package utils

import (
	"TicketSystem/config"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
)

//struct for mails
type Mail struct {
	ID                   int       `xml:"ID"`
	Mail                 string    `xml:"EMailAddress"`
	Subject              string    `xml:"Subject"`
	Message              string    `xml:"Message"`
	ReadAttemptCounter   int       `xml:"ReadAttempt"`
	FirstReadAttemptDate time.Time `xml:"FirstReadAttemptDate"`
}

type MailList struct {
	MailIDCounter int    `xml:"MailIDCounter"`
	MailList      []Mail `xml:"mails>mail"`
}

//creating or merging a ticket that was send by mail
func CreateTicketFromMail(mail string, reference string, message string) (Ticket, error) {
	tickets := GetTicketsByClient(mail)

	for _, actTicket := range tickets {
		if CheckStringsDeviation(2, strings.ToLower(actTicket.Reference), strings.ToLower(reference)) {
			newTicket, err := AddMessage(actTicket, mail, message)
			if err != nil {
				return newTicket, err
			}
			if newTicket.Status == TicketStatusClosed {
				err = ChangeStatus(newTicket.ID, TicketStatusInProcess)
				newTicket.Status = TicketStatusInProcess
				return newTicket, err
			}
			return newTicket, nil
		}
	}

	return CreateTicket(mail, reference, message)
}

//delete all mails in the xml file which are already sent
func DeleteMails(mailIds []int) error {
	maillist, err := ReadMailsFile()
	if err != nil {
		return err
	}

	mailIdCounter := maillist.MailIDCounter

	mailMap := make(map[int]Mail)
	for _, actMail := range maillist.MailList {
		mailMap[actMail.ID] = actMail
	}

	for _, actId := range mailIds {
		delete(mailMap, actId)
	}

	var newMaillist MailList
	for _, actMail := range mailMap {
		newMaillist.MailList = append(newMaillist.MailList, actMail)
	}

	newMaillist.MailIDCounter = mailIdCounter
	return WriteToXML(newMaillist, config.MailFilePath())
}

//store the message as a mail in the specific xml file
func SendMail(mail string, caption string, message string) error {
	maillist, err := ReadMailsFile()
	if err != nil {
		return err
	}
	nextMailId := maillist.MailIDCounter + 1
	newMail := Mail{Mail: mail, Subject: caption, Message: message, ID: nextMailId}
	maillist.MailList = append(maillist.MailList, newMail)
	maillist.MailIDCounter = nextMailId
	return WriteToXML(maillist, config.MailFilePath())
}

//get all mails from the xml file
func ReadMailsFile() (MailList, error) {
	file, err := ioutil.ReadFile(config.MailFilePath())
	if err != nil {
		return MailList{}, err
	}
	var maillist MailList
	err = xml.Unmarshal(file, &maillist)
	if err != nil {
		return MailList{}, err
	}
	return maillist, nil
}

func (m *Mail) IncrementReadAttemptsCounter() error {
	maillist, err := ReadMailsFile()
	if err != nil {
		return err
	}

	for i, mail := range maillist.MailList {
		if mail.ID == m.ID {
			maillist.MailList[i].ReadAttemptCounter = mail.ReadAttemptCounter + 1
			m.ReadAttemptCounter = maillist.MailList[i].ReadAttemptCounter
			if maillist.MailList[i].ReadAttemptCounter == 1 {
				maillist.MailList[i].FirstReadAttemptDate = time.Now()
				m.FirstReadAttemptDate = maillist.MailList[i].FirstReadAttemptDate
			}
			return WriteToXML(maillist, config.MailFilePath())
		}
	}
	return fmt.Errorf("couldn't find the email")
}

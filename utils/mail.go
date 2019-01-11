package utils

import (
	"TicketSystem/config"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

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

var mutexMailID = &sync.Mutex{}

// Creates or merges a ticket that was sent through the REST API (PUSH /mails)
func CreateTicketFromMail(mail string, reference string, message string) (Ticket, error) {
	tickets := GetTicketsByClient(mail)

	// Check if the ticket is referring to an existing ticket
	for _, actTicket := range tickets {
		if CheckStringsDeviation(2, strings.ToLower(actTicket.Reference), strings.ToLower(reference)) {
			newTicket, err := AddMessage(actTicket, mail, message)
			if err != nil {
				return newTicket, err
			}
			// Reopen closed tickets
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

// Deletes all mails in the xml file which are already sent
func DeleteMails(mailIds []int) error {
	// Synchronizing the change of the mail ID counter
	mutexMailID.Lock()
	defer mutexMailID.Unlock()

	mailList, err := ReadMailsFile()
	if err != nil {
		return err
	}

	mailIdCounter := mailList.MailIDCounter

	mailMap := make(map[int]Mail)
	for _, actMail := range mailList.MailList {
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

// Stores the input as a mail which needs to be sent
func SendMail(mail string, caption string, message string) error {
	// Synchronizing the change of the mail ID counter
	mutexMailID.Lock()
	defer mutexMailID.Unlock()

	mailList, err := ReadMailsFile()
	if err != nil {
		return err
	}

	nextMailId := mailList.MailIDCounter + 1
	newMail := Mail{Mail: mail, Subject: caption, Message: message, ID: nextMailId}
	mailList.MailList = append(mailList.MailList, newMail)
	mailList.MailIDCounter = nextMailId

	return WriteToXML(mailList, config.MailFilePath())
}

// Returns all mails which have to be sent
func ReadMailsFile() (MailList, error) {
	file, err := ioutil.ReadFile(config.MailFilePath())
	if err != nil {
		return MailList{}, err
	}

	var mailList MailList
	err = xml.Unmarshal(file, &mailList)
	if err != nil {
		return MailList{}, err
	}

	return mailList, nil
}

func (m *Mail) IncrementReadAttemptsCounter() error {
	mailList, err := ReadMailsFile()
	if err != nil {
		return err
	}

	for i, mail := range mailList.MailList {
		if mail.ID == m.ID {
			mailList.MailList[i].ReadAttemptCounter = mail.ReadAttemptCounter + 1
			m.ReadAttemptCounter = mailList.MailList[i].ReadAttemptCounter

			// Store the first read attempt only once
			if mailList.MailList[i].ReadAttemptCounter == 1 {
				mailList.MailList[i].FirstReadAttemptDate = time.Now()
				m.FirstReadAttemptDate = mailList.MailList[i].FirstReadAttemptDate
			}
			return WriteToXML(mailList, config.MailFilePath())
		}
	}
	return fmt.Errorf("couldn't find the email")
}

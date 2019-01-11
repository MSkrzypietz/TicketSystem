package main

import (
	"TicketSystem/utils"
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	url := flag.String("url", "https://localhost:4443", "URL of Website (root)")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Would you like to send a message to the ticket system? (Y/N)")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		input = strings.TrimSpace(input)

		if len(input) != 1 || !strings.ContainsAny(input, "ynYN") {
			fmt.Println("Invalid Input. Please try it again.")
			continue
		}

		if strings.ToUpper(input) == strings.ToUpper("n") {
			fmt.Println("Bye... :)")
			break
		}

		hasText := func(input string) bool {
			return len(strings.TrimSpace(input)) > 0
		}
		emailAddress := readInput(reader, "email address", utils.CheckMailFormal)
		subject := readInput(reader, "subject", hasText)
		message := readInput(reader, "message", hasText)

		if email, err := pushEmail(*url, emailAddress, subject, message); err == nil {
			fmt.Println("Successfully pushed the email with the following content:")
			fmt.Printf("\tE-Mail: %s\n", email.EMailAddress)
			fmt.Printf("\tSubject: %s\n", email.Subject)
			fmt.Printf("\tMessage: %s\n\n", email.Message)
		} else {
			fmt.Printf("We encountered the error '%v'. Please try it again!\n", err)
		}
	}
}

func readInput(reader *bufio.Reader, name string, checkInput func(string) bool) string {
	for {
		fmt.Printf("Enter %s: ", name)

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}

		input = strings.TrimSpace(input)
		if checkInput(input) {
			return input
		}

		fmt.Println("Invalid Input. Please try it again.")
	}
}

func pushEmail(url, emailAddress, subject, message string) (utils.MailData, error) {
	req := utils.Request{Mail: utils.MailData{EMailAddress: emailAddress, Subject: subject, Message: message}}
	buf, err := xml.Marshal(req)
	if err != nil {
		return utils.MailData{}, err
	}

	// Ignoring unauthorized certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Post(url+"/mails", "application/xml", bytes.NewBuffer(buf))
	if err != nil {
		return utils.MailData{}, err
	}
	defer res.Body.Close()

	return req.Mail, nil
}

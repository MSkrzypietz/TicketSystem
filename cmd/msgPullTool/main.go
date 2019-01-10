package main

import (
	"TicketSystem/utils"
	"bufio"
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
		fmt.Println("Would you like to pull the unsent emails of the ticket system? (Y/N)")
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

		mails, err := pullEmails(*url)
		if len(mails) == 1 {
			fmt.Printf("There is %d E-Mail to be sent:\n", len(mails))
		} else {
			fmt.Printf("There are %d E-Mails to be sent:\n", len(mails))
		}

		for _, mail := range mails {
			fmt.Println()
			fmt.Printf("E-Mail: %s\n", mail.EMailAddress)
			fmt.Printf("Subject: %s\n", mail.Subject)
			fmt.Printf("Message: %s\n", mail.Message)
		}
	}
}

func pullEmails(url string) ([]utils.MailData, error) {
	// Ignoring unauthorized certificate
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get(url + "/mails")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var mails utils.Response
	err = xml.NewDecoder(res.Body).Decode(&mails)
	if err != nil {
		return []utils.MailData{}, err
	}

	return mails.Data, nil
}

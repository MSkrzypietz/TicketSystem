package XML_IO

import (
	"encoding/xml"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

//struct that defines a ticket with the parameters ID, mail of the client, reference, actual status, editor and a list of all messages
type Ticket struct {
	XMLName     xml.Name  `xml:"Ticket"`
	Id          int       `xml:"ID"`
	Client      string    `xml:"Client"`
	Reference   string    `xml:"Reference"`
	Status      int       `xml:"Status"`
	Editor      int       `xml:"Editor"`
	MessageList []Message `xml:"MessageList>Message"`
}

//strcut that defines a message with the parameters date of creation, name or mail of the actor and the text of the message
type Message struct {
	CreationDate time.Time `xml:"CreationDate"`
	Actor        string    `xml:"Actor"`
	Text         string    `xml:Text`
}

var ticketMap map[int]Ticket = make(map[int]Ticket)

//function to create a ticket including the following parameters: mail of the client, reference and text of the ticket. Returns a bool whether the creation was successful.
func createTicket(client string, reference string, text string) bool {
	IDCounter := getTicketIDCounter() + 1
	newTicket := Ticket{Id: IDCounter, Client: client, Reference: reference, Status: 0, Editor: 0}
	return addMessage(newTicket, client, text) && writeToXML(IDCounter, "definitions")
}

//adds a message to a specified tickets. Functions includes the following parameters: specified ticket, the actor and the text of the message. Returns a bool whether it was successful.
func addMessage(ticket Ticket, actor string, text string) bool {
	newMessage := Message{CreationDate: time.Now(), Actor: actor, Text: text}
	ticket.MessageList = append(ticket.MessageList, newMessage)
	return storeTicket(ticket)
}

//stores a ticket to the cache (if there are too many tickets in the cache one will be written to the XML-file)
func storeTicket(ticket Ticket) bool {
	tmpBool := checkCache()
	ticketMap[ticket.Id] = ticket
	return tmpBool
}

//reads a specified ticket from the XML-file or from the cache. Function has as the parameter the ticket ID and returns the ticket
func readTicket(id int) Ticket {
	//returns the ticket from the cache if it is stored in there
	tempTicket := ticketMap[id]
	if tempTicket.Id != 0 {
		return tempTicket
	}

	//returns ticket from the XML-file and stores it to the cache
	file, err := ioutil.ReadFile("tickets/ticket" + strconv.Itoa(id) + ".xml")
	if err != nil {
		return Ticket{}
	}
	var ticket Ticket
	xml.Unmarshal(file, &ticket)
	checkCache()
	ticketMap[ticket.Id] = ticket
	return ticket
}

//deletes a ticket by its ID
func deleteTicket(id int) bool {
	delete(ticketMap, id)
	if id == getTicketIDCounter() {
		writeToXML(id-1, "definitions")
	}
	err := os.Remove("tickets/ticket" + strconv.Itoa(id) + ".xml")
	if err != nil {
		return false
	}
	return true
}

//changes the editor of a ticket
func changeEditor(id int, editor int) {
	ticket := readTicket(id)
	ticket.Editor = editor
	storeTicket(ticket)
}

//changes the status of a ticket
func changeStatus(id int, status int) {
	ticket := readTicket(id)
	ticket.Status = status
	storeTicket(ticket)
}

//returns a list of tickets by a specified ticket status. Status is specified in the parameters of the function.
func getTicketsByStatus(status int) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(); actualID++ {
		tmp := readTicket(actualID)
		if tmp.Status == status && tmp.Id != 0 {
			tickets = append(tickets, tmp)
		}
	}
	return tickets
}

//returns a list of tickets owned by one editor who is specified in the parameters of the function
func getTicketsByEditor(editor int) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(); actualID++ {
		tmp := readTicket(actualID)
		if tmp.Editor == editor && tmp.Id != 0 {
			tickets = append(tickets, tmp)
		}
	}
	return tickets
}

//returns the actual ticket ID in order to create a new ticket or to get to know the number of the stored tickets.
func getTicketIDCounter() int {
	file, err := ioutil.ReadFile("definitions.xml")
	if err != nil {
		panic(err)
	}
	var IDCounter int
	xml.Unmarshal(file, &IDCounter)
	return IDCounter
}

//merge two tickets, store them as one ticket and delete the other one
func mergeTickets(firstTicketID int, secondTicketID int) bool {
	firstTicket := readTicket(firstTicketID)
	secondTicket := readTicket(secondTicketID)
	if firstTicket.Editor != secondTicket.Editor {
		return false
	}
	for e := range secondTicket.MessageList {
		firstTicket.MessageList = append(firstTicket.MessageList, secondTicket.MessageList[e])
	}
	deleteTicket(secondTicketID)
	return storeTicket(firstTicket)
}

//functions writes an object to an specified xml file and returns a bool whether the writing was successful
func writeToXML(v interface{}, file string) bool {
	if xmlstring, err := xml.MarshalIndent(v, "", "    "); err == nil {
		xmlstring = []byte(xml.Header + string(xmlstring))
		err = ioutil.WriteFile(file+".xml", xmlstring, 0644)
		if err != nil {
			panic(err)
		}
		return true
	}
	return false
}

//function clears the cache
func clearCache() bool {
	tmpBool := true
	for e := range ticketMap {
		tmpBool = tmpBool && writeToXML(ticketMap[e], "tickets/ticket"+strconv.Itoa(ticketMap[e].Id))
		delete(ticketMap, e)
	}
	return tmpBool
}

//function checks if there are too many tickets in the cache and in the case of too many tickets one will be written to the XML-file.
func checkCache() bool {
	if len(ticketMap) > 10 {
		randNumber := rand.Intn(len(ticketMap))
		tmpInt := 1
		for e := range ticketMap {
			if tmpInt == randNumber {
				tmpTicket := ticketMap[e]
				delete(ticketMap, e)
				return writeToXML(tmpTicket, "tickets/ticket"+strconv.Itoa(e))
			}
			tmpInt++
		}
		return false
	} else {
		return true
	}
}

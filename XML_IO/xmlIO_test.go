package XML_IO

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strconv"
	"testing"
)

func TestTicketCreation(t *testing.T) {
	assert := assert.New(t)
	boolTicket := createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.True(boolTicket)

	ticketID := getTicketIDCounter()
	actTicket := ticketMap[ticketID]

	var expectedMsg []Message
	expectedMsg = append(expectedMsg, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Status: 0, Editor: 0, MessageList: expectedMsg}
	assert.Equal(expectedTicket, actTicket)

	_, err := ioutil.ReadFile("tickets/ticket" + strconv.Itoa(ticketID) + ".xml")
	assert.NotNil(err)

	deleteTicket(ticketID)
	clearCache()
	writeToXML(0, "definitions")
}

func TestAddMessage(t *testing.T) {
	assert := assert.New(t)
	boolTicket := createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.True(boolTicket)

	ticketID := getTicketIDCounter()
	addMessage(readTicket(ticketID), "4262", "please restart")
	addMessage(readTicket(ticketID), "client@dhbw.de", "Thank, it worked!")
	expectedMsgOne := Message{Actor: "4262", Text: "please restart"}
	expectedMsgTwo := Message{Actor: "client@dhbw.de", Text: "Thank, it worked!"}

	msgList := readTicket(ticketID).MessageList
	assert.Equal("client@dhbw.de", msgList[0].Actor)
	assert.Equal("PC does not start anymore. Any idea?", msgList[0].Text)
	assert.Equal(expectedMsgOne.Actor, msgList[1].Actor)
	assert.Equal(expectedMsgOne.Text, msgList[1].Text)
	assert.Equal(expectedMsgTwo.Actor, msgList[2].Actor)
	assert.Equal(expectedMsgTwo.Text, msgList[2].Text)

	deleteTicket(ticketID)
	clearCache()
}

func TestTicketStoring(t *testing.T) {
	assert := assert.New(t)
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		boolTicket := createTicket("client"+strconv.Itoa(tmpInt)+"@dhbw.de", "PC problem", "Pc does not start anymore")
		assert.True(boolTicket)
	}
	fmt.Println(ticketMap)
	actTicket := ticketMap[4]
	var expectedMsg []Message
	expectedMsg = append(expectedMsg, Message{Actor: "client4@dhbw.de", Text: "Pc does not start anymore", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicketFour := Ticket{XMLName: xml.Name{"", ""}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: 0, MessageList: expectedMsg}
	assert.Equal(expectedTicketFour, actTicket)

	_, err := ioutil.ReadFile("tickets/ticket4.xml")
	assert.NotNil(err)

	createTicket("client10@dhbw.de", "PC problem", "Pc does not start anymore")
	createTicket("client11@dhbw.de", "PC problem", "Pc does not start anymore")
	createTicket("client12@dhbw.de", "PC problem", "Pc does not start anymore")
	assert.Equal(11, len(ticketMap))

	clearCache()
	assert.Equal(0, len(ticketMap))
	actTicket = readTicket(4)
	expectedMsg = nil
	expectedMsg = append(expectedMsg, Message{Actor: "client4@dhbw.de", Text: "Pc does not start anymore", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicketFour = Ticket{XMLName: xml.Name{"", "Ticket"}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: 0, MessageList: expectedMsg}
	assert.Equal(expectedTicketFour, actTicket)

	clearCache()
	for tmpInt := 1; tmpInt <= 12; tmpInt++ {
		deleteTicket(tmpInt)
	}
	writeToXML(0, "definitions")
}

func TestTicketReading(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ticketID := getTicketIDCounter()
	actTicket := readTicket(ticketID)

	var msgList []Message
	msgList = append(msgList, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Status: 0, MessageList: msgList}
	assert.Equal(expectedTicket, actTicket)

	clearCache()
	actTicket = readTicket(ticketID)
	msgList = nil
	msgList = append(msgList, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket = Ticket{XMLName: xml.Name{"", "Ticket"}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Status: 0, MessageList: msgList}
	assert.Equal(expectedTicket, actTicket)

	clearCache()
	deleteTicket(ticketID)
	writeToXML(0, "definitions")
}

func TestIDCounter(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(2, getTicketIDCounter())

	clearCache()
	deleteTicket(getTicketIDCounter())
	deleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestTicketsByStatus(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	changeStatus(getTicketIDCounter(), 1)

	tickets := getTicketsByStatus(0)
	for _, element := range tickets {
		assert.Equal(0, element.Status)
	}

	tickets = getTicketsByStatus(1)
	for _, element := range tickets {
		assert.Equal(1, element.Status)
	}

	clearCache()
	deleteTicket(getTicketIDCounter())
	deleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestTicketByEditor(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	changeEditor(getTicketIDCounter()-1, 423)
	changeEditor(getTicketIDCounter(), 22)

	tickets := getTicketsByEditor(423)
	for _, element := range tickets {
		assert.Equal(423, element.Editor)
	}
	tickets = getTicketsByEditor(22)
	for _, element := range tickets {
		assert.Equal(22, element.Editor)
	}

	clearCache()
	deleteTicket(getTicketIDCounter())
	deleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestChangeEditor(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	changeEditor(getTicketIDCounter(), 4321)
	ticket := readTicket(getTicketIDCounter())
	assert.Equal(4321, ticket.Editor)

	clearCache()
	deleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestChangeStatus(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	changeStatus(getTicketIDCounter(), 2)
	ticket := readTicket(getTicketIDCounter())
	assert.Equal(2, ticket.Status)

	clearCache()
	deleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestDeleting(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "Computer", "PC not working")
	assert.Equal(1, len(ticketMap))
	deleteTicket(getTicketIDCounter())
	assert.Equal(0, len(ticketMap))
	//assert.Equal(0, getTicketIDCounter())

	createTicket("client@dhbw.de", "Computer", "PC not working")
	clearCache()
	assert.Equal(0, len(ticketMap))
	_, err := ioutil.ReadFile("tickets/ticket1.xml")
	assert.Nil(err)
	deleteTicket(1)
	_, err = ioutil.ReadFile("tickets/ticket1.xml")
	assert.NotNil(err)

	clearCache()
	writeToXML(0, "definitions")
}

func TestMerging(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Max Mustermann. Thanks.")
	createTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	firstTicket := readTicket(getTicketIDCounter() - 1)
	secondTicket := readTicket(getTicketIDCounter())
	changeStatus(firstTicket.Id, 1)
	changeStatus(secondTicket.Id, 1)
	changeEditor(firstTicket.Id, 202)
	changeEditor(secondTicket.Id, 202)
	firstTicket = readTicket(getTicketIDCounter() - 1)
	secondTicket = readTicket(getTicketIDCounter())

	//merge two tickets and test the function
	var msgList []Message
	msgList = firstTicket.MessageList
	for e := range secondTicket.MessageList {
		msgList = append(msgList, secondTicket.MessageList[e])
	}
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: firstTicket.Id, Client: firstTicket.Client, Reference: firstTicket.Reference, Status: firstTicket.Status, Editor: firstTicket.Editor, MessageList: msgList}

	boolMerged := mergeTickets(firstTicket.Id, secondTicket.Id)
	assert.True(boolMerged)
	assert.Equal(expectedTicket, readTicket(firstTicket.Id))

	//merge tickets with two different editors
	createTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	secondTicketID := getTicketIDCounter()
	changeEditor(secondTicketID, 412)
	assert.False(mergeTickets(firstTicket.Id, secondTicketID))

	clearCache()
	deleteTicket(firstTicket.Id)
	deleteTicket(secondTicket.Id)
	writeToXML(0, "definitions")
}

func TestClearCache(t *testing.T) {
	assert := assert.New(t)
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(3, len(ticketMap))

	clearCache()
	assert.Equal(0, len(ticketMap))

	_, err1 := ioutil.ReadFile("tickets/ticket1.xml")
	assert.Nil(err1)

	deleteTicket(1)
	deleteTicket(2)
	deleteTicket(3)
	clearCache()
	writeToXML(0, "definitions")
}

func TestCheckCache(t *testing.T) {
	assert := assert.New(t)
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	}

	assert.Equal(9, len(ticketMap))

	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(11, len(ticketMap))

	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	createTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(11, len(ticketMap))

	clearCache()
	for tmpInt := 1; tmpInt <= 13; tmpInt++ {
		deleteTicket(tmpInt)
	}
	writeToXML(0, "definitions")
}

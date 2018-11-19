package XML_IO

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestTicketCreation(t *testing.T) {
	boolTicket := CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.True(t, boolTicket)

	ticketID := getTicketIDCounter("definitions.xml")
	actTicket := ticketMap[ticketID]

	var expectedMsg []Message
	expectedMsg = append(expectedMsg, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0", MessageList: expectedMsg}
	assert.Equal(t, expectedTicket, actTicket)

	ClearCache("../data/tickets/ticket")

	f, err := ioutil.ReadFile("../data/tickets/ticket" + strconv.Itoa(ticketID) + ".xml")
	assert.NotNil(t, f)
	assert.Nil(t, err)

	DeleteTicket("../data/tickets/ticket", "definitions.xml", ticketID)
	ClearCache("../data/tickets/ticket")
	writeToXML(0, "definitions.xml")
}

func TestAddMessage(t *testing.T) {
	boolTicket := CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.True(t, boolTicket)

	ticketID := getTicketIDCounter("definitions.xml")
	AddMessage("../data/tickets/ticket", ReadTicket("../data/tickets/ticket", ticketID), "4262", "please restart")
	AddMessage("../data/tickets/ticket", ReadTicket("../data/tickets/ticket", ticketID), "client@dhbw.de", "Thank, it worked!")
	expectedMsgOne := Message{Actor: "4262", Text: "please restart"}
	expectedMsgTwo := Message{Actor: "client@dhbw.de", Text: "Thank, it worked!"}

	msgList := ReadTicket("../data/tickets/ticket", ticketID).MessageList
	assert.Equal(t, "client@dhbw.de", msgList[0].Actor)
	assert.Equal(t, "PC does not start anymore. Any idea?", msgList[0].Text)
	assert.Equal(t, expectedMsgOne.Actor, msgList[1].Actor)
	assert.Equal(t, expectedMsgOne.Text, msgList[1].Text)
	assert.Equal(t, expectedMsgTwo.Actor, msgList[2].Actor)
	assert.Equal(t, expectedMsgTwo.Text, msgList[2].Text)

	DeleteTicket("../data/tickets/ticket", "definitions.xml", ticketID)
	ClearCache("../data/tickets/ticket")
}

func TestTicketStoring(t *testing.T) {
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		boolTicket := CreateTicket("../data/tickets/ticket", "definitions.xml", "client"+strconv.Itoa(tmpInt)+"@dhbw.de", "PC problem", "Pc does not start anymore")
		assert.True(t, boolTicket)
	}

	actTicket := ticketMap[4]
	var expectedMsg []Message
	expectedMsg = append(expectedMsg, Message{Actor: "client4@dhbw.de", Text: "Pc does not start anymore", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicketFour := Ticket{XMLName: xml.Name{"", ""}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0", MessageList: expectedMsg}
	assert.Equal(t, expectedTicketFour, actTicket)

	_, err := ioutil.ReadFile("../data/tickets/ticket4.xml")
	assert.NotNil(t, err)

	CreateTicket("../data/tickets/ticket", "definitions.xml", "client10@dhbw.de", "PC problem", "Pc does not start anymore")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client11@dhbw.de", "PC problem", "Pc does not start anymore")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client12@dhbw.de", "PC problem", "Pc does not start anymore")
	assert.Equal(t, 11, len(ticketMap))

	ClearCache("../data/tickets/ticket")
	assert.Equal(t, 0, len(ticketMap))
	actTicket = ReadTicket("../data/tickets/ticket", 4)
	expectedMsg = nil
	expectedMsg = append(expectedMsg, Message{Actor: "client4@dhbw.de", Text: "Pc does not start anymore", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicketFour = Ticket{XMLName: xml.Name{"", "Ticket"}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0", MessageList: expectedMsg}
	assert.Equal(t, expectedTicketFour, actTicket)

	ClearCache("../data/tickets/ticket")
	for tmpInt := 1; tmpInt <= 12; tmpInt++ {
		DeleteTicket("../data/tickets/ticket", "definitions.xml", tmpInt)
	}
	writeToXML(0, "definitions.xml")
}

func TestTicketReading(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ticketID := getTicketIDCounter("definitions")
	actTicket := ReadTicket("../data/tickets/ticket", ticketID)

	var msgList []Message
	msgList = append(msgList, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Editor: "0", Status: 0, MessageList: msgList}
	assert.Equal(t, expectedTicket, actTicket)

	ClearCache("../data/tickets/ticket")
	actTicket = ReadTicket("../data/tickets/ticket", ticketID)
	msgList = nil
	msgList = append(msgList, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket = Ticket{XMLName: xml.Name{"", "Ticket"}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Editor: "0", Status: 0, MessageList: msgList}
	assert.Equal(t, expectedTicket, actTicket)

	assert.Equal(t, Ticket{}, ReadTicket("../data/tickets/ticket", -99))

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions", ticketID)
	writeToXML(0, "definitions")
}

func TestIDCounter(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 2, getTicketIDCounter("definitions"))

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	writeToXML(0, "definitions")
}

func TestTicketsByStatus(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus("../data/tickets/ticket", getTicketIDCounter("definitions"), 1)

	tickets := GetTicketsByStatus("../data/tickets/ticket", "definitions", 0)
	for _, element := range tickets {
		assert.Equal(t, 0, element.Status)
	}

	tickets = GetTicketsByStatus("../data/tickets/ticket", "definitions", 1)
	for _, element := range tickets {
		assert.Equal(t, 1, element.Status)
	}

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	writeToXML(0, "definitions")
}

func TestTicketByEditor(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeEditor("../data/tickets/ticket", getTicketIDCounter("definitions")-1, "423")
	ChangeEditor("../data/tickets/ticket", getTicketIDCounter("definitions"), "22")

	tickets := GetTicketsByEditor("../data/tickets/ticket", "definitions", "423")
	for _, element := range tickets {
		assert.Equal(t, "423", element.Editor)
	}
	tickets = GetTicketsByEditor("../data/tickets/ticket", "definitions", "22")
	for _, element := range tickets {
		assert.Equal(t, "22", element.Editor)
	}

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	writeToXML(0, "definitions")
}

func TestChangeEditor(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeEditor("../data/tickets/ticket", getTicketIDCounter("definitions"), "4321")
	ticket := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions"))
	assert.Equal(t, "4321", ticket.Editor)

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	writeToXML(0, "definitions")
}

func TestChangeStatus(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus("../data/tickets/ticket", getTicketIDCounter("definitions"), 2)
	ticket := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions"))
	assert.Equal(t, 2, ticket.Status)

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	writeToXML(0, "definitions")
}

func TestDeleting(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "Computer", "PC not working")
	assert.Equal(t, 1, len(ticketMap))
	DeleteTicket("../data/tickets/ticket", "definitions", getTicketIDCounter("definitions"))
	assert.Equal(t, 0, len(ticketMap))
	//assert.Equal(0, getTicketIDCounter())

	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "Computer", "PC not working")
	ClearCache("../data/tickets/ticket")
	assert.Equal(t, 0, len(ticketMap))
	_, err := ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.Nil(t, err)
	DeleteTicket("../data/tickets/ticket", "definitions", 1)
	_, err = ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.NotNil(t, err)

	ClearCache("../data/tickets/ticket")
	writeToXML(0, "definitions")
}

func TestMerging(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Max Mustermann. Thanks.")
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	firstTicket := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions")-1)
	secondTicket := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions"))
	ChangeStatus("../data/tickets/ticket", firstTicket.Id, 1)
	ChangeStatus("../data/tickets/ticket", secondTicket.Id, 1)
	ChangeEditor("../data/tickets/ticket", firstTicket.Id, "202")
	ChangeEditor("../data/tickets/ticket", secondTicket.Id, "202")
	firstTicket = ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions")-1)
	secondTicket = ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions"))

	//merge two tickets and test the function
	var msgList []Message
	msgList = firstTicket.MessageList
	for e := range secondTicket.MessageList {
		msgList = append(msgList, secondTicket.MessageList[e])
	}
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: firstTicket.Id, Client: firstTicket.Client, Reference: firstTicket.Reference, Status: firstTicket.Status, Editor: firstTicket.Editor, MessageList: msgList}

	boolMerged := MergeTickets("../data/tickets/ticket", "definitions", firstTicket.Id, secondTicket.Id)
	assert.True(t, boolMerged)
	assert.Equal(t, expectedTicket, ReadTicket("../data/tickets/ticket", firstTicket.Id))

	//merge tickets with two different editors
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	secondTicketID := getTicketIDCounter("definitions")
	ChangeEditor("../data/tickets/ticket", secondTicketID, "412")
	assert.False(t, MergeTickets("../data/tickets/ticket", "definitions", firstTicket.Id, secondTicketID))

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions", firstTicket.Id)
	DeleteTicket("../data/tickets/ticket", "definitions", secondTicket.Id)
	writeToXML(0, "definitions")
}

func TestClearCache(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 3, len(ticketMap))

	ClearCache("../data/tickets/ticket")
	assert.Equal(t, 0, len(ticketMap))

	_, err1 := ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.Nil(t, err1)

	DeleteTicket("../data/tickets/ticket", "definitions", 1)
	DeleteTicket("../data/tickets/ticket", "definitions", 2)
	DeleteTicket("../data/tickets/ticket", "definitions", 3)
	ClearCache("../data/tickets/ticket")
	writeToXML(0, "definitions")
}

func TestCheckCache(t *testing.T) {
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	}

	assert.Equal(t, 9, len(ticketMap))

	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 11, len(ticketMap))

	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 11, len(ticketMap))

	ClearCache("../data/tickets/ticket")
	for tmpInt := 1; tmpInt <= 13; tmpInt++ {
		DeleteTicket("../data/tickets/ticket", "definitions", tmpInt)
	}
	writeToXML(0, "definitions")
}

func TestCreateAndStoreUser(t *testing.T) {
	assert.True(t, CreateUser("../data/users/users.xml", "mustermann", "musterpasswort"))
	file, _ := ioutil.ReadFile("../data/users/users.xml")
	var userlist Userlist
	xml.Unmarshal(file, &userlist)

	var expectedUser []User
	expectedUser = append(expectedUser, User{Username: "mustermann", Password: "musterpasswort"})
	expected := Userlist{User: expectedUser}
	assert.Equal(t, expected, userlist)
}

func TestReadUser(t *testing.T) {
	assert.True(t, CreateUser("../data/users/users.xml", "mustermann", "musterpasswort"))
	expectedMap := make(map[string]User)
	expectedMap["mustermann"] = User{Username: "mustermann", Password: "musterpasswort"}
	assert.Equal(t, expectedMap, readUsers("../data/users/users"))
	os.Remove("../data/users/users")
	expectedMap = make(map[string]User)
	assert.Equal(t, expectedMap, readUsers("../data/users/users"))
	os.Create("../data/users/users")
}

func TestCheckUser(t *testing.T) {
	assert.True(t, CreateUser("../data/users/users", "mustermann", "musterpasswort"))
	assert.True(t, CheckUser("../data/users/users", "mustermann", "musterpasswort"))
	assert.False(t, CheckUser("../data/users/users", "mustermann", "falschespasswort"))
	assert.False(t, CheckUser("../data/users/users", "muster", "musterpasswort"))
}

func TestLoginUser(t *testing.T) {
	assert.True(t, CreateUser("../data/users/users", "mustermann", "musterpasswort"))
	assert.True(t, LoginUser("../data/users/users", "mustermann", "musterpasswort", "1234"))
	assert.False(t, LoginUser("../data/users/users", "mustermann", "falschespasswort", "1234"))
	usersMap := readUsers("../data/users/users")
	assert.Equal(t, "1234", usersMap["mustermann"].SessionID)
}

func TestLogoutUser(t *testing.T) {
	assert.True(t, CreateUser("../data/users/users", "mustermann", "musterpasswort"))
	assert.True(t, LoginUser("../data/users/users", "mustermann", "musterpasswort", "1234"))
	usersmap := readUsers("../data/users/users")
	assert.Equal(t, "1234", usersmap["mustermann"].SessionID)
	assert.True(t, LogoutUser("../data/users/users", "mustermann"))
	usersmap = readUsers("../data/users/users")
	assert.Equal(t, "", usersmap["mustermann"].SessionID)
	assert.False(t, LogoutUser("../data/users/users", "falscherName"))
}

func TestGetUserSession(t *testing.T) {
	assert.True(t, CreateUser("../data/users/users", "mustermann", "musterpasswort"))
	assert.Equal(t, "", GetUserSession("../data/users/users", "mustermann"))
	LoginUser("../data/users/users", "mustermann", "musterpasswort", "1234")
	assert.Equal(t, "1234", GetUserSession("../data/users/users", "mustermann"))
}

func TestGetUserBySession(t *testing.T) {
	assert.True(t, CreateUser("../data/users/users", "mustermann", "musterpasswort"))
	LoginUser("../data/users/users", "mustermann", "musterpasswort", "1234")
	expectedUser := User{Username: "mustermann", Password: "musterpasswort", SessionID: "1234"}
	assert.Equal(t, expectedUser, GetUserBySession("../data/users/users", "1234"))
	assert.Equal(t, User{}, GetUserBySession("../data/users/users", ""))
	assert.Equal(t, User{}, GetUserBySession("../data/users/users", "FalscheSession"))
}

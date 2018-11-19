package XML_IO

import (
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestTicketCreation(t *testing.T) {
	boolTicket := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.True(t, boolTicket)

	ticketID := getTicketIDCounter()
	actTicket := ticketMap[ticketID]

	var expectedMsg []Message
	expectedMsg = append(expectedMsg, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0", MessageList: expectedMsg}
	assert.Equal(t, expectedTicket, actTicket)

	_, err := ioutil.ReadFile("tickets/ticket" + strconv.Itoa(ticketID) + ".xml")
	assert.NotNil(t, err)

	DeleteTicket(ticketID)
	ClearCache()
	writeToXML(0, "definitions")
}

func TestAddMessage(t *testing.T) {
	boolTicket := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.True(t, boolTicket)

	ticketID := getTicketIDCounter()
	AddMessage(ReadTicket(ticketID), "4262", "please restart")
	AddMessage(ReadTicket(ticketID), "client@dhbw.de", "Thank, it worked!")
	expectedMsgOne := Message{Actor: "4262", Text: "please restart"}
	expectedMsgTwo := Message{Actor: "client@dhbw.de", Text: "Thank, it worked!"}

	msgList := ReadTicket(ticketID).MessageList
	assert.Equal(t, "client@dhbw.de", msgList[0].Actor)
	assert.Equal(t, "PC does not start anymore. Any idea?", msgList[0].Text)
	assert.Equal(t, expectedMsgOne.Actor, msgList[1].Actor)
	assert.Equal(t, expectedMsgOne.Text, msgList[1].Text)
	assert.Equal(t, expectedMsgTwo.Actor, msgList[2].Actor)
	assert.Equal(t, expectedMsgTwo.Text, msgList[2].Text)

	DeleteTicket(ticketID)
	ClearCache()
}

func TestTicketStoring(t *testing.T) {
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		boolTicket := CreateTicket("client"+strconv.Itoa(tmpInt)+"@dhbw.de", "PC problem", "Pc does not start anymore")
		assert.True(t, boolTicket)
	}
	fmt.Println(ticketMap)
	actTicket := ticketMap[4]
	var expectedMsg []Message
	expectedMsg = append(expectedMsg, Message{Actor: "client4@dhbw.de", Text: "Pc does not start anymore", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicketFour := Ticket{XMLName: xml.Name{"", ""}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0", MessageList: expectedMsg}
	assert.Equal(t, expectedTicketFour, actTicket)

	_, err := ioutil.ReadFile("tickets/ticket4.xml")
	assert.NotNil(t, err)

	CreateTicket("client10@dhbw.de", "PC problem", "Pc does not start anymore")
	CreateTicket("client11@dhbw.de", "PC problem", "Pc does not start anymore")
	CreateTicket("client12@dhbw.de", "PC problem", "Pc does not start anymore")
	assert.Equal(t, 11, len(ticketMap))

	ClearCache()
	assert.Equal(t, 0, len(ticketMap))
	actTicket = ReadTicket(4)
	expectedMsg = nil
	expectedMsg = append(expectedMsg, Message{Actor: "client4@dhbw.de", Text: "Pc does not start anymore", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicketFour = Ticket{XMLName: xml.Name{"", "Ticket"}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0", MessageList: expectedMsg}
	assert.Equal(t, expectedTicketFour, actTicket)

	ClearCache()
	for tmpInt := 1; tmpInt <= 12; tmpInt++ {
		DeleteTicket(tmpInt)
	}
	writeToXML(0, "definitions")
}

func TestTicketReading(t *testing.T) {
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ticketID := getTicketIDCounter()
	actTicket := ReadTicket(ticketID)

	var msgList []Message
	msgList = append(msgList, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Editor: "0", Status: 0, MessageList: msgList}
	assert.Equal(t, expectedTicket, actTicket)

	ClearCache()
	actTicket = ReadTicket(ticketID)
	msgList = nil
	msgList = append(msgList, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket = Ticket{XMLName: xml.Name{"", "Ticket"}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Editor: "0", Status: 0, MessageList: msgList}
	assert.Equal(t, expectedTicket, actTicket)

	assert.Equal(t, Ticket{}, ReadTicket(-99))

	ClearCache()
	DeleteTicket(ticketID)
	writeToXML(0, "definitions")
}

func TestIDCounter(t *testing.T) {
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 2, getTicketIDCounter())

	ClearCache()
	DeleteTicket(getTicketIDCounter())
	DeleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestTicketsByStatus(t *testing.T) {
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus(getTicketIDCounter(), 1)

	tickets := GetTicketsByStatus(0)
	for _, element := range tickets {
		assert.Equal(t, 0, element.Status)
	}

	tickets = GetTicketsByStatus(1)
	for _, element := range tickets {
		assert.Equal(t, 1, element.Status)
	}

	ClearCache()
	DeleteTicket(getTicketIDCounter())
	DeleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestTicketByEditor(t *testing.T) {
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeEditor(getTicketIDCounter()-1, "423")
	ChangeEditor(getTicketIDCounter(), "22")

	tickets := GetTicketsByEditor("423")
	for _, element := range tickets {
		assert.Equal(t, "423", element.Editor)
	}
	tickets = GetTicketsByEditor("22")
	for _, element := range tickets {
		assert.Equal(t, "22", element.Editor)
	}

	ClearCache()
	DeleteTicket(getTicketIDCounter())
	DeleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestChangeEditor(t *testing.T) {
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeEditor(getTicketIDCounter(), "4321")
	ticket := ReadTicket(getTicketIDCounter())
	assert.Equal(t, "4321", ticket.Editor)

	ClearCache()
	DeleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestChangeStatus(t *testing.T) {
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus(getTicketIDCounter(), 2)
	ticket := ReadTicket(getTicketIDCounter())
	assert.Equal(t, 2, ticket.Status)

	ClearCache()
	DeleteTicket(getTicketIDCounter())
	writeToXML(0, "definitions")
}

func TestDeleting(t *testing.T) {
	CreateTicket("client@dhbw.de", "Computer", "PC not working")
	assert.Equal(t, 1, len(ticketMap))
	DeleteTicket(getTicketIDCounter())
	assert.Equal(t, 0, len(ticketMap))
	//assert.Equal(0, getTicketIDCounter())

	CreateTicket("client@dhbw.de", "Computer", "PC not working")
	ClearCache()
	assert.Equal(t, 0, len(ticketMap))
	_, err := ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.Nil(t, err)
	DeleteTicket(1)
	_, err = ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.NotNil(t, err)

	ClearCache()
	writeToXML(0, "definitions")
}

func TestMerging(t *testing.T) {
	CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Max Mustermann. Thanks.")
	CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	firstTicket := ReadTicket(getTicketIDCounter() - 1)
	secondTicket := ReadTicket(getTicketIDCounter())
	ChangeStatus(firstTicket.Id, 1)
	ChangeStatus(secondTicket.Id, 1)
	ChangeEditor(firstTicket.Id, "202")
	ChangeEditor(secondTicket.Id, "202")
	firstTicket = ReadTicket(getTicketIDCounter() - 1)
	secondTicket = ReadTicket(getTicketIDCounter())

	//merge two tickets and test the function
	var msgList []Message
	msgList = firstTicket.MessageList
	for e := range secondTicket.MessageList {
		msgList = append(msgList, secondTicket.MessageList[e])
	}
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: firstTicket.Id, Client: firstTicket.Client, Reference: firstTicket.Reference, Status: firstTicket.Status, Editor: firstTicket.Editor, MessageList: msgList}

	boolMerged := MergeTickets(firstTicket.Id, secondTicket.Id)
	assert.True(t, boolMerged)
	assert.Equal(t, expectedTicket, ReadTicket(firstTicket.Id))

	//merge tickets with two different editors
	CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	secondTicketID := getTicketIDCounter()
	ChangeEditor(secondTicketID, "412")
	assert.False(t, MergeTickets(firstTicket.Id, secondTicketID))

	ClearCache()
	DeleteTicket(firstTicket.Id)
	DeleteTicket(secondTicket.Id)
	writeToXML(0, "definitions")
}

func TestClearCache(t *testing.T) {
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 3, len(ticketMap))

	ClearCache()
	assert.Equal(t, 0, len(ticketMap))

	_, err1 := ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.Nil(t, err1)

	DeleteTicket(1)
	DeleteTicket(2)
	DeleteTicket(3)
	ClearCache()
	writeToXML(0, "definitions")
}

func TestCheckCache(t *testing.T) {
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	}

	assert.Equal(t, 9, len(ticketMap))

	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 11, len(ticketMap))

	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 11, len(ticketMap))

	ClearCache()
	for tmpInt := 1; tmpInt <= 13; tmpInt++ {
		DeleteTicket(tmpInt)
	}
	writeToXML(0, "definitions")
}

func TestCreateAndStoreUser(t *testing.T) {
	assert.True(t, CreateUser("mustermann", "musterpasswort"))
	file, _ := ioutil.ReadFile("../data/users/users.xml")
	var userlist Userlist
	xml.Unmarshal(file, &userlist)

	var expectedUser []User
	expectedUser = append(expectedUser, User{Username: "mustermann", Password: "musterpasswort"})
	expected := Userlist{User: expectedUser}
	assert.Equal(t, expected, userlist)
}

func TestReadUser(t *testing.T) {
	assert.True(t, CreateUser("mustermann", "musterpasswort"))
	expectedMap := make(map[string]User)
	expectedMap["mustermann"] = User{Username: "mustermann", Password: "musterpasswort"}
	assert.Equal(t, expectedMap, readUsers())
	os.Remove("../data/users/users.xml")
	expectedMap = make(map[string]User)
	assert.Equal(t, expectedMap, readUsers())
	os.Create("../data/users/users.xml")
}

func TestCheckUser(t *testing.T) {
	assert.True(t, CreateUser("mustermann", "musterpasswort"))
	assert.True(t, CheckUser("mustermann", "musterpasswort"))
	assert.False(t, CheckUser("mustermann", "falschespasswort"))
	assert.False(t, CheckUser("muster", "musterpasswort"))
}

func TestLoginUser(t *testing.T) {
	assert.True(t, CreateUser("mustermann", "musterpasswort"))
	assert.True(t, LoginUser("mustermann", "musterpasswort", "1234"))
	assert.False(t, LoginUser("mustermann", "falschespasswort", "1234"))
	usersMap := readUsers()
	assert.Equal(t, "1234", usersMap["mustermann"].SessionID)
}

func TestLogoutUser(t *testing.T) {
	assert.True(t, CreateUser("mustermann", "musterpasswort"))
	assert.True(t, LoginUser("mustermann", "musterpasswort", "1234"))
	usersmap := readUsers()
	assert.Equal(t, "1234", usersmap["mustermann"].SessionID)
	assert.True(t, LogoutUser("mustermann"))
	usersmap = readUsers()
	assert.Equal(t, "", usersmap["mustermann"].SessionID)
	assert.False(t, LogoutUser("falscherName"))
}

func TestGetUserSession(t *testing.T) {
	assert.True(t, CreateUser("mustermann", "musterpasswort"))
	assert.Equal(t, "", GetUserSession("mustermann"))
	LoginUser("mustermann", "musterpasswort", "1234")
	assert.Equal(t, "1234", GetUserSession("mustermann"))
}

func TestGetUserBySession(t *testing.T) {
	assert.True(t, CreateUser("mustermann", "musterpasswort"))
	LoginUser("mustermann", "musterpasswort", "1234")
	expectedUser := User{Username: "mustermann", Password: "musterpasswort", SessionID: "1234"}
	assert.Equal(t, expectedUser, GetUserBySession("1234"))
	assert.Equal(t, User{}, GetUserBySession(""))
	assert.Equal(t, User{}, GetUserBySession("FalscheSession"))
}

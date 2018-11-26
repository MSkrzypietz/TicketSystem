package XML_IO

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestInitDataStorage(t *testing.T) {
	os.RemoveAll("../data")
	InitDataStorage()
}

func TestTicketCreation(t *testing.T) {
	expectedTicket, err := CreateTicket("../data/tickets", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	actTicket, _ := ReadTicket("../data/tickets", expectedTicket.Id)
	expectedTicket.XMLName.Local = "Ticket"
	expectedTicket.MessageList = nil
	actTicket.MessageList = nil
	assert.Equal(t, expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestAddMessage(t *testing.T) {
	createdTicket, err := CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	_, err1 := AddMessage("../data/tickets", createdTicket, "4262", "please restart")
	expectedTicket, err2 := AddMessage("../data/tickets", createdTicket, "client@dhbw.de", "Thank, it worked!")
	assert.Nil(t, err1)
	assert.Nil(t, err2)
	expectedTicket.XMLName.Local = "Ticket"
	actTicket, err := ReadTicket("../data/tickets", expectedTicket.Id)
	expectedTicket.MessageList[0].CreationDate = actTicket.MessageList[0].CreationDate
	expectedTicket.MessageList[1].CreationDate = actTicket.MessageList[1].CreationDate
	assert.Equal(t, expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestStoreTicket(t *testing.T) {
	expectedTicket := Ticket{XMLName: xml.Name{"", "Ticket"}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0"}
	ticketMap[expectedTicket.Id] = expectedTicket
	StoreTicket("../data/tickets", expectedTicket)
	assert.Equal(t, Ticket{}, ticketMap[expectedTicket.Id])
	actTicket, _ := ReadTicket("../data/tickets", expectedTicket.Id)
	assert.Equal(t, expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestReadTicket(t *testing.T) {
	expectedTicket := Ticket{XMLName: xml.Name{"", "Ticket"}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0"}
	StoreTicket("../data/tickets", expectedTicket)
	actTicket, _ := ReadTicket("../data/tickets", expectedTicket.Id)
	assert.Equal(t, expectedTicket, actTicket)
	actTicket, _ = ReadTicket("../data/tickets", expectedTicket.Id)
	assert.Equal(t, ticketMap[expectedTicket.Id], actTicket)
	_, err := ReadTicket("../data/tickets", -99)
	assert.NotNil(t, err)
	removeCompleteDataStorage()
}

func TestIDCounter(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 2, getTicketIDCounter("definitions.xml"))
	assert.Equal(t, -1, getTicketIDCounter("/\\!?$&"))
	removeCompleteDataStorage()
}

func TestGetTicketsByStatus(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus("../data/tickets/ticket", getTicketIDCounter("definitions.xml"), 1)

	tickets := GetTicketsByStatus("../data/tickets/ticket", "definitions.xml", 0)
	for _, element := range tickets {
		assert.Equal(t, 0, element.Status)
	}

	tickets = GetTicketsByStatus("../data/tickets/ticket", "definitions.xml", 1)
	for _, element := range tickets {
		assert.Equal(t, 1, element.Status)
	}
	removeCompleteDataStorage()
}

func TestGetTicketByEditor(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeEditor("../data/tickets/ticket", getTicketIDCounter("definitions.xml")-1, "423")
	ChangeEditor("../data/tickets/ticket", getTicketIDCounter("definitions.xml"), "22")

	tickets := GetTicketsByEditor("../data/tickets/ticket", "definitions.xml", "423")
	for _, element := range tickets {
		assert.Equal(t, "423", element.Editor)
	}
	tickets = GetTicketsByEditor("../data/tickets/ticket", "definitions.xml", "22")
	for _, element := range tickets {
		assert.Equal(t, "22", element.Editor)
	}
	removeCompleteDataStorage()
}

func TestChangeEditor(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeEditor("../data/tickets/ticket", getTicketIDCounter("definitions.xml"), "4321")
	ticket, _ := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml"))
	assert.Equal(t, "4321", ticket.Editor)
	assert.NotNil(t, ChangeEditor("../data/tickets/ticket", -99, "1234"))
	removeCompleteDataStorage()
}

func TestChangeStatus(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus("../data/tickets/ticket", getTicketIDCounter("definitions.xml"), 2)
	ticket, _ := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml"))
	assert.Equal(t, 2, ticket.Status)

	assert.NotNil(t, ChangeStatus("../data/tickets/ticket", -99, -1))
	removeCompleteDataStorage()
}

func TestDeleteTicket(t *testing.T) {
	expectedTicket := Ticket{XMLName: xml.Name{"", "Ticket"}, Id: 1, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0"}
	StoreTicket("../data/tickets", expectedTicket)
	ReadTicket("../data/tickets", expectedTicket.Id)
	assert.Equal(t, expectedTicket, ticketMap[expectedTicket.Id])
	assert.Nil(t, DeleteTicket("../data/tickets", "definitions.xml", expectedTicket.Id))
	assert.Equal(t, Ticket{}, ticketMap[1])
	assert.NotNil(t, DeleteTicket("../data/tickets", "definitinos.xml", -99))
	removeCompleteDataStorage()
}

func TestMerging(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Max Mustermann. Thanks.")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	firstTicket, _ := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml")-1)
	secondTicket, _ := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml"))
	ChangeStatus("../data/tickets/ticket", firstTicket.Id, 1)
	ChangeStatus("../data/tickets/ticket", secondTicket.Id, 1)
	ChangeEditor("../data/tickets/ticket", firstTicket.Id, "202")
	ChangeEditor("../data/tickets/ticket", secondTicket.Id, "202")
	firstTicket, _ = ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml")-1)
	secondTicket, _ = ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml"))

	//merge two tickets and test the function
	var msgList []Message
	msgList = firstTicket.MessageList
	for e := range secondTicket.MessageList {
		msgList = append(msgList, secondTicket.MessageList[e])
	}
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: firstTicket.Id, Client: firstTicket.Client, Reference: firstTicket.Reference, Status: firstTicket.Status, Editor: firstTicket.Editor, MessageList: msgList}

	assert.Nil(t, MergeTickets("../data/tickets/ticket", "definitions.xml", firstTicket.Id, secondTicket.Id))
	actTicket, _ := ReadTicket("../data/tickets/ticket", firstTicket.Id)
	expectedTicket.XMLName.Local = "Ticket"
	assert.Equal(t, expectedTicket, actTicket)

	//merge tickets with two different editors
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	secondTicketID := getTicketIDCounter("definitions.xml")
	ChangeEditor("../data/tickets/ticket", secondTicketID, "412")
	assert.NotNil(t, MergeTickets("../data/tickets/ticket", "definitions.xml", firstTicket.Id, secondTicketID))

	assert.NotNil(t, MergeTickets("../data/tickets/ticket", "definitions.xml", -1, 1))
	assert.NotNil(t, MergeTickets("../data/tickets/ticket", "definitions.xml", 1, -1))

	removeCompleteDataStorage()
}

func TestCheckCache(t *testing.T) {
	for tmpInt := 1; tmpInt <= 12; tmpInt++ {
		CreateTicket("../data/tickets", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	}
	assert.Equal(t, 0, len(ticketMap))
	for tmpInt := 1; tmpInt <= 2; tmpInt++ {
		ReadTicket("../data/tickets", tmpInt)
	}
	assert.Equal(t, 2, len(ticketMap))
	for tmpInt := 1; tmpInt <= 12; tmpInt++ {
		ReadTicket("../data/tickets", tmpInt)
	}
	assert.Equal(t, 11, len(ticketMap))
	removeCompleteDataStorage()
}

func TestCreateAndStoreUser(t *testing.T) {
	_, err := CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Nil(t, err)
	file, _ := ioutil.ReadFile("../data/users/users.xml")
	var userlist Userlist
	xml.Unmarshal(file, &userlist)

	var expectedUser []User
	expectedUser = append(expectedUser, User{Username: "mustermann", Password: "musterpasswort"})
	expected := Userlist{User: expectedUser}
	assert.Equal(t, expected, userlist)

	_, err = CreateUser("../data/users/u", "test", "test")
	assert.NotNil(t, err)
}

func TestReadUser(t *testing.T) {
	_, err := CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Nil(t, err)
	expectedMap := make(map[string]User)
	expectedMap["mustermann"] = User{Username: "mustermann", Password: "musterpasswort"}
	tmpMap, err := readUsers("../data/users/users.xml")
	assert.Nil(t, err)
	assert.Equal(t, expectedMap, tmpMap)
	os.Remove("../data/users/users.xml")
	expectedMap = make(map[string]User)
	tmpMap, err = readUsers("../data/users/users.xml")
	assert.NotNil(t, err)
	assert.Equal(t, expectedMap, tmpMap)
	os.Create("../data/users/users.xml")
}

func TestCheckUser(t *testing.T) {
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	tmpBool, err := CheckUser("../data/users/users.xml", "mustermann")
	assert.Nil(t, err)
	assert.False(t, tmpBool)
	tmpBool, err = CheckUser("../data/users/users.xml", "muster")
	assert.Nil(t, err)
	assert.True(t, tmpBool)

	_, err = CheckUser("../data/users/u", "test")
	assert.NotNil(t, err)
}

func TestVerifyUser(t *testing.T) {
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	tmpBool, err := VerifyUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.True(t, tmpBool)
	assert.Nil(t, err)
	tmpBool, err = VerifyUser("../data/users/users.xml", "mustermann", "xxx")
	assert.False(t, tmpBool)
	assert.Nil(t, err)

	_, err = VerifyUser("../data/users/u", "test", "test")
	assert.NotNil(t, err)
}

func TestLoginUser(t *testing.T) {
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Nil(t, LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234"))
	assert.NotNil(t, LoginUser("../data/users/users.xml", "mustermann", "falschespasswort", "1234"))
	usersMap, _ := readUsers("../data/users/users.xml")
	assert.Equal(t, "1234", usersMap["mustermann"].SessionID)

	assert.NotNil(t, LoginUser("../data/users/u", "test", "124", "1234"))
}

func TestLogoutUser(t *testing.T) {
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Nil(t, LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234"))
	usersmap, _ := readUsers("../data/users/users.xml")
	assert.Equal(t, "1234", usersmap["mustermann"].SessionID)
	assert.Nil(t, LogoutUser("../data/users/users.xml", "mustermann"))
	usersmap, _ = readUsers("../data/users/users.xml")
	assert.Equal(t, "", usersmap["mustermann"].SessionID)
	assert.NotNil(t, LogoutUser("../data/users/users.xml", "falscherName"))

	assert.NotNil(t, LogoutUser("../data/users/u", "test"))
}

func TestGetUserSession(t *testing.T) {
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Equal(t, "", GetUserSession("../data/users/users.xml", "mustermann"))
	LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234")
	assert.Equal(t, "1234", GetUserSession("../data/users/users.xml", "mustermann"))
}

func TestGetUserBySession(t *testing.T) {
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234")
	expectedUser := User{Username: "mustermann", Password: "musterpasswort", SessionID: "1234"}
	assert.Equal(t, expectedUser, GetUserBySession("../data/users/users.xml", "1234"))
	assert.Equal(t, User{}, GetUserBySession("../data/users/users.xml", ""))
	assert.Equal(t, User{}, GetUserBySession("../data/users/users.xml", "FalscheSession"))
}

func removeCompleteDataStorage() {
	os.RemoveAll("../data")
	InitDataStorage()
	ticketMap = make(map[int]Ticket)
	writeToXML(0, "definitions.xml")
}

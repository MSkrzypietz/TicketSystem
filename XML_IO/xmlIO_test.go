package XML_IO

import (
	"TicketSystem/config"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestInitDataStorage(t *testing.T) {
	assert := assert.New(t)
	os.RemoveAll("../data")
	InitDataStorage("../data/tickets", "../data/users")
	_, err := os.Stat(path.Join("../", config.UsersFilePath()))
	assert.Nil(err)
	_, err = os.Stat(path.Join("../", config.TicketsPath()))
	assert.Nil(err)
}

func TestTicketCreation(t *testing.T) {
	assert := assert.New(t)
	expectedTicket, err := CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(err)
	actTicket, _ := ReadTicket("../data/tickets/ticket", 1)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestAddMessage(t *testing.T) {
	assert := assert.New(t)
	tmpTicket, err := CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(err)
	expectedTicket, _ := AddMessage("../data/tickets/ticket", tmpTicket, "4262", "please restart")
	actTicket, err := ReadTicket("../data/tickets/ticket", expectedTicket.Id)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	actTicket.MessageList[1].CreationDate = expectedTicket.MessageList[1].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestTicketStoring(t *testing.T) {
	assert := assert.New(t)
	actTicket, err := CreateTicket("../data/tickets/ticket", "definitions.xml", "1234", "PC problem", "Pc does not start anymore")
	assert.Nil(err)
	ticketMap[1] = actTicket
	DeleteTicket("../data/tickets/ticket", "definitions.xml", 1)
	assert.Equal(Ticket{}, ticketMap[1])
	expectedTicket, err := ReadTicket("../data/tickets/ticket", 1)
	assert.NotNil(err)
	assert.Equal(Ticket{}, expectedTicket)
	removeCompleteDataStorage()
}

func TestTicketReading(t *testing.T) {
	assert := assert.New(t)
	tmpTicket := Ticket{Id: 1}
	ticketMap[1] = tmpTicket
	actTicket, _ := ReadTicket("../data/tickets/ticket", 1)
	assert.Equal(tmpTicket, actTicket)
	removeCompleteDataStorage()
	_, err := ReadTicket("../data/tickets/ticket", 1)
	assert.NotNil(err)
	expectedTicket, _ := CreateTicket("../data/tickets/ticket", "definitions.xml", "1234", "PC problem", "Pc does not start anymore")
	actTicket, _ = ReadTicket("../data/tickets/ticket", 1)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestIDCounter(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 2, getTicketIDCounter("definitions.xml"))

	removeCompleteDataStorage()
}

func TestTicketsByStatus(t *testing.T) {
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

func TestTicketByEditor(t *testing.T) {
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

	removeCompleteDataStorage()
}

func TestChangeStatus(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus("../data/tickets/ticket", getTicketIDCounter("definitions.xml"), 2)
	ticket, _ := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml"))
	assert.Equal(t, 2, ticket.Status)

	removeCompleteDataStorage()
}

func TestDeleting(t *testing.T) {
	assert := assert.New(t)
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "Computer", "PC not working")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", 1)
	assert.Equal(0, getTicketIDCounter("definitions.xml"))
	err := DeleteTicket("../data/tickets/ticket", "definitions.xml", 11)
	assert.NotNil(err)
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "Computer", "PC not working")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", 1)
	_, err = ReadTicket("../data/tickets/ticket", 1)
	assert.NotNil(err)
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
	actTicket.XMLName.Local = ""
	assert.Equal(t, expectedTicket, actTicket)

	//merge tickets with two different editors
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	secondTicketID := getTicketIDCounter("definitions.xml")
	ChangeEditor("../data/tickets/ticket", secondTicketID, "412")
	assert.NotNil(t, MergeTickets("../data/tickets/ticket", "definitions.xml", firstTicket.Id, secondTicketID))

	removeCompleteDataStorage()
}

func TestCheckCache(t *testing.T) {
	assert := assert.New(t)
	for tmpInt := 1; tmpInt <= 11; tmpInt++ {
		CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	}
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		ReadTicket("../data/tickets/ticket", tmpInt)
	}
	assert.Equal(9, len(ticketMap))
	ReadTicket("../data/tickets/ticket", 10)
	assert.Equal(10, len(ticketMap))
	ReadTicket("../data/tickets/ticket", 11)
	assert.Equal(10, len(ticketMap))
	removeCompleteDataStorage()
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)
	_, err := CreateUser("", "", "")
	assert.NotNil(err)
	expectedUser, _ := CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Equal(expectedUser.Username, "mustermann")
	removeCompleteDataStorage()
}

func TestStoreUser(t *testing.T) {
	assert := assert.New(t)
	tmpUserMap := make(map[string]User)
	tmpUserMap["mustermann"] = User{Username: "mustermann", Password: "musterpasswort"}
	storeUsers("../data/users/users.xml", tmpUserMap)
	actMap, _ := readUsers("../data/users/users.xml")
	assert.Equal(tmpUserMap, actMap)
	removeCompleteDataStorage()
}

//TODO: fix test for reading users
func TestReadUser(t *testing.T) {
	assert := assert.New(t)
	_, err := readUsers("falscherPfad")
	assert.NotNil(err)
	user1, _ := CreateUser("../data/users/users.xml", "mustermann1", "musterpasswort")
	user2, _ := CreateUser("../data/users/users.xml", "mustermann2", "musterpasswort")
	expectedMap := make(map[string]User)
	expectedMap[user1.Username] = user1
	expectedMap[user2.Username] = user2
	actMap, _ := readUsers("../data/users/users.xml")
	assert.Equal(expectedMap, actMap)
	removeCompleteDataStorage()
}

func TestCheckUser(t *testing.T) {
	assert := assert.New(t)
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	tmpBool, err := CheckUser("../data/users/users.xml", "mustermann")
	assert.Nil(err)
	assert.False(tmpBool)
	tmpBool, err = CheckUser("../data/users/users.xml", "muster")
	assert.Nil(err)
	assert.True(tmpBool)
	removeCompleteDataStorage()
}

func TestVerifyUser(t *testing.T) {
	assert := assert.New(t)
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	tmpBool, err := VerifyUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.True(tmpBool)
	assert.Nil(err)
	tmpBool, err = VerifyUser("../data/users/users.xml", "mustermann", "xxx")
	assert.False(tmpBool)
	assert.NotNil(err)
	removeCompleteDataStorage()
}

func TestLoginUser(t *testing.T) {
	assert := assert.New(t)
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Nil(LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234"))
	assert.NotNil(LoginUser("../data/users/users.xml", "mustermann", "falschespasswort", "1234"))
	usersMap, _ := readUsers("../data/users/users.xml")
	assert.Equal("1234", usersMap["mustermann"].SessionID)
	removeCompleteDataStorage()
}

func TestLogoutUser(t *testing.T) {
	assert := assert.New(t)
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Nil(LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234"))
	usersmap, _ := readUsers("../data/users/users.xml")
	assert.Equal("1234", usersmap["mustermann"].SessionID)
	assert.Nil(LogoutUser("../data/users/users.xml", "mustermann"))
	usersmap, _ = readUsers("../data/users/users.xml")
	assert.Equal("", usersmap["mustermann"].SessionID)
	assert.NotNil(LogoutUser("../data/users/users.xml", "falscherName"))
	removeCompleteDataStorage()
}

func TestGetUserSession(t *testing.T) {
	assert := assert.New(t)
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Equal("", GetUserSession("../data/users/users.xml", "mustermann"))
	LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234")
	assert.Equal("1234", GetUserSession("../data/users/users.xml", "mustermann"))
	removeCompleteDataStorage()
}

func TestGetUserBySession(t *testing.T) {
	assert := assert.New(t)
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234")
	actUser, err := GetUserBySession("../data/users/users.xml", "1234")
	assert.Equal("mustermann", actUser.Username)
	_, err = GetUserBySession("../data/users/users.xml", "")
	assert.NotNil(err)
	_, err = GetUserBySession("../data/users/users.xml", "FalscheSession")
	assert.NotNil(err)
	removeCompleteDataStorage()
}

func removeCompleteDataStorage() {
	os.RemoveAll("../data")
	InitDataStorage("../data/tickets", "../data/users")
	ticketMap = make(map[int]Ticket)
	writeToXML(0, "definitions.xml")
}

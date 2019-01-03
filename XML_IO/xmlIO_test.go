package XML_IO

import (
	"TicketSystem/config"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"os"
	"testing"
)

func TestInitDataStorage(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	os.RemoveAll(config.DataPath)
	InitDataStorage()
	_, err := os.Stat(config.UsersFilePath())
	assert.Nil(err)
	_, err = os.Stat(config.TicketsPath())
	assert.Nil(err)
}

func TestTicketCreation(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	expectedTicket, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(err)
	actTicket, _ := ReadTicket(1)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestAddMessage(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	tmpTicket, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(err)
	expectedTicket, _ := AddMessage(tmpTicket, "4262", "please restart")
	actTicket, err := ReadTicket(expectedTicket.Id)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	actTicket.MessageList[1].CreationDate = expectedTicket.MessageList[1].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestTicketStoring(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	actTicket, err := CreateTicket("1234", "PC problem", "Pc does not start anymore")
	assert.Nil(err)
	ticketMap[1] = actTicket
	DeleteTicket(1)
	assert.Equal(Ticket{}, ticketMap[1])
	expectedTicket, err := ReadTicket(1)
	assert.NotNil(err)
	assert.Equal(Ticket{}, expectedTicket)
	removeCompleteDataStorage()
}

func TestTicketReading(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	tmpTicket := Ticket{Id: 1}
	ticketMap[1] = tmpTicket
	actTicket, _ := ReadTicket(1)
	assert.Equal(tmpTicket, actTicket)
	removeCompleteDataStorage()
	_, err := ReadTicket(1)
	assert.NotNil(err)
	expectedTicket, _ := CreateTicket("1234", "PC problem", "Pc does not start anymore")
	actTicket, _ = ReadTicket(1)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestIDCounter(t *testing.T) {
	config.DataPath = "../data"
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 2, getTicketIDCounter())
	removeCompleteDataStorage()
}

func TestTicketsByStatus(t *testing.T) {
	config.DataPath = "../data"
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

	removeCompleteDataStorage()
}

func TestTicketByEditor(t *testing.T) {
	config.DataPath = "../data"
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

	removeCompleteDataStorage()
}

func TestChangeEditor(t *testing.T) {
	config.DataPath = "../data"
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeEditor(getTicketIDCounter(), "4321")
	ticket, _ := ReadTicket(getTicketIDCounter())
	assert.Equal(t, "4321", ticket.Editor)
	removeCompleteDataStorage()
}

func TestChangeStatus(t *testing.T) {
	config.DataPath = "../data"
	CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus(getTicketIDCounter(), 2)
	ticket, _ := ReadTicket(getTicketIDCounter())
	assert.Equal(t, 2, ticket.Status)
	removeCompleteDataStorage()
}

func TestDeleting(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	CreateTicket("client@dhbw.de", "Computer", "PC not working")
	DeleteTicket(1)
	assert.Equal(0, getTicketIDCounter())
	err := DeleteTicket(11)
	assert.NotNil(err)
	CreateTicket("client@dhbw.de", "Computer", "PC not working")
	DeleteTicket(1)
	_, err = ReadTicket(1)
	assert.NotNil(err)
	removeCompleteDataStorage()
}

func TestMerging(t *testing.T) {
	config.DataPath = "../data"
	CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Max Mustermann. Thanks.")
	CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	firstTicket, _ := ReadTicket(getTicketIDCounter() - 1)
	secondTicket, _ := ReadTicket(getTicketIDCounter())
	ChangeStatus(firstTicket.Id, 1)
	ChangeStatus(secondTicket.Id, 1)
	ChangeEditor(firstTicket.Id, "202")
	ChangeEditor(secondTicket.Id, "202")
	firstTicket, _ = ReadTicket(getTicketIDCounter() - 1)
	secondTicket, _ = ReadTicket(getTicketIDCounter())

	//merge two tickets and test the function
	var msgList []Message
	msgList = firstTicket.MessageList
	for e := range secondTicket.MessageList {
		msgList = append(msgList, secondTicket.MessageList[e])
	}
	expectedTicket := Ticket{XMLName: xml.Name{"", ""}, Id: firstTicket.Id, Client: firstTicket.Client, Reference: firstTicket.Reference, Status: firstTicket.Status, Editor: firstTicket.Editor, MessageList: msgList}

	assert.Nil(t, MergeTickets(firstTicket.Id, secondTicket.Id))
	actTicket, _ := ReadTicket(firstTicket.Id)
	actTicket.XMLName.Local = ""
	assert.Equal(t, expectedTicket, actTicket)

	//merge tickets with two different editors
	CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	secondTicketID := getTicketIDCounter()
	ChangeEditor(secondTicketID, "412")
	assert.NotNil(t, MergeTickets(firstTicket.Id, secondTicketID))

	removeCompleteDataStorage()
}

func TestCheckCache(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	for tmpInt := 1; tmpInt <= 11; tmpInt++ {
		CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	}
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		ReadTicket(tmpInt)
	}
	assert.Equal(9, len(ticketMap))
	ReadTicket(10)
	assert.Equal(10, len(ticketMap))
	ReadTicket(11)
	assert.Equal(10, len(ticketMap))
	removeCompleteDataStorage()
}

func TestCreateUser(t *testing.T) {
	config.DataPath = "wrongPath"
	assert := assert.New(t)
	_, err := CreateUser("", "")
	assert.NotNil(err)
	config.DataPath = "../data"
	expectedUser, _ := CreateUser("mustermann", "musterpasswort")
	assert.Equal(expectedUser.Username, "mustermann")
	removeCompleteDataStorage()
}

func TestStoreUser(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	tmpUserMap := make(map[string]User)
	tmpUserMap["mustermann"] = User{Username: "mustermann", Password: "musterpasswort"}
	storeUsers(tmpUserMap)
	actMap, _ := readUsers()
	assert.Equal(tmpUserMap, actMap)
	removeCompleteDataStorage()
}

func TestReadUser(t *testing.T) {
	config.DataPath = "wrongPath"
	assert := assert.New(t)
	_, err := readUsers()
	assert.NotNil(err)

	config.DataPath = "../data"
	CreateUser("testOne", "test")
	CreateUser("testTwo", "test")
	actMap, err := readUsers()
	assert.Nil(err)
	assert.Equal("testOne", actMap["testOne"].Username)
	assert.Equal("testTwo", actMap["testTwo"].Username)
	removeCompleteDataStorage()
}

func TestCheckUser(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	CreateUser("mustermann", "musterpasswort")
	tmpBool, err := CheckUser("mustermann")
	assert.Nil(err)
	assert.False(tmpBool)
	tmpBool, err = CheckUser("muster")
	assert.Nil(err)
	assert.True(tmpBool)
	removeCompleteDataStorage()
}

//TODO: add a positiv verification
func TestVerifyUser(t *testing.T) {
	config.DataPath = "wrongPath"
	assert := assert.New(t)
	_, err := VerifyUser("", "")
	assert.NotNil(err)
	config.DataPath = "../data"

	CreateUser("mustermann", "musterpasswort")
	tmpPsswd, _ := bcrypt.GenerateFromPassword([]byte("musterpasswort"), bcrypt.DefaultCost)
	tmpBool, err := VerifyUser("mustermann", string(tmpPsswd))

	tmpBool, err = VerifyUser("mustermann", "xxx")
	assert.False(tmpBool)
	assert.NotNil(err)
	removeCompleteDataStorage()
}

func TestLoginUser(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	CreateUser("mustermann", "musterpasswort")
	assert.Nil(LoginUser("mustermann", "musterpasswort", "1234"))
	assert.NotNil(LoginUser("mustermann", "falschespasswort", "1234"))
	usersMap, _ := readUsers()
	assert.Equal("1234", usersMap["mustermann"].SessionID)
	removeCompleteDataStorage()
}

func TestLogoutUser(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	CreateUser("mustermann", "musterpasswort")
	assert.Nil(LoginUser("mustermann", "musterpasswort", "1234"))
	usersmap, _ := readUsers()
	assert.Equal("1234", usersmap["mustermann"].SessionID)
	assert.Nil(LogoutUser("mustermann"))
	usersmap, _ = readUsers()
	assert.Equal("", usersmap["mustermann"].SessionID)
	assert.NotNil(LogoutUser("falscherName"))
	removeCompleteDataStorage()
}

func TestGetUserSession(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	CreateUser("mustermann", "musterpasswort")
	assert.Equal("", GetUserSession("mustermann"))
	LoginUser("mustermann", "musterpasswort", "1234")
	assert.Equal("1234", GetUserSession("mustermann"))
	removeCompleteDataStorage()
}

func TestGetUserBySession(t *testing.T) {
	config.DataPath = "../data"
	assert := assert.New(t)
	CreateUser("mustermann", "musterpasswort")
	LoginUser("mustermann", "musterpasswort", "1234")
	actUser, err := GetUserBySession("1234")
	assert.Equal("mustermann", actUser.Username)
	_, err = GetUserBySession("")
	assert.NotNil(err)
	_, err = GetUserBySession("FalscheSession")
	assert.NotNil(err)
	removeCompleteDataStorage()
}

func removeCompleteDataStorage() {
	config.DataPath = "../data"
	os.RemoveAll(config.DataPath)
	InitDataStorage()
	ticketMap = make(map[int]Ticket)
	writeToXML(0, config.DefinitionsFilePath())
}

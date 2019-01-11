package utils

import (
	"TicketSystem/config"
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func setup() {
	config.DataPath = "datatest"
	err := InitDataStorage()
	if err != nil {
		log.Println(err)
	}
}

func teardown() {
	err := os.RemoveAll(config.DataPath)
	if err != nil {
		log.Println(err)
	}
	ticketMap = make(map[int]Ticket)
}

func TestInitDataStorage(t *testing.T) {
	defer teardown()

	config.DataPath = "datatest"
	assert.Nil(t, os.RemoveAll(config.DataPath))
	assert.Nil(t, InitDataStorage())
	_, err := os.Stat(config.UsersFilePath())
	assert.Nil(t, err)
	_, err = os.Stat(config.TicketsPath())
	assert.Nil(t, err)
	_, err = os.Stat(config.MailFilePath())
	assert.Nil(t, err)
}

func TestTicketCreation(t *testing.T) {
	setup()
	defer teardown()

	expectedTicket, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	actTicket, _ := ReadTicket(1)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(t, expectedTicket, actTicket)
}

func TestAddMessage(t *testing.T) {
	setup()
	defer teardown()

	tmpTicket, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	expectedTicket, _ := AddMessage(tmpTicket, "4262", "please restart")
	actTicket, err := ReadTicket(expectedTicket.ID)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	actTicket.MessageList[1].CreationDate = expectedTicket.MessageList[1].CreationDate
	assert.Equal(t, expectedTicket, actTicket)
}

func TestTicketStoring(t *testing.T) {
	setup()
	defer teardown()

	actTicket, err := CreateTicket("1234", "PC problem", "Pc does not start anymore")
	assert.Nil(t, err)
	ticketMap[1] = actTicket
	assert.Nil(t, deleteTicket(1))
	assert.Equal(t, Ticket{}, ticketMap[1])
	expectedTicket, err := ReadTicket(1)
	assert.NotNil(t, err)
	assert.Equal(t, Ticket{}, expectedTicket)
}

func TestTicketReading(t *testing.T) {
	setup()
	defer teardown()

	tmpTicket := Ticket{ID: 1}
	ticketMap[1] = tmpTicket
	actTicket, _ := ReadTicket(1)
	assert.Equal(t, tmpTicket, actTicket)
	teardown()
	setup()

	_, err := ReadTicket(1)
	assert.NotNil(t, err)
	expectedTicket, _ := CreateTicket("1234", "PC problem", "Pc does not start anymore")
	actTicket, _ = ReadTicket(expectedTicket.ID)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(t, expectedTicket, actTicket)
}

func TestIDCounter(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	assert.NotNil(t, getTicketIDCounter())

	config.DataPath = "datatest"
	_, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	_, err = CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	assert.Equal(t, 2, getTicketIDCounter())
}

func TestTicketsByStatus(t *testing.T) {
	setup()
	defer teardown()

	_, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	_, err = CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	err = ChangeStatus(getTicketIDCounter(), 1)
	assert.Nil(t, err)

	tickets := GetTicketsByStatus(0)
	for _, element := range tickets {
		assert.Equal(t, 0, element.Status)
	}

	tickets = GetTicketsByStatus(1)
	for _, element := range tickets {
		assert.Equal(t, 1, element.Status)
	}
}

func TestTicketByEditor(t *testing.T) {
	setup()
	defer teardown()

	_, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	_, err = CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	err = ChangeEditor(getTicketIDCounter()-1, "423")
	assert.Nil(t, err)
	err = ChangeEditor(getTicketIDCounter(), "22")
	assert.Nil(t, err)

	tickets := GetTicketsByEditor("423")
	for _, element := range tickets {
		assert.Equal(t, "423", element.Editor)
	}
	tickets = GetTicketsByEditor("22")
	for _, element := range tickets {
		assert.Equal(t, "22", element.Editor)
	}
}

func TestTicketByClient(t *testing.T) {
	setup()
	defer teardown()

	_, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	_, err = CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	_, err = CreateTicket("clientTwo@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)

	tickets := GetTicketsByClient("client@dhbw.de")
	for _, element := range tickets {
		assert.Equal(t, "client@dhbw.de", element.Client)
	}
	tickets = GetTicketsByClient("clientTwo@dhbw.de")
	for _, element := range tickets {
		assert.Equal(t, "clientTwo@dhbw.de", element.Client)
	}
}

func TestChangeEditor(t *testing.T) {
	setup()
	defer teardown()

	_, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	err = ChangeEditor(getTicketIDCounter(), "4321")
	assert.Nil(t, err)
	ticket, _ := ReadTicket(getTicketIDCounter())
	assert.Equal(t, "4321", ticket.Editor)
}

func TestChangeStatus(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	assert.NotNil(t, ChangeStatus(1, TicketStatusClosed))

	config.DataPath = "datatest"
	_, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)
	err = ChangeStatus(getTicketIDCounter(), 2)
	assert.Nil(t, err)
	ticket, err := ReadTicket(getTicketIDCounter())
	assert.Nil(t, err)
	assert.Equal(t, 2, ticket.Status)
}

func TestDeleting(t *testing.T) {
	setup()
	defer teardown()

	_, err := CreateTicket("client@dhbw.de", "Computer", "PC not working")
	assert.Nil(t, err)
	err = deleteTicket(1)
	assert.Nil(t, err)
	err = deleteTicket(11)
	assert.NotNil(t, err)
	_, err = CreateTicket("client@dhbw.de", "Computer", "PC not working")
	assert.Nil(t, err)
	err = deleteTicket(2)
	assert.Nil(t, err)
	_, err = ReadTicket(2)
	assert.NotNil(t, err)
}

func TestMerging(t *testing.T) {
	setup()
	defer teardown()

	ticket, err := CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Max Mustermann. Thanks.")
	assert.Nil(t, err)
	assert.NotNil(t, MergeTickets(ticket.ID, 1337))
	assert.NotNil(t, MergeTickets(1337, ticket.ID))

	_, err = CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Max Mustermann. Thanks.")
	assert.Nil(t, err)
	_, err = CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	assert.Nil(t, err)
	firstTicket, _ := ReadTicket(getTicketIDCounter() - 1)
	secondTicket, _ := ReadTicket(getTicketIDCounter())
	err = ChangeStatus(firstTicket.ID, 1)
	assert.Nil(t, err)
	err = ChangeStatus(secondTicket.ID, 1)
	assert.Nil(t, err)
	err = ChangeEditor(firstTicket.ID, "202")
	assert.Nil(t, err)
	err = ChangeEditor(secondTicket.ID, "202")
	assert.Nil(t, err)
	firstTicket, _ = ReadTicket(getTicketIDCounter() - 1)
	secondTicket, _ = ReadTicket(getTicketIDCounter())

	//merge two tickets and test the function
	var msgList []Message
	msgList = firstTicket.MessageList
	for e := range secondTicket.MessageList {
		msgList = append(msgList, secondTicket.MessageList[e])
	}
	expectedTicket := Ticket{XMLName: xml.Name{Space: "", Local: ""}, ID: firstTicket.ID, Client: firstTicket.Client, Reference: firstTicket.Reference, Status: firstTicket.Status, Editor: firstTicket.Editor, MessageList: msgList}

	assert.Nil(t, MergeTickets(firstTicket.ID, secondTicket.ID))
	actTicket, _ := ReadTicket(firstTicket.ID)
	actTicket.XMLName.Local = ""
	assert.Equal(t, expectedTicket, actTicket)

	//merge tickets with two different editors
	_, err = CreateTicket("client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	assert.Nil(t, err)
	secondTicketID := getTicketIDCounter()
	err = ChangeEditor(secondTicketID, "412")
	assert.Nil(t, err)
	assert.NotNil(t, MergeTickets(firstTicket.ID, secondTicketID))
}

func TestCheckCache(t *testing.T) {
	setup()
	defer teardown()

	for tmpInt := 1; tmpInt <= 11; tmpInt++ {
		_, err := CreateTicket("client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
		assert.Nil(t, err)
	}
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		_, err := ReadTicket(tmpInt)
		assert.Nil(t, err)
	}
	assert.Equal(t, 9, len(ticketMap))
	_, err := ReadTicket(10)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(ticketMap))
	_, err = ReadTicket(11)
	assert.Nil(t, err)
	assert.Equal(t, 10, len(ticketMap))
}

func TestCreateUser(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	_, err := CreateUser("", "")
	assert.NotNil(t, err)

	config.DataPath = "datatest"
	_, err = CreateUser("mustermann", "musterpasswort")
	assert.Nil(t, err)

	_, err = CreateUser("mustermann", "musterpasswort")
	assert.NotNil(t, err)
}

func TestStoreUser(t *testing.T) {
	setup()
	defer teardown()

	tmpUserMap := make(map[string]User)
	tmpUserMap["mustermann"] = User{Username: "mustermann", Password: "musterpasswort"}
	err := storeUsers(tmpUserMap)
	assert.Nil(t, err)
	actMap, _ := ReadUsers()
	assert.Equal(t, tmpUserMap, actMap)
}

func TestReadUser(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	_, err := ReadUsers()
	assert.NotNil(t, err)

	config.DataPath = "datatest"
	_, err = CreateUser("testOne", "test")
	assert.Nil(t, err)
	_, err = CreateUser("testTwo", "test")
	assert.Nil(t, err)
	actMap, err := ReadUsers()
	assert.Nil(t, err)
	assert.Equal(t, "testOne", actMap["testOne"].Username)
	assert.Equal(t, "testTwo", actMap["testTwo"].Username)
}

func TestCheckUser(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	_, err := CheckUser("")
	assert.NotNil(t, err)

	config.DataPath = "datatest"
	_, err = CreateUser("mustermann", "musterpasswort")
	assert.Nil(t, err)
	tmpBool, err := CheckUser("mustermann")
	assert.Nil(t, err)
	assert.False(t, tmpBool)
	tmpBool, err = CheckUser("muster")
	assert.Nil(t, err)
	assert.True(t, tmpBool)
}

func TestVerifySessionCookie(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	err := VerifySessionCookie("", "")
	assert.NotNil(t, err)

	config.DataPath = "datatest"
	user, err := CreateUser("mustermann", "musterpasswort")
	assert.Nil(t, err)

	err = VerifySessionCookie(user.Username, user.Password)
	assert.Nil(t, err)

	err = VerifySessionCookie(user.Username, "xxx")
	assert.NotNil(t, err)
}

func TestLoginUser(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	err := LoginUser("", "", "")
	assert.NotNil(t, err)

	config.DataPath = "datatest"
	_, err = CreateUser("mustermann", "musterpasswort")
	assert.Nil(t, err)
	assert.Nil(t, LoginUser("mustermann", "musterpasswort", "1234"))
	assert.NotNil(t, LoginUser("mustermann", "falschespasswort", "1234"))
	usersMap, _ := ReadUsers()
	assert.Equal(t, "1234", usersMap["mustermann"].SessionID)
}

func TestLogoutUser(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	err := LogoutUser("")
	assert.NotNil(t, err)

	config.DataPath = "datatest"
	_, err = CreateUser("mustermann", "musterpasswort")
	assert.Nil(t, err)
	assert.Nil(t, LoginUser("mustermann", "musterpasswort", "1234"))
	usersMap, _ := ReadUsers()
	assert.Equal(t, "1234", usersMap["mustermann"].SessionID)
	assert.Nil(t, LogoutUser("mustermann"))
	usersMap, _ = ReadUsers()
	assert.Equal(t, "", usersMap["mustermann"].SessionID)
	assert.NotNil(t, LogoutUser("falscherName"))
}

func TestGetUserSession(t *testing.T) {
	setup()
	defer teardown()

	_, err := CreateUser("mustermann", "musterpasswort")
	assert.Nil(t, err)
	assert.Equal(t, "", GetUserSession("mustermann"))
	err = LoginUser("mustermann", "musterpasswort", "1234")
	assert.Nil(t, err)
	assert.Equal(t, "1234", GetUserSession("mustermann"))
}

func TestGetUserBySession(t *testing.T) {
	setup()
	defer teardown()

	_, err := CreateUser("mustermann", "musterpasswort")
	assert.Nil(t, err)
	err = LoginUser("mustermann", "musterpasswort", "1234")
	assert.Nil(t, err)
	actUser, err := GetUserBySession("1234")
	assert.Equal(t, "mustermann", actUser.Username)
	_, err = GetUserBySession("")
	assert.NotNil(t, err)
	_, err = GetUserBySession("FalscheSession")
	assert.NotNil(t, err)
}

func TestSetUserHolidayMode(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	assert.NotNil(t, SetUserHolidayMode("testUser", true))
	config.DataPath = "datatest"
	assert.NotNil(t, SetUserHolidayMode("testUser", true))
	_, _ = CreateUser("testUser", "test")
	assert.Nil(t, SetUserHolidayMode("testUser", true))
	tmpUsers, _ := ReadUsers()
	assert.True(t, tmpUsers["testUser"].HolidayMode)
	assert.Nil(t, SetUserHolidayMode("testUser", false))
	tmpUsers, _ = ReadUsers()
	assert.False(t, tmpUsers["testUser"].HolidayMode)
}

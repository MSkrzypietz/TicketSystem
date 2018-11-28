package XML_IO

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func TestInitDataStorage(t *testing.T) {
	os.RemoveAll("../data")
	InitDataStorage("../data/tickets", "../data/users")
}

func TestTicketCreation(t *testing.T) {
	expectedTicket, err := CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)

	ticketID := getTicketIDCounter("definitions.xml")
	actTicket := ticketMap[ticketID]

	assert.Equal(t, expectedTicket, actTicket)

	assert.Nil(t, ClearCache("../data/tickets/ticket"))

	f, err := ioutil.ReadFile("../data/tickets/ticket" + strconv.Itoa(ticketID) + ".xml")
	assert.NotNil(t, f)
	assert.Nil(t, err)

	DeleteTicket("../data/tickets/ticket", "definitions.xml", ticketID)
	ClearCache("../data/tickets/ticket")
	writeToXML(0, "definitions.xml")
}

func TestAddMessage(t *testing.T) {
	createdTicket, err := CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Nil(t, err)

	_, err1 := AddMessage("../data/tickets/ticket", createdTicket, "4262", "please restart")
	expectedTicket, err2 := AddMessage("../data/tickets/ticket", createdTicket, "client@dhbw.de", "Thank, it worked!")
	assert.Nil(t, err1)
	assert.Nil(t, err2)

	actTicket, err := ReadTicket("../data/tickets/ticket", expectedTicket.Id)
	assert.Equal(t, expectedTicket, actTicket)

	DeleteTicket("../data/tickets/ticket", "definitions.xml", expectedTicket.Id)
	ClearCache("../data/tickets/ticket")
}

func TestTicketStoring(t *testing.T) {
	expectedTicketFour := Ticket{}
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		actTicket, err := CreateTicket("../data/tickets/ticket", "definitions.xml", "client"+strconv.Itoa(tmpInt)+"@dhbw.de", "PC problem", "Pc does not start anymore")
		assert.Nil(t, err)
		if tmpInt == 4 {
			expectedTicketFour = actTicket
		}
	}

	assert.Equal(t, expectedTicketFour, ticketMap[4])

	_, err := ioutil.ReadFile("../data/tickets/ticket4.xml")
	assert.NotNil(t, err)

	CreateTicket("../data/tickets/ticket", "definitions.xml", "client10@dhbw.de", "PC problem", "Pc does not start anymore")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client11@dhbw.de", "PC problem", "Pc does not start anymore")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client12@dhbw.de", "PC problem", "Pc does not start anymore")
	assert.Equal(t, 11, len(ticketMap))

	ClearCache("../data/tickets/ticket")
	assert.Equal(t, 0, len(ticketMap))

	actTicket, err := ReadTicket("../data/tickets/ticket", 4)
	var expectedMsg []Message
	expectedMsg = append(expectedMsg, Message{Actor: "client4@dhbw.de", Text: "Pc does not start anymore", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicketFour = Ticket{XMLName: xml.Name{"", "Ticket"}, Id: 4, Client: "client4@dhbw.de", Reference: "PC problem", Status: 0, Editor: "0", MessageList: expectedMsg}
	expectedTicketFour.XMLName.Local = "Ticket"
	assert.Equal(t, expectedTicketFour, actTicket)

	ClearCache("../data/tickets/ticket")
	for tmpInt := 1; tmpInt <= 12; tmpInt++ {
		DeleteTicket("../data/tickets/ticket", "definitions.xml", tmpInt)
	}
	writeToXML(0, "definitions.xml")
}

func TestTicketReading(t *testing.T) {
	expectedTicket, _ := CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ticketID := getTicketIDCounter("definitions.xml")
	actTicket, err := ReadTicket("../data/tickets/ticket", ticketID)
	assert.Nil(t, err)

	assert.Equal(t, expectedTicket, actTicket)

	ClearCache("../data/tickets/ticket")
	actTicket, err = ReadTicket("../data/tickets/ticket", ticketID)
	assert.Nil(t, err)
	var msgList []Message
	msgList = append(msgList, Message{Actor: "client@dhbw.de", Text: "PC does not start anymore. Any idea?", CreationDate: actTicket.MessageList[0].CreationDate})
	expectedTicket = Ticket{XMLName: xml.Name{"", "Ticket"}, Id: ticketID, Client: "client@dhbw.de", Reference: "PC problem", Editor: "0", Status: 0, MessageList: msgList}
	assert.Equal(t, expectedTicket, actTicket)

	tmpTicket, err := ReadTicket("../data/tickets/ticket", -99)
	assert.NotNil(t, err)
	assert.Equal(t, tmpTicket, Ticket{})

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", ticketID)
	writeToXML(0, "definitions.xml")
}

/*
func TestIDCounter(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 2, getTicketIDCounter("definitions.xml"))

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	writeToXML(0, "definitions.xml")
}
*/

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

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	writeToXML(0, "definitions.xml")
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

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	writeToXML(0, "definitions.xml")
}

func TestChangeEditor(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeEditor("../data/tickets/ticket", getTicketIDCounter("definitions.xml"), "4321")
	ticket, _ := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml"))
	assert.Equal(t, "4321", ticket.Editor)

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	writeToXML(0, "definitions.xml")
}

func TestChangeStatus(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	ChangeStatus("../data/tickets/ticket", getTicketIDCounter("definitions.xml"), 2)
	ticket, _ := ReadTicket("../data/tickets/ticket", getTicketIDCounter("definitions.xml"))
	assert.Equal(t, 2, ticket.Status)

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	writeToXML(0, "definitions.xml")
}

func TestDeleting(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "Computer", "PC not working")
	assert.Equal(t, 1, len(ticketMap))
	DeleteTicket("../data/tickets/ticket", "definitions.xml", getTicketIDCounter("definitions.xml"))
	assert.Equal(t, 0, len(ticketMap))

	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "Computer", "PC not working")
	ClearCache("../data/tickets/ticket")
	assert.Equal(t, 0, len(ticketMap))
	_, err := ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.Nil(t, err)
	DeleteTicket("../data/tickets/ticket", "definitions.xml", 1)
	_, err = ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.NotNil(t, err)

	ClearCache("../data/tickets/ticket")
	writeToXML(0, "definitions.xml")
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
	assert.Equal(t, expectedTicket, actTicket)

	//merge tickets with two different editors
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "New employee", "Hello, please create a new login account for our new employee Erika Musterfrau. Thank you.")
	secondTicketID := getTicketIDCounter("definitions.xml")
	ChangeEditor("../data/tickets/ticket", secondTicketID, "412")
	assert.NotNil(t, MergeTickets("../data/tickets/ticket", "definitions.xml", firstTicket.Id, secondTicketID))

	ClearCache("../data/tickets/ticket")
	DeleteTicket("../data/tickets/ticket", "definitions.xml", firstTicket.Id)
	DeleteTicket("../data/tickets/ticket", "definitions.xml", secondTicket.Id)
	writeToXML(0, "definitions.xml")
}

func TestClearCache(t *testing.T) {
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 3, len(ticketMap))

	ClearCache("../data/tickets/ticket")
	assert.Equal(t, 0, len(ticketMap))

	_, err1 := ioutil.ReadFile("../data/tickets/ticket1.xml")
	assert.Nil(t, err1)

	DeleteTicket("../data/tickets/ticket", "definitions.xml", 1)
	DeleteTicket("../data/tickets/ticket", "definitions.xml", 2)
	DeleteTicket("../data/tickets/ticket", "definitions.xml", 3)
	ClearCache("../data/tickets/ticket")
	writeToXML(0, "definitions.xml")
}

// TODO: Wieso bleibt len(ticketMap) = 11 nachdem 2 weitere tickets erstellt werden? + Es werden bei mir nicht alle tickets gelöscht  -> Groese des Caches ist auf 11 festgelegt
func TestCheckCache(t *testing.T) {
	for tmpInt := 1; tmpInt <= 9; tmpInt++ {
		CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	}

	assert.Equal(t, 9, len(ticketMap))

	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 11, len(ticketMap))

	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	CreateTicket("../data/tickets/ticket", "definitions.xml", "client@dhbw.de", "PC problem", "PC does not start anymore. Any idea?")
	assert.Equal(t, 11, len(ticketMap))

	ClearCache("../data/tickets/ticket")
	for tmpInt := 1; tmpInt <= 13; tmpInt++ {
		DeleteTicket("../data/tickets/ticket", "definitions.xml", tmpInt)
	}
	writeToXML(0, "definitions.xml")
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
}

func TestVerifyUser(t *testing.T) {
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	tmpBool, err := VerifyUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.True(t, tmpBool)
	assert.Nil(t, err)
	tmpBool, err = VerifyUser("../data/users/users.xml", "mustermann", "xxx")
	assert.False(t, tmpBool)
	assert.Nil(t, err)
}

func TestLoginUser(t *testing.T) {
	CreateUser("../data/users/users.xml", "mustermann", "musterpasswort")
	assert.Nil(t, LoginUser("../data/users/users.xml", "mustermann", "musterpasswort", "1234"))
	assert.NotNil(t, LoginUser("../data/users/users.xml", "mustermann", "falschespasswort", "1234"))
	usersMap, _ := readUsers("../data/users/users.xml")
	assert.Equal(t, "1234", usersMap["mustermann"].SessionID)
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
	InitDataStorage("../data/tickets", "../data/users")
	ticketMap = make(map[int]Ticket)
	writeToXML(0, "definitions.xml")
}

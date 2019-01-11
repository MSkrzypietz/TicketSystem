package utils

import (
	"TicketSystem/config"
	"encoding/xml"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Ticket struct {
	XMLName     xml.Name  `xml:"Ticket"`
	ID          int       `xml:"ID"`
	Client      string    `xml:"ClientAddress"`
	Reference   string    `xml:"Subject"`
	Status      int       `xml:"Status"`
	Editor      string    `xml:"Editor"`
	MessageList []Message `xml:"MessageList>Message"`
}

type Message struct {
	CreationDate time.Time `xml:"CreationDate"`
	Actor        string    `xml:"Actor"`
	Text         string    `xml:"Text"`
}

const (
	TicketStatusOpen = iota
	TicketStatusInProcess
	TicketStatusClosed
)

var ticketMap = make(map[int]Ticket)

var mutexTicketID = &sync.Mutex{}

type User struct {
	Username    string `xml:"Username"`
	Password    string `xml:"Password"`
	SessionID   string `xml:"SessionID"`
	HolidayMode bool   `xml:"HolidayMode"`
}

type UserList struct {
	User []User `xml:"users>user"`
}

// Creates directory for the data storage if it does not exist
func InitDataStorage() error {
	_, err := os.Stat(config.TicketsPath())
	if err != nil {
		if os.IsNotExist(err) {
			tmpErr := os.MkdirAll(config.TicketsPath(), 0777)
			if tmpErr != nil {
				return tmpErr
			}
		}
	}

	_, err = os.Stat(config.UsersFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			tmpErr := os.MkdirAll(config.UsersPath(), 0777)
			if tmpErr != nil {
				return tmpErr
			}
		}
		f, err := os.Create(config.UsersFilePath())
		if err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(config.DefinitionsFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(config.DefinitionsFilePath())
			if err != nil {
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}
			err = WriteToXML(0, config.DefinitionsFilePath())
			if err != nil {
				return err
			}
		}
	}

	_, err = os.Stat(config.MailFilePath())
	if err != nil {
		if os.IsNotExist(err) {
			f, err := os.Create(config.MailFilePath())
			if err != nil {
				return err
			}
			err = f.Close()
			if err != nil {
				return err
			}
			return WriteToXML(MailList{}, config.MailFilePath())
		}
	}

	return err
}

// Creates a ticket from the inputs
func CreateTicket(client string, reference string, text string) (Ticket, error) {
	// Synchronizing this method to prevent multiple tickets with the same ID
	mutexTicketID.Lock()
	defer mutexTicketID.Unlock()

	IDCounter := getTicketIDCounter() + 1
	newTicket := Ticket{ID: IDCounter, Client: client, Reference: reference, Status: TicketStatusOpen}
	err := WriteToXML(IDCounter, config.DefinitionsFilePath())
	if err != nil {
		return Ticket{}, err
	}

	return AddMessage(newTicket, client, text)
}

// Adds a message to a specified tickets
func AddMessage(ticket Ticket, actor string, text string) (Ticket, error) {
	newMessage := Message{CreationDate: time.Now(), Actor: actor, Text: text}
	ticket.MessageList = append(ticket.MessageList, newMessage)
	return ticket, StoreTicket(ticket)
}

// Stores a ticket as a xml file
func StoreTicket(ticket Ticket) error {
	delete(ticketMap, ticket.ID)
	return WriteToXML(ticket, config.TicketXMLPath(ticket.ID))
}

// Returns the requested ticket from the cache or from the corresponding xml file
func ReadTicket(id int) (Ticket, error) {
	if ticketMap[id].ID != 0 {
		return ticketMap[id], nil
	}

	file, err := ioutil.ReadFile(config.TicketXMLPath(id))
	if err != nil {
		return Ticket{}, err
	}

	var ticket Ticket
	err = xml.Unmarshal(file, &ticket)
	if err != nil {
		return Ticket{}, err
	}

	err = checkCache()
	if err != nil {
		return Ticket{}, err
	}

	ticketMap[ticket.ID] = ticket
	return ticket, nil
}

// Deletes a ticket by its ID
func deleteTicket(id int) error {
	delete(ticketMap, id)

	err := os.Remove(config.TicketXMLPath(id))
	if err != nil {
		return err
	}

	return nil
}

// Changes the editor of a ticket
func ChangeEditor(id int, editor string) error {
	ticket, err := ReadTicket(id)
	if err != nil {
		return err
	}

	ticket.Editor = editor
	return StoreTicket(ticket)
}

// Changes the status of a ticket
func ChangeStatus(id int, status int) error {
	ticket, err := ReadTicket(id)
	if err != nil {
		return err
	}

	ticket.Status = status
	return StoreTicket(ticket)
}

// Returns a list of tickets by a specified ticket status
func GetTicketsByStatus(status int) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(); actualID++ {
		tmp, _ := ReadTicket(actualID)
		if tmp.Status == status && tmp.ID != 0 {
			tickets = append(tickets, tmp)
		}
	}

	return tickets
}

// Returns a list of tickets owned by the specified editor
func GetTicketsByEditor(editor string) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(); actualID++ {
		tmp, _ := ReadTicket(actualID)
		if tmp.Editor == editor && tmp.ID != 0 {
			tickets = append(tickets, tmp)
		}
	}

	return tickets
}

// Returns a list of tickets owned by the specified client
func GetTicketsByClient(client string) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(); actualID++ {
		tmp, _ := ReadTicket(actualID)
		if tmp.Client == client && tmp.ID != 0 {
			tickets = append(tickets, tmp)
		}
	}

	return tickets
}

// Returns the current ticket ID. Unexpected errors will return -1
func getTicketIDCounter() int {
	file, err := ioutil.ReadFile(config.DefinitionsFilePath())
	if err != nil {
		return -1
	}

	var IDCounter int
	err = xml.Unmarshal(file, &IDCounter)
	if err != nil {
		return -1
	}

	return IDCounter
}

// Merges two tickets, store them as one ticket and delete the other one
func MergeTickets(firstTicketID int, secondTicketID int) error {
	firstTicket, err := ReadTicket(firstTicketID)
	if err != nil {
		return err
	}

	secondTicket, err := ReadTicket(secondTicketID)
	if err != nil {
		return err
	}

	if firstTicket.Editor != secondTicket.Editor {
		return fmt.Errorf("the two tickets for the merging process do not have the same editors")
	}

	for _, msgList := range secondTicket.MessageList {
		firstTicket.MessageList = append(firstTicket.MessageList, msgList)
	}

	err = ChangeStatus(firstTicketID, TicketStatusInProcess)
	if err != nil {
		return err
	}

	err = deleteTicket(secondTicketID)
	if err != nil {
		return err
	}

	return StoreTicket(firstTicket)
}

// Writes an object to the specified xml file
func WriteToXML(v interface{}, path string) error {
	content, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}

	content = []byte(xml.Header + string(content))
	err = ioutil.WriteFile(path, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Checks if there are too many tickets in the cache and in the case of too many tickets one will be deleted
func checkCache() error {
	if len(ticketMap) > 9 {
		randNumber := rand.Intn(len(ticketMap))
		tmpInt := 1
		for _, tmpTicket := range ticketMap {
			if tmpInt == randNumber {
				delete(ticketMap, tmpTicket.ID)
				return nil
			}
			tmpInt++
		}
		return fmt.Errorf("no ticket found in the cache")
	}

	return nil
}

// Creates a new user
func CreateUser(name string, password string) (User, error) {
	usersMap, err := ReadUsers()
	if err != nil {
		return User{}, err
	}

	if _, ok := usersMap[name]; ok {
		return User{}, fmt.Errorf("a user with the same name already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return User{}, err
	}

	usersMap[name] = User{Username: name, Password: string(hash), SessionID: "", HolidayMode: false}
	err = storeUsers(usersMap)
	if err != nil {
		return User{}, err
	}

	return usersMap[name], nil
}

// Returns all users
func ReadUsers() (map[string]User, error) {
	usersMap := make(map[string]User)

	file, err := ioutil.ReadFile(config.UsersFilePath())
	if err != nil {
		return usersMap, err
	}

	var userList UserList
	err = xml.Unmarshal(file, &userList)
	if err != nil && err != io.EOF {
		return usersMap, err
	}

	for _, tmpUser := range userList.User {
		usersMap[tmpUser.Username] = tmpUser
	}

	return usersMap, nil
}

// Stores all users from the map to the xml file
func storeUsers(userMap map[string]User) error {
	var users []User
	for _, tmpUser := range userMap {
		users = append(users, tmpUser)
	}

	return WriteToXML(UserList{User: users}, config.UsersFilePath())
}

// Checks if the user is registered
func CheckUser(name string) (bool, error) {
	usersMap, err := ReadUsers()
	if err != nil {
		return false, err
	}
	if usersMap[name].Username == name {
		return false, nil
	}
	return true, nil
}

// Checks if the hashed password is the correct password
func VerifySessionCookie(name string, password string) error {
	usersMap, err := ReadUsers()
	if err != nil {
		return err
	}

	if password != usersMap[name].Password {
		return fmt.Errorf("invalid session cookie")
	}

	return nil
}

// Login of a user to the ticket system
func LoginUser(name string, password string, session string) error {
	usersMap, err := ReadUsers()
	if err != nil {
		return fmt.Errorf("wrong path to user file")
	}

	err = bcrypt.CompareHashAndPassword([]byte(usersMap[name].Password), []byte(password))
	if err != nil {
		return err
	}

	tmpUser := usersMap[name]
	tmpUser.SessionID = session
	usersMap[name] = tmpUser
	return storeUsers(usersMap)
}

// Logout of a user and deletes the session id
func LogoutUser(name string) error {
	usersMap, err := ReadUsers()
	if err != nil {
		return err
	}

	if usersMap[name].Username != name {
		return fmt.Errorf("user does not exist")
	}

	tmpUser := usersMap[name]
	tmpUser.SessionID = ""
	usersMap[name] = tmpUser
	return storeUsers(usersMap)
}

// Returns the current session id of an user
func GetUserSession(name string) string {
	usersMap, _ := ReadUsers()
	return usersMap[name].SessionID
}

// Returns a user by the specified session id
func GetUserBySession(session string) (User, error) {
	if session == "" {
		return User{}, fmt.Errorf("session is not set")
	}

	usersMap, _ := ReadUsers()
	for _, tmpUser := range usersMap {
		if tmpUser.SessionID == session {
			return tmpUser, nil
		}
	}

	return User{}, fmt.Errorf("user does not exist")
}

// Sets the holiday mode of the specified user
func SetUserHolidayMode(name string, holidayMode bool) error {
	tmpUsers, err := ReadUsers()
	if err != nil {
		return err
	}

	user := tmpUsers[name]
	if user.Username == "" {
		return fmt.Errorf("user does not exist")
	}

	user.HolidayMode = holidayMode
	tmpUsers[name] = user
	err = storeUsers(tmpUsers)
	return err
}

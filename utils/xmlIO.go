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
	"time"
)

//struct that defines a ticket with the parameters ID, mail of the client, reference, actual status, editor and a list of all messages
type Ticket struct {
	XMLName     xml.Name  `xml:"Ticket"`
	Id          int       `xml:"ID"`
	Client      string    `xml:"Client"`
	Reference   string    `xml:"Reference"`
	Status      int       `xml:"Status"`
	Editor      string    `xml:"Editor"`
	MessageList []Message `xml:"MessageList>Message"`
}

//struct that defines a message with the parameters date of creation, name or mail of the actor and the text of the message
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

//struct for one user
type User struct {
	Username  string `xml:"Username"`
	Password  string `xml:"Password"`
	SessionID string `xml:"SessionID"`
}

type Userlist struct {
	User []User `xml:"users>user"`
}

//creates directory for the data storage if it doesn´t exist
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
			return WriteToXML(Maillist{}, config.MailFilePath())
		}
	}

	return err
}

//function to create a ticket including the following parameters: mail of the client, reference and text of the ticket. Returns the ticket struct and an error whether the creation was successful.
func CreateTicket(client string, reference string, text string) (Ticket, error) {
	IDCounter := getTicketIDCounter() + 1
	newTicket := Ticket{Id: IDCounter, Client: client, Reference: reference, Status: TicketStatusOpen}
	err := WriteToXML(IDCounter, config.DefinitionsFilePath())
	if err != nil {
		return Ticket{}, err
	}
	return AddMessage(newTicket, client, text)
}

//adds a message to a specified tickets. Functions includes the following parameters: specified ticket, the actor and the text of the message. Returns the new ticket and an error whether it was successful.
func AddMessage(ticket Ticket, actor string, text string) (Ticket, error) {
	newMessage := Message{CreationDate: time.Now(), Actor: actor, Text: text}
	ticket.MessageList = append(ticket.MessageList, newMessage)
	return ticket, StoreTicket(ticket)
}

//stores a ticket as xml file
func StoreTicket(ticket Ticket) error {
	delete(ticketMap, ticket.Id)
	return WriteToXML(ticket, config.TicketXMLPath(ticket.Id))
}

//returns a ticket from the cache or from the corresponding xml file.
func ReadTicket(id int) (Ticket, error) {
	if ticketMap[id].Id != 0 {
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
	ticketMap[ticket.Id] = ticket
	return ticket, nil
}

//deletes a ticket by its ID and returns an error whether it was successful.
func DeleteTicket(id int) error {
	delete(ticketMap, id)
	if id == getTicketIDCounter() {
		err := WriteToXML(id-1, config.DefinitionsFilePath())
		if err != nil {
			return err
		}
	}
	err := os.Remove(config.TicketXMLPath(id))
	if err != nil {
		return err
	}
	return nil
}

//changes the editor of a ticket and returns an error whether the change was successful.
func ChangeEditor(id int, editor string) error {
	ticket, err := ReadTicket(id)
	if err != nil {
		return err
	}
	ticket.Editor = editor
	return StoreTicket(ticket)
}

//changes the status of a ticket and returns an error whether the change was successful.
func ChangeStatus(id int, status int) error {
	ticket, err := ReadTicket(id)
	if err != nil {
		return err
	}
	ticket.Status = status
	return StoreTicket(ticket)
}

//returns a list of tickets by a specified ticket status. Status is specified in the parameters of the function.
func GetTicketsByStatus(status int) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(); actualID++ {
		tmp, _ := ReadTicket(actualID)
		if tmp.Status == status && tmp.Id != 0 {
			tickets = append(tickets, tmp)
		}
	}
	return tickets
}

//returns a list of tickets owned by one editor who is specified in the parameters of the function
func GetTicketsByEditor(editor string) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(); actualID++ {
		tmp, _ := ReadTicket(actualID)
		if tmp.Editor == editor && tmp.Id != 0 {
			tickets = append(tickets, tmp)
		}
	}
	return tickets
}

//returns a list of tickets owned by a client who is specified in the parameters of the function
func GetTicketsByClient(client string) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(); actualID++ {
		tmp, _ := ReadTicket(actualID)
		if tmp.Client == client && tmp.Id != 0 {
			tickets = append(tickets, tmp)
		}
	}
	return tickets
}

//returns the actual ticket ID in order to create a new ticket or to get to know the number of the stored tickets.
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

//merge two tickets, store them as one ticket and delete the other one. Returns an error whether the merge was successful.
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
	err = DeleteTicket(secondTicketID)
	if err != nil {
		return err
	}
	return StoreTicket(firstTicket)
}

//functions writes an object to an specified xml file and returns an error whether the writing was successful.
func WriteToXML(v interface{}, path string) error {
	xmlstring, err := xml.MarshalIndent(v, "", "    ")
	if err != nil {
		return err
	}
	xmlstring = []byte(xml.Header + string(xmlstring))
	err = ioutil.WriteFile(path, xmlstring, 0644)
	if err != nil {
		return err
	}
	return nil
}

//function checks if there are too many tickets in the cache and in the case of too many tickets one will be deleted. Returns an error whether it was successful.
func checkCache() error {
	if len(ticketMap) > 9 {
		randNumber := rand.Intn(len(ticketMap))
		tmpInt := 1
		for _, tmpTicket := range ticketMap {
			if tmpInt == randNumber {
				delete(ticketMap, tmpTicket.Id)
				return nil
			}
			tmpInt++
		}
		return fmt.Errorf("no ticket found in the cache")
	} else {
		return nil
	}
}

//creates a new user and returns the user and an error whether the creation was successful.
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
	usersMap[name] = User{Username: name, Password: string(hash), SessionID: ""}
	err = storeUsers(usersMap)
	if err != nil {
		return User{}, err
	}
	return usersMap[name], nil
}

//reads all users from the xml-file and returns the users and an error whether the reading process was successful.
func ReadUsers() (map[string]User, error) {
	usersMap := make(map[string]User)
	file, err := ioutil.ReadFile(config.UsersFilePath())
	if err != nil {
		return usersMap, err
	}
	var userlist Userlist
	err = xml.Unmarshal(file, &userlist)
	if err != nil && err != io.EOF {
		return usersMap, err
	}
	for _, tmpUser := range userlist.User {
		usersMap[tmpUser.Username] = tmpUser
	}
	return usersMap, nil
}

//stores all users from the map to the xml file and returns an error whether the storing process was successful.
func storeUsers(usermap map[string]User) error {
	var users []User
	for _, tmpUser := range usermap {
		users = append(users, tmpUser)
	}
	return WriteToXML(Userlist{User: users}, config.UsersFilePath())
}

//checks if the user is registrated and returns a bool. The bool value is false if there is already a user with that name.
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

//checks if the username and the password is correct. Returns a bool whether it is correct.
func VerifyUser(name string, password string) (bool, error) {
	usersMap, err := ReadUsers()
	if err != nil {
		return false, err
	}

	if password != usersMap[name].Password {
		return false, fmt.Errorf("passwords don't match. Cannot verify the user")
	}

	return true, nil
}

//Login of a user to the ticket system. Returns an error if an error occurs.
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

//Logout of a user and deletes the session id. Returns an error if an error occurs.
func LogoutUser(name string) error {
	usersmap, err := ReadUsers()
	if err != nil {
		return err
	}
	if usersmap[name].Username != name {
		return fmt.Errorf("user does not exist")
	}
	tmpUser := usersmap[name]
	tmpUser.SessionID = ""
	usersmap[name] = tmpUser
	return storeUsers(usersmap)
}

//gets the actual session id of an user
func GetUserSession(name string) string {
	usersMap, _ := ReadUsers()
	return usersMap[name].SessionID
}

//returns an user by a specified session id
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

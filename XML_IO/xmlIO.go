package XML_IO

import (
	"TicketSystem/config"
	"encoding/xml"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

//TODO: (Kleinigkeit - nicht dringend): In users.xml gibt es 2 "root" Elemente? <UserList> oder <users> kann entfernt werden

//TODO: Error handling bei allen Warnungen

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

//TODO: Handling of the errors of mkdirall and writetoxml
//creates directory for the data storage if it doesn´t exist
func InitDataStorage(ticketPath string, usersPath string) {
	_, err := os.Stat(ticketPath)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(ticketPath, 0777)
	}
	_, err = os.Stat(path.Join(usersPath, "users.xml"))
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(usersPath, 0777)
		writeToXML(nil, config.UsersFilePath())
	}
}

//function to create a ticket including the following parameters: mail of the client, reference and text of the ticket. Returns the ticket struct and an error whether the creation was successful.
func CreateTicket(path string, definitionsPath string, client string, reference string, text string) (Ticket, error) {
	IDCounter := getTicketIDCounter(definitionsPath) + 1
	newTicket := Ticket{Id: IDCounter, Client: client, Reference: reference, Status: 0, Editor: "0"}
	err := writeToXML(IDCounter, definitionsPath)
	if err != nil {
		return Ticket{}, err
	}
	return AddMessage(path, newTicket, client, text)
}

//adds a message to a specified tickets. Functions includes the following parameters: specified ticket, the actor and the text of the message. Returns the new ticket and an error whether it was successful.
func AddMessage(path string, ticket Ticket, actor string, text string) (Ticket, error) {
	newMessage := Message{CreationDate: time.Now(), Actor: actor, Text: text}
	ticket.MessageList = append(ticket.MessageList, newMessage)
	return ticket, StoreTicket(path, ticket)
}

//stores a ticket as xml file
func StoreTicket(path string, ticket Ticket) error {
	delete(ticketMap, ticket.Id)
	return writeToXML(ticket, path+strconv.Itoa(ticket.Id)+".xml")
}

//returns a ticket from the cache or from the corresponding xml file.
func ReadTicket(path string, id int) (Ticket, error) {
	if ticketMap[id].Id != 0 {
		return ticketMap[id], nil
	}

	file, err := ioutil.ReadFile(path + strconv.Itoa(id) + ".xml")
	if err != nil {
		return Ticket{}, err
	}
	var ticket Ticket
	xml.Unmarshal(file, &ticket)
	checkCache()
	ticketMap[ticket.Id] = ticket
	return ticket, nil
}

//deletes a ticket by its ID and returns an error whether it was successful.
func DeleteTicket(path string, definitionsPath string, id int) error {
	delete(ticketMap, id)
	if id == getTicketIDCounter(definitionsPath) {
		writeToXML(id-1, definitionsPath)
	}
	err := os.Remove(path + strconv.Itoa(id) + ".xml")
	if err != nil {
		return err
	}
	return nil
}

//changes the editor of a ticket and returns an error whether the change was successful.
func ChangeEditor(path string, id int, editor string) error {
	ticket, err := ReadTicket(path, id)
	if err != nil {
		return err
	}
	ticket.Editor = editor
	return StoreTicket(path, ticket)
}

//changes the status of a ticket and returns an error whether the change was successful.
func ChangeStatus(path string, id int, status int) error {
	ticket, err := ReadTicket(path, id)
	if err != nil {
		return err
	}
	ticket.Status = status
	return StoreTicket(path, ticket)
}

//returns a list of tickets by a specified ticket status. Status is specified in the parameters of the function.
func GetTicketsByStatus(path string, definitionsPath string, status int) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(definitionsPath); actualID++ {
		tmp, _ := ReadTicket(path, actualID)
		if tmp.Status == status && tmp.Id != 0 {
			tickets = append(tickets, tmp)
		}
	}
	return tickets
}

//returns a list of tickets owned by one editor who is specified in the parameters of the function
func GetTicketsByEditor(path string, definitionsPath string, editor string) []Ticket {
	var tickets []Ticket
	for actualID := 1; actualID <= getTicketIDCounter(definitionsPath); actualID++ {
		tmp, _ := ReadTicket(path, actualID)
		if tmp.Editor == editor && tmp.Id != 0 {
			tickets = append(tickets, tmp)
		}
	}
	return tickets
}

//returns the actual ticket ID in order to create a new ticket or to get to know the number of the stored tickets.
func getTicketIDCounter(definitionsPath string) int {
	file, err := ioutil.ReadFile(definitionsPath)
	if err != nil {
		return -1
	}
	var IDCounter int
	xml.Unmarshal(file, &IDCounter)
	return IDCounter
}

//merge two tickets, store them as one ticket and delete the other one. Returns an error whether the merge was successful.
func MergeTickets(path string, definitionsPath string, firstTicketID int, secondTicketID int) error {
	firstTicket, err1 := ReadTicket(path, firstTicketID)
	secondTicket, err2 := ReadTicket(path, secondTicketID)
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	if firstTicket.Editor != secondTicket.Editor {
		return errors.New("the two tickets for the merging process do not have the same editors")
	}
	for _, msgList := range secondTicket.MessageList {
		firstTicket.MessageList = append(firstTicket.MessageList, msgList)
	}
	DeleteTicket(path, definitionsPath, secondTicketID)
	return StoreTicket(path, firstTicket)
}

//functions writes an object to an specified xml file and returns an error whether the writing was successful.
func writeToXML(v interface{}, path string) error {
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
	if len(ticketMap) > 10 {
		randNumber := rand.Intn(len(ticketMap))
		tmpInt := 1
		for _, tmpTicket := range ticketMap {
			if tmpInt == randNumber {
				delete(ticketMap, tmpTicket.Id)
				return nil
			}
			tmpInt++
		}
		return errors.New("no ticket found in the cache")
	} else {
		return nil
	}
}

//creates a new user and returns the user and an error whether the creation was successful.
func CreateUser(path string, name string, password string) (User, error) {
	usersMap, err := readUsers(path)
	if err != nil {
		return User{}, err
	}
	usersMap[name] = User{Username: name, Password: password, SessionID: ""}
	err = storeUsers(path, usersMap)
	if err != nil {
		return User{}, err
	}
	return usersMap[name], nil
}

//reads all users from the xml-file and returns the users and an error whether the reading process was successful.
func readUsers(path string) (map[string]User, error) {
	usersMap := make(map[string]User)
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return usersMap, err
	}
	var userlist Userlist
	xml.Unmarshal(file, &userlist)
	for _, tmpUser := range userlist.User {
		usersMap[tmpUser.Username] = tmpUser
	}
	return usersMap, nil
}

//stores all users from the map to the xml file and returns an error whether the storing process was successful.
func storeUsers(path string, usermap map[string]User) error {
	var users []User
	for _, tmpUser := range usermap {
		users = append(users, tmpUser)
	}
	return writeToXML(Userlist{User: users}, path)
}

//checks if the user is registrated and returns a bool. The bool value is false if there is already a user with that name.
func CheckUser(path string, name string) (bool, error) {
	usersMap, err := readUsers(path)
	if err != nil {
		return false, err
	}
	if usersMap[name].Username == name {
		return false, nil
	}
	return true, nil
}

//checks if the username and the password is correct. Returns a bool whether it is correct.
func VerifyUser(path string, name string, password string) (bool, error) {
	usersMap, err := readUsers(path)
	if err != nil {
		return false, err
	}
	if usersMap[name].Password == password {
		return true, nil
	}
	return false, nil
}

//Login of a user to the ticket system. Returns an error if an error occurs.
func LoginUser(path string, name string, password string, session string) error {
	usersMap, err := readUsers(path)
	if err != nil {
		return errors.New("wrong path to user file")
	}
	err = bcrypt.CompareHashAndPassword([]byte(usersMap[name].Password), []byte(password))
	if err != nil {
		return err
	}
	tmpUser := usersMap[name]
	tmpUser.SessionID = session
	usersMap[name] = tmpUser
	return storeUsers(path, usersMap)
}

//Logout of a user and deletes the session id. Returns an error if an error occurs.
func LogoutUser(path string, name string) error {
	usersmap, err := readUsers(path)
	if err != nil {
		return err
	}
	if usersmap[name].Username != name {
		return errors.New("user does not exist")
	}
	tmpUser := usersmap[name]
	tmpUser.SessionID = ""
	usersmap[name] = tmpUser
	return storeUsers(path, usersmap)
}

//gets the actual session id of an user
func GetUserSession(path string, name string) string {
	usersMap, _ := readUsers(path)
	return usersMap[name].SessionID
}

//returns an user by a specified session id
func GetUserBySession(path string, session string) User {
	if session == "" {
		return User{}
	}
	usersMap, _ := readUsers(path)
	for _, tmpUser := range usersMap {
		if tmpUser.SessionID == session {
			return tmpUser
		}
	}
	return User{}
}

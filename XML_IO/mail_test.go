package XML_IO

import (
	"TicketSystem/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTicketFromMail(t *testing.T) {
	removeCompleteDataStorage()
	assert := assert.New(t)
	config.DataPath = "../datatest"

	expectedTicket, err := CreateTicketFromMail("mail@test", "testCaption", "testMsg")
	assert.Nil(err)
	actTicket, err := ReadTicket(expectedTicket.Id)
	assert.Nil(err)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()

	tmpTicket, _ := CreateTicket("test@mail", "testCaption", "testMsgOne")
	ChangeStatus(tmpTicket.Id, Closed)
	expectedTicket, err = CreateTicketFromMail("test@mail", "testCaption", "testMsgTwo")
	assert.Nil(err)
	actTicket, err = ReadTicket(expectedTicket.Id)
	expectedTicket.XMLName.Local = ""
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	actTicket.MessageList[1].CreationDate = expectedTicket.MessageList[1].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()

	tmpTicket, _ = CreateTicket("test@mail", "testCaption", "testMsgOne")
	expectedTicket, err = CreateTicketFromMail("test@mail", "testCaption", "testMsgTwo")
	assert.Nil(err)
	actTicket, err = ReadTicket(expectedTicket.Id)
	expectedTicket.XMLName.Local = ""
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	actTicket.MessageList[1].CreationDate = expectedTicket.MessageList[1].CreationDate
	assert.Equal(expectedTicket, actTicket)
	removeCompleteDataStorage()
}

func TestDeleteMails(t *testing.T) {
	assert := assert.New(t)
	config.DataPath = "../datatest"
	SendMail("mail@test", "captionOne", "test")
	SendMail("mail@test", "captionTwo", "test")
	var idField []int
	idField = append(idField, 1)
	DeleteMails(idField)
	actMaillist, _ := readMailsFile()
	var expectedMaillist []Mail
	expectedMaillist = append(expectedMaillist, Mail{"mail@test", "captionTwo", "test", 2})
	assert.Equal(Maillist{1, expectedMaillist}, actMaillist)
	removeCompleteDataStorage()
}

func TestSendMail(t *testing.T) {
	removeCompleteDataStorage()
	assert := assert.New(t)
	config.DataPath = "wrongPath"
	assert.NotNil(SendMail("", "", ""))
	config.DataPath = "../datatest"
	assert.Nil(SendMail("test@test", "testCaption", "testMsg"))
	var expectedMaillist []Mail
	expectedMaillist = append(expectedMaillist, Mail{"test@test", "testCaption", "testMsg", 1})
	actMaillist, err := readMailsFile()
	assert.Nil(err)
	assert.Equal(Maillist{1, expectedMaillist}, actMaillist)
	removeCompleteDataStorage()
}

func TestReadMailsFile(t *testing.T) {
	assert := assert.New(t)
	config.DataPath = "wrongPath"
	_, err := readMailsFile()
	assert.NotNil(err)
	config.DataPath = "../datatest"

	var mails []Mail
	mails = append(mails, Mail{"test@test", "testOne", "testOne", 1})
	mails = append(mails, Mail{"test@test", "testTwo", "testTwo", 2})
	expectedMaillist := Maillist{1, mails}
	WriteToXML(expectedMaillist, config.MailFilePath())
	actMaillist, err := readMailsFile()
	assert.Nil(err)
	assert.Equal(expectedMaillist, actMaillist)
	removeCompleteDataStorage()
}

package XML_IO

import (
	"TicketSystem/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateTicketFromMail(t *testing.T) {
	setup()
	defer teardown()

	expectedTicket, err := CreateTicketFromMail("mail@test", "testCaption", "testMsg")
	assert.Nil(t, err)
	actTicket, err := ReadTicket(expectedTicket.Id)
	assert.Nil(t, err)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(t, expectedTicket, actTicket)
	teardown()
	setup()

	tmpTicket, _ := CreateTicket("test@mail", "testCaption", "testMsgOne")
	assert.Nil(t, ChangeStatus(tmpTicket.Id, Closed))
	expectedTicket, err = CreateTicketFromMail("test@mail", "testCaption", "testMsgTwo")
	assert.Nil(t, err)
	actTicket, err = ReadTicket(expectedTicket.Id)
	expectedTicket.XMLName.Local = ""
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	actTicket.MessageList[1].CreationDate = expectedTicket.MessageList[1].CreationDate
	assert.Equal(t, expectedTicket, actTicket)
	teardown()
	setup()

	tmpTicket, _ = CreateTicket("test@mail", "testCaption", "testMsgOne")
	expectedTicket, err = CreateTicketFromMail("test@mail", "testCaption", "testMsgTwo")
	assert.Nil(t, err)
	actTicket, err = ReadTicket(expectedTicket.Id)
	expectedTicket.XMLName.Local = ""
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	actTicket.MessageList[1].CreationDate = expectedTicket.MessageList[1].CreationDate
	assert.Equal(t, expectedTicket, actTicket)
}

func TestDeleteMails(t *testing.T) {
	setup()
	defer teardown()

	assert.Nil(t, SendMail("mail@test", "captionOne", "test"))
	assert.Nil(t, SendMail("mail@test", "captionTwo", "test"))
	var idField []int
	idField = append(idField, 1)
	assert.Nil(t, DeleteMails(idField))
	actMaillist, _ := readMailsFile()
	var expectedMaillist []Mail
	expectedMaillist = append(expectedMaillist, Mail{"mail@test", "captionTwo", "test", 2})
	assert.Equal(t, Maillist{1, expectedMaillist}, actMaillist)
}

func TestSendMail(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	assert.NotNil(t, SendMail("", "", ""))

	config.DataPath = "datatest"
	assert.Nil(t, SendMail("test@test", "testCaption", "testMsg"))

	var expectedMaillist []Mail
	expectedMaillist = append(expectedMaillist, Mail{"test@test", "testCaption", "testMsg", 1})
	actMaillist, err := readMailsFile()
	assert.Nil(t, err)
	assert.Equal(t, Maillist{1, expectedMaillist}, actMaillist)
}

func TestReadMailsFile(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	_, err := readMailsFile()
	assert.NotNil(t, err)
	config.DataPath = "datatest"

	var mails []Mail
	mails = append(mails, Mail{"test@test", "testOne", "testOne", 1})
	mails = append(mails, Mail{"test@test", "testTwo", "testTwo", 2})
	expectedMaillist := Maillist{1, mails}
	assert.Nil(t, WriteToXML(expectedMaillist, config.MailFilePath()))
	actMaillist, err := readMailsFile()
	assert.Nil(t, err)
	assert.Equal(t, expectedMaillist, actMaillist)
}

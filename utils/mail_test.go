package utils

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
	actTicket, err := ReadTicket(expectedTicket.ID)
	assert.Nil(t, err)
	actTicket.XMLName.Local = ""
	actTicket.MessageList[0].CreationDate = expectedTicket.MessageList[0].CreationDate
	assert.Equal(t, expectedTicket, actTicket)
	teardown()
	setup()

	tmpTicket, _ := CreateTicket("test@mail", "testCaption", "testMsgOne")
	assert.Nil(t, ChangeStatus(tmpTicket.ID, TicketStatusClosed))
	expectedTicket, err = CreateTicketFromMail("test@mail", "testCaption", "testMsgTwo")
	assert.Nil(t, err)
	actTicket, err = ReadTicket(expectedTicket.ID)
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
	actTicket, err = ReadTicket(expectedTicket.ID)
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
	actMailList, _ := ReadMailsFile()
	var expectedMailList []Mail
	expectedMailList = append(expectedMailList, Mail{Mail: "mail@test", Subject: "captionTwo", Message: "test", ID: 2})
	assert.Equal(t, MailList{2, expectedMailList}, actMailList)
}

func TestSendMail(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	assert.NotNil(t, SendMail("", "", ""))

	config.DataPath = "datatest"
	assert.Nil(t, SendMail("test@test", "testCaption", "testMsg"))

	var expectedMailList []Mail
	expectedMailList = append(expectedMailList, Mail{Mail: "test@test", Subject: "testCaption", Message: "testMsg", ID: 1})
	actMailList, err := ReadMailsFile()
	assert.Nil(t, err)
	assert.Equal(t, MailList{1, expectedMailList}, actMailList)
}

func TestReadMailsFile(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	_, err := ReadMailsFile()
	assert.NotNil(t, err)
	config.DataPath = "datatest"

	var mails []Mail
	mails = append(mails, Mail{Mail: "test@test", Subject: "testOne", Message: "testOne", ID: 1})
	mails = append(mails, Mail{Mail: "test@test", Subject: "testTwo", Message: "testTwo", ID: 2})
	expectedMailList := MailList{1, mails}
	assert.Nil(t, WriteToXML(expectedMailList, config.MailFilePath()))
	actMailList, err := ReadMailsFile()
	assert.Nil(t, err)
	assert.Equal(t, expectedMailList, actMailList)
}

func TestIncrementReadAttemptsCounter(t *testing.T) {
	setup()
	defer teardown()

	config.DataPath = "wrongPath"
	err := (&Mail{}).IncrementReadAttemptsCounter()
	assert.NotNil(t, err)
	config.DataPath = "datatest"

	testMail := Mail{Mail: "test@test", Subject: "testOne", Message: "testOne", ID: 1}
	err = testMail.IncrementReadAttemptsCounter()
	assert.NotNil(t, err)

	mailList := MailList{1, []Mail{testMail}}
	err = WriteToXML(mailList, config.MailFilePath())
	assert.Nil(t, err)
	err = testMail.IncrementReadAttemptsCounter()
	assert.Nil(t, err)
	assert.Equal(t, 1, testMail.ReadAttemptCounter)
}

package test

import (
	"bytes"
	"github.com/teamhanko/hanko/backend/mail"
	"gopkg.in/gomail.v2"
)

var emailsArray [][]gomail.Message

type TestMailer struct {
	index  int
	mailer mail.Mailer
}

func NewTestMailer() *TestMailer {
	emailsArray = append(emailsArray, []gomail.Message{})
	mailer := newMailer(len(emailsArray) - 1)
	return &TestMailer{mailer: mailer, index: len(emailsArray) - 1}
}

func newMailer(index int) mail.Mailer {
	return &mailer{index: index}
}

type mailer struct {
	index int
}

func (m mailer) Send(message *gomail.Message) error {
	emailsArray[m.index] = append(emailsArray[m.index], *message)
	return nil
}

func (m TestMailer) GetMailer() mail.Mailer {
	return m.mailer
}

func (m TestMailer) GetEmailCount() int {
	return len(emailsArray[m.index])
}

func (m TestMailer) GetLatestEmail() *gomail.Message {
	emails := emailsArray[m.index]
	if len(emails) == 0 {
		return nil
	}
	return &emails[len(emails)-1]
}

func (m TestMailer) GetEmails() []gomail.Message {
	return emailsArray[m.index]
}

type EmailAttributes struct {
	To       string
	From     string
	Subject  string
	Contents string
}

func (m TestMailer) GetEmailAttributes(message gomail.Message) EmailAttributes {
	buf := new(bytes.Buffer)
	message.WriteTo(buf)
	attributes := EmailAttributes{
		To:       message.GetHeader("To")[0],
		From:     message.GetHeader("From")[0],
		Subject:  message.GetHeader("Subject")[0],
		Contents: buf.String(),
	}

	return attributes
}

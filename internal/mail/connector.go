package mail

import (
	"ToDoProject/internal/mail/smtp"
	"ToDoProject/internal/mail/structures"
	"github.com/labstack/gommon/log"
	"os"
)

var (
	mailDriver = os.Getenv("MAIL_DRIVER")
)

type Mailer struct {
	SendMailFunc   func(to string, subject string, body structures.MessagesData) (err error)
	MailRecipients chan structures.EmailRecipient
}

func (m *Mailer) Init() {
	m.MailRecipients = make(chan structures.EmailRecipient)
	switch mailDriver {
	case "smtp":
		m.SendMailFunc = smtp.SendMail
	default:
		panic("unrecognized mail driver")
	}
}

func (m *Mailer) StartHandling() {
	for r := range m.MailRecipients {
		go func(r structures.EmailRecipient) {
			err := m.SendMailFunc(r.Email, r.Subject, r.Messages)
			if err != nil {
				log.Errorf("Error sending email to %s: %v\n", r.Email, err)
			}
		}(r)
	}
}

func (m *Mailer) QueueEmail(to string, subject string, messages structures.MessagesData) {
	m.MailRecipients <- structures.EmailRecipient{to, subject, messages}
}

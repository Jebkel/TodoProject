package smtp

import (
	"ToDoProject/internal/mail/structures"
	"bytes"
	"crypto/tls"
	"github.com/labstack/gommon/log"
	"gopkg.in/gomail.v2"
	"html/template"
	"os"
	"strconv"
	"time"
)

var (
	host       = os.Getenv("MAIL_HOST")
	port       = os.Getenv("MAIL_PORT")
	login      = os.Getenv("MAIL_LOGIN")
	password   = os.Getenv("MAIL_PASSWORD")
	cryptType  = os.Getenv("MAIL_CRYPTO_TYPE")
	fromMailer = os.Getenv("MAIL_FROM")
)

func SendMail(to string, subject string, messages structures.MessagesData) (err error) {
	port, err := strconv.Atoi(port)
	if err != nil {
		log.Error(err)
		return err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", fromMailer)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)

	t := template.New("template.html")
	t, err = t.ParseFiles("assets/email/template.html")
	if err != nil {
		log.Error(err)
		return err
	}
	var tpl bytes.Buffer

	if err := t.Execute(&tpl, messages); err != nil {
		log.Error(err)
		return err
	}

	msg.SetBody("text/html", tpl.String())

	d := gomail.NewDialer(host, port, login, password)
	if cryptType == "tls" {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	for i := 0; i < 3; i++ {
		if err = d.DialAndSend(msg); err == nil {
			return nil
		}
		time.Sleep(time.Minute)
	}
	log.Error(err)
	return err
}

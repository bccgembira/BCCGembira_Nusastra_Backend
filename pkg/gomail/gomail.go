package gomail

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/CRobinDev/BCCGembira_Nusastra/internal/dto"
	"github.com/CRobinDev/BCCGembira_Nusastra/pkg/response"
	"gopkg.in/gomail.v2"
)

type Gomail struct {
	message  *gomail.Message
	dialer   *gomail.Dialer
	htmlPath string
}

func NewGomail() *Gomail {
	port := os.Getenv("EMAIL_PORT")
	intPort, err := strconv.Atoi(port)
	if err != nil {
		log.Printf("Error converting port to integer: %v", err)
	}

	return &Gomail{
		message:  gomail.NewMessage(),
		dialer:   gomail.NewDialer(os.Getenv("EMAIL_HOST"), intPort, os.Getenv("SENDER"), os.Getenv("PASSWORD")),
		htmlPath: os.Getenv("HTML_PATH"),
	}
}

func (g *Gomail) SetBodyHTML(path string, data interface{}) (string, error) {
	var body bytes.Buffer
	t, err := template.ParseFiles(g.htmlPath + path)
	if err != nil {
		return "", &response.ErrSetHTML
	}

	err = t.Execute(&body, data)
	if err != nil {
		return "", &response.ErrExecuteHTML
	}

	return body.String(), nil
}

func (g *Gomail) SendNotification(req dto.NotificationRequest) error {
	emailBody, err := g.SetBodyHTML("notification.html", req)
	if err != nil {
		return &response.ErrSetHTML
	}

	g.message.SetHeader("From", os.Getenv("SENDER"))
	g.message.SetHeader("To", req.Email)
	g.message.SetHeader("Subject", "🎉 Fitur Baru Hadir di NusaGo! Eksklusif untuk Anda 🚀")
	g.message.SetBody("text/html", emailBody)

	if err := g.dialer.DialAndSend(g.message); err != nil {
		return &response.ErrFailedSendNotification
	}

	return nil
}

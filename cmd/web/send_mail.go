package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/TranQuocToan1996/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

// listenForMail creates go routine to make this function running in backgound asynchronous with other function in bookings app
func listenForMail() {
	go func() {
		// For loop for continuous listen for mails
		for {
			message := <-app.MailChan
			sendMessage(message)

		}
	}()
}

func sendMessage(mailData models.MailData) {
	// Create dummy sending mail server, and receiver is mailHog
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false // Active only when needed to send an email
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	// server.Encryption = mail.EncryptionSSLTLS

	// Client connect to server
	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	// Set information for email
	email := mail.NewMSG()
	email.SetFrom(mailData.From).AddTo(mailData.To).SetSubject(mailData.Subject)
	if mailData.Template == "" {
		email.SetBody(mail.TextHTML, mailData.Content)
	} else {
		// read template into memory
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", mailData.Template))
		if err != nil {
			errorLog.Println(err)
		}

		mailTemplate := string(data)
		// Replace in the content mailTemplate: "[%emailContent%]" -> m.Content with nolimit replacement
		msgToSend := strings.Replace(mailTemplate, "[%emailContent%]", mailData.Content, -1)
		email.SetBody(mail.TextHTML, msgToSend)
	}

	// Sending email
	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
	} else {
		infoLog.Println("email sent!")
	}

}

// func main() {

// 	// Sender data.
// 	from := "from@gmail.com"
// 	password := "<Email Password>"

// 	// Receiver email address.
// 	to := []string{
// 		"sender@example.com",
// 	}

// 	// smtp server configuration.
// 	smtpHost := "smtp.gmail.com"
// 	smtpPort := "587"

// 	// Message.
// 	message := []byte("This is a test email message.")

// 	// Authentication.
// 	auth := smtp.PlainAuth("", from, password, smtpHost)

// 	// Sending email.
// 	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println("Email Sent Successfully!")
// }

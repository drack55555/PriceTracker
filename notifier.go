package main

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

type NotifierCofig struct {
	SMTPServer     string
	SMTPPort       int
	SenderEmail    string
	SenderPassword string
}

type Notifier struct {
	config NotifierCofig
	dialer *gomail.Dialer
}

func NewNotifier(config NotifierCofig) *Notifier {
	d := gomail.NewDialer(
		config.SMTPServer,
		config.SMTPPort,
		config.SenderEmail,
		config.SenderPassword,
	)
	return &Notifier{
		config: config,
		dialer: d,
	}
}

func (n *Notifier) SendEmailAlert(toEmail string, productURL string, price float64) error {

	//Create a new blank Email Message using 'gomail.NewMessage()'
	m := gomail.NewMessage()

	// We set who it's from (us), who it's to (the user), and the subject.
	m.SetHeader("From", n.config.SenderEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Price Drop Alert! 📉")

	// Set the Email Body.
	body := fmt.Sprintf(`
		<h1>Price Drop Alert!</h1>
		<p>Good news! The price for the following product has dropped:</p>
		<p><a href="%s">Click here to view product</a></p>
		<h2>New Price: $%.2f</h2>
	`, productURL, price)

	m.SetBody("text/html", body)

	// 'DialAndSend' opens the connection to the SMTP server,
	// authenticates, sends the email, and closes the connection.
	if err := n.dialer.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

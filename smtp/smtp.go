package main

import (
	"flag"
	"fmt"
	"log"
	"net/smtp"
)

var (
	host    string
	from    string
	to      string
	subject string
	body    string
)

func main() {
	flag.StringVar(&host, "host", "", "host with port (mail.example.com:25)")
	flag.StringVar(&from, "from", "", "envelope FROM")
	flag.StringVar(&to, "to", "", "envelope RCPT TO")
	flag.StringVar(&subject, "subject", "Test subject", "actual email subject")
	flag.StringVar(&body, "body", "Test body", "actual email body")

	flag.Parse()

	if len(host) == 0 || len(from) == 0 || len(to) == 0 {
		flag.Usage()
		panic("Invalid arguments")
	}

	// Connect to the remote SMTP server.
	c, err := smtp.Dial(host)
	if err != nil {
		log.Fatal(err)
	}

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		log.Fatal(err)
	}
	if err := c.Rcpt(to); err != nil {
		log.Fatal(err)
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		log.Fatal(err)
	}
	_, err = fmt.Fprintf(wc, "From: %s\r\nTo: %s\r\nSubject: %subject\r\n\r\n%s", from, to, subject, body)
	if err != nil {
		log.Fatal(err)
	}
	err = wc.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
	"os"
)

var (
	host     string
	username string
	password string
	clean    bool
	ssl      bool
)

func main() {
	flag.StringVar(&host, "host", "", "host with port (mail.example.com:993)")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.BoolVar(&clean, "clean", false, "clean whole INBOX afterwards")
	flag.BoolVar(&ssl, "ssl", true, "start SSL connection (cert errors are ignored)")

	flag.Parse()

	if len(host) == 0 || len(username) == 0 || len(password) == 0 {
		flag.Usage()
		panic("Invalid arguments")
	}

	log.Println("Connecting to server...")

	// Connect to server
	var c *client.Client
	var err error

	if ssl == false {
		c, err = client.Dial(host)
	} else {
		c, err = client.DialTLS(host, &tls.Config{InsecureSkipVerify: true})
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()

	// Login
	if err := c.Login(username, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	log.Println("Mailboxes:")
	for m := range mailboxes {
		log.Println("* " + m.Name)
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 3 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	done = make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)

		if clean {
			item := imap.FormatFlagsOp(imap.AddFlags, true)
			flags := []interface{}{imap.DeletedFlag}

			if err := c.Store(seqset, item, flags, nil); err != nil {
				fmt.Println("IMAP Message Flag Update Failed")
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}()

	if clean {
		if err := c.Expunge(nil); err != nil {
			fmt.Println("IMAP Message Delete Failed")
			os.Exit(1)
		}
	}

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)

	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	log.Println("Done!")
}

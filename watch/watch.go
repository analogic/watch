package watch

import (
	"crypto/tls"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"
)

func ImapRetrieve(host string, username string, password string, clean bool, ssl bool, awaitTimeout int, awaitSubject string) {
	// Connect to server
	var c *client.Client
	var err error
	var done chan error

	if ssl == false {
		c, err = client.Dial(host)
	} else {
		c, err = client.DialTLS(host, &tls.Config{InsecureSkipVerify: true})
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")
	//c.SetDebug(os.Stdout)

	// Don't forget to logout
	//defer c.Logout()

	// Login
	if err := c.Login(username, password); err != nil {
		log.Println("Login error:")
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done = make(chan error, 1)
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

	if awaitTimeout > 0 {
		done := make(chan bool)
		go func() {
			for {
				subjects := imapList(c, clean)
				for _, subject := range subjects {
					if strings.Contains(subject, awaitSubject) {
						done <- true
					}
				}
				time.Sleep(time.Millisecond * 200)
			}
		}()

		select {
		case <-time.After(time.Duration(awaitTimeout) * time.Second):
			log.Println("Timeouted")
			os.Exit(1)
		case <-done:
			log.Println("DONE")
		}

	} else {
		imapList(c, clean)
	}
}

func imapList(c *client.Client, clean bool) []string {
	result := make([]string, 0)

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Flags for INBOX:", mbox.Flags)

	// Get the last 4 messages
	from := uint32(1)
	to := mbox.Messages

	if mbox.Messages == 0 {
		return result
	}

	if mbox.Messages > 3 {
		// We're using unsigned integers here, only substract if the result is > 0
		from = mbox.Messages - 3
	}
	seqset := new(imap.SeqSet)
	seqset.AddRange(to, from)

	messages := make(chan *imap.Message, 10)
	if err = c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages); err != nil {
		log.Println("IMAP fetch failed")
		os.Exit(1)
	}

	if clean {
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.DeletedFlag}

		if err := c.Store(seqset, item, flags, nil); err != nil {
			log.Println("IMAP Message Flag Update Failed")
			log.Println(err)
			//os.Exit(1)
		}
	}

	log.Println("Last 4 messages:")
	for msg := range messages {
		log.Println("* " + msg.Envelope.Subject)
		result = append(result, msg.Envelope.Subject)
	}

	if clean {
		if err := c.Expunge(nil); err != nil {
			log.Println("IMAP Message Delete Failed")
			//os.Exit(1)
		}
	}

	return result
}

func SMTPSend(host string, from string, to string, subject string, body string) {
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
	_, err = fmt.Fprintf(wc, "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)
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

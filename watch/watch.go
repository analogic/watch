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

func ImapRetrieve(host string, username string, password string, clean bool, ssl bool, awaitTimeout int, awaitSubject string) error {
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
		return err
	}
	log.Println("Connected")
	//c.SetDebug(os.Stdout)

	// Don't forget to logout
	//defer c.Logout()

	// Login
	if err := c.Login(username, password); err != nil {
		log.Println("Login error:")
		return err
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
		return err
	}

	if awaitTimeout > 0 {
		done := make(chan bool)
		errch := make(chan error)

		go func() {
			for {
				err, subjects := imapList(c, clean)
				if err != nil {
					errch <- err
				}
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
		case <-errch:
			return err
		}

	} else {
		err, _ := imapList(c, clean)
		if err != nil {
			return err
		}
	}

	return nil
}

func imapList(c *client.Client, clean bool) (error, []string) {
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
		return nil, result
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
		return err, result
	}

	if clean {
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.DeletedFlag}

		if err := c.Store(seqset, item, flags, nil); err != nil {
			log.Println("IMAP Message Flag Update Failed")
			return err, result
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
			return err, result
		}
	}

	return nil, result
}

func SMTPSend(host string, from string, to string, subject string, body string) error {
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(host)
	if err != nil {
		return err
	}

	// Set the sender and recipient first
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return err
	}

	// Send the email body.
	wc, err := c.Data()
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(wc, "From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, body)
	if err != nil {
		return err
	}
	err = wc.Close()
	if err != nil {
		return err
	}

	// Send the QUIT command and close the connection.
	err = c.Quit()
	if err != nil {
		return err
	}

	return nil
}

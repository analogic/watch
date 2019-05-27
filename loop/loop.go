package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"
	"watch/watch"
)

var (
	smtp  string
	imap1 string
	imap2 string
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	flag.StringVar(&smtp, "smtp", "", "SMTP send definition host:port:from_address:to_address")
	flag.StringVar(&imap1, "imap1", "", "IMAP retrieve definition host:port:username:password")
	flag.StringVar(&imap2, "imap2", "", "IMAP retrieve definition host:port:username:password")

	flag.Parse()

	if len(smtp) == 0 || len(imap1) == 0 || len(imap2) == 0 {
		flag.Usage()
		panic("Invalid arguments")
	}

	subject := fmt.Sprintf("Mailserver loop test %s", randSeq(10))
	s := strings.Split(smtp, ":")

	log.Println("Sending email through SMTP")
	watch.SMTPSend(s[0]+":"+s[1], s[2], s[3], subject, "Loop test body")

	log.Println("IMAP 1")
	i1 := strings.Split(imap1, ":")
	watch.ImapRetrieve(i1[0]+":"+i1[1], i1[2], i1[3], true, true, 10, subject)

	log.Println("IMAP 2")
	i2 := strings.Split(imap2, ":")
	watch.ImapRetrieve(i2[0]+":"+i2[1], i2[2], i2[3], true, true, 10, subject)
}

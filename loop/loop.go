package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"strings"
	"time"
	"watch/watch"
)

var (
	smtp       string
	imap1      string
	imap2      string
	resultOnly bool
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

	start := time.Now()

	flag.StringVar(&smtp, "smtp", "", "SMTP send definition host:port:from_address:to_address")
	flag.StringVar(&imap1, "imap1", "", "IMAP retrieve definition host:port:username:password")
	flag.StringVar(&imap2, "imap2", "", "IMAP retrieve definition host:port:username:password")
	flag.BoolVar(&resultOnly, "result-only", false, "Show loop time only")

	flag.Parse()

	if len(smtp) == 0 || len(imap1) == 0 || len(imap2) == 0 {
		flag.Usage()
		panic("Invalid arguments")
	}

	if resultOnly {
		log.SetOutput(ioutil.Discard)
	}

	subject := fmt.Sprintf("Mailserver loop test %s", randSeq(10))
	s := strings.Split(smtp, ":")

	log.Println("-------------------------------------------------------------------------------------------------")

	log.Println("Sending email through SMTP")
	watch.SMTPSend(s[0]+":"+s[1], s[2], s[3], subject, "Loop test body")

	log.Println("-------------------------------------------------------------------------------------------------")

	log.Println("IMAP 1")
	i1 := strings.Split(imap1, ":")
	watch.ImapRetrieve(i1[0]+":"+i1[1], i1[2], i1[3], true, true, 10, subject)

	log.Println("-------------------------------------------------------------------------------------------------")

	log.Println("IMAP 2")
	i2 := strings.Split(imap2, ":")
	watch.ImapRetrieve(i2[0]+":"+i2[1], i2[2], i2[3], true, true, 10, subject)

	elapsed := time.Since(start)
	log.Printf("Send&Receive took %fs", elapsed.Seconds())

	if resultOnly {
		fmt.Println(elapsed.Seconds())
	}
}

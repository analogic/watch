package main

import (
	"flag"
	"log"
	"os"
	"watch/watch"
)

var (
	host     string
	username string
	password string
	clean    bool
	ssl      bool

	awaitTimeout int
	awaitSubject string
)

func main() {
	flag.StringVar(&host, "host", "", "host with port (mail.example.com:993)")
	flag.StringVar(&username, "username", "", "username")
	flag.StringVar(&password, "password", "", "password")
	flag.BoolVar(&clean, "clean", false, "clean whole INBOX afterwards")
	flag.BoolVar(&ssl, "ssl", true, "start SSL connection (cert errors are ignored)")

	flag.IntVar(&awaitTimeout, "await-timeout", 0, "number of seconds to wait for email with searched subject")
	flag.StringVar(&awaitSubject, "await-subject", "", "part of searched subject")

	flag.Parse()

	if len(host) == 0 || len(username) == 0 || len(password) == 0 {
		flag.Usage()
		panic("Invalid arguments")
	}

	log.Println("Connecting to server...")

	watch.ImapRetrieve(host, username, password, clean, ssl, awaitTimeout, awaitSubject)

	log.Println("Done!")
	os.Exit(0)
}

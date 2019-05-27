package main

import (
	"flag"
	"watch/watch"
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

	watch.SMTPSend(host, from, to, subject, body)
}

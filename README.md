# Custom smoke test for mailserver availability

Simple script in GO which should test: 
1. receiving from remote host (SMTP relay)
2. delivery to local box (LDA)
3. fetch from local box (IMAP)
4. copy filter to remote box (Sieve)
5. delivery to remote box (MTA)

#### Requirements: 
 - two mailservers
 - IMAPs, SMTP
 
 
#### How to:

- create mailbox watchdog@tested.com at mail.tested.com with some password
- add copy filter to this box with target watchdog-tested@our.com 
- create mailbox watchdog-tested@our.com with some password
- create test for your NMS with command like this:

```
$ ./loop \
    -smtp mail.tested.com:25:watchdog-tested@our.com:watchdog@tested.com \
    -imap1 mail.tested.com:993:watchdog@tested.com:password \
    -imap2 mail.our.com:993:watchdog-tested@our.com:password \
    -result-only
1.743403509
```

#### Result:
When success then process time is returned. If execution fail then process should return 1 as error code and 0 as processing time
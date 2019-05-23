#!/bin/sh

set -ex

from=$1
via=$2
viaPassword=$3
viaHost=$4
to=$5
toPassword=$6
toHost=$7

echo "Sending email through SMTP"
smtp/smtp -from $from -to $via -host $viaHost:25 -subject "Test123"

echo "Waiting for incoming mail at mailserver"
imap/imap -host $viaHost:993 -username $via -password $viaPassword -await-timeout 5 -await-subject "Test123" -clean

echo "Waiting for redirected email at poste"
imap/imap -host $toHost:993 -username $to -password $toPassword -await-timeout 5 -await-subject "Test123" -clean

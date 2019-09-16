#!/bin/bash
openssl genrsa -out myCA.key 2048
openssl req -x509 -new -key myCA.key  -out myCA.cer -days 3650 -subj="/CN=\"$1\"/O=\"DNSFidder untrusted MITM proxy Inc\""

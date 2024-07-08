#!/usr/bin/sh

HOST="your-fileserver.tld"

curl --connect-timeout 5 \
      --max-time 10 \
      --retry 5 \
      --retry-delay 5 \
      --retry-max-time 40 \
      -X POST "http://${HOST}:9999/webhook"

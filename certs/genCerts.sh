#! /bin/bash

# Use the below lines to generate the certs
openssl req -x509 -new -nodes -key server.key -sha256 -days 1024 -out server.crt
openssl req -x509 -new -nodes -key client.key -sha256 -days 1024 -out client.crt

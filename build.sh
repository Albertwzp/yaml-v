#! /bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o yaml-v
tar -zcvf  yaml-v.tgz yaml-v && rm yaml-v

#docker build -t yaml-v:v0.1 .

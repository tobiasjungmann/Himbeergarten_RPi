#!/bin/sh

cd testClient || exit
go clean -modcache
go mod tidy
go get github.com/tobiasjungmann/Himbeergarten_RPi/server@latest

cd ../humidity_forwarder|| exit
go clean -modcache
go mod tidy
go get github.com/tobiasjungmann/Himbeergarten_RPi/server@latest

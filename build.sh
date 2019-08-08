#!/bin/sh
go build -o build/linux_x64/doorparty-connector
GOARM=6 GOARCH=arm go build -o build/linux_arm6/doorparty-connector
GOOS=windows GOARCH=386 go build -o build/win32/doorparty-connector.exe

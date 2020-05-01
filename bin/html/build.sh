#!/bin/bash

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o html_wallet_windows.exe
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o html_wallet_linux
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o html_wallet_mac
rm html_wallet -rf
mkdir html_wallet
copy conf.json html_wallet
mv html_wallet_windows.exe html_wallet
mv html_wallet_linux html_wallet
mv html_wallet_mac html_wallet
cp static html_wallet -rf
tar zcvf html_wallet_$(date +'%y%m%d_%H%M%S').tar.gz html_wallet
rm html_wallet -rf
echo Enter to exit
read k

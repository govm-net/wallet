#!/bin/bash

CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o html_wallet_windows.exe
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o html_wallet_linux
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o html_wallet_mac
folder=govm_html_wallet
# echo $folder "$folder"
rm $folder -rf
mkdir $folder
mv html_wallet_windows.exe $folder
mv html_wallet_linux $folder
mv html_wallet_mac $folder
cp static $folder -rf
#tar zcvf "$folder"_$(date +'%Y%m%d_%H%M%S').tar.gz $folder
zip -r "$folder"_$(date +'%Y%m%d_%H%M%S').tar.gz $folder
rm $folder -rf
echo Enter to exit
read k

#!/bin/bash

fyne.exe package -os windows -name govm.net -icon ./assets/govm.png
mv gui.exe govm.exe
folder=govm_windows_wallet
rm $folder -rf
mkdir $folder
mv govm.exe $folder
cp assets $folder -rf
cp conf.json $folder -rf
zip -r govm_windows_wallet_$(date +'%Y%m%d_%H%M%S').tar.gz $folder
echo Enter to exit
read k
rm $folder -rf

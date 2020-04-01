package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/lengzhao/wallet/trans"
	"golang.org/x/crypto/ssh/terminal"
)

// TConfig config of app
type TConfig struct {
	trans.Wallet
	Password   string `json:"password,omitempty"`
	WalletFile string `json:"wallet_file,omitempty"`
	APIServer  string `json:"api_server,omitempty"`
}

// DebugMod debug mode
const (
	CreateFristTrans = false
)

var (
	conf TConfig
)

func init() {
	err := loadConfig()
	if err != nil {
		log.Println("fail to load config,", err)
		time.Sleep(time.Second * 3)
		os.Exit(2)
	}
}

func loadConfig() error {
	writeConf := false
	data, err := ioutil.ReadFile("conf.json")
	if err == nil {
		err = json.Unmarshal(data, &conf)
		if err != nil {
			log.Println("fail to Unmarshal configure,conf.json")
			return err
		}
	} else {
		// log.Println("fail to read file,conf.json")
		writeConf = true
	}

	if conf.WalletFile == "" {
		conf.WalletFile = "wallet.key"
	}

	if conf.Password == "" {
		fmt.Println("please enter the password of wallet:(default is govm_pwd@2019)")
		password, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return err
		}
		conf.Password = string(password)
		if conf.Password == "" {
			return fmt.Errorf("need password")
		}
	}
	if conf.APIServer == "" {
		conf.APIServer = "http://govm.net"
	}
	if _, err := os.Stat(conf.WalletFile); os.IsNotExist(err) {
		conf.Wallet.New(conf.Password)
		out := conf.Wallet.String()
		ioutil.WriteFile(conf.WalletFile, []byte(out), 666)
	} else {
		data, err = ioutil.ReadFile(conf.WalletFile)
		if err != nil {
			log.Println("fail to read wallet file:", conf.WalletFile, err)
			return err
		}
		err = conf.Wallet.Load(conf.Password, string(data))
		if err != nil {
			log.Println("fail to load wallet.", err)
			return err
		}
	}
	if writeConf {
		conf.Password = ""
		data, _ = json.MarshalIndent(conf, "", "  ")
		ioutil.WriteFile(conf.WalletFile, data, 666)
	}
	return nil
}

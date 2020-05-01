package conf

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Xuanwo/go-locale"
	"github.com/lengzhao/wallet/trans"
)

// Vserion version of wallet
const Vserion = "v0.1.1"

var (
	conf       map[string]string
	myWallet   trans.Wallet
	myPassword string
)

// Public key of config
const (
	APIServer   = "APIServer"
	WalletFile  = "WalletFile"
	Langure     = "Langure"
	CoinUnit    = "Unit"
	FyneScale   = "FYNE_SCALE"
	LangureList = "langure_list"

	confFile = "conf.json"
)

// Languages languages
var Languages []string

func init() {
	writeConf := false
	var c map[string]string

	data, err := ioutil.ReadFile(confFile)
	if err == nil {
		err = json.Unmarshal(data, &c)
		if err != nil {
			log.Println("fail to Unmarshal configure,conf.json")
			writeConf = true
		}
	} else {
		writeConf = true
		c = make(map[string]string)
	}

	if c[WalletFile] == "" {
		c[WalletFile] = "wallet.key"
	}

	if c[APIServer] == "" {
		c[APIServer] = "http://govm.net"
	}
	if c[LangureList] == "" {
		c[LangureList] = "en,zh"
	}
	if writeConf {
		data, _ = json.MarshalIndent(c, "", "  ")
		ioutil.WriteFile(confFile, data, 666)
	}
	Languages = strings.Split(c[LangureList], ",")
	conf = c
}

// Load load wallet
func Load(password string) error {
	if _, err := os.Stat(conf["WalletFile"]); os.IsNotExist(err) {
		myWallet.New(password)
		out := myWallet.String()
		ioutil.WriteFile(conf["WalletFile"], []byte(out), 666)
		log.Println("new wallet:", myWallet.AddressStr)
	} else {
		data, err := ioutil.ReadFile(conf["WalletFile"])
		if err != nil {
			log.Println("fail to read wallet file:", conf["WalletFile"], err)
			return err
		}

		err = myWallet.Load(password, string(data))
		if err != nil {
			log.Println("fail to load wallet.", err)
			return err
		}
		log.Println("load wallet:", myWallet.AddressStr)
	}

	myPassword = password
	return nil
}

// CheckPassword check password, return true when password is right
func CheckPassword(in string) bool {
	if in == myPassword {
		return true
	}
	return false
}

// Get get config
func Get(key string) string {
	out, ok := conf[key]
	if !ok {
		switch key {
		case Langure:
			tag, err := locale.Detect()
			if err != nil {
				log.Println("fail to get lang")
				return "en"
			}
			lang := tag.String()
			out = Languages[0]
			for _, it := range Languages {
				if strings.HasPrefix(lang, it) {
					out = it
				}
			}
			log.Println("lang", lang)
		case APIServer:
			out = "http://govm.net"
		case CoinUnit:
			out = "tc"
		}
	}
	return out
}

// Set set
func Set(key, value string) {
	conf[key] = value
	data, _ := json.MarshalIndent(conf, "", "  ")
	ioutil.WriteFile(confFile, data, 666)
}

// GetWallet get wallet
func GetWallet() trans.Wallet {
	return myWallet
}

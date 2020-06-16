package conf

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/Xuanwo/go-locale"
	"github.com/lengzhao/wallet/trans"
)

// Vserion version of wallet
const Vserion = "v0.5.1"

// Config config
type Config struct {
	APIServer    string   `json:"api_server,omitempty"`
	WalletFile   string   `json:"wallet_file,omitempty"`
	Langure      string   `json:"langure,omitempty"`
	LangureList  []string `json:"langure_list,omitempty"`
	CoinUnit     string   `json:"coin_unit,omitempty"`
	DefaultChain string   `json:"default_chain,omitempty"`
	Chains       []string `json:"chains,omitempty"`
}

// Public key of config
const confFile = "conf.json"

var (
	conf       Config
	myWallet   trans.Wallet
	myPassword string
)

func init() {
	writeConf := false
	conf.APIServer = "http://govm.net:9090"
	conf.Chains = []string{"1", "2"}
	conf.WalletFile = "wallet.key"
	conf.CoinUnit = "govm"
	conf.LangureList = []string{"en", "zh"}
	conf.DefaultChain = "1"

	data, err := ioutil.ReadFile(confFile)
	if err == nil {
		err = json.Unmarshal(data, &conf)
		if err != nil {
			log.Println("fail to Unmarshal configure,conf.json")
			writeConf = true
		}
	} else {
		writeConf = true
	}
	if conf.Langure == "" {
		conf.Langure = getLangure()
	}

	if writeConf {
		data, _ = json.MarshalIndent(conf, "", "  ")
		ioutil.WriteFile(confFile, data, 666)
	}
}

func getLangure() string {
	tag, err := locale.Detect()
	if err != nil {
		log.Println("fail to get lang")
		return conf.LangureList[0]
	}
	lang := tag.String()
	out := conf.LangureList[0]
	for _, it := range conf.LangureList {
		if strings.HasPrefix(lang, it) {
			out = it
		}
	}
	return out
}

// Load load wallet
func Load(password string) error {
	if _, err := os.Stat(conf.WalletFile); os.IsNotExist(err) {
		myWallet.New(password)
		out := myWallet.String()
		ioutil.WriteFile(conf.WalletFile, []byte(out), 666)
		log.Println("new wallet:", myWallet.AddressStr)
	} else {
		data, err := ioutil.ReadFile(conf.WalletFile)
		if err != nil {
			log.Println("fail to read wallet file:", conf.WalletFile, err)
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

// SaveConf save conf to file
func SaveConf() error {
	data, _ := json.MarshalIndent(conf, "", "  ")
	ioutil.WriteFile(confFile, data, 666)
	return nil
}

// CheckPassword check password, return true when password is right
func CheckPassword(in string) bool {
	if in == myPassword {
		return true
	}
	return false
}

// ChangePassword change password
func ChangePassword(pwd string) error {
	if !myWallet.ChangePwd(myPassword, pwd) {
		return errors.New("error old password")
	}
	out := myWallet.String()
	os.Rename(conf.WalletFile, "old_"+conf.WalletFile)
	err := ioutil.WriteFile(conf.WalletFile, []byte(out), 666)
	if err != nil {
		log.Println("fail to save wallet")
		return err
	}
	myPassword = pwd
	return nil
}

// GetWallet get wallet
func GetWallet() trans.Wallet {
	return myWallet
}

// Get get conf
func Get() *Config {
	return &conf
}

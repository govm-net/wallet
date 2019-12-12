package main

import (
	"encoding/hex"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	"github.com/lengzhao/wallet/trans"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"
)

// Conf configure
type Conf struct {
	WalletFile string `json:"wallet_file"`
	Operation  string `json:"operation"`
	Descript   string `json:"descript"`
	Server     string `json:"server"`
	Chain      uint64 `json:"chain"`
	Energy     uint64 `json:"energy"`
	Cost       uint64 `json:"cost"`
	TransOps   byte   `json:"trans_ops"`
	TransFile  string `json:"trans_file"`
	Transfer   struct {
		Peer string `json:"peer"`
		Msg  string `json:"msg"`
	} `json:"transfer"`
	Move struct {
		DstChain uint64 `json:"dst_chain"`
	} `json:"move"`
	RunApp struct {
		AppName string `json:"app_name"`
		Data    string `json:"data"`
		IsHex   bool   `json:"is_hex"`
	} `json:"run_app"`
	UpdateLife struct {
		AppName string `json:"app_name"`
		Life    uint64 `json:"life"`
	} `json:"update_life"`
	pwd string
}

const (
	opsNewWallet   = "wallet"
	opsShow        = "show"
	opsTransaction = "transaction"
	opsSendTrans   = "sendTrans"
)

var conf Conf

func main() {
	var wait bool
	var pwd, confFile string
	flag.BoolVar(&wait, "w", true, "wait 1 second when exit")
	flag.StringVar(&pwd, "pwd", "", "the password of the wallet")
	flag.StringVar(&confFile, "f", "conf.json", "configure file")

	flag.Parse()

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		conf.WalletFile = "wallet.key"
		conf.TransFile = "transaction.data"
		conf.Chain = 1
		conf.Server = "http://govm.net"
		conf.Energy = 10000000
		conf.Operation = opsNewWallet
		conf.Descript = "operation:wallet/show/transaction/sendTrans"
		data, _ := json.MarshalIndent(conf, "  ", "  ")
		ioutil.WriteFile(confFile, data, 0666)
	} else {
		data, err := ioutil.ReadFile(confFile)
		if err != nil {
			fmt.Println("fail to read config:", confFile, err)
			return
		}
		err = json.Unmarshal(data, &conf)
		if err != nil {
			fmt.Println("fail to unmarshal config:", err)
			return
		}
	}
	conf.pwd = pwd

	if wait {
		defer func() {
			fmt.Println("exit 1 second later")
			time.Sleep(time.Second)
		}()
	}

	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "http://govm.net")
	case "darwin":
		cmd = exec.Command("open", "http://govm.net")
	case "linux":
		cmd = exec.Command("xdg-open", "http://govm.net")
	}
	err := cmd.Start()
	if err != nil {
		fmt.Println("fail to open baidu.com. ", err)
		return
	}

	switch conf.Operation {
	case opsNewWallet:
		createWallet()
	case opsShow:
		showWallet()
	case opsTransaction:
		newTransaction()
	case opsSendTrans:
		sendTransaction()
	default:
		fmt.Println("error:not support.", conf.Operation)
	}
}

func sendTransaction() {
	urlStr := fmt.Sprintf("%s/api/v1/%d/transaction/new", conf.Server, conf.Chain)
	data, err := ioutil.ReadFile(conf.TransFile)
	if err != nil {
		fmt.Println("fail to read transaction data from", conf.TransFile, err)
		return
	}
	reader := bytes.NewReader(data)
	req, err := http.NewRequest(http.MethodPost, urlStr, reader)
	if err != nil {
		fmt.Println("fail to create request.", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("fail to post transactin.", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode > 299 {
		data, _ = ioutil.ReadAll(resp.Body)
		fmt.Println("fail to send transaction,response:", resp.StatusCode)
		fmt.Println("error msg:", string(data))
		return
	}
	fmt.Println("success to send transaction")
}

func newTransaction() {
	if conf.pwd == "" {
		fmt.Printf("Enter password: ")
		data, err := gopass.GetPasswdMasked()
		if err != nil {
			fmt.Println("fail to get password:", err)
			return
		}
		conf.pwd = string(data)
		if conf.pwd == "" {
			fmt.Println("fail to get password.")
			return
		}
	}
	data, err := ioutil.ReadFile(conf.WalletFile)
	if err != nil {
		fmt.Println("fail to read wallet file:", conf.WalletFile, err)
		return
	}
	var w trans.Wallet
	err = w.Load(conf.pwd, string(data))
	if err != nil {
		fmt.Println("fail to load wallet:", conf.WalletFile, err)
		return
	}
	t := trans.NewTransaction(conf.Chain, w.Address)
	switch conf.TransOps {
	case trans.OpsTransfer:
		in := conf.Transfer
		err := t.CreateTransfer(in.Peer, in.Msg, conf.Cost, conf.Energy)
		if err != nil {
			fmt.Println("error:", err)
			return
		}
	case trans.OpsMove:
		t.CreateMove(conf.Move.DstChain, conf.Cost, conf.Energy)
	case trans.OpsRunApp:
		in := conf.RunApp
		var d = []byte(in.Data)
		if in.IsHex{
			d,err = hex.DecodeString(in.Data)
			if err != nil{
				fmt.Println("error hex data:",err)
				return
			}
		}
		err = t.RunApp(in.AppName, conf.Cost, conf.Energy, d)
		if err != nil {
			fmt.Println("fail to RunApp:", err)
			return
		}
	case trans.OpsUpdateAppLife:
		in := conf.UpdateLife
		err = t.UpdateAppLife(in.AppName, in.Life, conf.Energy)
		if err != nil {
			fmt.Println("fail to RunApp:", err)
			return
		}
	default:
		fmt.Println("error:not support,", conf.TransOps)
	}
	signData := t.GetSignData()
	t.SetSign(w.Sign(signData))
	tBytes := t.Output()
	ioutil.WriteFile(conf.TransFile, tBytes, 0666)
}

func createWallet() {
	if _, err := os.Stat(conf.WalletFile); !os.IsNotExist(err) {
		fmt.Println("error:already exist the wallet file:", conf.WalletFile)
		return
	}
	if conf.pwd == "" {
		fmt.Printf("Enter password: ")
		data, err := gopass.GetPasswdMasked()
		if err != nil {
			fmt.Println("fail to get password:", err)
			return
		}
		conf.pwd = string(data)
		if conf.pwd == "" {
			fmt.Println("fail to get password.")
			return
		}
	}
	var w trans.Wallet
	w.New(conf.pwd)
	data := w.String()
	f, err := os.Create(conf.WalletFile)
	if err != nil {
		fmt.Println("fail to create file:", conf.WalletFile, err)
		return
	}
	defer f.Close()
	f.Write([]byte(data))
	fmt.Println("new wallet address:", w.AddressStr)
	return
}

func showWallet() {
	var info struct {
		Tag        string `json:"tag"`
		AddressStr string `json:"address_str"`
		Address    []byte `json:"address"`
	}
	data, err := ioutil.ReadFile(conf.WalletFile)
	if err != nil {
		fmt.Println("fail to read wallet file:", conf.WalletFile, err)
		return
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		fmt.Println("fail to Unmarshal wallet file:", conf.WalletFile, err)
		return
	}
	if len(info.Address) == 0 {
		fmt.Println("error address of wallet")
		return
	}
	fmt.Printf("wallet address:%x\n", info.Address)
	fmt.Println("try to get account info from server.")
	urlStr := fmt.Sprintf("%s/api/v1/%d/account?address=%x", conf.Server, conf.Chain, info.Address)
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("fail to get info from server:", urlStr, err)
		return
	}
	defer func() { resp.Body.Close() }()
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("fail to read body:", err)
		return
	}
	fmt.Println("info:", string(data))
}

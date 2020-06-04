package main

import (
	"encoding/hex"
	"fmt"
	"syscall/js"

	"github.com/lengzhao/wallet/trans"
)

var w trans.Wallet

func main() {
	js.Global().Set("loadWallet",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) != 2 {
				return "error parament,hope password+walletString"
			}
			err := w.Load(args[0].String(), args[1].String())
			if err != nil {
				fmt.Println("fail to load wallet")
				return fmt.Sprintf("error:%s", err)
			}
			return "ok"
		}),
	)
	js.Global().Set("newWallet",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) != 1 {
				return "error parament,need password"
			}
			err := w.New(args[0].String())
			if err != nil {
				fmt.Println("fail to new wallet")
				return fmt.Sprintf("error:%s", err)
			}
			return "ok"
		}),
	)
	js.Global().Set("changePwd",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) != 2 {
				return "error parament,need old and new password"
			}
			if w.ChangePwd(args[0].String(), args[1].String()) {
				return "ok"
			}
			return "error:password"
		}),
	)
	js.Global().Set("getAddress",
		js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return w.AddressStr
		}),
	)
	js.Global().Set("newTransaction", js.FuncOf(newTransaction))

	defer recover()
	println := js.Global().Get("wasmLoaded")
	if println == js.Undefined() {
		fmt.Println("undefine wasmLoaded")
	} else {
		println.Invoke()
	}

	done := make(chan int, 0)
	<-done
}

func newTransaction(this js.Value, args []js.Value) interface{} {
	if len(args) < 3 {
		fmt.Println("need parament")
		return "error:need paramement"
	}
	if w.AddressStr == "" {
		return "error:privateKey"
	}
	chain := uint64(args[1].Float())
	cost := uint64(args[2].Float())
	t := trans.NewTransaction(chain, w.Address, cost)
	ops := args[0].Int()
	switch uint8(ops) {
	// OpsTransfer pTransfer
	case trans.OpsTransfer:
		if len(args) < 4 {
			fmt.Println("OpsTransfer,peerAddress and cost")
			return "error:need paramement"
		}
		peerStr := args[3].String()
		var energy uint64
		if len(args) > 4 {
			energy = uint64(args[4].Float())
		}
		var msg string
		if len(args) > 5 {
			msg = args[5].String()
		}
		err := t.CreateTransfer(peerStr, msg, cost)
		if err != nil {
			fmt.Println("error:", err)
			return fmt.Sprintf("error:%s", err)
		}
		t.SetEnergy(energy)
		fmt.Printf("transfer. to:%s,cost:%d,energy:%d", peerStr, cost, energy)
		// OpsMove Move out of coin, move from this chain to adjacent chains
	case trans.OpsMove:
		if len(args) < 4 {
			fmt.Println("OpsMove,dstChain and cost")
			return "error:need paramement"
		}
		dstChain := uint64(args[3].Float())
		var energy uint64
		if len(args) > 4 {
			energy = uint64(args[4].Float())
		}
		t.CreateMove(dstChain)
		t.SetEnergy(energy)
		fmt.Printf("move. to:%d,cost:%d,energy:%d", dstChain, cost, energy)
		// OpsRunApp run app
	case trans.OpsRunApp:
		if len(args) < 4 {
			fmt.Println("OpsRunApp,dstChain and cost")
			return "error:need paramement"
		}
		appStr := args[3].String()
		var data string
		if len(args) > 4 {
			data = args[4].String()
		}
		var energy uint64
		if len(args) > 5 {
			energy = uint64(args[5].Float())
		}
		err := t.RunApp(appStr, []byte(data))
		t.SetEnergy(energy)
		if err != nil {
			fmt.Println("error:", err)
			return fmt.Sprintf("error:%s", err)
		}
		fmt.Printf("OpsRunApp. app:%s,cost:%d,energy:%d,data:%s", appStr, cost, energy, data)
		// OpsUpdateAppLife update app life
	case trans.OpsUpdateAppLife:
		if len(args) < 4 {
			fmt.Println("OpsUpdateAppLife,peerAddress and cost")
			return "error:need paramement"
		}
		appStr := args[4].String()
		life := uint64(args[5].Float())
		var energy uint64
		if len(args) > 5 {
			energy = uint64(args[6].Float())
		}
		err := t.UpdateAppLife(appStr, life)
		if err != nil {
			fmt.Println("error:", err)
			return fmt.Sprintf("error:%s", err)
		}
		fmt.Printf("OpsUpdateAppLife. app:%s,life:%d,energy:%d", appStr, life, energy)
	default:
		fmt.Println("not support option:", ops)
		return fmt.Sprintf("error:not support %d", ops)
	}
	signData := t.GetSignData()
	t.SetTheSign(w.Sign(signData))
	tBytes := t.Output()
	return hex.EncodeToString(tBytes)
}

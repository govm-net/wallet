package trans

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/lengzhao/govm/encrypt"
	"github.com/lengzhao/govm/wallet"
	"time"
)

// Wallet wallet
type Wallet struct {
	Tag        string `json:"tag,omitempty"`
	AddressStr string `json:"address_str,omitempty"`
	Address    []byte `json:"address,omitempty"`
	Key        []byte `json:"key,omitempty"`
	SignPrefix []byte `json:"sign_prefix,omitempty"`
	pwd        string
}

const (
	// TimeDuration one year/12
	TimeDuration = 31558150000 / 12
)

// New new wallet
func (w *Wallet) New(pwd string) error {
	priv := wallet.NewPrivateKey()
	priv = append(priv, []byte(pwd)...)
	priv = wallet.GetHash(priv)

	pubK := wallet.GetPublicKey(priv)
	w.Address = wallet.PublicKeyToAddress(pubK, EAddrTypeDefault)
	w.Key = priv
	w.pwd = pwd
	w.AddressStr = hex.EncodeToString(w.Address)
	return nil
}

// ChangePwd change password
func (w *Wallet) ChangePwd(old, pwd string) bool {
	if old != w.pwd {
		return false
	}
	w.pwd = pwd
	return true
}

// Load load wallet by json data
func (w *Wallet) Load(pwd, data string) error {
	err := json.Unmarshal([]byte(data), w)
	if err != nil {
		fmt.Println("fail to json.Unmarshal:", err)
		return err
	}

	aesEnc := encrypt.AesEncrypt{}
	aesEnc.Key = pwd
	privKey, err := aesEnc.Decrypt(w.Key)
	if err != nil {
		fmt.Println("error password:", err)
		w.Address = nil
		w.AddressStr = ""
		return err
	}
	w.Key = privKey
	w.pwd = pwd
	return nil
}

const (
	// EAddrTypeDefault the type of default public address
	EAddrTypeDefault = byte(iota + 1)
	// EAddrTypeIBS identity-based signature 基于身份的签名，不同时间，使用不同私钥签名(签名时间是消息的前8个字节)
	EAddrTypeIBS
)

// String return json string
func (w *Wallet) String() string {
	aesEnc := encrypt.AesEncrypt{}
	aesEnc.Key = w.pwd

	arrEncrypt, err := aesEnc.Encrypt(w.Key)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	info := Wallet{}
	info.AddressStr = w.AddressStr
	info.Address = w.Address
	info.SignPrefix = w.SignPrefix
	info.Key = arrEncrypt
	if info.Address[0] == EAddrTypeIBS {
		t := wallet.GetDeadlineOfIBS(info.Address)
		info.Tag = time.Unix(int64(t/1000), 0).Format(time.RFC3339)
	} else {
		info.Tag = time.Now().Add(time.Hour * 24 * 365 * 85).Format(time.RFC3339)
	}

	sd, _ := json.Marshal(info)
	return string(sd)
}

// Sign sign
func (w *Wallet) Sign(data []byte) []byte {
	var out []byte
	if len(w.SignPrefix) > 0 {
		out = append(out, w.SignPrefix...)
	}
	sign := wallet.Sign(w.Key, data)
	out = append(out, sign...)
	return out
}

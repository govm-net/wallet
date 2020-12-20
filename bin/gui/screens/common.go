package screens

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/trans"
)

//statCoinLock
//statCoinUnlock
//statMiningCount
//statMinerHit
//statMinerReg
//statTransferIn
//statTransferOut
//statMove
//statAPPRun

const coreName = "ff0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f"

func newFormItemWithUnit(in fyne.CanvasObject, right string) fyne.CanvasObject {
	unit := widget.NewLabel(right)
	borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
	return fyne.NewContainerWithLayout(borderLayout, unit, in)
}
func newFormItemWithObj(in, right fyne.CanvasObject) fyne.CanvasObject {
	borderLayout := layout.NewBorderLayout(nil, nil, nil, right)
	return fyne.NewContainerWithLayout(borderLayout, right, in)
}

func runApp(chain, cost, energy uint64, app, prefix string, param []byte) (string, error) {
	myWlt := conf.GetWallet()
	trans := trans.NewTransaction(chain, myWlt.Address, cost)
	var body []byte
	var err error
	if prefix != "" {
		body, err = hex.DecodeString(prefix)
		if err != nil {
			return "", err
		}
	}
	if len(param) > 0 {
		body = append(body, param...)
	}
	trans.Energy = energy
	err = trans.RunApp(app, body)
	if err != nil {
		return "", err
	}
	td := trans.GetSignData()
	sign := myWlt.Sign(td)
	trans.SetTheSign(sign)
	td = trans.Output()
	key := trans.Key[:]

	err = postTrans(chain, td)
	if err != nil {
		// result.SetText(fmt.Sprintf("%s", err))
		log.Println("fail to run app:", err)
		return "", err
	}
	return hex.EncodeToString(key), nil
}

type dataInfo struct {
	AppName    string `json:"app_name,omitempty"`
	StructName string `json:"struct_name,omitempty"`
	IsDBData   bool   `json:"is_db_data,omitempty"`
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
	Life       uint64 `json:"life,omitempty"`
}

func getNextKeyOfDB(chain, appName, structName, preKey string) string {
	if appName == "" {
		appName = coreName
	}
	urlStr := conf.Get().APIServer
	urlStr += "/api/v1/" + chain + "/data/visit?app_name=" + appName
	urlStr += "&is_db_data=true&struct_name=" + structName
	urlStr += "&pre_key=" + preKey
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println("fail to get db.", urlStr, err)
		return ""
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || len(data) == 0 {
		log.Println("fail to get db.", urlStr, resp.Status, string(data))
		return ""
	}
	info := dataInfo{}
	json.Unmarshal(data, &info)
	return info.Key
}
func getIntOfDB(chain, appName, structName, address string) uint64 {
	value, _ := getStringOfDB(chain, appName, structName, address)
	if value == "" {
		return 0
	}
	out, err := strconv.ParseUint(value, 16, 64)
	if err != nil {
		log.Println("parse error.", value, err)
	}

	return out
}

func getStringOfDB(chain, appName, structName, address string) (string, uint64) {
	if appName == "" {
		appName = coreName
	}
	urlStr := conf.Get().APIServer
	urlStr += "/api/v1/" + chain + "/data?app_name=" + appName
	urlStr += "&is_db_data=true&struct_name=" + structName
	urlStr += "&key=" + address
	resp, err := http.Get(urlStr)
	if err != nil {
		log.Println("fail to get db.", urlStr, err)
		return "", 0
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || len(data) == 0 {
		log.Println("fail to get db.", urlStr, resp.Status, string(data))
		return "", 0
	}
	info := dataInfo{}
	err = json.Unmarshal(data, &info)
	if err != nil || info.Value == "" {
		log.Println("not value.", urlStr)
		return "", 0
	}
	// log.Println("success to get:", urlStr, info.Value)
	return info.Value, info.Life
}

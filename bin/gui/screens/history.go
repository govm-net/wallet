package screens

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/res"
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

type dataInfo struct {
	AppName    string `json:"app_name,omitempty"`
	StructName string `json:"struct_name,omitempty"`
	IsDBData   bool   `json:"is_db_data,omitempty"`
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
	Life       uint64 `json:"life,omitempty"`
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

func makeTransferOutList(w fyne.Window) fyne.Widget {
	var his = make(map[string]string)
	var numbers = make(map[string]uint64)
	// var oldCount uint64
	c := conf.Get()
	chain := widget.NewSelect(c.Chains, nil)
	chain.SetSelected(c.DefaultChain)
	number := widget.NewEntry()
	number.Disable()
	number.SetText("0")
	updateTime := widget.NewEntry()
	updateTime.Disable()
	box := widget.NewVBox()
	btn := widget.NewButton(res.GetLocalString("Update"), func() {
		address := conf.GetWallet().AddressStr
		// address := "01853433fb23a8e55663bc2b3cba0db2a8530acd60540fd9"
		// newCount := getIntOfDB(chain.Selected, "statTransferOut", address)
		newCount := getIntOfDB(chain.Selected, "", "statTransList", address)
		log.Println("update:", "statTransferOut", address, newCount)
		updateTime.SetText(time.Now().Local().String())

		number.SetText(fmt.Sprintf("%d", newCount))
		numbers[chain.Selected] = newCount
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			_, ok := his[key]
			if ok {
				break
			}
			// transKey := getStringOfDB(chain.Selected, "statTransferOut", key)
			transKey, _ := getStringOfDB(chain.Selected, "", "statTransList", key)
			if transKey != "" {
				his[key] = transKey
			}
		}
		box.Children = nil
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			item := his[key]
			if item == "" {
				item = fmt.Sprintf("error,index%d", i)
			}
			ent := widget.NewEntry()
			ent.SetText(item)
			ent.Disable()
			box.Children = append(box.Children, ent)
		}
		box.Refresh()
	})
	layout2 := layout.NewBorderLayout(nil, nil, nil, btn)
	updateItem := fyne.NewContainerWithLayout(layout2, btn, updateTime)
	form := &widget.Form{}
	form.Append(res.GetLocalString("Chain"), chain)
	form.Append(res.GetLocalString("Number"), number)
	form.Append(res.GetLocalString("UpdateTime"), updateItem)

	return widget.NewVBox(form, widget.NewGroup(res.GetLocalString("HistoryList"), box))
}

func makeTransferInList(w fyne.Window) fyne.Widget {
	var his = make(map[string]string)
	var numbers = make(map[string]uint64)
	c := conf.Get()
	chain := widget.NewSelect(c.Chains, nil)
	chain.SetSelected(c.DefaultChain)
	number := widget.NewEntry()
	number.Disable()
	number.SetText("0")
	updateTime := widget.NewEntry()
	updateTime.Disable()
	box := widget.NewVBox()
	btn := widget.NewButton(res.GetLocalString("Update"), func() {
		address := conf.GetWallet().AddressStr
		// address := "01853433fb23a8e55663bc2b3cba0db2a8530acd60540fd9"
		newCount := getIntOfDB(chain.Selected, "", "statTransferIn", address)
		log.Println("update:", "statTransferIn", address, newCount)
		updateTime.SetText(time.Now().Local().String())

		number.SetText(fmt.Sprintf("%d", newCount))
		numbers[chain.Selected] = newCount
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			_, ok := his[key]
			if ok {
				break
			}
			transKey, _ := getStringOfDB(chain.Selected, "", "statTransferIn", key)
			if transKey != "" {
				his[key] = transKey
			}
		}
		box.Children = nil
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			item := his[key]
			if item == "" {
				item = "error"
			}
			ent := widget.NewEntry()
			ent.SetText(item)
			ent.Disable()
			box.Children = append(box.Children, ent)
		}
		box.Refresh()
	})
	layout2 := layout.NewBorderLayout(nil, nil, nil, btn)
	updateItem := fyne.NewContainerWithLayout(layout2, btn, updateTime)
	form := &widget.Form{}
	form.Append(res.GetLocalString("Chain"), chain)
	form.Append(res.GetLocalString("Number"), number)
	form.Append(res.GetLocalString("UpdateTime"), updateItem)
	// btn.Tapped()
	return widget.NewVBox(form, widget.NewGroup(res.GetLocalString("HistoryList"), box))
}

// move coins to other chain
func makeMoveList(w fyne.Window) fyne.Widget {
	var his = make(map[string]string)
	var numbers = make(map[string]uint64)
	c := conf.Get()
	chain := widget.NewSelect(c.Chains, nil)
	chain.SetSelected(c.DefaultChain)
	number := widget.NewEntry()
	number.Disable()
	number.SetText("0")
	updateTime := widget.NewEntry()
	updateTime.Disable()
	// addr := widget.NewEntry()
	// addr.SetText("01f66d878a44f1a0ea8ec699c842ce9783b6df590a1584f7")
	box := widget.NewVBox()
	btn := widget.NewButton(res.GetLocalString("Update"), func() {
		address := conf.GetWallet().AddressStr
		// address := addr.Text
		newCount := getIntOfDB(chain.Selected, "", "statMove", address)
		log.Println("update:", "statMove", address, newCount)
		updateTime.SetText(time.Now().Local().String())

		number.SetText(fmt.Sprintf("%d", newCount))
		numbers[chain.Selected] = newCount
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			_, ok := his[key]
			if ok {
				break
			}
			transKey, _ := getStringOfDB(chain.Selected, "", "statMove", key)
			if transKey != "" {
				his[key] = transKey
			}
		}
		box.Children = nil
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			item := his[key]
			if item == "" {
				item = "error"
			}
			ent := widget.NewEntry()
			ent.SetText(item)
			ent.Disable()
			box.Children = append(box.Children, ent)
		}
		box.Refresh()
	})
	layout2 := layout.NewBorderLayout(nil, nil, nil, btn)
	updateItem := fyne.NewContainerWithLayout(layout2, btn, updateTime)
	form := &widget.Form{}
	form.Append(res.GetLocalString("Chain"), chain)
	form.Append(res.GetLocalString("Number"), number)
	form.Append(res.GetLocalString("UpdateTime"), updateItem)

	return widget.NewVBox(form, widget.NewGroup(res.GetLocalString("HistoryList"), box))
}

// HistoryScreen history of transaction
func HistoryScreen(w fyne.Window) fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem(res.GetLocalString("TransferIn"), makeTransferInList(w)),
		widget.NewTabItem(res.GetLocalString("TransferOut"), makeTransferOutList(w)),
		widget.NewTabItem(res.GetLocalString("Move"), makeMoveList(w)),
	)
}

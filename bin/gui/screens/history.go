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

const coreName = "e4a05b2b8a4de21d9e6f26e9d7992f7f33e89689f3015f3fc8a3a3278815e28c"

type dataInfo struct {
	AppName    string `json:"app_name,omitempty"`
	StructName string `json:"struct_name,omitempty"`
	IsDBData   bool   `json:"is_db_data,omitempty"`
	Key        string `json:"key,omitempty"`
	Value      string `json:"value,omitempty"`
	Life       uint64 `json:"life,omitempty"`
}

func getIntOfDB(chain, structName, address string) uint64 {
	value := getStringOfDB(chain, structName, address)
	if value == "" {
		return 0
	}
	out, err := strconv.ParseUint(value, 16, 64)
	if err != nil {
		log.Println("parse error.", value, err)
	}

	return out
}
func getStringOfDB(chain, structName, address string) string {
	urlStr := conf.Get(conf.APIServer)
	urlStr += "/api/v1/" + chain + "/data?app_name=" + coreName
	urlStr += "&is_db_data=true&struct_name=" + structName
	urlStr += "&key=" + address
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
	err = json.Unmarshal(data, &info)
	if err != nil || info.Value == "" {
		log.Println("not value.", urlStr)
		return ""
	}
	// log.Println("success to get:", urlStr, info.Value)
	return info.Value
}

func makeTransferOutList(w fyne.Window) fyne.Widget {
	var his = make(map[string]string)
	var numbers = make(map[string]uint64)
	// var oldCount uint64
	chain := widget.NewSelect([]string{"1", "2"}, nil)
	chain.SetSelected("1")
	number := widget.NewEntry()
	number.Disable()
	number.SetText("0")
	updateTime := widget.NewEntry()
	updateTime.Disable()
	box := widget.NewVBox()
	btn := widget.NewButton(res.GetLocalString("Update"), func() {
		address := conf.GetWallet().AddressStr
		// address := "01853433fb23a8e55663bc2b3cba0db2a8530acd60540fd9"
		newCount := getIntOfDB(chain.Selected, "statTransferOut", address)
		log.Println("update:", "statTransferOut", address, newCount)
		updateTime.SetText(time.Now().Local().String())
		old := numbers[chain.Selected]
		if old >= newCount {
			return
		}
		number.SetText(fmt.Sprintf("%d", newCount))
		numbers[chain.Selected] = newCount
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			_, ok := his[key]
			if ok {
				break
			}
			transKey := getStringOfDB(chain.Selected, "statTransferOut", key)
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
	form.Append(res.GetLocalString("TransferNumber"), number)
	form.Append(res.GetLocalString("UpdateTime"), updateItem)

	return widget.NewVBox(form, widget.NewGroup(res.GetLocalString("HistoryList"), box))
}

func makeTransferInList(w fyne.Window) fyne.Widget {
	var his = make(map[string]string)
	var numbers = make(map[string]uint64)
	// var oldCount uint64
	chain := widget.NewSelect([]string{"1", "2"}, nil)
	chain.SetSelected("1")
	number := widget.NewEntry()
	number.Disable()
	number.SetText("0")
	updateTime := widget.NewEntry()
	updateTime.Disable()
	box := widget.NewVBox()
	btn := widget.NewButton(res.GetLocalString("Update"), func() {
		address := conf.GetWallet().AddressStr
		// address := "01853433fb23a8e55663bc2b3cba0db2a8530acd60540fd9"
		newCount := getIntOfDB(chain.Selected, "statTransferIn", address)
		log.Println("update:", "statTransferIn", address, newCount)
		updateTime.SetText(time.Now().Local().String())
		old := numbers[chain.Selected]
		if old >= newCount {
			return
		}
		number.SetText(fmt.Sprintf("%d", newCount))
		numbers[chain.Selected] = newCount
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			_, ok := his[key]
			if ok {
				break
			}
			transKey := getStringOfDB(chain.Selected, "statTransferIn", key)
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
	form.Append(res.GetLocalString("TransferNumber"), number)
	form.Append(res.GetLocalString("UpdateTime"), updateItem)
	// btn.Tapped()
	return widget.NewVBox(form, widget.NewGroup(res.GetLocalString("HistoryList"), box))
}

// move coins to other chain
func makeMoveList(w fyne.Window) fyne.Widget {
	var his = make(map[string]string)
	var numbers = make(map[string]uint64)
	// var oldCount uint64
	chain := widget.NewSelect([]string{"1", "2"}, nil)
	chain.SetSelected("1")
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
		newCount := getIntOfDB(chain.Selected, "statMove", address)
		log.Println("update:", "statMove", address, newCount)
		updateTime.SetText(time.Now().Local().String())
		old := numbers[chain.Selected]
		if old >= newCount {
			return
		}
		number.SetText(fmt.Sprintf("%d", newCount))
		numbers[chain.Selected] = newCount
		for i := newCount; i > 0 && i+10 > newCount; i-- {
			id := fmt.Sprintf("%016x", i)
			key := address + id
			_, ok := his[key]
			if ok {
				break
			}
			transKey := getStringOfDB(chain.Selected, "statMove", key)
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
			} else if len(item) == 16 {
				val, _ := strconv.ParseUint(item, 16, 64)
				item = fmt.Sprintf("(In, Block ID)%d", val)
			} else {
				item = "(Out)" + item
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
	form.Append(res.GetLocalString("TransferNumber"), number)
	// form.Append(res.GetLocalString("Address"), addr)
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

package screens

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/event"
	"github.com/lengzhao/wallet/bin/gui/res"
)

type blockInfo struct {
	Time       uint64 `json:"time,omitempty"`
	Previous   string `json:"previous,omitempty"`
	Parent     string `json:"parent,omitempty"`
	LeftChild  string `json:"left_child,omitempty"`
	RightChild string `json:"right_child,omitempty"`
	Producer   string `json:"producer,omitempty"`
	Chain      uint64 `json:"chain,omitempty"`
	Index      uint64 `json:"index,omitempty"`
	Key        string `json:"key,omitempty"`
}

func searchBlock(w fyne.Window) fyne.Widget {
	srcChain := widget.NewSelect([]string{"1", "2"}, nil)
	srcChain.SetSelected("1")
	srcIndex := widget.NewEntry()
	srcKey := widget.NewEntry()

	eChain := widget.NewEntry()
	eChain.Disable()
	eKey := widget.NewEntry()
	eKey.Disable()
	eTime := widget.NewEntry()
	eTime.Disable()
	eProducer := widget.NewEntry()
	eProducer.Disable()
	eIndex := widget.NewEntry()
	eIndex.Disable()
	ePrevious := widget.NewEntry()
	ePrevious.Disable()
	eParent := widget.NewEntry()
	eParent.Disable()
	eLeftChild := widget.NewEntry()
	eLeftChild.Disable()
	eRightChild := widget.NewEntry()
	eRightChild.Disable()

	reqBlock := func(chain, index uint64, key string) error {
		apiServer := conf.Get(conf.APIServer)
		urlStr1 := fmt.Sprintf("%s/api/v1/%d/block/info?index=%d&key=%s", apiServer, chain, index, key)
		resp, err := http.Get(urlStr1)
		if err != nil {
			log.Print("fail to do get block,", err)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Print("error response when getting block,", resp.Status)
			return fmt.Errorf("error Status:%s", resp.Status)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print("fail to read body(block),", err)
			return err
		}
		info := blockInfo{}
		err = json.Unmarshal(data, &info)
		if err != nil {
			log.Println("fail to unmarshal,", string(data), err)
			return err
		}
		eChain.SetText(fmt.Sprintf("%d", info.Chain))
		eKey.SetText(info.Key)
		t := time.Unix(int64(info.Time/1000), 0)
		eTime.SetText(t.Local().String())
		eProducer.SetText(info.Producer)
		eIndex.SetText(fmt.Sprintf("%d", info.Index))
		ePrevious.SetText(info.Previous)
		eParent.SetText(info.Parent)
		eLeftChild.SetText(info.LeftChild)
		eRightChild.SetText(info.RightChild)
		return nil
	}

	srcform := &widget.Form{OnCancel: func() {
		srcIndex.SetText("")
		srcKey.SetText("")
	},
		OnSubmit: func() {
			cid, err := strconv.ParseUint(srcChain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			var index uint64
			var key string
			if srcIndex.Text != "" {
				index, err = strconv.ParseUint(srcIndex.Text, 10, 64)
				if err != nil {
					dialog.ShowError(fmt.Errorf("error index"), w)
					return
				}
			}
			if srcKey.Text != "" {
				key = srcKey.Text
			}
			err = reqBlock(cid, index, key)
			if err != nil {
				dialog.ShowError(err, w)
			}
		}}
	srcform.Append(res.GetLocalString("Chain"), srcChain)
	srcform.Append(res.GetLocalString("Index"), srcIndex)
	srcform.Append(res.GetLocalString("Key"), srcKey)

	showForm := &widget.Form{}
	showForm.Append(res.GetLocalString("Chain"), eChain)
	showForm.Append(res.GetLocalString("Key"), eKey)
	showForm.Append(res.GetLocalString("Time"), eTime)
	showForm.Append(res.GetLocalString("Producer"), eProducer)
	showForm.Append(res.GetLocalString("Index"), eIndex)
	showForm.Append(res.GetLocalString("Previous"), ePrevious)
	showForm.Append(res.GetLocalString("Parent"), eParent)
	showForm.Append(res.GetLocalString("LeftChild"), eLeftChild)
	showForm.Append(res.GetLocalString("RightChild"), eRightChild)

	info := widget.NewGroup("Info", showForm)

	win := widget.NewVBox(srcform, info)
	return win
}

type transactionInfo struct {
	Time   uint64
	User   []byte
	Chain  uint64
	Energy uint64
	Cost   uint64
	Ops    uint8
	Key    []byte
	Size   int
	Others map[string]interface{}
}

func searchTransaction(w fyne.Window) fyne.Widget {
	srcChain := widget.NewSelect([]string{"1", "2"}, nil)
	srcChain.SetSelected("1")
	srcKey := widget.NewEntry()
	srcKey.SetText("0c64e484f3b329fea41a03be2677161eaac92741105cb0548b6ec4a5529efc71")

	eChain := widget.NewEntry()
	eChain.Disable()
	eKey := widget.NewEntry()
	eKey.Disable()
	eTime := widget.NewEntry()
	eTime.Disable()
	eUser := widget.NewEntry()
	eUser.Disable()
	eCost := widget.NewEntry()
	eCost.Disable()
	unit := conf.Get(conf.CoinUnit)
	eCUnit := widget.NewLabel(unit)
	costLayout := layout.NewBorderLayout(nil, nil, nil, eCUnit)
	eCostUnit := fyne.NewContainerWithLayout(costLayout, eCUnit, eCost)
	eEnergy := widget.NewEntry()
	eEnergy.Disable()
	eEUnit := widget.NewLabel(unit)
	energyLayout := layout.NewBorderLayout(nil, nil, nil, eEUnit)
	eEnergyUnit := fyne.NewContainerWithLayout(energyLayout, eEUnit, eEnergy)
	eBlockID := widget.NewEntry()
	eBlockID.Disable()
	eOpcode := widget.NewEntry()
	eOpcode.Disable()
	eOthers := widget.NewMultiLineEntry()
	eOthers.Disable()

	event.RegisterConsumer(event.EChangeUnit, func(e string, param ...interface{}) error {
		eCost.SetText("")
		eEnergy.SetText("")
		eCUnit.SetText(conf.Get(conf.CoinUnit))
		eEUnit.SetText(conf.Get(conf.CoinUnit))
		return nil
	})

	reqTrans := func(chain uint64, key string) error {
		apiServer := conf.Get(conf.APIServer)
		urlStr1 := fmt.Sprintf("%s/api/v1/%d/transaction/info?key=%s", apiServer, chain, key)
		resp, err := http.Get(urlStr1)
		if err != nil {
			log.Print("fail to do get trans,", err)
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			log.Print("error response when getting trans,", resp.Status)
			return fmt.Errorf("error Status:%s", resp.Status)
		}
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print("fail to read body(trans),", err)
			return err
		}
		info := transactionInfo{}
		json.Unmarshal(data, &info)
		eChain.SetText(fmt.Sprintf("%d", info.Chain))
		eKey.SetText(hex.EncodeToString(info.Key))
		t := time.Unix(int64(info.Time/1000), 0)
		eTime.SetText(t.Local().String())
		eUser.SetText(hex.EncodeToString(info.User))
		base := res.GetBaseOfUnit(eCUnit.Text)
		eCost.SetText(fmt.Sprintf("%.3f", float64(info.Cost)/float64(base)))
		eEnergy.SetText(fmt.Sprintf("%.3f", float64(info.Energy)/float64(base)))
		eBlockID.SetText(fmt.Sprintf("%v", info.Others["BlockID"]))
		eOpcode.SetText(fmt.Sprintf("%d", info.Ops))
		switch info.Ops {
		case 0:
			var peer []byte
			d, err := json.Marshal(info.Others["peer"])
			if err != nil {
				log.Println("fail to marshal app name.", err)
				return nil
			}
			json.Unmarshal(d, &peer)
			eOthers.SetText("Peer:" + hex.EncodeToString(peer))
		case 1:
			eOthers.SetText(fmt.Sprintf("DstChain:%v", info.Others["peer"]))
		case 2:
			eOthers.SetText(fmt.Sprintf("New Chain:%v", info.Others["peer"]))
		case 3:
			data := fmt.Sprintf("App Name:%v\nPublic:%v\nEnable Run:%v\nEnable Import:%v",
				info.Others["app_name"], info.Others["is_public"], info.Others["enable_run"],
				info.Others["enable_import"])
			eOthers.SetText(data)
		case 4:
			var name []byte
			d, err := json.Marshal(info.Others["name"])
			if err != nil {
				log.Println("fail to marshal app name.", err)
				return nil
			}
			json.Unmarshal(d, &name)
			eOthers.SetText(fmt.Sprintf("App Name:%x", name))
		case 5:
		case 6:
			eOthers.SetText(fmt.Sprintf("Index:%v", info.Others["index"]))
		}
		// eOthers.Hide()
		return nil
	}

	srcform := &widget.Form{OnCancel: func() {
		srcKey.SetText("")
	},
		OnSubmit: func() {
			cid, err := strconv.ParseUint(srcChain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			if srcKey.Text == "" {
				dialog.ShowError(fmt.Errorf("error transaction key"), w)
				return
			}
			err = reqTrans(cid, srcKey.Text)
			if err != nil {
				dialog.ShowError(err, w)
			}
		}}
	srcform.Append(res.GetLocalString("Chain"), srcChain)
	srcform.Append(res.GetLocalString("Key"), srcKey)

	showForm := &widget.Form{}
	showForm.Append(res.GetLocalString("Chain"), eChain)
	showForm.Append(res.GetLocalString("Key"), eKey)
	showForm.Append(res.GetLocalString("Time"), eTime)
	showForm.Append(res.GetLocalString("User"), eUser)
	showForm.Append(res.GetLocalString("Cost"), eCostUnit)
	showForm.Append(res.GetLocalString("Energy"), eEnergyUnit)
	showForm.Append(res.GetLocalString("BlockID"), eBlockID)
	showForm.Append(res.GetLocalString("Opcode"), eOpcode)
	showForm.Append(res.GetLocalString("Other"), eOthers)

	info := widget.NewGroup("Info", showForm)

	win := widget.NewVBox(srcform, info)
	return win
}

// SearchScreen search transaction or block
func SearchScreen(w fyne.Window) fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem(res.GetLocalString("Block"), searchBlock(w)),
		widget.NewTabItem(res.GetLocalString("Transaction"), searchTransaction(w)),
	)
}

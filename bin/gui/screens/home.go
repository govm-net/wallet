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
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/event"
	"github.com/lengzhao/wallet/bin/gui/res"
)

// Account account
type account struct {
	Chain   uint64 `json:"chain,omitempty"`
	Address string `json:"address,omitempty"`
	Cost    uint64 `json:"cost,omitempty"`
}

// VoteInfo vote info
type VoteInfo struct {
	Admin    [24]byte
	Cost     uint64
	StartDay uint64
}

func getAccount(chain uint64, address string) account {
	apiServer := conf.Get().APIServer
	info := account{}
	urlStr1 := fmt.Sprintf("%s/api/v1/%d/account?address=%s", apiServer, chain, address)
	resp, err := http.Get(urlStr1)
	if err != nil {
		log.Print("fail to do get account,", err)
		return info
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Print("error response when getting account,", resp.Status)
		return info
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("fail to read body(account),", err)
		return info
	}
	err = json.Unmarshal(data, &info)
	if err != nil {
		log.Println("fail to unmarshal,", string(data), err)
	}
	// log.Println("get account:", info)
	return info
}

// AccountScreen account info
func AccountScreen(w fyne.Window) fyne.Widget {
	address := widget.NewEntry()
	address.Disable()
	addrCpy := fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widget.NewButtonWithIcon(res.GetLocalString("Copy"), theme.ContentCopyIcon(), func() {
			w.Clipboard().SetContent(conf.GetWallet().AddressStr)
		}))

	bl0 := layout.NewBorderLayout(nil, nil, nil, addrCpy)
	addrItem := fyne.NewContainerWithLayout(bl0, addrCpy, address)

	form := &widget.Form{}
	form.Append(res.GetLocalString("Address"), addrItem)

	c := conf.Get()
	allWidget := make(map[string]*widget.Entry)
	for i, chain := range c.Chains {
		fmt.Println(i, chain)
		chainW := widget.NewEntry()
		chainW.Disable()
		allWidget["chain"+chain] = chainW
		form.Append("Chain"+chain, chainW)
		balanceW := widget.NewEntry()
		balanceW.Disable()
		allWidget["balance"+chain] = balanceW
		unit := widget.NewLabel(c.CoinUnit)
		borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
		showBalance := fyne.NewContainerWithLayout(borderLayout, unit, balanceW)
		form.Append("Balance"+chain, showBalance)

		votes := widget.NewEntry()
		votes.Disable()
		allWidget["votes"+chain] = votes
		form.Append("Votes"+chain, votes)
		// voteReward := widget.NewEntry()
		// voteReward.Disable()
		// allWidget["voteReward"+chain] = voteReward
		// form.Append("VoteReward"+chain, voteReward)
	}

	eTime := widget.NewEntry()
	eTime.Disable()
	form.Append("Time", eTime)

	updateAccount := func() {
		addr := conf.GetWallet().AddressStr
		address.SetText(addr)
		unit := c.CoinUnit
		base := res.GetBaseOfUnit(unit)
		for _, it := range c.Chains {
			chain, _ := strconv.ParseUint(it, 10, 64)
			info := getAccount(chain, addr)
			cost := float64(info.Cost) / float64(base)
			chainW := allWidget["chain"+it]
			chainW.SetText(fmt.Sprintf("%d", info.Chain))
			chainW.Refresh()
			balanceW := allWidget["balance"+it]
			balanceW.SetText(fmt.Sprintf("%.3f", cost))
			balanceW.Refresh()

			data := getDataFromServer(chain, c.APIServer, "", "dbVote", addr)
			var vInfo1 VoteInfo
			if len(data) > 0 {
				Decode(data, &vInfo1)
				adminOfVote = fmt.Sprintf("%x", vInfo1.Admin)
				fmt.Printf("set admin of vote:%x\n", vInfo1.Admin)
				event.Send(event.EAdminOfVote, w)
			}
			votesW := allWidget["votes"+it]
			votesW.Text = fmt.Sprintf("%d", vInfo1.Cost/1000000000)
			votesW.Refresh()
		}
		eTime.SetText(time.Now().Local().String())
		eTime.Refresh()
	}

	event.RegisterConsumer(event.EShowHome, func(e string, param ...interface{}) error {
		address.SetText(conf.GetWallet().AddressStr)
		updateAccount()
		return nil
	})

	btn := fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widget.NewButtonWithIcon(res.GetLocalString("Update"), theme.ViewRefreshIcon(), func() {
			updateAccount()
		}))

	win := widget.NewVBox(form, btn)
	return win
}

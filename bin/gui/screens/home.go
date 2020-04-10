package screens

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func getAccount(chain uint64, address string) account {
	apiServer := conf.Get(conf.APIServer)
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
	chain1 := widget.NewEntry()
	chain1.Disable()
	blance1 := widget.NewEntry()
	blance1.Disable()
	unit1 := widget.NewLabel(conf.Get(conf.CoinUnit))
	borderLayout1 := layout.NewBorderLayout(nil, nil, nil, unit1)
	showBlance1 := fyne.NewContainerWithLayout(borderLayout1, unit1, blance1)

	chain2 := widget.NewEntry()
	chain2.Disable()
	blance2 := widget.NewEntry()
	blance2.Disable()
	unit2 := widget.NewLabel(conf.Get(conf.CoinUnit))
	borderLayout2 := layout.NewBorderLayout(nil, nil, nil, unit2)
	showBlance2 := fyne.NewContainerWithLayout(borderLayout2, unit2, blance2)

	eTime := widget.NewEntry()
	eTime.Disable()

	form := &widget.Form{}
	form.Append(res.GetLocalString("Address"), address)
	form.Append("Chain1", chain1)
	form.Append("Blance1", showBlance1)
	form.Append("Chain2", chain2)
	form.Append("Blance2", showBlance2)
	form.Append("Time", eTime)

	updateAccount := func() {
		addr := conf.GetWallet().AddressStr
		address.SetText(addr)
		info := getAccount(1, addr)
		unit := conf.Get(conf.CoinUnit)
		unit1.SetText(unit)
		unit2.SetText(unit)
		base := res.GetBaseOfUnit(unit)
		cost1 := float64(info.Cost) / float64(base)
		chain1.SetText(fmt.Sprintf("%d", info.Chain))
		blance1.SetText(fmt.Sprintf("%.3f", cost1))
		// log.Println("get account1:", info, base)
		info2 := getAccount(2, addr)
		cost2 := float64(info2.Cost) / float64(base)
		chain2.SetText(fmt.Sprintf("%d", info2.Chain))
		blance2.SetText(fmt.Sprintf("%.3f", cost2))
		log.Println("get account:", base, info, info2)
		eTime.SetText(time.Now().Local().String())
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

package screens

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/event"
	"github.com/lengzhao/wallet/bin/gui/res"
	"github.com/lengzhao/wallet/trans"
)

func makeAppRunTab(w fyne.Window) fyne.Widget {
	c := conf.Get()
	chain := widget.NewSelect(c.Chains, nil)
	chain.SetSelected(c.DefaultChain)
	app := widget.NewEntry()
	app.SetPlaceHolder("app name")
	var dataType = "string"
	dType := widget.NewSelect([]string{"string", "hex"}, func(nv string) {
		dataType = nv
	})
	dType.SetSelected(dataType)
	dEntry := widget.NewEntry()
	amount := widget.NewEntry()
	unit := widget.NewLabel(c.CoinUnit)
	energy := widget.NewEntry()
	energy.SetText("0.1")
	unit2 := widget.NewLabel(c.CoinUnit)
	result := widget.NewEntry()
	result.Disable()

	event.RegisterConsumer(event.EChangeUnit, func(e string, param ...interface{}) error {
		unit.SetText(c.CoinUnit)
		return nil
	})

	form := &widget.Form{
		OnCancel: func() {
			app.SetText("")
			dEntry.SetText("")
			amount.SetText("")
			result.SetText("")
		},
		OnSubmit: func() {
			result.SetText("")
			costF, err := strconv.ParseFloat(amount.Text, 10)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error amount"), w)
				return
			}
			amount.SetPlaceHolder(amount.Text)
			amount.SetText("")
			cid, err := strconv.ParseUint(chain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			if app.Text == "" {
				dialog.ShowError(fmt.Errorf("require app name"), w)
				return
			}
			var param []byte
			if dEntry.Text != "" {
				switch dataType {
				case "hex":
					param, err = hex.DecodeString(dEntry.Text)
					if err != nil {
						dialog.ShowError(fmt.Errorf("error hex data.%s", err), w)
						return
					}
				default:
					param = []byte(dEntry.Text)
				}
			}

			base := res.GetBaseOfUnit(c.CoinUnit)
			cost := uint64(costF * float64(base))
			myWlt := conf.GetWallet()
			trans := trans.NewTransaction(cid, myWlt.Address, cost)
			err = trans.RunApp(app.Text, param)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}

			td := trans.GetSignData()
			sign := myWlt.Sign(td)
			trans.SetTheSign(sign)
			td = trans.Output()
			key := trans.Key[:]

			err = postTrans(cid, td)
			if err != nil {
				// result.SetText(fmt.Sprintf("%s", err))
				dialog.ShowError(err, w)
				return
			}
			log.Printf("new transfer:%x\n", key)
			// dialog.ShowInformation("transaction", fmt.Sprintf("%x", key), w)
			result.SetText(fmt.Sprintf("%x", key))
			app.SetPlaceHolder(app.Text)
			app.SetText("")
		},
	}
	form.Append(res.GetLocalString("Chain"), chain)
	form.Append(res.GetLocalString("App"), app)
	form.Append(res.GetLocalString("Data Type"), dType)
	form.Append(res.GetLocalString("Data"), dEntry)
	borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
	form.Append(res.GetLocalString("Amount"), fyne.NewContainerWithLayout(borderLayout, unit, amount))
	borderLayout2 := layout.NewBorderLayout(nil, nil, nil, unit2)
	form.Append(res.GetLocalString("Energy"), fyne.NewContainerWithLayout(borderLayout2, unit2, energy))

	return widget.NewVBox(form, result)
}

type wGovmBody struct {
	EthAddr string `json:"eth_addr,omitempty"`
}

func makeWGOVMTab(w fyne.Window) fyne.Widget {
	c := conf.Get()
	desc := widget.NewMultiLineEntry()
	desc.Disable()
	desc.SetText(res.GetLocalString("wgovm_desc"))
	ethEntry := widget.NewEntry()
	amount := widget.NewEntry()
	unit := widget.NewLabel(c.CoinUnit)
	energy := widget.NewEntry()
	energy.SetText("0.1")
	unit2 := widget.NewLabel(c.CoinUnit)
	result := widget.NewEntry()
	result.Disable()

	event.RegisterConsumer(event.EChangeUnit, func(e string, param ...interface{}) error {
		unit.SetText(c.CoinUnit)
		return nil
	})

	form := &widget.Form{
		OnCancel: func() {
			ethEntry.SetText("")
			amount.SetText("")
			result.SetText("")
		},
		OnSubmit: func() {
			result.SetText("")
			costF, err := strconv.ParseFloat(amount.Text, 10)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error amount"), w)
				return
			}
			amount.SetPlaceHolder(amount.Text)
			amount.SetText("")

			engF, err := strconv.ParseFloat(energy.Text, 10)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error energy"), w)
				return
			}

			if costF < 5000 {
				dialog.ShowError(fmt.Errorf("require amount >= 5000govm"), w)
				return
			}

			if ethEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("require eth address"), w)
				return
			}
			param := []byte{0}
			info := wGovmBody{ethEntry.Text}
			d, _ := json.Marshal(info)
			param = append(param, d...)
			base := res.GetBaseOfUnit(c.CoinUnit)
			cost := uint64(costF * float64(base))
			myWlt := conf.GetWallet()
			trans := trans.NewTransaction(1, myWlt.Address, cost)
			err = trans.RunApp("a99b97a411b45c91779e1609386fa18484b8e50016379bd98c6822b491247b9b", param)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			trans.Energy = uint64(engF * float64(base))

			td := trans.GetSignData()
			sign := myWlt.Sign(td)
			trans.SetTheSign(sign)
			td = trans.Output()
			key := trans.Key[:]

			err = postTrans(1, td)
			if err != nil {
				// result.SetText(fmt.Sprintf("%s", err))
				log.Println(ethEntry.Text, trans.Cost, trans.Energy, err)
				dialog.ShowError(err, w)
				return
			}
			log.Printf("new transfer:%x\n", key)
			// dialog.ShowInformation("transaction", fmt.Sprintf("%x", key), w)
			result.SetText(fmt.Sprintf("%x", key))
		},
	}

	form.Append(res.GetLocalString("Description"), desc)
	form.Append(res.GetLocalString("Eth Address"), ethEntry)
	borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
	form.Append(res.GetLocalString("Amount"), fyne.NewContainerWithLayout(borderLayout, unit, amount))
	borderLayout2 := layout.NewBorderLayout(nil, nil, nil, unit2)
	form.Append(res.GetLocalString("Energy"), fyne.NewContainerWithLayout(borderLayout2, unit2, energy))

	return widget.NewVBox(form, result)
}

// AppScreen shows a panel containing widget demos
func AppScreen(w fyne.Window) fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem(res.GetLocalString("Run APP"), makeAppRunTab(w)),
		widget.NewTabItem(res.GetLocalString("wGOVM"), makeWGOVMTab(w)),
	)
}

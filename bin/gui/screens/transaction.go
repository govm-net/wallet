package screens

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
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

func postTrans(chain uint64, data []byte) error {
	apiServer := conf.Get(conf.APIServer)
	buf := bytes.NewBuffer(data)
	urlStr := fmt.Sprintf("%s/api/v1/%d/transaction/new", apiServer, chain)
	req, err := http.NewRequest(http.MethodPost, urlStr, buf)
	if err != nil {
		log.Println("fail to new http request:", urlStr, err)
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("fail to do request:", urlStr, err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Println("error status:", resp.Status)
		data, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error response status:%s,msg:%s", resp.Status, data)
	}
	return nil
}

func makeTransferTab(w fyne.Window) fyne.Widget {
	desc := widget.NewLabel(res.GetLocalString("transfer.desc"))
	chain := widget.NewSelect([]string{"1", "2"}, nil)
	chain.SetSelected("1")
	peer := widget.NewEntry()
	peer.SetPlaceHolder("peer address")
	// peer.SetText("01853433fb23a8e55663bc2b3cba0db2a8530acd60540fd9")
	amount := widget.NewEntry()
	unit := widget.NewLabel(conf.Get(conf.CoinUnit))
	result := widget.NewLabel("")

	event.RegisterConsumer(event.EChangeUnit, func(e string, param ...interface{}) error {
		unit.SetText(conf.Get(conf.CoinUnit))
		return nil
	})

	form := &widget.Form{
		OnCancel: func() {
			peer.SetText("")
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
			cid, err := strconv.ParseUint(chain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			base := res.GetBaseOfUnit(conf.Get(conf.CoinUnit))
			cost := uint64(costF * float64(base))
			myWlt := conf.GetWallet()
			trans := trans.NewTransaction(cid, myWlt.Address)
			err = trans.CreateTransfer(peer.Text, "", cost, 0)
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			var eid int
			rid := rand.Int()
			eid = event.RegisterConsumer(event.EResponsePwd, func(e string, param ...interface{}) error {
				event.Unregister(eid)
				if len(param) < 2 {
					result.SetText("fail to get password")
					return nil
				}
				if param[0].(int) != rid {
					peer.SetText("")
					amount.SetText("")
					result.SetText("fail to get password")
					return nil
				}
				pwd := param[1].(string)
				if pwd == "" {
					peer.SetText("")
					amount.SetText("")
					result.SetText("cancel")
					return nil
				}
				if !conf.CheckPassword(pwd) {
					result.SetText("error password")
					return nil
				}
				td := trans.GetSignData()
				myWlt := conf.GetWallet()
				sign := myWlt.Sign(td)
				trans.SetSign(sign)
				td = trans.Output()
				key := trans.Key[:]

				err = postTrans(cid, td)
				if err != nil {
					result.SetText(fmt.Sprintf("fail to send trans:%s", err))
					return nil
				}
				log.Printf("new transfer:%x\n", key)
				// dialog.ShowInformation("transaction", fmt.Sprintf("%x", key), w)
				result.SetText(fmt.Sprintf("trans:%x", key))
				peer.SetText("")
				amount.SetText("")
				return nil
			})
			event.Send(event.ERequsetPwd, rid)
		},
	}
	form.Append(res.GetLocalString("Description"), desc)
	form.Append(res.GetLocalString("Chain"), chain)
	form.Append(res.GetLocalString("transfer.peer"), peer)
	borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
	form.Append(res.GetLocalString("Amount"), fyne.NewContainerWithLayout(borderLayout, unit, amount))

	return widget.NewVBox(form, result)
}

func makeMoveTransTab(w fyne.Window) fyne.Widget {
	desc := widget.NewLabel(res.GetLocalString("move.desc"))
	srcChain := widget.NewSelect([]string{"1", "2"}, nil)
	srcChain.SetSelected("1")
	dstChain := widget.NewSelect([]string{"1", "2"}, nil)
	dstChain.SetSelected("2")
	amount := widget.NewEntry()
	// amount.SetText("1.1")
	unit := widget.NewLabel(conf.Get(conf.CoinUnit))
	unit.TextStyle.Bold = true
	result := widget.NewLabel("")

	event.RegisterConsumer(event.EChangeUnit, func(e string, param ...interface{}) error {
		unit.SetText(conf.Get(conf.CoinUnit))
		return nil
	})

	form := &widget.Form{
		OnCancel: func() {
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
			srcid, err := strconv.ParseUint(srcChain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			dstid, err := strconv.ParseUint(dstChain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			if srcid == dstid {
				dialog.ShowError(fmt.Errorf("same chain id"), w)
				return
			}
			base := res.GetBaseOfUnit(conf.Get(conf.CoinUnit))
			cost := uint64(costF * float64(base))
			if cost == 0 {
				dialog.ShowError(fmt.Errorf("error amount"), w)
				return
			}
			trans := trans.NewTransaction(srcid, conf.GetWallet().Address)
			trans.CreateMove(dstid, cost, 0)
			var eid int
			rid := rand.Int()
			eid = event.RegisterConsumer(event.EResponsePwd, func(e string, param ...interface{}) error {
				event.Unregister(eid)
				if len(param) < 2 {
					result.SetText("fail to get password")
					return nil
				}
				if param[0].(int) != rid {
					amount.SetText("")
					result.SetText("fail to get password")
					return nil
				}
				pwd := param[1].(string)
				if pwd == "" {
					amount.SetText("")
					result.SetText("cancel")
					return nil
				}
				if !conf.CheckPassword(pwd) {
					result.SetText("error password")
					return nil
				}
				td := trans.GetSignData()
				myWlt := conf.GetWallet()
				sign := myWlt.Sign(td)
				trans.SetSign(sign)
				td = trans.Output()
				key := trans.Key[:]

				err = postTrans(srcid, td)
				if err != nil {
					result.SetText(fmt.Sprintf("fail to send trans:%s", err))
					return nil
				}
				log.Printf("new transfer:%x\n", key)
				// dialog.ShowInformation("transaction", fmt.Sprintf("%x", key), w)
				result.SetText(fmt.Sprintf("trans:%x", key))
				amount.SetText("")
				return nil
			})
			event.Send(event.ERequsetPwd, rid)
		},
	}
	form.Append(res.GetLocalString("Description"), desc)
	form.Append(res.GetLocalString("move.from_chain"), srcChain)
	form.Append(res.GetLocalString("move.to_chain"), dstChain)
	borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
	form.Append(res.GetLocalString("Amount"), fyne.NewContainerWithLayout(borderLayout, unit, amount))

	return widget.NewVBox(form, result)
}

// TransactionScreen shows a panel containing widget demos
func TransactionScreen(w fyne.Window) fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem(res.GetLocalString("Transfer"), makeTransferTab(w)),
		widget.NewTabItem(res.GetLocalString("Move"), makeMoveTransTab(w)),
	)
}

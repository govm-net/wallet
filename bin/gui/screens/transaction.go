package screens

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
		return fmt.Errorf("error,status:%s,msg:%s", resp.Status, data)
	}
	return nil
}

func makeTransferTab(w fyne.Window) fyne.Widget {
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
			amount.SetText("")
			cid, err := strconv.ParseUint(chain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			base := res.GetBaseOfUnit(conf.Get(conf.CoinUnit))
			cost := uint64(costF * float64(base))
			myWlt := conf.GetWallet()
			trans := trans.NewTransaction(cid, myWlt.Address, cost)
			err = trans.CreateTransfer(peer.Text, "")
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
			result.SetText(fmt.Sprintf("trans:%x", key))
			peer.SetText("")
		},
	}
	form.Append(res.GetLocalString("Chain"), chain)
	form.Append(res.GetLocalString("transfer.peer"), peer)
	borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
	form.Append(res.GetLocalString("Amount"), fyne.NewContainerWithLayout(borderLayout, unit, amount))

	return widget.NewVBox(form, result)
}

func makeMoveTransTab(w fyne.Window) fyne.Widget {
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
			amount.SetText("")
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
			trans := trans.NewTransaction(srcid, conf.GetWallet().Address, cost)
			trans.CreateMove(dstid)
			td := trans.GetSignData()
			myWlt := conf.GetWallet()
			sign := myWlt.Sign(td)
			trans.SetTheSign(sign)
			td = trans.Output()
			key := trans.Key[:]

			err = postTrans(srcid, td)
			if err != nil {
				// result.SetText(fmt.Sprintf("%s", err))
				dialog.ShowError(err, w)
				return
			}
			log.Printf("new transfer:%x\n", key)
			result.SetText(fmt.Sprintf("trans:%x", key))
		},
	}
	form.Append(res.GetLocalString("move.from_chain"), srcChain)
	form.Append(res.GetLocalString("move.to_chain"), dstChain)
	borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
	form.Append(res.GetLocalString("Amount"), fyne.NewContainerWithLayout(borderLayout, unit, amount))

	return widget.NewVBox(form, result)
}

func makeVoteTab(w fyne.Window) fyne.Widget {
	chain := widget.NewSelect([]string{"1", "2"}, nil)
	chain.SetSelected("1")
	peer := widget.NewEntry()
	peer.SetPlaceHolder("admin address")
	peer.SetText("01ccaf415a3a6dc8964bf935a1f40e55654a4243ae99c709")
	votes := widget.NewEntry()
	result := widget.NewLabel("")

	form := &widget.Form{
		OnCancel: func() {
			votes.SetText("")
			result.SetText("")
		},
		OnSubmit: func() {
			result.SetText("")
			v, err := strconv.ParseUint(votes.Text, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error amount"), w)
				return
			}
			votes.SetText("")
			cid, err := strconv.ParseUint(chain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			base := res.GetBaseOfUnit("govm")
			cost := uint64(v * base)
			myWlt := conf.GetWallet()
			trans := trans.NewTransaction(cid, myWlt.Address, cost)
			err = trans.CreateVote(peer.Text)
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
			result.SetText(fmt.Sprintf("trans:%x", key))
		},
	}
	form.Append(res.GetLocalString("Chain"), chain)
	form.Append(res.GetLocalString("transfer.peer"), peer)
	form.Append(res.GetLocalString("Votes"), votes)

	return widget.NewVBox(form, result)
}

func makeUnvoteTab(w fyne.Window) fyne.Widget {
	chain := widget.NewSelect([]string{"1", "2"}, nil)
	chain.SetSelected("1")
	result := widget.NewLabel("")

	form := &widget.Form{
		OnSubmit: func() {
			result.SetText("")
			cid, err := strconv.ParseUint(chain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			myWlt := conf.GetWallet()
			trans := trans.NewTransaction(cid, myWlt.Address, 0)
			err = trans.Unvote()
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
			result.SetText(fmt.Sprintf("trans:%x", key))
		},
	}
	form.Append(res.GetLocalString("Chain"), chain)

	return widget.NewVBox(form, result)
}

// TransactionScreen shows a panel containing widget demos
func TransactionScreen(w fyne.Window) fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem(res.GetLocalString("Transfer"), makeTransferTab(w)),
		widget.NewTabItem(res.GetLocalString("Move"), makeMoveTransTab(w)),
		widget.NewTabItem(res.GetLocalString("Vote"), makeVoteTab(w)),
		widget.NewTabItem(res.GetLocalString("Unvote"), makeUnvoteTab(w)),
	)
}

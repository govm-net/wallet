package screens

import (
	"fmt"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/res"
)

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

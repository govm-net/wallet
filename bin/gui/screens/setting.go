package screens

import (
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/event"
	"github.com/lengzhao/wallet/bin/gui/res"
)

// SettingScreen setting
func SettingScreen(win fyne.Window) fyne.CanvasObject {
	desc := widget.NewLabel(res.GetLocalString("setting.desc"))
	unit := widget.NewSelect([]string{"govm", "t6", "t3", "t0"}, func(in string) {
		conf.Set(conf.CoinUnit, in)
		event.Send(event.EChangeUnit)
	})
	unit.SetSelected(conf.Get(conf.CoinUnit))

	var list = []string{"en", "zh"}
	selected := conf.Get(conf.Langure)
	if selected == "" {
		selected = list[0]
	}
	lng := widget.NewSelect(list, func(in string) {
		conf.Set(conf.Langure, in)
	})
	lng.SetSelected(selected)

	entry1 := widget.NewEntry()
	entry1.SetText(conf.Get(conf.APIServer))
	btn := widget.NewButton(res.GetLocalString("Set"), func() {
		conf.Set(conf.APIServer, entry1.Text)
	})
	borderLayout := layout.NewBorderLayout(nil, nil, nil, btn)
	apiServer := fyne.NewContainerWithLayout(borderLayout, btn, entry1)
	version := widget.NewLabel(conf.Vserion)
	link, _ := url.Parse("http://govm.net")
	website := widget.NewHyperlink("http://govm.net", link)

	form := &widget.Form{}
	form.Append(res.GetLocalString("Description"), desc)
	form.Append(res.GetLocalString("setting.unit"), unit)
	form.Append(res.GetLocalString("setting.language"), lng)
	form.Append(res.GetLocalString("setting.server"), apiServer)
	form.Append(res.GetLocalString("Version"), version)
	form.Append(res.GetLocalString("Website"), website)
	// form.Append(res.GetLocalString("Window Scale"), itemScale)

	return widget.NewVBox(form, layout.NewSpacer())
}

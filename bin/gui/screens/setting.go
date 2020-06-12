package screens

import (
	"fmt"
	"net/url"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
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

	oldPwd := widget.NewPasswordEntry()
	newPwd1 := widget.NewPasswordEntry()
	newPwd2 := widget.NewPasswordEntry()
	pwdForm := &widget.Form{
		OnSubmit: func() {
			if !conf.CheckPassword(oldPwd.Text) {
				dialog.ShowError(fmt.Errorf("error password"), win)
				return
			}
			if newPwd1.Text == "" {
				dialog.ShowError(fmt.Errorf("not new password"), win)
				return
			}
			if newPwd1.Text != newPwd2.Text {
				dialog.ShowError(fmt.Errorf("Two passwords are different"), win)
				return
			}
			err := conf.ChangePassword(newPwd1.Text)
			if err != nil {
				dialog.ShowError(err, win)
			} else {
				dialog.ShowInformation("Change Password", "Success", win)
				oldPwd.Text = ""
				newPwd1.Text = ""
				newPwd2.Text = ""
			}
		},
	}
	pwdForm.Append(res.GetLocalString("Old Password"), oldPwd)
	pwdForm.Append(res.GetLocalString("New Password"), newPwd1)
	pwdForm.Append(res.GetLocalString("New Password"), newPwd2)

	privBtn := fyne.NewContainerWithLayout(layout.NewCenterLayout(),
		widget.NewButtonWithIcon(res.GetLocalString("View Private Key"), theme.ViewFullScreenIcon(), func() {
			logo := canvas.NewImageFromResource(res.GetResource("point.svg"))
			logo.SetMinSize(fyne.NewSize(300, 5))
			content := widget.NewPasswordEntry()

			dialog.ShowCustomConfirm(res.GetLocalString("Password"), "", "",
				widget.NewVBox(logo, content), func(rst bool) {
					pwd := content.Text
					if pwd == "" {
						return
					}
					if !conf.CheckPassword(pwd) {
						dialog.ShowInformation("Error", "error password", win)
						return
					}
					wal := conf.GetWallet()
					l := len(wal.Key)
					info := fmt.Sprintf("%x\n%x\n%x\n%x\n",
						wal.Key[:l/4], wal.Key[l/4:l/2], wal.Key[l/2:l*3/4], wal.Key[l*3/4:])
					dialog.ShowInformation(res.GetLocalString("Privage Key"), info, win)
				}, win,
			)
		}))

	return widget.NewVBox(form, widget.NewGroup(res.GetLocalString("Change Password"), pwdForm), privBtn)
}

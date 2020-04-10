package screens

import (
	"log"

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

// Master master window
func Master(a fyne.App) fyne.Window {
	w := a.NewWindow(res.GetLocalString("GOVM"))
	w.SetMaster()

	tabs := widget.NewTabContainer(
		widget.NewTabItemWithIcon(res.GetLocalString("Home"), theme.HomeIcon(), widget.NewScrollContainer(AccountScreen(w))),
		widget.NewTabItemWithIcon(res.GetLocalString("Transaction"), theme.MailSendIcon(), widget.NewScrollContainer(TransactionScreen(w))),
		widget.NewTabItemWithIcon(res.GetLocalString("Search"), theme.SearchIcon(), widget.NewScrollContainer(SearchScreen(w))),
		widget.NewTabItemWithIcon(res.GetLocalString("Setting"), theme.SettingsIcon(), widget.NewScrollContainer(SettingScreen(w))))
	tabs.SetTabLocation(widget.TabLocationLeading)

	logo := canvas.NewImageFromResource(res.GetStaticResource("point"))
	logo.SetMinSize(fyne.NewSize(300, 5))
	pwd := widget.NewPasswordEntry()
	// pwd.SetText("govm_pwd@2019")

	btn := widget.NewButton("Ok", func() {
		if pwd.Text == "" {
			dialog.ShowInformation("Error", res.GetLocalString("login.empty_pwd"), w)
			return
		}
		err := conf.Load(pwd.Text)
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
			return
		}
		w.SetContent(tabs)
		min := w.Content().MinSize()
		sizeW := a.Preferences().Int("size_w")
		sizeH := a.Preferences().Int("size_h")
		if sizeW > min.Width && sizeH > min.Height {
			w.Resize(fyne.NewSize(sizeW, sizeH))
		} else {
			w.Resize(min)
		}
		w.CenterOnScreen()
		event.Send(event.EShowHome, w)
		w.SetOnClosed(func() {
			size := w.Canvas().Size()
			a.Preferences().SetInt("size_w", size.Width)
			a.Preferences().SetInt("size_h", size.Height)
		})
	})
	login := widget.NewVBox(
		logo,
		widget.NewLabel(res.GetLocalString("login.msg")),
		pwd,
		widget.NewHBox(layout.NewSpacer(), btn, layout.NewSpacer()),
	)

	w.SetContent(login)

	w.Show()

	event.RegisterConsumer(event.ERequsetPwd, func(e string, param ...interface{}) error {
		// var pwd string
		logo := canvas.NewImageFromResource(res.GetResource("resource/img/point.svg"))
		logo.SetMinSize(fyne.NewSize(300, 5))
		content := widget.NewPasswordEntry()
		var p interface{}
		if len(param) > 0 {
			p = param[0]
		}

		dialog.ShowCustomConfirm(res.GetLocalString("Password"), "", "",
			widget.NewVBox(logo, content), func(rst bool) {
				pwd := content.Text
				log.Println("ERequsetPwd:", rst, pwd)
				if !rst {
					pwd = ""
				}
				event.Send(event.EResponsePwd, p, pwd)
			}, w,
		)
		return nil
	})

	return w
}

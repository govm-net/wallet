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
		widget.NewTabItemWithIcon(res.GetLocalString("Home"), theme.HomeIcon(), AccountScreen(w)),
		widget.NewTabItemWithIcon(res.GetLocalString("Transaction"), theme.MailSendIcon(), TransactionScreen(w)),
		widget.NewTabItemWithIcon(res.GetLocalString("APP"), theme.MenuIcon(), AppScreen(w)),
		widget.NewTabItemWithIcon(res.GetLocalString("Search"), theme.SearchIcon(), SearchScreen(w)),
		widget.NewTabItemWithIcon(res.GetLocalString("History"), theme.ContentPasteIcon(), HistoryScreen(w)),
		widget.NewTabItemWithIcon(res.GetLocalString("Setting"), theme.SettingsIcon(), SettingScreen(w)))
	tabs.SetTabLocation(widget.TabLocationLeading)

	widthLimit := canvas.NewImageFromResource(res.GetResource("point.svg"))
	widthLimit.SetMinSize(fyne.NewSize(300, 1))
	logo := canvas.NewImageFromResource(res.GetResource("govm.png"))
	logo.SetMinSize(fyne.NewSize(100, 100))
	pwd := widget.NewPasswordEntry()

	btn := widget.NewButton("Ok", func() {
		if pwd.Text == "" {
			pwd.Text = "govm_pwd@2019"
		}
		err := conf.Load(pwd.Text)
		if err != nil {
			dialog.ShowInformation("Error", err.Error(), w)
			return
		}
		w.SetContent(tabs)
		w.CenterOnScreen()
		event.Send(event.EShowHome, w)
	})
	login := widget.NewVBox(
		widthLimit,
		widget.NewHBox(layout.NewSpacer(), logo, layout.NewSpacer()),
		widget.NewLabel(res.GetLocalString("login.msg")),
		pwd,
		widget.NewHBox(layout.NewSpacer(), btn, layout.NewSpacer()),
	)

	w.SetContent(login)
	w.CenterOnScreen()

	event.RegisterConsumer(event.ERequsetPwd, func(e string, param ...interface{}) error {
		// var pwd string
		logo := canvas.NewImageFromResource(res.GetResource("point.svg"))
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

	w.ShowAndRun()

	return w
}

// Package main provides various examples of Fyne API capabilities
package main

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"github.com/lengzhao/wallet/bin/gui/event"
	"github.com/lengzhao/wallet/bin/gui/res"
	"github.com/lengzhao/wallet/bin/gui/screens"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(&lumberjack.Logger{
		Filename:   "./log/gui.log",
		MaxSize:    50, // megabytes
		MaxBackups: 5,
		MaxAge:     10,   //days
		Compress:   true, // disabled by default
	})
	res.LoadLanguage()
	a := app.NewWithID("net.govm")
	a.SetIcon(res.GetResourceW("govm.png", theme.CancelIcon()))
	// a.Settings().SetTheme(theme.LightTheme())
	a.Settings().SetTheme(screens.NewCustomTheme())

	event.Send(event.ELogin)
	// screens.Login(a)
	screens.Master(a)
}

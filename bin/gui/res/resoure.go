package res

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"

	"fyne.io/fyne"
	"fyne.io/fyne/theme"
	"github.com/lengzhao/wallet/bin/gui/conf"
)

var localI18n map[string]string

var i18n = map[string]string{
	"GOVM":             "GOVM",
	"login.msg":        "Please entry password",
	"login.empty_pwd":  "not entry password",
	"transfer.peer":    "Peer",
	"transfer.desc":    "After the transfer, wait until the block index of the transaction is not 0 to be considered a success. Normally less than 2 minutes.",
	"move.from_chain":  "From Chain",
	"move.to_chain":    "To Chain",
	"move.desc":        "Move coin to other chain,need wait 10min after BlockID not equal 0.",
	"setting.unit":     "Coin Unit:",
	"setting.language": "Language:",
	"setting.server":   "API Server",
	"setting.desc":     "Restart after setting",
}

// LoadLanguage load language
func LoadLanguage() error {
	lang := conf.Get(conf.Langure)
	dir := path.Join("resource", "i18n", lang)
	data, err := ioutil.ReadFile(path.Join(dir, "local.json"))
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &localI18n)
	if err != nil {
		return err
	}
	ttf := path.Join(dir, "font.ttf")
	if _, err := os.Stat(ttf); !os.IsNotExist(err) {
		os.Setenv("FYNE_FONT", ttf)
		// fmt.Println("set FYNE_FONT:", ttf)
	}

	return nil
}

// GetResource returns a resource by name
func GetResource(path string) fyne.Resource {
	out, err := fyne.LoadResourceFromPath(path)
	if err != nil {
		log.Println("fail to load resource.", path, err)
		return theme.CancelIcon()
	}
	return out
}

var pointIconRes = &fyne.StaticResource{
	StaticName:    "point.svg",
	StaticContent: []byte("<svg><circle r='1'/></svg>")}

// GetStaticResource get static resource
func GetStaticResource(id string) fyne.Resource {
	switch id {
	case "point":
		return pointIconRes
	}
	return theme.CancelIcon()
}

// GetLocalString get local string
func GetLocalString(id string) string {
	out := localI18n[id]
	if out != "" {
		return out
	}
	out = i18n[id]
	if out != "" {
		return out
	}
	return id
}

// GetBaseOfUnit get base of unit
func GetBaseOfUnit(in string) uint64 {
	switch in {
	case "tc":
		return 1000000000000
	case "t9":
		return 1000000000
	case "t6":
		return 1000000
	case "t3":
		return 1000
	case "t0":
		return 1
	}
	return 1000000000000
}

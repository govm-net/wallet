package screens

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/govm-net/govm/wallet"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/res"
	"github.com/lengzhao/wallet/trans"
)

type dataInMode int
type dataOutMode int

const (
	diHex = dataInMode(iota + 1)
	diAddress
	diString
	diUint64
	diUint32
	diUint16
	diUint8
	diFloat64
	diHash
	diMap
)
const (
	doHex = dataOutMode(iota + 100)
	doMap
	doList
	doString
	doNumber
)

var enableEmptyType = map[dataInMode]bool{
	diHash: true, diMap: true,
}

func (a *dataInMode) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		return fmt.Errorf("unknow input mode:%s", s)
	case "hex":
		*a = diHex
	case "address":
		*a = diAddress
	case "string":
		*a = diString
	case "uint64":
		*a = diUint64
	case "uint32":
		*a = diUint32
	case "uint16":
		*a = diUint16
	case "uint8":
		*a = diUint8
	case "float":
		*a = diFloat64
	case "hash":
		*a = diHash
	case "map":
		*a = diMap
	}
	return nil
}

func (a *dataOutMode) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	default:
		return fmt.Errorf("unknow view mode:%s", s)
	case "hex":
		*a = doHex
	case "string":
		*a = doString
	case "number":
		*a = doNumber
	case "map":
		*a = doMap
	case "list":
		*a = doList
	}
	return nil
}

type uiInput struct {
	Mode        dataInMode `json:"mode,omitempty"`
	Title       string     `json:"title,omitempty"`
	Key         string     `json:"key,omitempty"`
	Value       string     `json:"value,omitempty"`
	Hide        bool       `json:"hide,omitempty"`
	MultiLine   bool       `json:"multi_line,omitempty"`
	EmptyEnable bool       `json:"empty_enable,omitempty"`
	Sub         []uiInput  `json:"sub,omitempty"`
	win         *widget.Entry
}

type viewOutput struct {
	Mode      dataOutMode  `json:"mode,omitempty"`
	Title     string       `json:"title,omitempty"`
	Length    int          `json:"length,omitempty"`
	Key       string       `json:"key,omitempty"`
	MultiLine bool         `json:"multi_line,omitempty"`
	Sub       []viewOutput `json:"sub,omitempty"`
	win       *widget.Entry
}

type dReadUI struct {
	Name           string       `json:"name,omitempty"`
	App            string       `json:"app,omitempty"`
	Description    string       `json:"description,omitempty"`
	Struct         string       `json:"struct,omitempty"`
	IsLog          bool         `json:"is_log,omitempty"`
	Chain          uint64       `json:"chain,omitempty"`
	ShowExpiration bool         `json:"show_expiration,omitempty"`
	ShowKey        bool         `json:"show_key,omitempty"`
	Hide           bool         `json:"hide,omitempty"`
	Input          []uiInput    `json:"input,omitempty"`
	View           []viewOutput `json:"view,omitempty"`
	winExpiration  *widget.Entry
	winKey         *widget.Entry
}

func addInputItem(form *widget.Form, lst []uiInput) {
	for i, it := range lst {
		if it.Title != "" {
			ent := widget.NewEntry()
			if it.MultiLine {
				ent = widget.NewMultiLineEntry()
			}
			if it.Value != "" {
				ent.SetText(it.Value)
			} else if it.Mode == diAddress {
				ent.SetText(conf.GetWallet().AddressStr)
			}
			it.win = ent
			lst[i] = it
			if it.Hide {
				continue
			}
			form.Append(it.Title, ent)
		}
		if len(it.Sub) > 0 {
			addInputItem(form, it.Sub)
		}
	}
}

func addViewItem(form *widget.Form, lst []viewOutput) {
	for i, it := range lst {
		if it.Title != "" {
			ent := widget.NewEntry()
			if it.MultiLine {
				ent = widget.NewMultiLineEntry()
			}
			ent.Disable()
			form.Append(it.Title, ent)
			it.win = ent
			lst[i] = it
		}
		if len(it.Sub) > 0 {
			addViewItem(form, it.Sub)
		}
	}
}
func showViewItemFromMap(lst []viewOutput, data map[string]interface{}) error {
	for _, it := range lst {
		switch it.Mode {
		case doHex:
			v := data[it.Key]
			var in []byte
			if v != nil {
				d, _ := json.Marshal(v)
				err := json.Unmarshal(d, &in)
				if err != nil {
					return err
				}
			}

			if it.win != nil {
				it.win.SetText(hex.EncodeToString(in))
			}
			if len(it.Sub) > 0 {
				err := showViewItem(it.Sub, in)
				if err != nil {
					return err
				}
			}
		case doMap:
			v := data[it.Key]
			var in map[string]interface{}
			if v != nil {
				d, _ := json.Marshal(v)
				err := json.Unmarshal(d, &in)
				if err != nil {
					return err
				}
			}

			if it.win != nil {
				d, _ := json.MarshalIndent(in, "", "  ")
				it.win.SetText(string(d))
			}
			if len(it.Sub) > 0 {
				err := showViewItemFromMap(it.Sub, in)
				if err != nil {
					return err
				}
			}
		case doList:
			v := data[it.Key]
			var in []interface{}
			if v != nil {
				d, _ := json.Marshal(v)
				err := json.Unmarshal(d, &in)
				if err != nil {
					return err
				}
			}

			if it.win != nil {
				d, _ := json.MarshalIndent(in, "", "  ")
				it.win.SetText(string(d))
			}
			if len(it.Sub) > 0 {
				return fmt.Errorf("not support sub of list")
			}
		case doString:
			v := data[it.Key]
			var in string
			if v != nil {
				d, _ := json.Marshal(v)
				err := json.Unmarshal(d, &in)
				if err != nil {
					return err
				}
			}
			if it.win != nil {
				it.win.SetText(in)
			}
			if len(it.Sub) > 0 {
				err := showViewItem(it.Sub, []byte(in))
				if err != nil {
					return err
				}
			}
		case doNumber:
			v := data[it.Key]
			var in int64
			if v != nil {
				d, _ := json.Marshal(v)
				err := json.Unmarshal(d, &in)
				if err != nil {
					return err
				}
			}
			if it.win != nil {
				it.win.SetText(fmt.Sprintf("%d", in))
			}
			if len(it.Sub) > 0 {
				return fmt.Errorf("not support")
			}
		}
	}
	return nil
}

func showViewItem(lst []viewOutput, data []byte) error {
	for _, it := range lst {
		switch it.Mode {
		case doHex:
			ld := data
			if it.Length > 0 {
				ld = data[:it.Length]
				data = data[it.Length:]
			}
			if it.win != nil {
				it.win.SetText(hex.EncodeToString(ld))
			}
			if len(it.Sub) > 0 {
				err := showViewItem(it.Sub, ld)
				if err != nil {
					return err
				}
			}
		case doMap:
			ld := data
			if it.Length > 0 {
				ld = data[:it.Length]
				data = data[it.Length:]
			}
			var info map[string]interface{}
			err := json.Unmarshal(ld, &info)
			if err != nil {
				return err
			}
			if it.win != nil {
				d, _ := json.MarshalIndent(info, "", "  ")
				it.win.SetText(string(d))
			}
			if len(it.Sub) > 0 {
				err = showViewItemFromMap(it.Sub, info)
				if err != nil {
					return err
				}
			}
		case doList:
			ld := data
			if it.Length > 0 {
				ld = data[:it.Length]
				data = data[it.Length:]
			}
			var info []interface{}
			err := json.Unmarshal(ld, &info)
			if err != nil {
				return err
			}
			if it.win != nil {
				d, _ := json.MarshalIndent(info, "", "  ")
				it.win.SetText(string(d))
			}
			if len(it.Sub) > 0 {
				return fmt.Errorf("not support")
			}
		case doString:
			ld := data
			if it.Length > 0 {
				ld = data[:it.Length]
				data = data[it.Length:]
			}
			if it.win != nil {
				it.win.SetText(string(ld))
			}
			if len(it.Sub) > 0 {
				err := showViewItem(it.Sub, ld)
				if err != nil {
					return err
				}
			}
		case doNumber:
			ld := data
			if it.Length > 0 {
				ld = data[:it.Length]
				data = data[it.Length:]
			}
			if it.win != nil {
				var num big.Int
				num.SetBytes(ld)
				it.win.SetText(num.String())
			}
			if len(it.Sub) > 0 {
				err := showViewItem(it.Sub, ld)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func newReadUI(w fyne.Window, in dReadUI) fyne.Widget {
	inputForm := &widget.Form{}
	viewForm := &widget.Form{}
	var lastKey string
	chainStr := fmt.Sprintf("%d", in.Chain)
	if in.ShowKey {
		in.winKey = widget.NewEntry()
		viewForm.Append(res.GetLocalString("Key:"), in.winKey)
	}
	if in.ShowExpiration {
		in.winExpiration = widget.NewEntry()
		viewForm.Append(res.GetLocalString("Expiration:"), in.winExpiration)
	}
	if in.Description != "" {
		inputForm.Append(res.GetLocalString("Description"),
			widget.NewLabel(in.Description))
	}

	addInputItem(inputForm, in.Input)

	empty := widget.NewLabel("")
	viewData := func() {
		// var localKey []byte
		localKey, err := encodeRunInput(in.Input)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		if len(localKey) == 0 {
			dialog.ShowError(fmt.Errorf("request key"), w)
			return
		}
		k := hex.EncodeToString(localKey)
		data, life := getStringOfDB(chainStr, in.App, in.Struct, k)
		if life == 0 {
			dialog.ShowError(fmt.Errorf("not found"), w)
			return
		}

		log.Println("key:", localKey, data, life)
		if in.ShowExpiration {
			t := time.Unix(int64(life/1000), 0)
			in.winExpiration.SetText(t.Local().String())
		}
		if in.ShowKey {
			in.winKey.SetText(k)
		}

		hexData, _ := hex.DecodeString(data)
		showViewItem(in.View, hexData)
		lastKey = k
	}
	btnSearch := widget.NewButton(res.GetLocalString("Search"), viewData)
	btnNext := widget.NewButton(res.GetLocalString("Next"), func() {
		key := getNextKeyOfDB(chainStr, in.App, in.Struct, lastKey)
		lastKey = key
		if in.ShowKey {
			in.winKey.SetText(key)
		}
		data, life := getStringOfDB(chainStr, in.App, in.Struct, key)
		if life == 0 {
			dialog.ShowError(fmt.Errorf("not found"), w)
			return
		}

		log.Println("key:", key, data, life)
		if in.ShowExpiration {
			t := time.Unix(int64(life/1000), 0)
			in.winExpiration.SetText(t.Local().String())
		}
		hexData, _ := hex.DecodeString(data)
		showViewItem(in.View, hexData)
	})
	btns := widget.NewHBox(btnNext, empty, btnSearch)
	searchItem := newFormItemWithObj(empty, btns)

	addViewItem(viewForm, in.View)

	return widget.NewVBox(inputForm, searchItem, viewForm)
}

type appRunUI struct {
	Name        string    `json:"name,omitempty"`
	App         string    `json:"app,omitempty"`
	Description string    `json:"description,omitempty"`
	Chains      []uint64  `json:"chains,omitempty"`
	Hide        bool      `json:"hide,omitempty"`
	Cost        float64   `json:"cost,omitempty"`
	Energy      float64   `json:"energy,omitempty"`
	Input       []uiInput `json:"input,omitempty"`
}

type uiConf struct {
	ViewUI []dReadUI  `json:"view_ui,omitempty"`
	RunUI  []appRunUI `json:"run_ui,omitempty"`
}

func encodeRunInputMap(in []uiInput, param map[string]interface{}) error {
	for _, it := range in {
		var data string
		if it.win != nil {
			data = it.win.Text
		}
		switch it.Mode {
		case diHex:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return fmt.Errorf("empty value:%s", it.Title)
			}
			d, err := hex.DecodeString(data)
			if err != nil {
				return err
			}
			param[it.Key] = d
		case diAddress:
			fallthrough
		case diString:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return fmt.Errorf("empty value:%s", it.Title)
			}
			param[it.Key] = data
		case diFloat64:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseFloat(data, 64)
			if err != nil {
				return err
			}
			param[it.Key] = number
		case diUint64:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseUint(data, 10, 64)
			if err != nil {
				return err
			}
			param[it.Key] = number
		case diUint32:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseUint(data, 10, 32)
			if err != nil {
				return err
			}
			param[it.Key] = number
		case diUint16:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseUint(data, 10, 16)
			if err != nil {
				return err
			}
			param[it.Key] = number
		case diUint8:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseUint(data, 10, 8)
			if err != nil {
				return err
			}
			param[it.Key] = number
		case diMap:
			info := make(map[string]interface{})
			err := encodeRunInputMap(it.Sub, info)
			if err != nil {
				return err
			}
			param[it.Key] = info
		default:
			return fmt.Errorf("unsupport:%d", it.Mode)
		}
	}
	return nil
}

func encodeRunInput(in []uiInput) ([]byte, error) {
	var out []byte
	for _, it := range in {
		var lv []byte
		var data string
		if it.win != nil {
			data = it.win.Text
		}
		switch it.Mode {
		case diAddress:
			fallthrough
		case diHex:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return nil, fmt.Errorf("empty value:%s", it.Title)
			}
			d, err := hex.DecodeString(data)
			if err != nil {
				return nil, err
			}
			lv = d
		case diString:
			lv = []byte(data)
		case diFloat64:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return nil, fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseFloat(data, 64)
			if err != nil {
				return nil, err
			}
			lv = trans.Encode(number)
		case diUint64:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return nil, fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseUint(data, 10, 64)
			if err != nil {
				return nil, err
			}
			lv = trans.Encode(number)
		case diUint32:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return nil, fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseUint(data, 10, 32)
			if err != nil {
				return nil, err
			}
			lv = trans.Encode(uint32(number))
		case diUint16:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return nil, fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseUint(data, 10, 16)
			if err != nil {
				return nil, err
			}
			lv = trans.Encode(uint16(number))
		case diUint8:
			if data == "" {
				if it.EmptyEnable || enableEmptyType[it.Mode] {
					continue
				}
				return nil, fmt.Errorf("empty value:%s", it.Title)
			}
			number, err := strconv.ParseUint(data, 10, 8)
			if err != nil {
				return nil, err
			}
			lv = trans.Encode(uint8(number))
		case diMap:
			info := make(map[string]interface{})
			err := encodeRunInputMap(it.Sub, info)
			if err != nil {
				return nil, err
			}
			lv, _ = json.Marshal(info)
		case diHash:
			lv = wallet.GetHash(out)
			out = nil
		default:
			return nil, fmt.Errorf("unsupport:%d", it.Mode)
		}
		out = append(out, lv...)
	}
	return out, nil
}

func newRunUI(w fyne.Window, in appRunUI) fyne.Widget {
	infoForm := &widget.Form{}
	inputForm := &widget.Form{}
	result := widget.NewEntry()
	result.Disable()

	// chainStr := fmt.Sprintf("%d", in.Chain)
	if len(in.Chains) == 0 {
		in.Chains = append(in.Chains, 1)
	}
	var chainArr []string
	for _, c := range in.Chains {
		chainArr = append(chainArr, fmt.Sprintf("%d", c))
	}
	chainID := in.Chains[0]
	chain := widget.NewSelect(chainArr, func(s string) {
		chainID, _ = strconv.ParseUint(s, 10, 64)
	})
	chain.SetSelected(chainArr[0])

	amount := widget.NewEntry()
	amount.SetText("0")
	eng := widget.NewEntry()
	eng.SetText("0")
	infoForm.Append(res.GetLocalString("Description"),
		widget.NewLabel(in.Description))
	infoForm.Append(res.GetLocalString("Chain"), chain)

	if in.Cost > 0 {
		amount.SetText(fmt.Sprintf("%f", in.Cost))
		infoForm.Append(res.GetLocalString("Cost"),
			newFormItemWithUnit(amount, "govm"))
	}
	if in.Energy > 0 {
		eng.SetText(fmt.Sprintf("%f", in.Energy))
		infoForm.Append(res.GetLocalString("Energy"),
			newFormItemWithUnit(eng, "govm"))
	}
	addInputItem(inputForm, in.Input)

	inputForm.SubmitText = res.GetLocalString("Run")
	inputForm.OnSubmit = func() {
		result.SetText("")
		var cost, energy uint64
		if amount.Text != "" {
			costF, err := strconv.ParseFloat(amount.Text, 10)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error cost"), w)
				return
			}
			amount.SetPlaceHolder(amount.Text)
			amount.SetText("")
			cost = uint64(costF * 1e9)
		}
		if eng.Text != "" {
			engF, err := strconv.ParseFloat(eng.Text, 10)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error cost"), w)
				return
			}
			energy = uint64(engF * 1e9)
		}

		data, err := encodeRunInput(in.Input)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		log.Printf("run. chain:%d,app:%s,msg:%x\n", chainID, in.App, data)
		key, err := runApp(chainID, cost, energy, in.App, "", data)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		result.SetText(key)
	}
	inputForm.Refresh()

	return widget.NewVBox(infoForm, inputForm, result)
}

// CustomizeScreen Customize dApp ui
func CustomizeScreen(w fyne.Window) fyne.CanvasObject {
	c := conf.Get()
	tb := widget.NewTabContainer()
	d, _ := ioutil.ReadFile(c.DynamicUIFile)
	var info uiConf
	err := json.Unmarshal(d, &info)
	if err != nil {
		// fmt.Println("dReadUI:", err)
		tb.Append(widget.NewTabItem("error", widget.NewLabel(fmt.Sprint(err))))
		tb.SelectTabIndex(0)
		return tb
	}
	for _, it := range info.RunUI {
		if it.Hide {
			continue
		}
		tb.Append(widget.NewTabItem("R_"+it.Name, newRunUI(w, it)))
	}
	for _, it := range info.ViewUI {
		if it.Hide {
			continue
		}
		tb.Append(widget.NewTabItem("V_"+it.Name, newReadUI(w, it)))
	}
	tb.SelectTabIndex(0)
	return tb
}

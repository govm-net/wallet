package screens

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/lengzhao/wallet/bin/gui/conf"
	"github.com/lengzhao/wallet/bin/gui/res"
	"github.com/lengzhao/wallet/trans"
)

var taskApp = "c11b3b8aa630a7fbfccec9e023c363749e9c60db43d7678f43c96075e5c2ddc0"

type taskHead struct {
	ID          uint64 `json:"id,omitempty"`
	Title       string `json:"title,omitempty"`
	Life        uint64 `json:"life,omitempty"`
	Description string `json:"desc,omitempty"`
	Number      uint32 `json:"number,omitempty"`
}

// Task task
type taskData struct {
	taskHead
	// AcceptRule  trans.Hash             `json:"accept_rule,omitempty"`
	// CommitRule  trans.Hash             `json:"commit_rule,omitempty"`
	// RewardRule  trans.Hash             `json:"reward_rule,omitempty"`
	// Others      map[string]interface{} `json:"others,omitempty"`
	Owner      trans.Address `json:"owner,omitempty"`
	Bounty     uint64        `json:"bounty,omitempty"`
	AcceptNum  uint32        `json:"accept_num,omitempty"`
	Status     int           `json:"status,omitempty"`
	Rewarded   uint64        `json:"rewarded,omitempty"`
	Message    string        `json:"message,omitempty"`
	Expiration int64         `json:"expiration,omitempty"`
}

type actionHead struct {
	ID     uint64 `json:"id,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Reward uint64 `json:"reward,omitempty"`
}
type actionData struct {
	actionHead
	TaskID  uint64        `json:"task_id,omitempty"`
	User    trans.Address `json:"user,omitempty"`
	Index   uint32        `json:"index,omitempty"`
	Status  int           `json:"status,omitempty"`
	Message string
}

func statusToString(in int) string {
	switch in {
	case 0:
		return "init"
	case 1:
		return "accept"
	case 2:
		return "commit"
	case 3:
		return "refuse"
	case 4:
		return "closed"
	default:
		return ""
	}
}

func runApp(chain, cost, energy uint64, app, prefix string, param []byte) (string, error) {
	myWlt := conf.GetWallet()
	trans := trans.NewTransaction(chain, myWlt.Address, cost)
	var body []byte
	var err error
	if prefix != "" {
		body, err = hex.DecodeString(prefix)
		if err != nil {
			return "", err
		}
	}
	if len(param) > 0 {
		body = append(body, param...)
	}
	trans.Energy = energy
	err = trans.RunApp(app, body)
	if err != nil {
		return "", err
	}
	td := trans.GetSignData()
	sign := myWlt.Sign(td)
	trans.SetTheSign(sign)
	td = trans.Output()
	key := trans.Key[:]

	err = postTrans(chain, td)
	if err != nil {
		// result.SetText(fmt.Sprintf("%s", err))
		log.Println("fail to run app:", err)
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func makeTaskInfoTab(w fyne.Window) fyne.Widget {
	c := conf.Get()
	chain := widget.NewSelect(c.Chains, nil)
	chain.SetSelected("1")

	number := widget.NewEntry()
	number.Disable()
	number.SetText("0")
	updateTime := widget.NewEntry()
	updateTime.Disable()
	index := widget.NewEntry()

	showForm := &widget.Form{}
	taskID := widget.NewLabel("")
	title := widget.NewLabel("")
	lifeLab := widget.NewLabel("")
	bounty := widget.NewLabel("")
	description := widget.NewMultiLineEntry()
	description.SetReadOnly(true)
	description.SetText("")
	owner := widget.NewLabel("")
	accepted := widget.NewLabel("")
	rewarded := widget.NewLabel("")
	status := widget.NewLabel("")

	showForm.Append(res.GetLocalString("Task ID"), taskID)
	showForm.Append(res.GetLocalString("Title"), title)
	showForm.Append(res.GetLocalString("Time"), lifeLab)
	showForm.Append(res.GetLocalString("Bounty"), bounty)
	showForm.Append(res.GetLocalString("Owner"), owner)
	showForm.Append(res.GetLocalString("Accepted"), accepted)
	showForm.Append(res.GetLocalString("Rewarded"), rewarded)
	showForm.Append(res.GetLocalString("Status"), status)
	showForm.Append(res.GetLocalString("Description"), description)
	showForm.SubmitText = "Accept"
	showForm.CancelText = "View Actions"
	result := widget.NewEntry()
	result.Disable()

	updateEvent := func() {
		num := getIntOfDB(chain.Selected, taskApp, "tApp", "0000")
		number.SetText(fmt.Sprintf("%d", num))
		updateTime.SetText(time.Now().Local().String())
		if index.Text == "" {
			index.SetText(number.Text)
		}
	}

	clickEvent := func() {
		if index.Text == "" {
			updateEvent()
		}
		result.SetText("")
		id, err := strconv.ParseUint(index.Text, 10, 64)
		if err != nil || id == 0 {
			result.SetText("error task id")
			return
		}
		key := fmt.Sprintf("%016x", id)
		taskStr, life := getStringOfDB(chain.Selected, taskApp, "tTask", key)
		if life == 0 || taskStr == "" {
			result.SetText("not found the task")
			return
		}
		var info taskData
		d, _ := hex.DecodeString(taskStr)
		json.Unmarshal(d, &info)
		info.Expiration = int64(life / 1000)

		taskID.SetText(index.Text)
		title.SetText(info.Title)
		lifeLab.SetText(time.Unix(info.Expiration, 0).String())
		bounty.SetText(fmt.Sprintf("%.3f govm", float64(info.Bounty)/1e9))

		description.SetText(info.Description)
		owner.SetText(info.Owner.String())
		accepted.SetText(fmt.Sprintf("%d", info.AcceptNum))
		rewarded.SetText(fmt.Sprintf("%.3f govm", float64(info.Rewarded)/1e9))
		status.SetText(statusToString(info.Status))
		showForm.OnSubmit = nil
		showForm.OnCancel = nil
		log.Println("task:", string(d), info)
		if time.Now().Unix() < info.Expiration &&
			info.Bounty > info.Rewarded && info.Number > info.AcceptNum {
			statKey := fmt.Sprintf("%s%s", conf.GetWallet().AddressStr, key)
			n := getIntOfDB(chain.Selected, taskApp, "tActionStatus", statKey)
			log.Println("accepted:", n)
			if n == 0 {
				showForm.OnSubmit = func() {
					cid, err := strconv.ParseUint(chain.Selected, 10, 64)
					if err != nil {
						dialog.ShowError(fmt.Errorf("error chain id"), w)
						return
					}

					trans, err := runApp(cid, 0, 1e8, taskApp, "02"+key, nil)
					if err != nil {
						dialog.ShowError(err, w)
						return
					}
					log.Println("accept task:", id)
					result.SetText(trans)
				}

			}
		}
		// if info.AcceptNum > 0 {
		// 	showForm.OnCancel = func() {
		// 		fmt.Println("view action")
		// 		taskTB.SelectTabIndex(1)
		// 	}
		// }
		showForm.Refresh()
	}
	previousEvent := func() {
		num := getIntOfDB(chain.Selected, taskApp, "tApp", "0000")
		id, err := strconv.ParseUint(index.Text, 10, 64)
		if err != nil || id <= 1 {
			id = num + 1
		}
		id--
		index.SetText(fmt.Sprintf("%d", id))
		clickEvent()
	}
	btn := widget.NewButton(res.GetLocalString("Update"), updateEvent)
	layout2 := layout.NewBorderLayout(nil, nil, nil, btn)
	updateItem := fyne.NewContainerWithLayout(layout2, btn, number)

	btnSearch := widget.NewButton(res.GetLocalString("Go"), clickEvent)
	btnPrev := widget.NewButton(res.GetLocalString("Previous"), previousEvent)
	btns := widget.NewHBox(btnPrev, btnSearch)
	layout3 := layout.NewBorderLayout(nil, nil, nil, btns)
	searchItem := fyne.NewContainerWithLayout(layout3, btns, index)

	form := &widget.Form{}
	form.Append(res.GetLocalString("Task Number"), updateItem)
	form.Append(res.GetLocalString("UpdateTime"), updateTime)
	form.Append(res.GetLocalString("Search"), searchItem)
	go clickEvent()

	return widget.NewVBox(form, widget.NewGroup(
		res.GetLocalString("Info"), showForm), result)
}

func makeTaskTab(w fyne.Window) fyne.Widget {
	c := conf.Get()
	chain := widget.NewSelect(c.Chains, nil)
	chain.SetSelected("1")
	title := widget.NewEntry()
	title.SetPlaceHolder("task title")

	desc := widget.NewMultiLineEntry()
	bounty := widget.NewEntry()
	bounty.SetText("100")
	unit := widget.NewLabel("govm")
	energy := widget.NewEntry()
	energy.SetText("0.02")
	unit2 := widget.NewLabel("govm")
	lifeE := widget.NewEntry()
	lifeE.SetText("30")
	unit3 := widget.NewLabel("Days")
	result := widget.NewEntry()
	result.Disable()

	form := &widget.Form{
		OnCancel: func() {
			title.SetText("")
			desc.SetText("")
			bounty.SetText("")
			result.SetText("")
		},
		OnSubmit: func() {
			result.SetText("")
			costF, err := strconv.ParseFloat(bounty.Text, 10)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error bounty"), w)
				return
			}
			energyF, err := strconv.ParseFloat(energy.Text, 10)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error bounty"), w)
				return
			}
			bounty.SetPlaceHolder(bounty.Text)
			bounty.SetText("")
			cid, err := strconv.ParseUint(chain.Selected, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error chain id"), w)
				return
			}
			life, err := strconv.ParseUint(lifeE.Text, 10, 64)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error expiration"), w)
				return
			}
			if title.Text == "" {
				dialog.ShowError(fmt.Errorf("request title"), w)
				return
			}
			if costF < 1.0 {
				dialog.ShowError(fmt.Errorf("not enough bounty"), w)
				return
			}

			var task taskHead
			task.Title = title.Text
			task.Description = desc.Text
			// task.Expiration
			task.Life = life * 24 * 3600 * 1000
			d, _ := json.Marshal(task)
			title.SetPlaceHolder(title.Text)
			title.SetText("")
			desc.SetPlaceHolder(desc.Text)
			desc.SetText("")

			cost := uint64(costF * float64(1e9))
			eng := uint64(energyF * float64(1e9))
			key, err := runApp(cid, cost, eng, taskApp, "01", d)
			if err != nil {
				dialog.ShowError(err, w)
				log.Println("fail to create transaction(run app):", err)
				return
			}

			log.Printf("new transfer:%x\n", key)
			// dialog.ShowInformation("transaction", fmt.Sprintf("%x", key), w)
			result.SetText(key)
		},
	}
	// form.Append(res.GetLocalString("Chain"), chain)
	form.Append(res.GetLocalString("Title"), title)
	borderLayout3 := layout.NewBorderLayout(nil, nil, nil, unit3)
	form.Append(res.GetLocalString("Life"), fyne.NewContainerWithLayout(borderLayout3, unit3, lifeE))
	form.Append(res.GetLocalString("Description"), desc)
	borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
	form.Append(res.GetLocalString("Bounty"), fyne.NewContainerWithLayout(borderLayout, unit, bounty))
	borderLayout2 := layout.NewBorderLayout(nil, nil, nil, unit2)
	form.Append(res.GetLocalString("Energy"), fyne.NewContainerWithLayout(borderLayout2, unit2, energy))

	return widget.NewVBox(form, result)
}

func makeActionTab(w fyne.Window) fyne.Widget {
	c := conf.Get()
	chain := widget.NewSelect(c.Chains, nil)
	chain.SetSelected("1")

	taskIndex := widget.NewEntry()
	taskIndex.SetText("1")
	offset := widget.NewEntry()
	offset.SetText("1")

	showForm := &widget.Form{}
	taskID := widget.NewLabel("")
	title := widget.NewLabel("")
	bounty := widget.NewLabel("")
	actionID := widget.NewLabel("")
	actionUser := widget.NewLabel("")
	actionStatus := widget.NewLabel("")
	actionReward := widget.NewLabel("")
	actionMsg := widget.NewMultiLineEntry()
	actionMsg.Disable()

	showForm.Append(res.GetLocalString("Task ID"), taskID)
	showForm.Append(res.GetLocalString("Title"), title)
	showForm.Append(res.GetLocalString("Bounty"), bounty)
	showForm.Append(res.GetLocalString("ID"), actionID)
	showForm.Append(res.GetLocalString("User"), actionUser)
	showForm.Append(res.GetLocalString("Status"), actionStatus)
	showForm.Append(res.GetLocalString("Reward"), actionReward)
	showForm.Append(res.GetLocalString("Message"), actionMsg)
	showForm.SubmitText = "Accept"
	showForm.CancelText = "View Actions"

	result := widget.NewEntry()
	result.Disable()

	var info taskData

	clickEvent := func() {
		result.SetText("")
		showForm.Hide()
		id, err := strconv.ParseUint(taskIndex.Text, 10, 64)
		if err != nil || id == 0 {
			result.SetText("error task id")
			return
		}
		ofs, err := strconv.ParseUint(offset.Text, 10, 64)
		if err != nil || ofs == 0 {
			result.SetText("error action offset")
			return
		}
		if info.ID != id || uint32(ofs) > info.AcceptNum {
			info = taskData{}
			key := fmt.Sprintf("%016x", id)
			taskStr, life := getStringOfDB(chain.Selected, taskApp, "tTask", key)
			if life == 0 || taskStr == "" {
				result.SetText("not found the task")
				return
			}
			d, _ := hex.DecodeString(taskStr)
			json.Unmarshal(d, &info)
			info.Expiration = int64(life / 1000)
			info.ID = id
		}
		if uint32(ofs) > info.AcceptNum {
			ofs = 1
			offset.SetText(fmt.Sprintf("%d", ofs))
		}
		taskID.SetText(fmt.Sprintf("%d", info.ID))
		title.SetText(info.Title)
		bounty.SetText(fmt.Sprintf("%.3f govm", float64(info.Bounty)/1e9))
		showForm.OnSubmit = nil
		key := fmt.Sprintf("%016x%08x", id, ofs)
		aid := getIntOfDB(chain.Selected, taskApp, "tApp", key)
		if aid == 0 {
			result.SetText("not found the action")
			return
		}
		key1 := fmt.Sprintf("%016x", aid)
		actionStr, life := getStringOfDB(chain.Selected, taskApp, "tAction", key1)
		if life == 0 || actionStr == "" {
			result.SetText("not found the action")
			return
		}
		var action actionData
		d, _ := hex.DecodeString(actionStr)
		json.Unmarshal(d, &action)
		action.ID = aid

		actionID.SetText(fmt.Sprintf("%d", aid))
		actionUser.SetText(action.User.String())
		actionStatus.SetText(statusToString(action.Status))
		actionReward.SetText(fmt.Sprintf("%.3f govm", float64(action.Reward)/1e9))
		actionMsg.SetText(action.Message)

		if time.Now().Unix() < info.Expiration {
			if action.Status == 1 &&
				conf.GetWallet().AddressStr == action.User.String() {
				showForm.SubmitText = "Commit"
				showForm.OnSubmit = func() {
					logo := canvas.NewImageFromResource(res.GetResource("point.svg"))
					logo.SetMinSize(fyne.NewSize(300, 5))
					title := widget.NewLabel(res.GetLocalString("Proof"))
					content := widget.NewMultiLineEntry()
					dialog.ShowCustomConfirm(res.GetLocalString("Commit Action"), "", "",
						widget.NewVBox(logo, title, content), func(rst bool) {
							if !rst {
								fmt.Println("cancel content.Text")
								return
							}
							// fmt.Println("content.Text:", content.Text)
							var newAction actionHead
							newAction.ID = aid
							newAction.Msg = content.Text
							cid, _ := strconv.ParseUint(chain.Selected, 10, 64)
							d1, _ := json.Marshal(newAction)
							transKey, err := runApp(cid, 0, 0, taskApp, "03", d1)
							if err != nil {
								// result.SetText("error:" + err.Error())
								dialog.ShowError(err, w)
								return
							}
							// runApp
							result.SetText(transKey)
						}, w,
					)
				}
			} else if action.Status == 2 &&
				conf.GetWallet().AddressStr == info.Owner.String() {
				showForm.SubmitText = "Reward"
				showForm.OnSubmit = func() {
					logo := canvas.NewImageFromResource(res.GetResource("point.svg"))
					logo.SetMinSize(fyne.NewSize(300, 5))
					title := widget.NewLabel(res.GetLocalString("Reward"))
					content := widget.NewEntry()
					v := float64(info.Bounty-info.Rewarded) / 1e9
					fmt.Printf("Bounty:%d,Rewarded:%d,v:%f\n",
						info.Bounty, info.Rewarded, v)
					content.SetText(fmt.Sprintf("%.3f", v))
					unit := widget.NewLabel("govm")
					borderLayout := layout.NewBorderLayout(nil, nil, nil, unit)
					dialog.ShowCustomConfirm(res.GetLocalString("Reward Action"), "", "",
						widget.NewVBox(logo, title,
							fyne.NewContainerWithLayout(borderLayout, unit, content)), func(rst bool) {
							if !rst {
								fmt.Println("cancel content.Text")
								return
							}
							costF, err := strconv.ParseFloat(content.Text, 10)
							if err != nil {
								dialog.ShowError(fmt.Errorf("error reward"), w)
								return
							}
							var newAction actionHead
							newAction.ID = aid
							newAction.Reward = uint64(costF * float64(1e9))
							cid, _ := strconv.ParseUint(chain.Selected, 10, 64)
							d1, _ := json.Marshal(newAction)
							transKey, err := runApp(cid, 0, 1e8, taskApp, "04", d1)
							if err != nil {
								// result.SetText("error:" + err.Error())
								dialog.ShowError(err, w)
								return
							}
							result.SetText(transKey)
						}, w,
					)
				}
			}
		}
		showForm.Show()
		showForm.Refresh()
	}
	NextEvent := func() {
		id, _ := strconv.ParseUint(offset.Text, 10, 64)
		id++
		offset.SetText(fmt.Sprintf("%d", id))
		clickEvent()
	}

	btnSearch := widget.NewButton(res.GetLocalString("Go"), clickEvent)
	btnPrev := widget.NewButton(res.GetLocalString("Next"), NextEvent)
	btns := widget.NewHBox(btnPrev, btnSearch)
	layout3 := layout.NewBorderLayout(nil, nil, nil, btns)
	searchItem := fyne.NewContainerWithLayout(layout3, btns, offset)

	form := &widget.Form{}
	form.Append(res.GetLocalString("Task ID"), taskIndex)
	form.Append(res.GetLocalString("Offset"), searchItem)
	go clickEvent()

	return widget.NewVBox(form,
		widget.NewGroup(res.GetLocalString("Action Info"), showForm), result)
}

// TaskScreen shows a panel containing widget demos
func TaskScreen(w fyne.Window) fyne.CanvasObject {
	return widget.NewTabContainer(
		widget.NewTabItem(res.GetLocalString("Task Info"), makeTaskInfoTab(w)),
		widget.NewTabItem(res.GetLocalString("Action Info"), makeActionTab(w)),
		widget.NewTabItem(res.GetLocalString("New Task"), makeTaskTab(w)),
	)
}

package panel

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/function"
)

// StartPlayPanel ..
func (p *Panel) StartPlayPanel() (exitGame bool, err error) {
	for {
		var msg *function.MessageFromBacked

		minecraftForm := form.ActionForm{
			Title: "游玩租赁服",
			Buttons: []form.ActionFormElement{
				{
					Text: "按上次的配置进服",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "选择租赁服进行游玩",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "添加新租赁服配置",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "修改已有租赁服配置",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "删除已有租赁服配置",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "返回上一级菜单",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		if len(p.f.EulogistUserData().RentalServerConfig) == 0 {
			minecraftForm.Content = "您目前没有添加任何租赁服, 请点击§r§e添加新租赁服配置§r以添加。\n请选择您要进行的操作。"
		} else {
			minecraftForm.Content = fmt.Sprintf(
				"您目前已存有 §r§b%d§r 个租赁服配置。\n请选择您要进行的操作。",
				len(p.f.EulogistUserData().RentalServerConfig),
			)
		}

		resp, isUserCancel, err := p.f.Interact().SendFormAndWaitResponse(minecraftForm)
		if err != nil {
			return false, fmt.Errorf("StartPlayPanel: %v", err)
		}
		if isUserCancel {
			return false, nil
		}

		idx := resp.(int32)
		switch idx {
		case 0:
			msg, err = p.f.StartPlayByLastConfig()
		case 1:
			msg, err = p.f.SelectRentalServerAndStartPlay()
		case 2:
			err = p.f.AddNewRentalServerConfig()
		case 3:
			err = p.f.ModifyRentalServerConfig()
		case 4:
			err = p.f.RemoveRentalServerConfig()
		case 5:
			return false, nil
		}
		if err != nil {
			return false, fmt.Errorf("StartPlayPanel: %v", err)
		}

		if msg != nil {
			if msg.CanTerminate {
				return true, nil
			}
			if !msg.LoginServerMeetError {
				continue
			}
			_, _, err := p.f.Interact().SendFormAndWaitResponse(form.MessageForm{
				Title:   "进服失败",
				Content: msg.LoginServrErrorInfo,
				Button1: "确定",
				Button2: "取消",
			})
			if err != nil {
				return false, fmt.Errorf("StartPlayPanel: %v", err)
			}
		}
	}
}

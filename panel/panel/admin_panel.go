package panel

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// AdminPanel ..
func (p *Panel) AdminPanel() error {
	switch p.f.EulogistUserData().UserPermissionLevel {
	case define.UserPermissionSystem:
	case define.UserPermissionAdmin:
	default:
		_, _, err := p.f.Interact().SendFormAndWaitResponse(form.MessageForm{
			Title:   "权限不足",
			Content: "您不是§r§a赞颂者管理人员§r, 因此您无法打开此选项卡。",
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("AdminPanel: %v", err)
		}
		return nil
	}

	resp, isUserCancel, err := p.f.Interact().SendFormAndWaitResponse(form.MessageForm{
		Title:   "通知",
		Content: "您将需要§r§e搜索§r需要修改权限的赞颂者用户。",
		Button1: "继续",
		Button2: "返回上一级菜单",
	})
	if err != nil {
		return fmt.Errorf("AdminPanel: %v", err)
	}
	if isUserCancel || !resp.(bool) {
		return nil
	}

	eulogistUserName, isUserCancel, err := p.f.ShowEulogistSearch()
	if err != nil {
		return fmt.Errorf("AdminPanel: %v", err)
	}
	if isUserCancel {
		return nil
	}

	for {
		actionForm := form.ActionForm{
			Title: "选择",
			Content: fmt.Sprintf(
				"您将修改 §r§e%s§r 的用户信息, 请从以下选项中选择一个。",
				eulogistUserName,
			),
			Buttons: []form.ActionFormElement{
				{
					Text: "修改主要部分",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "修改该用户可以管理的租赁服",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "返回上一级菜单",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		resp, isUserCancel, err := p.f.Interact().SendFormAndWaitResponse(actionForm)
		if err != nil {
			return fmt.Errorf("AdminPanel: %v", err)
		}
		if isUserCancel {
			return nil
		}

		switch resp.(int32) {
		case 0:
			err = p.f.AdminChangeMainConfig(eulogistUserName)
		case 1:
			err = p.f.AdminChangeManager(eulogistUserName)
		case 2:
			return nil
		}
		if err != nil {
			return fmt.Errorf("AdminPanel: %v", err)
		}
	}
}

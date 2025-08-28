package panel

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// DevelopPanel ..
func (p *Panel) DevelopPanel() error {
	switch p.f.EulogistUserData().UserPermissionLevel {
	case define.UserPermissionSystem:
	case define.UserPermissionAdmin:
	case define.UserPermissionAdvance:
	default:
		_, _, err := p.f.Interact().SendFormAndWaitResponse(form.MessageForm{
			Title:   "权限不足",
			Content: "命令拒绝。",
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("DevelopPanel: %v", err)
		}
		return nil
	}

	for {
		var success bool

		actionForm := form.ActionForm{
			Title: "Developer features",
			Content: "" +
				"This page only open for those who have advance access.\n\n" +
				"" +
				" - Select §r§eUse built in skin§r would lead you to set a built in skin.\n" +
				" - Select §r§eUse custom skin§r would lead you to set a custom skin.\n" +
				" - Select §r§eUse skin from standard Minecraft§r would lead you to\n" +
				"select a skin from standard Minecraft without exit current game.\n\n" +
				"" +
				"Note that this setting is §r§eone-time§r, \n" +
				"and also you must ensure you §r§eequipped§r a verified NetEase skin.",
			Buttons: []form.ActionFormElement{
				{
					Text: "Use built in skin",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "Use custom skin",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "Use skin from standard Minecraft",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "Back to main panel",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		resp, isUserCancel, err := p.f.Interact().SendFormAndWaitResponse(actionForm)
		if err != nil {
			return fmt.Errorf("DevelopPanel: %v", err)
		}
		if isUserCancel {
			return nil
		}

		switch resp.(int32) {
		case 0:
			success, err = p.f.UseBuiltInSkin()
		case 1:
			success, err = p.f.UseCustomSkin()
		case 2:
			err = p.f.UseSkinFromMinecraft()
		case 3:
			return nil
		}

		if err != nil {
			return fmt.Errorf("DevelopPanel: %v", err)
		}
		if success {
			return nil
		}
	}
}

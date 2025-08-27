package panel

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// SkinCachePanel ..
func (p *Panel) SkinCachePanel() error {
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
			return fmt.Errorf("SkinCachePanel: %v", err)
		}
		return nil
	}

	for {
		actionForm := form.ActionForm{
			Title: "Choose an operation",
			Content: "" +
				"This page only open for those who have advance access.\n\n" +
				"Choose §r§eUse Built In Skin§r would lead you to set a built in skin, \n" +
				"and choose §r§eUse Custom Skin§r would lead you to set a custom skin.\n\n" +
				"Note that this setting is §r§eone-time§r, and also you must ensure you §r§eequipped§r a verified NetEase skin.",
			Buttons: []form.ActionFormElement{
				{
					Text: "Use Built In Skin",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "Use Custom Skin",
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
			return fmt.Errorf("SkinCachePanel: %v", err)
		}
		if isUserCancel {
			return nil
		}

		switch resp.(int32) {
		case 0:
			isUserCancel, err = p.f.UseBuiltInSkin()
		case 1:
			isUserCancel, err = p.f.UseCustomSkin()
		case 2:
			return nil
		}

		if err != nil {
			return fmt.Errorf("SkinCachePanel: %v", err)
		}
		if !isUserCancel {
			return nil
		}
	}
}

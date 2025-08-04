package panel

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/form"
)

// ServerOwnerPanel 向用户展示租赁服服主的控制面板。
// 从设计上讲，ServerOwnerPanel 应该从 MainPanel 调用
func (p *Panel) ServerOwnerPanel() error {
	for {
		minecraftForm := form.ActionForm{
			Title: "服主管理面板",
			Buttons: []form.ActionFormElement{
				{
					Text: "新增授权",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "修改已有授权",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "移除已有授权",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "返回上一级菜单",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		if len(p.f.EulogistUserData().RentalServerCanManage) == 0 {
			minecraftForm.Content = "您旗下没有可管理的租赁服, 请联系§r§a管理员§r添加。\n请选择您要进行的操作。"
		} else {
			minecraftForm.Content = fmt.Sprintf(
				"您目前总共可以管理 §r§b%d§r 个租赁服。\n请选择您要进行的操作。",
				len(p.f.EulogistUserData().RentalServerCanManage),
			)
		}

		resp, isUserCancel, err := p.f.Interact().SendFormAndWaitResponse(minecraftForm)
		if err != nil {
			return fmt.Errorf("ServerOwnerPanel: %v", err)
		}
		if isUserCancel {
			return nil
		}

		idx := resp.(int32)
		switch idx {
		case 0:
			_, err = p.f.RentalServerManage(p.f.RentalServerManageAddNew)
		case 1:
			_, err = p.f.RentalServerManage(p.f.RentalServerManageModify)
		case 2:
			_, err = p.f.RentalServerManage(p.f.RentalServerManageDelete)
		case 3:
			return nil
		}
		if err != nil {
			return fmt.Errorf("ServerOwnerPanel: %v", err)
		}
	}
}

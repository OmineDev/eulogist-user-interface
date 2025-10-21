package panel

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/form"
)

// AuthHelperPanel 向用户展示验证服务账户的管理面板。
// 从设计上讲，AuthHelperPanel 应该从 MainPanel 调用
func (p *Panel) AuthHelperPanel() error {
	for {
		minecraftForm := form.ActionForm{
			Title: "MC 账号管理面板",
			Buttons: []form.ActionFormElement{
				{
					Text: "查看正在使用的 MC 账号的详细信息",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "切换正在使用的 MC 账号",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "添加新 MC 账号",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "修改已添加的第三方验证服务账户",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "删除已添加的 MC 账号",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "下载标准验证服务的账户凭据",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "设置一次性 PE Auth 凭据",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "返回上一级菜单",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		if len(p.f.EulogistUserData().MultipleAuthServerAccounts) == 0 {
			minecraftForm.Content = "您目前还没有 MC 账号, 请先§r§e添加§r一个 MC 账号。\n请选择您要进行的操作。"
		} else {
			account, ok := p.f.EulogistUserData().CurrentAuthServerAccount.Value()
			if ok {
				minecraftForm.Content = fmt.Sprintf(
					"当前正在使用的账号是: \n  %s\n\n请选择您要进行的操作。",
					account.FormatInGame(),
				)
			} else {
				minecraftForm.Content = "当前还没有设置使用的 MC 账号, 请先§r§e设置§r一个 MC 账号。\n请选择您要进行的操作。"
			}
		}

		resp, isUserCancel, err := p.f.Interact().SendFormAndWaitResponse(minecraftForm)
		if err != nil {
			return fmt.Errorf("ManageAuthHelper: %v", err)
		}
		if isUserCancel {
			return nil
		}

		idx := resp.(int32)
		switch idx {
		case 0:
			err = p.f.ShowHelperDetails()
		case 1:
			err = p.f.ChangeCurrentHelper()
		case 2:
			err = p.f.AddNewHelper()
		case 3:
			err = p.f.ModifyCustomHelper()
		case 4:
			err = p.f.DeleteHelper()
		case 5:
			err = p.f.DevAskToken()
		case 6:
			err = p.f.SetAuthData()
		case 7:
			return nil
		}
		if err != nil {
			return fmt.Errorf("ManageAuthHelper: %v", err)
		}
	}
}

package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// HelperAddRequest ..
type HelperDeleteRequest struct {
	Token string `json:"token,omitempty"`
	Index uint   `json:"index"`
}

// HelperAddResponse ..
type HelperDeleteResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// DeleteHelper 打开表单以使用户其已有的删除 MC 账户
func (f *Function) DeleteHelper() error {
	for {
		actionForm := form.ActionForm{
			Title: "切换账户",
		}
		if len(f.userData.MultipleAuthServerAccounts) == 0 {
			actionForm.Content = "" +
				"您目前没有添加任何 MC 账号, \n" +
				"因此您无法删除任何 MC 账号。\n" +
				"您可以点击右上方的 §r§e[x]§r 退出本级菜单。"
		} else {
			actionForm.Content = fmt.Sprintf(
				""+
					"您目前总计已有 §r§b%d§r 个账号, \n"+
					"请从中选择一个以便于§r§e删除§r。",
				len(f.userData.MultipleAuthServerAccounts),
			)
		}
		for _, value := range f.userData.MultipleAuthServerAccounts {
			actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
				Text: value.FormatInGame(),
				Icon: form.ActionFormIconNone{},
			})
		}

		resp, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
		if err != nil {
			return fmt.Errorf("DeleteHelper: %v", err)
		}
		if isUserCancel {
			return nil
		}

		ensure, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "二重确认",
			Content: "您确定要§r§c删除§r目标 MC 账号吗？",
			Button1: "确定",
			Button2: "取消",
		})
		if err != nil {
			return fmt.Errorf("DeleteHelper: %v", err)
		}
		if isUserCancel || !ensure.(bool) {
			continue
		}

		helperDeleteResponse, err := SendAndGetHttpResponse[HelperDeleteResponse](
			fmt.Sprintf("%s/delete_helper", define.StdAuthServerAddress),
			HelperDeleteRequest{
				Token: f.config.EulogistToken,
				Index: uint(resp),
			},
		)
		if err != nil {
			return fmt.Errorf("DeleteHelper: %v", err)
		}
		if !helperDeleteResponse.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: helperDeleteResponse.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("DeleteHelper: %v", err)
			}
			return nil
		}

		newAccounts := make([]define.AuthServerAccount, 0)
		for index, value := range f.userData.MultipleAuthServerAccounts {
			if index != int(resp) {
				newAccounts = append(newAccounts, value)
			}
		}
		f.userData.MultipleAuthServerAccounts = newAccounts

		return nil
	}
}

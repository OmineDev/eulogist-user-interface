package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
	"github.com/ToolDelta-Basic/gophertunnel/minecraft/protocol"
)

// HelperChangeRequest ..
type HelperChangeRequest struct {
	Token string `json:"token,omitempty"`
	Index uint   `json:"index"`
}

// HelperChangeResponse ..
type HelperChangeResponse struct {
	ErrorInfo    string `json:"error_info"`
	Success      bool   `json:"success"`
	GameNickName string `json:"game_nick_name"`
	G79UserUID   string `json:"g79_user_uid"`
}

// ChangeCurrentHelper 向用户展示其已有的所有 MC 账户，
// 以供用户切换当前正在使用的 MC 账户
func (f *Function) ChangeCurrentHelper() error {
	var helperChangeRespose HelperChangeResponse

	actionForm := form.ActionForm{
		Title: "设置账户",
	}
	if len(f.userData.MultipleAuthServerAccounts) == 0 {
		actionForm.Content = "" +
			"您目前没有添加任何 MC 账号, \n" +
			"请回到点击 §r§e[x]§r 回到上一级菜单, \n" +
			"然后§r§e添加§r一个 MC 账号。"
	} else {
		actionForm.Content = fmt.Sprintf(
			""+
				"您目前总计已添加 §r§b%d§r 个账号, \n"+
				"请§r§e从中选择一个§r作为当前使用的 MC 账号。",
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
		return fmt.Errorf("ChangeCurrentHelper: %v", err)
	}
	if isUserCancel {
		return nil
	}

	helperChangeRespose, err = utils.SendAndGetHttpResponse[HelperChangeResponse](
		fmt.Sprintf("%s/change_current_helper", define.AddressEulogistAPI),
		HelperChangeRequest{
			Token: f.config.EulogistToken,
			Index: uint(resp),
		},
	)
	if err != nil {
		return fmt.Errorf("ChangeCurrentHelper: %v", err)
	}
	if !helperChangeRespose.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: helperChangeRespose.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("ChangeCurrentHelper: %v", err)
		}
		return nil
	}

	account := f.userData.MultipleAuthServerAccounts[resp]
	f.userData.CurrentAuthServerAccount = protocol.Option(account)

	if account.IsStdAccount() {
		account.UpdateData(map[string]any{
			"gameNickName":       helperChangeRespose.GameNickName,
			"g79UserUID":         helperChangeRespose.G79UserUID,
			"authHelperUniqueID": account.AuthServerSecret(),
		})
		f.userData.CurrentAuthServerAccount = protocol.Option(account)

		for index, value := range f.userData.MultipleAuthServerAccounts {
			if !value.IsStdAccount() {
				continue
			}
			if value.AuthServerSecret() == account.AuthServerSecret() {
				f.userData.MultipleAuthServerAccounts[index] = account
			}
		}
	}

	return nil
}

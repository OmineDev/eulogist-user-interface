package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
	"github.com/ToolDelta-Basic/gophertunnel/minecraft/protocol"
)

// CustomHelperModifyRequest ..
type CustomHelperModifyRequest struct {
	Token             string `json:"token,omitempty"`
	Index             uint   `json:"index"`
	AuthServerAddress string `json:"auth_server_address"`
	AuthServerToken   string `json:"auth_server_token"`
}

// CustomHelperModifyResponse ..
type CustomHelperModifyResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// ModifyCustomHelper 向用户打开表单，
// 用于修改已添加的第三方验证服务账户
func (f *Function) ModifyCustomHelper() error {
	for {
		toGlobalIndex := make([]uint, 0)
		actionForm := form.ActionForm{
			Title: "修改第三方验证服务账户",
		}

		for globalIndex, value := range f.userData.MultipleAuthServerAccounts {
			if !value.IsStdAccount() {
				actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
					Text: value.FormatInGame(),
					Icon: form.ActionFormIconNone{},
				})
				toGlobalIndex = append(toGlobalIndex, uint(globalIndex))
			}
		}
		if len(actionForm.Buttons) == 0 {
			actionForm.Content = "" +
				"您目前没有添加任何第三方验证服务的 MC 账号, \n" +
				"请回到点击 §r§e[x]§r 回到上一级菜单, \n" +
				"然后§r§e添加§r一个第三方验证服务的 MC 账号。"
		} else {
			actionForm.Content = fmt.Sprintf(
				""+
					"您目前总计已添加 §r§b%d§r 个第三方验证服务的账号, \n"+
					"请§r§e从中选择一个§r, 然后您可以开始修改它的数据。",
				len(actionForm.Buttons),
			)
		}

		resp, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
		if err != nil {
			return fmt.Errorf("ModifyCustomHelper: %v", err)
		}
		if isUserCancel {
			return nil
		}
		idx := toGlobalIndex[resp]

		src := f.userData.MultipleAuthServerAccounts[idx]
		minecraftForm := form.ModalForm{
			Title: "修改第三方验证服务账户",
			Contents: []form.ModalFormElement{
				form.ModalFormElementLabel{
					Text: "请§r§e修改§r第三方验证服务账户的数据。",
				},
				form.ModalFormElementInput{
					Text:        "验证服务地址",
					Default:     src.AuthServerAddress(),
					PlaceHolder: "Auth Service Address",
				},
				form.ModalFormElementInput{
					Text:        "验证服务令牌",
					Default:     src.AuthServerSecret(),
					PlaceHolder: "Auth Service Token",
				},
			},
		}
		answer, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
		if err != nil {
			return fmt.Errorf("ModifyCustomHelper: %v", err)
		}
		if isUserCancel {
			continue
		}

		answerList := answer.([]any)
		authServerAddress := answerList[1].(string)
		authServerToken := answerList[2].(string)

		customHelperModifyResp, err := utils.SendAndGetHttpResponse[CustomHelperModifyResponse](
			fmt.Sprintf("%s/modify_custom_helper", define.StdAuthServerAddress),
			CustomHelperModifyRequest{
				Token:             f.config.EulogistToken,
				Index:             idx,
				AuthServerAddress: authServerAddress,
				AuthServerToken:   authServerToken,
			},
		)
		if err != nil {
			return fmt.Errorf("ModifyCustomHelper: %v", err)
		}
		if !customHelperModifyResp.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: customHelperModifyResp.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("ModifyCustomHelper: %v", err)
			}
			continue
		}

		newAccount := define.CustomAuthServerAccount{}
		newAccount.UpdateData(map[string]any{
			"internalAccountID": src.(*define.CustomAuthServerAccount).InternalAccountID(),
			"authServerAddress": authServerAddress,
			"authServerToken":   authServerToken,
		})
		f.userData.MultipleAuthServerAccounts[idx] = &newAccount

		dst, ok := f.userData.CurrentAuthServerAccount.Value()
		if ok {
			if !dst.IsStdAccount() && src.AuthServerAddress() == dst.AuthServerAddress() && src.AuthServerSecret() == dst.AuthServerSecret() {
				f.userData.CurrentAuthServerAccount = protocol.Option(f.userData.MultipleAuthServerAccounts[idx])
			}
		}

		return nil
	}
}

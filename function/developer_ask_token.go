package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

// DevTokenAsk ..
type DevTokenAsk struct {
	Token          string `json:"token,omitempty"`
	HelperUniqueID string `json:"helper_unique_id,omitempty"`
}

// DevTokenResp ..
type DevTokenResp struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
	Token     string `json:"token"`
}

// DevAskToken ..
func (f *Function) DevAskToken() error {
	if !f.userData.CanGetHelperToken {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "权限不足",
			Content: "该功能是为§r§e关键的开发者§r启用的, 请联系管理员以取得授权。",
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("DevAskToken: %v", err)
		}
		return nil
	}

	actionForm := form.ActionForm{
		Title: "下载辅助用户的凭据 (Token)",
	}
	if len(f.userData.MultipleAuthServerAccounts) == 0 {
		actionForm.Content = "" +
			"您目前没有添加任何 MC 账号, \n" +
			"因此您无法获取任何 MC 账号的 Token。\n" +
			"您可以点击右上方的 §r§e[x]§r 退出本级菜单。"
	} else {
		actionForm.Content = fmt.Sprintf(
			""+
				"您目前总计已有 §r§b%d§r 个账号, \n"+
				"请从中选择一个, 然后便可§r§e取得§r它的 Token。",
			len(f.userData.MultipleAuthServerAccounts),
		)
	}

	idxToHelperUniqueID := make([]string, 0)
	for _, value := range f.userData.MultipleAuthServerAccounts {
		if value.IsStdAccount() {
			actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
				Text: value.FormatInGame(),
				Icon: form.ActionFormIconNone{},
			})
			idxToHelperUniqueID = append(idxToHelperUniqueID, value.AuthServerSecret())
		}
	}

	resp, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
	if err != nil {
		return fmt.Errorf("DevAskToken: %v", err)
	}
	if isUserCancel {
		return nil
	}

	devTokenResp, err := utils.SendAndGetHttpResponse[DevTokenResp](
		fmt.Sprintf("%s/dev_ask_token", define.StdAuthServerAddress),
		DevTokenAsk{
			Token:          f.config.EulogistToken,
			HelperUniqueID: idxToHelperUniqueID[resp],
		},
	)
	if err != nil {
		return fmt.Errorf("DevAskToken: %v", err)
	}
	if !devTokenResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: devTokenResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("DevAskToken: %v", err)
		}
		return nil
	}

	_, _, err = f.interact.SendFormAndWaitResponse(form.ModalForm{
		Title: "Auth Server Token",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "" +
					"The following string is your token of the corresponding std auth helper.\n" +
					"Please §r§l§eKEEP IT ON YOUR OWN§r, because it can be directly used for bot, \n" +
					"like ToolDelta or Neo-Omega.\n" +
					"Note that the token will expire if §r§aEulogist Admins§r disable your permission.",
			},
			form.ModalFormElementInput{
				Text:        "Standard Auth Server Address",
				Default:     define.StdAuthServerPhoenixAPI,
				PlaceHolder: "Standard Auth Server Address",
			},
			form.ModalFormElementInput{
				Text:        "Standard Auth Server Token",
				Default:     devTokenResp.Token,
				PlaceHolder: "Standard Auth Server Token",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("DevAskToken: %v", err)
	}

	return nil
}

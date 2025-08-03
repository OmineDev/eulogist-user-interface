package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// PEAuthSetRequest ..
type PEAuthSetRequest struct {
	Token  string `json:"token,omitempty"`
	PEAuth string `json:"pe_auth,omitempty"`
}

// PEAuthSetResponse ..
type PEAuthSetResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// SetPEAuth 打开表单以使用户设置其所使用的 PE Auth
func (f *Function) SetPEAuth() error {
	minecraftForm := form.ModalForm{
		Title: "设置 PE Auth",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "" +
					"这是一个§r§e高级选项§r, 用于§r§e渠道服登录§r。\n" +
					"设置后, 下次将使用它对应的 MC 账户进服。\n" +
					"如果您不知道它的使用方法, 请§r§e返回到上一级菜单§r, 否则您的下次进服可能会§r§c出现未知问题§r。",
			},
			form.ModalFormElementInput{
				Text:        "PE Auth",
				Default:     "",
				PlaceHolder: "Encrypt hex or JSON string",
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("SetPEAuth: %v", err)
	}
	if isUserCancel {
		return nil
	}
	peAuth := resp.([]any)[1].(string)

	peAuthResponse, err := SendAndGetHttpResponse[PEAuthSetResponse](
		fmt.Sprintf("%s/set_pe_auth", define.StdAuthServerAddress),
		PEAuthSetRequest{
			Token:  f.config.EulogistToken,
			PEAuth: resp.([]any)[1].(string),
		},
	)
	if err != nil {
		return fmt.Errorf("SetPEAuth: %v", err)
	}
	if !peAuthResponse.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: peAuthResponse.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("SetPEAuth: %v", err)
		}
		return nil
	}

	f.userData.ProvidedPeAuthData = peAuth
	return nil
}

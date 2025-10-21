package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

// AuthDataSetRequest ..
type AuthDataSetRequest struct {
	Token   string `json:"token,omitempty"`
	DoClean bool   `json:"do_clean,omitempty"`
	PeAuth  string `json:"pe_auth,omitempty"`
	SaAuth  string `json:"sa_auth,omitempty"`
}

// AuthDataSetResponse ..
type AuthDataSetResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// SetAuthData 打开表单以使用户设置其所使用的 Pe Auth 或 Sa Auth
func (f *Function) SetAuthData() error {
	for {
		var doClean bool
		var peAuth string
		var saAuth string

		minecraftForm := form.ActionForm{
			Title: "设置账号登录状态",
			Content: "" +
				"这是一个§r§e高级选项§r, 用于§r§e渠道服登录§r。\n" +
				"设置后, 未来一段时间内将通过该 MC 账户进服。\n" +
				"请确保您§r§e知道§r它们的用法, 否则可能会遇见§r§e出现未知问题§r。",
			Buttons: []form.ActionFormElement{
				{
					Text: "设置 Pe Auth",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "设置 Sa Auth",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "清除登录状态",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "返回上一级菜单",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
		if err != nil {
			return fmt.Errorf("SetAuthData: %v", err)
		}
		if isUserCancel || resp.(int32) == 3 {
			return nil
		}

		if resp.(int32) == 0 {
			modalForm := form.ModalForm{
				Title: "设置 Pe Auth",
				Contents: []form.ModalFormElement{
					form.ModalFormElementLabel{
						Text: "请填写 Pe Auth 字符串。",
					},
					form.ModalFormElementInput{
						Text:        "Pe Auth",
						Default:     "",
						PlaceHolder: "Encrypt hex or JSON string",
					},
				},
			}

			resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(modalForm)
			if err != nil {
				return fmt.Errorf("SetAuthData: %v", err)
			}
			if isUserCancel {
				continue
			}

			peAuth = resp.([]any)[1].(string)
		}
		if resp.(int32) == 1 {
			modalForm := form.ModalForm{
				Title: "设置 Sa Auth",
				Contents: []form.ModalFormElement{
					form.ModalFormElementLabel{
						Text: "请填写 Sa Auth 字符串。",
					},
					form.ModalFormElementInput{
						Text:        "Sa Auth",
						Default:     "",
						PlaceHolder: `JSON string`,
					},
				},
			}

			resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(modalForm)
			if err != nil {
				return fmt.Errorf("SetAuthData: %v", err)
			}
			if isUserCancel {
				continue
			}

			saAuth = resp.([]any)[1].(string)
		}
		if resp.(int32) == 2 {
			messageForm := form.MessageForm{
				Title:   "二重确认",
				Content: "您确定要§r§c清除登录状态§r吗？",
				Button1: "确定",
				Button2: "取消",
			}

			resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(messageForm)
			if err != nil {
				return fmt.Errorf("SetAuthData: %v", err)
			}
			if isUserCancel || !resp.(bool) {
				continue
			}

			doClean = true
		}

		authDataSetResp, err := utils.SendAndGetHttpResponse[AuthDataSetResponse](
			fmt.Sprintf("%s/set_auth_data", define.StdAuthServerAddress),
			AuthDataSetRequest{
				Token:   f.config.EulogistToken,
				DoClean: doClean,
				PeAuth:  peAuth,
				SaAuth:  saAuth,
			},
		)
		if err != nil {
			return fmt.Errorf("SetAuthData: %v", err)
		}
		if !authDataSetResp.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: authDataSetResp.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("SetAuthData: %v", err)
			}
			continue
		}

		return nil
	}
}

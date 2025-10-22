package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

const (
	RequestTypeRegisterSession uint8 = iota
	RequestTypeCleanUpSession
	RequestTypeEnableVitality
	RequestTypeDisableVitality
)

// LoginStatesSetRequest ..
type LoginStatesSetRequest struct {
	Token       string `json:"token,omitempty"`
	RequestType uint8  `json:"request_type,omitempty"`
	PeAuth      string `json:"pe_auth,omitempty"`
	SaAuth      string `json:"sa_auth,omitempty"`
}

// LoginStatesSetResponse ..
type LoginStatesSetResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// SetLoginStates 打开表单以使用户设置其 mc 账号的登录状态
func (f *Function) SetLoginStates() error {
	for {
		var requestType uint8
		var peAuth string
		var saAuth string
		var doClean bool

		minecraftForm := form.ActionForm{
			Title:   "设置账号的登录状态",
			Content: "请选择您要进行的操作。",
			Buttons: []form.ActionFormElement{
				{
					Text: "设置会话状态",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "设置会话维持器",
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
			return fmt.Errorf("SetLoginStates: %v", err)
		}
		if isUserCancel || resp.(int32) == 3 {
			return nil
		}

		switch resp.(int32) {
		case 0:
			peAuth, saAuth, doClean, isUserCancel, err = f.SetSessionStates()
			if err != nil {
				return fmt.Errorf("SetLoginStates: %v", err)
			}
			if isUserCancel {
				continue
			}
			if doClean {
				requestType = RequestTypeCleanUpSession
			} else {
				requestType = RequestTypeRegisterSession
			}
		case 1:
			enableVitality, isUserCancel, err := f.SetSessionMaintainer()
			if err != nil {
				return fmt.Errorf("SetLoginStates: %v", err)
			}
			if isUserCancel {
				continue
			}
			if enableVitality {
				requestType = RequestTypeEnableVitality
			} else {
				requestType = RequestTypeDisableVitality
			}
		case 2:
			return nil
		}

		loginStatesSetResp, err := utils.SendAndGetHttpResponse[LoginStatesSetResponse](
			fmt.Sprintf("%s/set_login_states", define.AddressEulogistAPI),
			LoginStatesSetRequest{
				Token:       f.config.EulogistToken,
				RequestType: requestType,
				PeAuth:      peAuth,
				SaAuth:      saAuth,
			},
		)
		if err != nil {
			return fmt.Errorf("SetLoginStates: %v", err)
		}
		if !loginStatesSetResp.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: loginStatesSetResp.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("SetLoginStates: %v", err)
			}
			continue
		}

		return nil
	}
}

func (f *Function) SetSessionStates() (peAuth string, saAuth string, doClean bool, isUserCancel bool, err error) {
	for {
		minecraftForm := form.ActionForm{
			Title: "设置会话状态",
			Content: "" +
				"设置会话状态被用于§r§e渠道服登录§r。\n\n" +
				"设置 §r§ePe Auth§r 或 §r§eSa Auth§r 后, \n" +
				"未来一段时间内将通过它们对应的 MC 账户进服。\n\n" +
				"您可以通过§r§e清除账户的登录状态§r来撤销这些设置。\n\n" +
				"请选择您要进行的操作。",
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
					Text: "清除账户的登录状态",
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
			return "", "", false, false, fmt.Errorf("SetSessionStates: %v", err)
		}
		if isUserCancel {
			return "", "", false, true, nil
		}

		switch resp.(int32) {
		case 0:
			modalForm := form.ModalForm{
				Title: "设置 Pe Auth",
				Contents: []form.ModalFormElement{
					form.ModalFormElementLabel{
						Text: "请填写 §r§ePe Auth§r 字符串。",
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
				return "", "", false, false, fmt.Errorf("SetSessionStates: %v", err)
			}
			if isUserCancel {
				break
			}

			peAuth = resp.([]any)[1].(string)
			return peAuth, "", false, false, nil
		case 1:
			modalForm := form.ModalForm{
				Title: "设置 Sa Auth",
				Contents: []form.ModalFormElement{
					form.ModalFormElementLabel{
						Text: "请填写 §r§eSa Auth§r 字符串。",
					},
					form.ModalFormElementInput{
						Text:        "Sa Auth",
						Default:     "",
						PlaceHolder: "JSON string",
					},
				},
			}

			resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(modalForm)
			if err != nil {
				return "", "", false, false, fmt.Errorf("SetSessionStates: %v", err)
			}
			if isUserCancel {
				break
			}

			saAuth = resp.([]any)[1].(string)
			return "", saAuth, false, false, nil
		case 2:
			messageForm := form.MessageForm{
				Title:   "二重确认",
				Content: "您确定要§r§c清除登录状态§r吗？",
				Button1: "确定",
				Button2: "取消",
			}

			resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(messageForm)
			if err != nil {
				return "", "", false, false, fmt.Errorf("SetSessionStates: %v", err)
			}
			if isUserCancel || !resp.(bool) {
				break
			}

			doClean = true
			return "", "", doClean, false, nil
		default:
			return "", "", false, true, nil
		}
	}
}

// SetSessionMaintainer ..
func (f *Function) SetSessionMaintainer() (enableVitality bool, isUserCancel bool, err error) {
	minecraftForm := form.ActionForm{
		Title: "设置会话维持器",
		Content: "" +
			"会话维持器用于维持当前所\n" +
			"使用的 MC 账号的§r§e在线状态§r。\n\n" +
			"" +
			"该维持器会在您通过赞颂者进入游戏后启动，\n" +
			"直到您关闭赞颂者时会终止运行。\n\n" +
			"" +
			"在会话维持器运行期间, \n" +
			"当前账号将得到对应的§r§e在线经验§r。\n\n" +
			"" +
			"请确保会话维持器运行期间, " +
			"\n该账号§r§e不会被重新登录§r, \n" +
			"或§r§e存在更多的会话维持器在运行§r。\n\n" +
			"" +
			"如果您需要使用同一个 MC 账号服务于§r§e多个§r服务器, \n" +
			"请§r§e关闭§r会话维持器, \n" +
			"并运行§r§e独立的会话维持器§r以确保维持器是唯一的。\n\n" +
			"" +
			"另外, 该功能, 即会话维持器, \n" +
			"在§r§e默认情况§r下呈现§r§e开启§r状态。\n\n" +
			"" +
			"请选择您要进行的操作。",
		Buttons: []form.ActionFormElement{
			{
				Text: "启用会话维持器",
				Icon: form.ActionFormIconNone{},
			},
			{
				Text: "关闭会话维持器",
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
		return false, false, fmt.Errorf("SetSessionMaintainer: %v", err)
	}
	if isUserCancel {
		return false, true, nil
	}

	switch resp.(int32) {
	case 0:
		return true, false, nil
	case 1:
		return false, false, nil
	default:
		return false, true, nil
	}
}

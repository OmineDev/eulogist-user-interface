package function

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/google/uuid"
)

const (
	ActionTypeAddStdEmailHelper uint8 = iota
	ActionTypeAddCustomHelper
)

// HelperAddRequest ..
type HelperAddRequest struct {
	Token      string `json:"token,omitempty"`
	ActionType uint8  `json:"action_type"`

	Email       string `json:"email,omitempty"`
	MD5Password string `json:"md5_password,omitempty"`

	AuthServerAddress string `json:"auth_server_address,omitempty"`
	AuthServerToken   string `json:"auth_server_token,omitempty"`
}

// HelperAddResponse ..
type HelperAddResponse struct {
	ErrorInfo            string `json:"error_info"`
	NetEaseRequireVerify bool   `json:"netease_require_verify"`
	VerifyURL            string `json:"verify_url"`
	Success              bool   `json:"success"`
	HelperUniqueID       string `json:"helper_unique_id"`
	GameNickName         string `json:"game_nick_name"`
	G79UserUID           string `json:"g79_user_uid"`
}

// AddNewHelper 为客户端打开添加验证服务账户的表单
func (f *Function) AddNewHelper() error {
	actionForm := form.ActionForm{
		Title:   "添加账户",
		Content: "请从以下选项中选择一个。",
		Buttons: []form.ActionFormElement{
			{
				Text: "使用内置验证服务 (§r§a推荐§r)",
				Icon: form.ActionFormIconNone{},
			},
			{
				Text: "使用第三方验证服务",
				Icon: form.ActionFormIconNone{},
			},
			{
				Text: "返回上一级菜单",
				Icon: form.ActionFormIconNone{},
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(actionForm)
	if err != nil {
		return fmt.Errorf("AddNewHelper: %v", err)
	}
	if isUserCancel {
		return nil
	}

	switch resp.(int32) {
	case 0:
		err = f.addStdHelper()
	case 1:
		err = f.addCustomHelper()
	}
	if err != nil {
		return fmt.Errorf("AddNewHelper: %v", err)
	}

	return nil
}

// addCustomHelper ..
func (f *Function) addCustomHelper() error {
	var helperAddResponse HelperAddResponse

	minecraftForm := form.ModalForm{
		Title: "添加第三方验证服务账户",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "请输入第三方验证服务账户的§r§e验证服务地址§r和§r§e验证服务令牌§r。",
			},
			form.ModalFormElementInput{
				Text:        "验证服务地址",
				Default:     "https://nv1.nethard.pro",
				PlaceHolder: "Auth Service Address",
			},
			form.ModalFormElementInput{
				Text:        "验证服务令牌",
				Default:     "",
				PlaceHolder: "Auth Service Token",
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("addCustomHelper: %v", err)
	}
	if isUserCancel {
		return nil
	}

	respList := resp.([]any)
	authServerAddress := respList[1].(string)
	authServerToken := respList[2].(string)

	helperAddResponse, err = SendAndGetHttpResponse[HelperAddResponse](
		fmt.Sprintf("%s/add_helper_normal", define.StdAuthServerAddress),
		HelperAddRequest{
			Token:             f.config.EulogistToken,
			ActionType:        ActionTypeAddCustomHelper,
			AuthServerAddress: authServerAddress,
			AuthServerToken:   authServerToken,
		},
	)
	if err != nil {
		return fmt.Errorf("addCustomHelper: %v", err)
	}
	if !helperAddResponse.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: helperAddResponse.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("addCustomHelper: %v", err)
		}
		return nil
	}

	f.userData.InternalIncreasingAccountID++
	account := define.CustomAuthServerAccount{}
	account.UpdateData(map[string]any{
		"internalAccountID": f.userData.InternalIncreasingAccountID,
		"authServerAddress": authServerAddress,
		"authServerToken":   authServerToken,
	})
	f.userData.MultipleAuthServerAccounts = append(f.userData.MultipleAuthServerAccounts, &account)

	return nil
}

// addStdHelper ..
func (f *Function) addStdHelper() error {
	actionForm := form.ActionForm{
		Title:   "添加账户",
		Content: "我们目前只支持官服, 渠道服登录请返回上一级菜单, 并使用 §r§ePE Auth§r 凭据。",
		Buttons: []form.ActionFormElement{
			{
				Text: "添加网易邮箱对应的 MC 账号",
				Icon: form.ActionFormIconNone{},
			},
			{
				Text: "添加手机号对应的 MC 账号",
				Icon: form.ActionFormIconNone{},
			},
			{
				Text: "返回上一级菜单",
				Icon: form.ActionFormIconNone{},
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(actionForm)
	if err != nil {
		return fmt.Errorf("addStdHelper: %v", err)
	}
	if isUserCancel {
		return nil
	}

	switch resp.(int32) {
	case 0:
		err = f.addStdHelperByEmail()
	case 1:
		err = f.addStdHelperByMobile()
	}
	if err != nil {
		return fmt.Errorf("addStdHelper: %v", err)
	}

	return nil
}

// addStdHelperByEmail ..
func (f *Function) addStdHelperByEmail() error {
	var helperAddResponse HelperAddResponse

	minecraftForm := form.ModalForm{
		Title: "通过邮箱添加账户",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "请输入您网易邮箱的§r§e邮箱地址§r和§r§e邮箱密码§r。",
			},
			form.ModalFormElementInput{
				Text:        "邮箱地址",
				Default:     "",
				PlaceHolder: "xxx@163.com",
			},
			form.ModalFormElementInput{
				Text:        "邮箱密码",
				Default:     "",
				PlaceHolder: "123456",
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("addStdHelperByEmail: %v", err)
	}
	if isUserCancel {
		return nil
	}

	respList := resp.([]any)
	emailAddress := respList[1].(string)
	emailPassword := respList[2].(string)
	emailPasswordMD5 := md5.Sum([]byte(emailPassword))

	for {
		helperAddResponse, err = SendAndGetHttpResponse[HelperAddResponse](
			fmt.Sprintf("%s/add_helper_normal", define.StdAuthServerAddress),
			HelperAddRequest{
				Token:       f.config.EulogistToken,
				ActionType:  ActionTypeAddStdEmailHelper,
				Email:       emailAddress,
				MD5Password: hex.EncodeToString(emailPasswordMD5[:]),
			},
		)
		if err != nil {
			return fmt.Errorf("addStdHelperByEmail: %v", err)
		}
		if !helperAddResponse.Success {
			err = f.ShowAuthServerError(
				helperAddResponse.NetEaseRequireVerify,
				helperAddResponse.VerifyURL,
				helperAddResponse.ErrorInfo,
			)
			if err != nil {
				return fmt.Errorf("addStdHelperByEmail: %v", err)
			}
			if helperAddResponse.NetEaseRequireVerify {
				continue
			} else {
				return nil
			}
		}
		break
	}

	account := define.StdAuthServerAccount{}
	account.UpdateData(map[string]any{
		"gameNickName":       helperAddResponse.GameNickName,
		"g79UserUID":         helperAddResponse.G79UserUID,
		"authHelperUniqueID": helperAddResponse.HelperUniqueID,
	})
	f.userData.MultipleAuthServerAccounts = append(f.userData.MultipleAuthServerAccounts, &account)

	return nil
}

const (
	ActionTypeOpenNewTransaction uint8 = iota // Open new transaction
	ActionTypeFinishVerify                    // User send/receive SMS to/from NetEase
)

const (
	ResponseTypeClientNeedSendSMS    uint8 = iota // User need send SMS to NetEase
	ResponseTypeClientNeedReceiveSMS              // User need receive SMS from NetEase
	ResponseTypeLoginSuccess                      // SMS Login success
	ResponseTypeMeetError                         // SMS Login meet error
)

// SMSHelperAddRequest ..
type SMSHelperAddRequest struct {
	Token           string `json:"token,omitempty"`
	TransactionUUID string `json:"transaction_uuid"`
	ActionType      uint8  `json:"action_type"`
	Mobile          string `json:"mobile,omitempty"`
	SMSVerifyCode   string `json:"sms_verify_code,omitempty"`
}

// SMSHelperAddResponse ..
type SMSHelperAddResponse struct {
	ErrorInfo      string `json:"error_info"`
	ResponseType   uint8  `json:"response_type"`
	VerifyURL      string `json:"verify_url"`
	HelperUniqueID string `json:"helper_unique_id"`
	GameNickName   string `json:"game_nick_name"`
	G79UserUID     string `json:"g79_user_uid"`
}

// addStdHelperByMobile ..
func (f *Function) addStdHelperByMobile() error {
	var verifyCode string

	minecraftForm := form.ModalForm{
		Title: "通过手机号添加账户",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "请输入您的§r§e手机号§r。",
			},
			form.ModalFormElementInput{
				Text:        "手机号",
				Default:     "",
				PlaceHolder: "12345678901",
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("addStdHelperByMobile: %v", err)
	}
	if isUserCancel {
		return nil
	}

	respList := resp.([]any)
	mobile := respList[1].(string)

	tranUUID := uuid.NewString()
	smsHelperAddResp, err := SendAndGetHttpResponse[SMSHelperAddResponse](
		fmt.Sprintf("%s/add_std_helper_sms", define.StdAuthServerAddress),
		SMSHelperAddRequest{
			Token:           f.config.EulogistToken,
			TransactionUUID: tranUUID,
			ActionType:      ActionTypeOpenNewTransaction,
			Mobile:          mobile,
		},
	)
	if err != nil {
		return fmt.Errorf("addStdHelperByMobile: %v", err)
	}

	switch smsHelperAddResp.ResponseType {
	case ResponseTypeClientNeedSendSMS:
		resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "短信验证",
			Content: fmt.Sprintf("此操作需要§r§e短信验证§r, 内容为: \n  %s", smsHelperAddResp.ErrorInfo),
			Button1: "我已完成验证",
			Button2: "取消",
		})
		if err != nil {
			return fmt.Errorf("addStdHelperByMobile: %v", err)
		}
		if isUserCancel || !resp.(bool) {
			return nil
		}
	case ResponseTypeClientNeedReceiveSMS:
		resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.ModalForm{
			Title: "短信验证",
			Contents: []form.ModalFormElement{
				form.ModalFormElementLabel{
					Text: "此操作需要§r§e短信验证§r, 请从手机上接收来自网易的短信, 并填写在下方。",
				},
				form.ModalFormElementInput{
					Text:        "短信验证码",
					Default:     "",
					PlaceHolder: "6 位数验证码",
				},
			},
		})
		if err != nil {
			return fmt.Errorf("addStdHelperByMobile: %v", err)
		}
		if isUserCancel {
			return nil
		}
		verifyCode = resp.([]any)[1].(string)
	}

	smsHelperAddResp, err = SendAndGetHttpResponse[SMSHelperAddResponse](
		fmt.Sprintf("%s/add_std_helper_sms", define.StdAuthServerAddress),
		SMSHelperAddRequest{
			Token:           f.config.EulogistToken,
			TransactionUUID: tranUUID,
			ActionType:      ActionTypeFinishVerify,
			SMSVerifyCode:   verifyCode,
		},
	)
	if err != nil {
		return fmt.Errorf("addStdHelperByMobile: %v", err)
	}

	if smsHelperAddResp.ResponseType != ResponseTypeLoginSuccess {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: smsHelperAddResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("addStdHelperByMobile: %v", err)
		}
		return nil
	}

	account := define.StdAuthServerAccount{}
	account.UpdateData(map[string]any{
		"gameNickName":       smsHelperAddResp.GameNickName,
		"g79UserUID":         smsHelperAddResp.G79UserUID,
		"authHelperUniqueID": smsHelperAddResp.HelperUniqueID,
	})
	f.userData.MultipleAuthServerAccounts = append(f.userData.MultipleAuthServerAccounts, &account)

	return nil
}

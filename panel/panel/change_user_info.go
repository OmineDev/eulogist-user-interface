package panel

import (
	"crypto/sha256"
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

// UserInfoChangeRequest ..
type UserInfoChangeRequest struct {
	Token             string `json:"token,omitempty"`
	NewName           string `json:"new_name,omitempty"`
	NewPasswordSum256 []byte `json:"new_password_sum256,omitempty"`
}

// UserInfoChangeResponse ..
type UserInfoChangeResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
	NewToken  string `json:"new_token"`
}

// ChangeUserInfo ..
func (p *Panel) ChangeUserInfo() error {
	var newPasswordSum256Bytes []byte

	minecraftForm := form.ModalForm{
		Title: "编辑赞颂者账户信息",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "" +
					"您将§r§e编辑§r您自己的赞颂者账户的用户档案。\n" +
					"如果您§r§e不需要§r修改密码, 请将密码输入框§r§e置空§r。\n" +
					"§r§e一旦§r密码更改, 则您的§r§e赞颂者令牌§r将会发生变化。",
			},
			form.ModalFormElementInput{
				Text:        "新用户名",
				Default:     p.f.EulogistUserData().UserName,
				PlaceHolder: "您的赞颂者用户名",
			},
			form.ModalFormElementInput{
				Text:        "新密码",
				Default:     "",
				PlaceHolder: "如果无需修改, 请置空",
			},
		},
	}

	resp, isUserCancel, err := p.f.Interact().SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("ChangeUserInfo: %v", err)
	}
	if isUserCancel {
		return nil
	}

	respList := resp.([]any)
	newUserName := respList[1].(string)
	newPassword := respList[2].(string)
	if len(newPassword) > 0 {
		newPasswordSum256 := sha256.Sum256([]byte(newPassword + define.UserPasswordSlat))
		newPasswordSum256Bytes = newPasswordSum256[:]
	}

	userInfoChangeResp, err := utils.SendAndGetHttpResponse[UserInfoChangeResponse](
		fmt.Sprintf("%s/change_user_info", define.StdAuthServerAddress),
		UserInfoChangeRequest{
			Token:             p.f.EulogistConfig().EulogistToken,
			NewName:           newUserName,
			NewPasswordSum256: newPasswordSum256Bytes,
		},
	)
	if err != nil {
		return fmt.Errorf("ChangeUserInfo: %v", err)
	}

	if !userInfoChangeResp.Success {
		_, _, err := p.f.Interact().SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: userInfoChangeResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("ChangeUserInfo: %v", err)
		}
		return nil
	}

	p.f.EulogistUserData().UserName = newUserName
	p.f.EulogistConfig().EulogistToken = userInfoChangeResp.NewToken

	err = define.WriteConfig(p.f.EulogistConfig())
	if err != nil {
		return fmt.Errorf("ChangeUserInfo: %v", err)
	}

	return nil
}

package function

import (
	"crypto/sha256"
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// LoginRequest ..
type LoginRequest struct {
	IsRegister         bool   `json:"is_register"`
	UserName           string `json:"user_name,omitempty"`
	UserPasswordSum256 []byte `json:"user_password_sum256,omitempty"`
}

// LoginResponse ..
type LoginResponse struct {
	ErrorInfo     string `json:"error_info"`
	Success       bool   `json:"success"`
	EulogistToken string `json:"eulogist_token"`
}

// RegisterOrLogin 使用户注册或登录赞颂者账户。
// exitGame 指示用户是否希望关闭赞颂者并退出假 MC 服务器
func (f *Function) RegisterOrLogin() (exitGame bool, err error) {
	for {
		minecraftForm := form.ActionForm{
			Title:   "欢迎使用",
			Content: "您还没有注册或登录赞颂者账号, 请先注册或登录！",
			Buttons: []form.ActionFormElement{
				{
					Text: "注册新的赞颂者账户",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "登录已有的赞颂者账户",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "退出游戏",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		resp, isUserCacel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
		if err != nil {
			return false, fmt.Errorf("RegisterOrLogin: %v", err)
		}
		if isUserCacel {
			continue
		}

		idx := resp.(int32)
		switch idx {
		case 0:
			err = f.register()
		case 1:
			err = f.login()
		case 2:
			return true, nil
		}

		if err != nil {
			return false, fmt.Errorf("RegisterOrLogin: %v", err)
		}
		if len(f.config.EulogistToken) > 0 {
			return false, nil
		}
	}
}

// register ..
func (f *Function) register() error {
	minecraftForm := form.ModalForm{
		Title: "注册",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "您将开始注册一个新的赞颂者账户。",
			},
			form.ModalFormElementInput{
				Text:        "账户名",
				Default:     "",
				PlaceHolder: "赞颂者账户之名, 在 6 到 16 个字符之间",
			},
			form.ModalFormElementInput{
				Text:        "密码",
				Default:     "",
				PlaceHolder: "赞颂者账户的密码, 不得为空",
			},
		},
	}

	resp, isUserCacel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("register: %v", err)
	}
	if isUserCacel {
		return nil
	}

	respList := resp.([]any)
	userName := respList[1].(string)
	userPassword := respList[2].(string)
	userPasswordSum256 := sha256.Sum256([]byte(userPassword + define.UserPasswordSlat))

	loginResponse, err := SendAndGetHttpResponse[LoginResponse](
		fmt.Sprintf("%s/register_or_login", define.StdAuthServerAddress),
		LoginRequest{
			IsRegister:         true,
			UserName:           userName,
			UserPasswordSum256: userPasswordSum256[:],
		},
	)
	if err != nil {
		return fmt.Errorf("register: %v", err)
	}

	if !loginResponse.Success {
		_, _, err = f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "注册失败",
			Content: loginResponse.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("register: %v", err)
		}
		return nil
	}

	f.config.EulogistToken = loginResponse.EulogistToken
	err = define.WriteConfig(f.config)
	if err != nil {
		return fmt.Errorf("register: %v", err)
	}

	return nil
}

// login ..
func (f *Function) login() error {
	minecraftForm := form.ModalForm{
		Title: "登录",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "您将登录一个已经注册好的赞颂者用户。",
			},
			form.ModalFormElementInput{
				Text:        "账户名",
				Default:     "",
				PlaceHolder: "",
			},
			form.ModalFormElementInput{
				Text:        "密码",
				Default:     "",
				PlaceHolder: "",
			},
		},
	}

	resp, isUserCacel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("login: %v", err)
	}
	if isUserCacel {
		return nil
	}

	respList := resp.([]any)
	userName := respList[1].(string)
	userPassword := respList[2].(string)
	userPasswordSum256 := sha256.Sum256([]byte(userPassword + define.UserPasswordSlat))

	loginResponse, err := SendAndGetHttpResponse[LoginResponse](
		fmt.Sprintf("%s/register_or_login", define.StdAuthServerAddress),
		LoginRequest{
			IsRegister:         false,
			UserName:           userName,
			UserPasswordSum256: userPasswordSum256[:],
		},
	)
	if err != nil {
		return fmt.Errorf("login: %v", err)
	}

	if !loginResponse.Success {
		_, _, err = f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "登录失败",
			Content: loginResponse.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("login: %v", err)
		}
		return nil
	}

	f.config.EulogistToken = loginResponse.EulogistToken
	err = define.WriteConfig(f.config)
	if err != nil {
		return fmt.Errorf("login: %v", err)
	}

	return nil
}

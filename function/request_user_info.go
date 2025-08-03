package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

const (
	RequestTypeGetUserInfoNormal uint8 = iota
	RequestTypeGetUserInfoAdmin
)

// UserInfoRequest ..
type UserInfoRequest struct {
	Token            string `json:"token,omitempty"`
	RequestType      uint8  `json:"request_type"`
	EulogistUserName string `json:"eulogist_user_name,omitempty"`
}

// UserInfoResponse ..
type UserInfoResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
	Payload   []byte `json:"payload"`
}

// RequestUserInfo 请求用户的赞颂者账户数据。
// 如果 isReGet 为假且之前已经成功得到了数据，
// 则不进行任何操作
func (f *Function) RequestUserInfo(isReGet bool) error {
	for {
		if !isReGet && f.userData != nil {
			return nil
		}

		userInfoResponse, err := SendAndGetHttpResponse[UserInfoResponse](
			fmt.Sprintf("%s/request_user_info", define.StdAuthServerAddress),
			UserInfoRequest{
				Token:       f.config.EulogistToken,
				RequestType: RequestTypeGetUserInfoNormal,
			},
		)
		if err != nil {
			return fmt.Errorf("RequestUserInfo: %v", err)
		}

		if !userInfoResponse.Success {
			resp, isUserCacel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title: "错误",
				Content: fmt.Sprintf(
					"请求赞颂者信息时失败 (%s), 可能是赞颂者令牌已经过期, 要重新登录吗？",
					userInfoResponse.ErrorInfo,
				),
				Button1: "重新登录",
				Button2: "重试",
			})
			if err != nil {
				return fmt.Errorf("RequestUserInfo: %v", err)
			}
			if !isUserCacel && resp.(bool) {
				err = f.login()
				if err != nil {
					return fmt.Errorf("RequestUserInfo: %v", err)
				}
			}
			continue
		}

		f.userData = new(define.EulogistUser)
		*f.userData = define.DecodeEulogistUser(userInfoResponse.Payload)
		return nil
	}
}

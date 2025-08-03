package function

import (
	"fmt"
	"strconv"
	"time"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// ChangeMainConfigRequest ..
type ChangeMainConfigRequest struct {
	Token            string `json:"token,omitempty"`
	EulogistUserName string `json:"eulogist_user_name,omitempty"`
	NewUserData      []byte `json:"new_user_data"`
}

// ChangeMainConfigResponse ..
type ChangeMainConfigResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// AdminChangeMainConfig ..
func (f *Function) AdminChangeMainConfig(eulogistUserName string) error {
	userInfoResponse, err := SendAndGetHttpResponse[UserInfoResponse](
		fmt.Sprintf("%s/request_user_info", define.StdAuthServerAddress),
		UserInfoRequest{
			Token:            f.config.EulogistToken,
			RequestType:      RequestTypeGetUserInfoAdmin,
			EulogistUserName: eulogistUserName,
		},
	)
	if err != nil {
		return fmt.Errorf("AdminChangeMainConfig: %v", err)
	}

	if !userInfoResponse.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: userInfoResponse.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("AdminChangeMainConfig: %v", err)
		}
		return nil
	}
	userData := define.DecodeEulogistUser(userInfoResponse.Payload)

	minecraftForm := form.ModalForm{
		Title: "修改用户信息",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "以下列出了该用户的信息, 标有§r§e仅复制§r的仅限用于复制, 修改它们是无效的。",
			},
			form.ModalFormElementInput{
				Text:        "用户唯一标识符 (仅复制)",
				Default:     userData.UserUniqueID,
				PlaceHolder: "UserUniqueID",
			},
			form.ModalFormElementInput{
				Text:        "用户名",
				Default:     userData.UserName,
				PlaceHolder: "UserName",
			},
			form.ModalFormElementDropdown{
				Text:    "权限等级",
				Options: []string{"系统", "赞颂者管理人员", "租赁服管理人员", "普通用户", "访客"},
				Default: uint(userData.UserPermissionLevel),
			},
			form.ModalFormElementInput{
				Text: fmt.Sprintf(
					"解封时间戳 (封禁至 %s)",
					time.Unix(userData.UnbanUnixTime, 0).Format("2006-01-02 15:04:05"),
				),
				Default:     fmt.Sprintf("%d", userData.UnbanUnixTime),
				PlaceHolder: "UnbanUnixTime",
			},
			form.ModalFormElementToggle{
				Text:    "允许无权限进入任意租赁服",
				Default: userData.DisableGlobalOpertorVerify,
			},
			form.ModalFormElementToggle{
				Text:    "可以使用内置验证服务账户进入任意租赁服",
				Default: userData.CanAccessAnyRentalServer,
			},
			form.ModalFormElementToggle{
				Text:    "可以下载任何存档的解密密钥",
				Default: userData.CanGetGameSavesKeyCipher,
			},
			form.ModalFormElementToggle{
				Text:    "可以下载 Token 用于机器人进服",
				Default: userData.CanGetHelperToken,
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("AdminChangeMainConfig: %v", err)
	}
	if isUserCancel {
		return nil
	}

	respList := resp.([]any)
	userName := respList[2].(string)
	userPermissionLevel := uint8(respList[3].(int32))
	unbanUnixTime, err := strconv.ParseInt(respList[4].(string), 10, 64)
	if err != nil {
		unbanUnixTime = userData.UnbanUnixTime
	}
	disableGlobalOpertorVerify := respList[5].(bool)
	canAccessAnyRentalServer := respList[6].(bool)
	canGetGameSavesKeyCipher := respList[7].(bool)
	canGetHelperToken := respList[8].(bool)

	changeMainConfigResp, err := SendAndGetHttpResponse[ChangeMainConfigResponse](
		fmt.Sprintf("%s/admin_change_main_config", define.StdAuthServerAddress),
		ChangeMainConfigRequest{
			Token:            f.config.EulogistToken,
			EulogistUserName: eulogistUserName,
			NewUserData: define.EncodeEulogistUser(define.EulogistUser{
				UserName:                   userName,
				UserPermissionLevel:        userPermissionLevel,
				UnbanUnixTime:              unbanUnixTime,
				DisableGlobalOpertorVerify: disableGlobalOpertorVerify,
				CanAccessAnyRentalServer:   canAccessAnyRentalServer,
				CanGetGameSavesKeyCipher:   canGetGameSavesKeyCipher,
				CanGetHelperToken:          canGetHelperToken,
			}),
		},
	)
	if err != nil {
		return fmt.Errorf("AdminChangeMainConfig: %v", err)
	}

	if !changeMainConfigResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: changeMainConfigResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("AdminChangeMainConfig: %v", err)
		}
		return nil
	}

	if eulogistUserName == f.userData.UserName {
		if err = f.RequestUserInfo(true); err != nil {
			return fmt.Errorf("AdminChangeMainConfig: %v", err)
		}
	}
	return nil
}

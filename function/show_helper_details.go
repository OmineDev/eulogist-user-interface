package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// HelperInfoRequest ..
type HelperInfoRequest struct {
	Token string `json:"token,omitempty"`
}

// UserInfoResponse ..
type HelperInfoResponse struct {
	ErrorInfo            string `json:"error_info"`
	NetEaseRequireVerify bool   `json:"netease_require_verify"`
	VerifyURL            string `json:"verify_url"`
	Success              bool   `json:"success"`
	GameNickName         string `json:"game_nick_name"`
	G79UserUID           string `json:"g79_user_uid"`
}

// ShowHelperDetails 向客户端展示其当前正在使用的 MC 账户的信息。
// 即便该用户当前没有正在使用的 MC 账号，调用 ShowHelperDetails 也是安全的
func (f *Function) ShowHelperDetails() error {
	for {
		account, ok := f.userData.CurrentAuthServerAccount.Value()
		if !ok {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "账号详细信息",
				Content: "您目前没有正在使用的 MC 账号, 请§r§e回到上一级菜单§r选择一个 MC 账号作为当前使用的 MC 账号。",
				Button1: "确定",
				Button2: "回到上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("ShowHelperDetails: %v", err)
			}
			return nil
		}

		if !account.IsStdAccount() {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title: "账号详细信息",
				Content: fmt.Sprintf(
					"● 显示名称: %s\n●账户类型: 第三方验证服务账户\n● 验证服务地址: %s\n● 验证服务令牌: %s\n",
					account.FormatInGame(),
					account.AuthServerAddress(),
					account.AuthServerSecret(),
				),
				Button1: "确定",
				Button2: "取消",
			})
			if err != nil {
				return fmt.Errorf("ShowHelperDetails: %v", err)
			}
			return nil
		}

		helpInfoResponse, err := SendAndGetHttpResponse[HelperInfoResponse](
			fmt.Sprintf("%s/get_std_helper_info", define.StdAuthServerAddress),
			HelperInfoRequest{
				Token: f.config.EulogistToken,
			},
		)
		if err != nil {
			return fmt.Errorf("ShowHelperDetails: %v", err)
		}

		if !helpInfoResponse.Success {
			err = f.ShowAuthServerError(
				helpInfoResponse.NetEaseRequireVerify,
				helpInfoResponse.VerifyURL,
				helpInfoResponse.ErrorInfo,
			)
			if err != nil {
				return fmt.Errorf("ShowHelperDetails: %v", err)
			}
			if helpInfoResponse.NetEaseRequireVerify {
				continue
			} else {
				return nil
			}
		}

		account.UpdateData(map[string]any{
			"gameNickName":       helpInfoResponse.GameNickName,
			"g79UserUID":         helpInfoResponse.G79UserUID,
			"authHelperUniqueID": account.AuthServerSecret(),
		})
		f.userData.CurrentAuthServerAccount = protocol.Option(account)

		for index, value := range f.userData.MultipleAuthServerAccounts {
			if !value.IsStdAccount() {
				continue
			}
			if value.AuthServerSecret() == account.AuthServerSecret() {
				f.userData.MultipleAuthServerAccounts[index] = account
			}
		}

		_, _, err = f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title: "账号详细信息",
			Content: fmt.Sprintf(
				"● 游戏昵称: %s\n●账户类型: 内置验证服务账户",
				account.FormatInGame(),
			),
			Button1: "确定",
			Button2: "取消",
		})
		if err != nil {
			return fmt.Errorf("ShowHelperDetails: %v", err)
		}
		return nil
	}
}

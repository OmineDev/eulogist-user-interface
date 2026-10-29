package function

import (
	"fmt"
	"time"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
	"github.com/YingLunTown-DreamLand/gophertunnel/minecraft/protocol"
)

const (
	AccountTypeMpayUser uint8 = iota
	AccountTypePeAuth
	AccountTypeSaAuth
)

// HelperInfoRequest ..
type HelperInfoRequest struct {
	Token string `json:"token,omitempty"`
}

// HelperInfoResponse ..
type HelperInfoResponse struct {
	ErrorInfo         string `json:"error_info"`
	Success           bool   `json:"success"`
	AccountType       uint8  `json:"account_type"`
	AccountExpireTime int64  `json:"account_expire_time"`
	GameNickName      string `json:"game_nick_name"`
	G79UserUID        string `json:"g79_user_uid"`
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
					"● 显示名称: %s\n● 账户类型: 第三方验证服务账户\n● 验证服务地址: %s\n● 验证服务令牌: %s\n",
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

		helpInfoResponse, err := utils.SendAndGetHttpResponse[HelperInfoResponse](
			fmt.Sprintf("%s/get_std_helper_info", define.AddressEulogistAPI),
			HelperInfoRequest{
				Token: f.config.EulogistToken,
			},
		)
		if err != nil {
			return fmt.Errorf("ShowHelperDetails: %v", err)
		}

		if !helpInfoResponse.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: helpInfoResponse.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("ShowHelperDetails: %v", err)
			}
			return nil
		}

		if helpInfoResponse.AccountType == AccountTypeMpayUser {
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
		}

		content := ""
		switch helpInfoResponse.AccountType {
		case AccountTypeMpayUser:
			content = fmt.Sprintf(
				""+
					"● 游戏昵称: %s\n"+
					"● 账户类型: 内置验证服务账户\n"+
					"● 会话过期时间: §r§b%s§r",
				account.FormatInGame(),
				time.Unix(helpInfoResponse.AccountExpireTime, 0).Format(time.DateTime),
			)
		case AccountTypePeAuth:
			content = fmt.Sprintf(
				""+
					"● 游戏昵称: %s\n"+
					"● 账户类型: 内置验证服务账户 (Pe Auth)\n"+
					"● 会话过期时间: §r§b%s§r",
				fmt.Sprintf(
					"§r§l§e%s §r§l(§b%s§r§l)§r",
					helpInfoResponse.GameNickName,
					helpInfoResponse.G79UserUID,
				),
				time.Unix(helpInfoResponse.AccountExpireTime, 0).Format(time.DateTime),
			)
		case AccountTypeSaAuth:
			content = fmt.Sprintf(
				""+
					"● 游戏昵称: %s\n"+
					"● 账户类型: 内置验证服务账户 (Sa Auth)\n"+
					"● 会话过期时间: §r§b%s§r",
				fmt.Sprintf(
					"§r§l§e%s §r§l(§b%s§r§l)§r",
					helpInfoResponse.GameNickName,
					helpInfoResponse.G79UserUID,
				),
				time.Unix(helpInfoResponse.AccountExpireTime, 0).Format(time.DateTime),
			)
		}

		_, _, err = f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "账号详细信息",
			Content: content,
			Button1: "确定",
			Button2: "取消",
		})
		if err != nil {
			return fmt.Errorf("ShowHelperDetails: %v", err)
		}
		return nil
	}
}

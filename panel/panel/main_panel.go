package panel

import (
	"fmt"
	"time"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// MainPanel 向用户展示赞颂者的控制面板
func (p *Panel) MainPanel() (exitGame bool, err error) {
	conn := p.f.Interact().Server().MinecraftConn()

	for {
		// Register or login
		for {
			if len(p.f.EulogistConfig().EulogistToken) != 0 {
				break
			}

			exitGame, err := p.f.RegisterOrLogin()
			select {
			case <-conn.Context().Done():
				return false, nil
			default:
			}

			if err != nil {
				_, _, _ = p.f.Interact().SendFormAndWaitResponse(form.MessageForm{
					Title:   "注册或登录失败",
					Content: fmt.Sprintf("%v", err),
					Button1: "确定",
					Button2: "继续",
				})
				continue
			}

			if exitGame {
				return true, nil
			}
		}

		// Request user info
		for {
			if err = p.f.RequestUserInfo(false); err == nil {
				break
			}

			select {
			case <-conn.Context().Done():
				return false, nil
			default:
			}

			_, _, _ = p.f.Interact().SendFormAndWaitResponse(form.MessageForm{
				Title:   "请求赞颂者账户信息失败",
				Content: fmt.Sprintf("%v", err),
				Button1: "确定",
				Button2: "继续",
			})
		}

		// Show main panel
		minecraftForm := form.ActionForm{
			Title:   "欢迎您",
			Content: fmt.Sprintf("%s请选择您要进行的操作。", p.formatUserData()),
			Buttons: []form.ActionFormElement{
				{
					Text: "游玩租赁服",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "管理您自己的网易 MC 账号",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "编辑您自己的赞颂者账号信息",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "下载您存档的解密密钥",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "进入服主管理面板",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "进入赞颂者管理面板",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "退出登录",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "退出游戏",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		resp, isUserCancel, err := p.f.Interact().SendFormAndWaitResponse(minecraftForm)
		if err != nil {
			return false, fmt.Errorf("MainPanel: %v", err)
		}
		if isUserCancel {
			continue
		}

		switch resp.(int32) {
		case 0:
			exitGame, err = p.StartPlayPanel()
		case 1:
			err = p.AuthHelperPanel()
		case 2:
			err = p.ChangeUserInfo()
		case 3:
			err = p.f.GetGameSavesKey()
		case 4:
			err = p.ServerOwnerPanel()
		case 5:
			err = p.AdminPanel()
		case 6:
			err = p.f.Logout()
		case 7:
			exitGame = true
		}
		if err != nil {
			return false, fmt.Errorf("MainPanel: %v", err)
		}

		if exitGame {
			return true, nil
		}
	}
}

// formatUserData ..
func (p *Panel) formatUserData() (content string) {
	userData := p.f.EulogistUserData()
	content = fmt.Sprintf(
		"基本用户信息\n  - 用户名: §r§e%s§r\n",
		userData.UserName,
	)

	switch userData.UserPermissionLevel {
	case define.UserPermissionSystem:
		content += "  - 用户权限: §r§a系统§r\n"
	case define.UserPermissionAdmin:
		content += "  - 用户权限: §r§a管理员§r\n"
	case define.UserPermissionManager:
		content += "  - 用户权限: §r§e租赁服管理员§r\n"
	case define.UserPermissionNormal:
		content += "  - 用户权限: §r§r普通用户§r\n"
	case define.UserPermissionNone:
		content += "  - 用户权限: §r§c访客§r\n"
	default:
		content += fmt.Sprintf("  - 用户权限: 未知 §r(§b%d§r)\n", userData.UserPermissionLevel)
	}

	currentTime := time.Now()
	if currentTime.Unix() < userData.UnbanUnixTime {
		totalSeconds := int64(time.Unix(userData.UnbanUnixTime, 0).Sub(currentTime).Seconds())
		days := totalSeconds / 86400
		hours := totalSeconds/3600 - days*24
		minutes := totalSeconds/60 - hours*60 - days*1440
		content += fmt.Sprintf(
			"  - 封禁状态: §r§c正被封禁 §r(还剩 §b%d §r天 §b%d §r时 §b%d §r分)§r\n",
			days, hours, minutes,
		)
	} else {
		content += "  - 封禁状态: §r§a未封禁§r\n"
	}

	if len(userData.ProvidedPeAuthData) > 0 {
		content += "  - PE Auth: §r§a正在使用§r\n"
	} else {
		content += "  - PE Auth: §r§7未使用§r\n"
	}

	if userData.DisableGlobalOpertorVerify {
		content += "  - 无权限进入任何租赁服: §r§a已授权§r\n"
	} else {
		content += "  - 无权限进入任何租赁服: §r§c未授权§r\n"
	}

	if userData.CanAccessAnyRentalServer {
		content += "  - 可访问任意租赁服: §r§a已授权§r\n"
	} else {
		content += "  - 可访问任意租赁服: §r§c未授权§r\n"
	}

	if userData.CanGetGameSavesKeyCipher {
		content += "  - 可取得任意租赁服的存档密钥: §r§a已授权§r\n"
	} else {
		content += "  - 可取得任意租赁服的存档密钥: §r§c未授权§r\n"
	}

	if userData.CanGetHelperToken {
		content += "  - 可取得辅助用户令牌: §r§a已授权§r\n"
	} else {
		content += "  - 可取得辅助用户令牌: §r§c未授权§r\n"
	}

	return
}

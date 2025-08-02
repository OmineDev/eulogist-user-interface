package function

import (
	"fmt"
	"time"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// ShowPanel 向用户展示赞颂者的控制面板。
// 在调用它前必须先调用 RequestUserInfo 获取用户信息
func (f *Function) ShowPanel() (exitGame bool, err error) {
	for {
		if f.userData == nil {
			return false, fmt.Errorf("ShowPanel: Needs call RequestUserInfo first")
		}

		minecraftForm := form.ActionForm{
			Title:   "欢迎您",
			Content: fmt.Sprintf("%s请选择您要进行的操作。", f.formatUserData()),
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
					Text: "进入服主管理面板",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "进入赞颂者管理面板",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "退出游戏",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		f.interact.SendFormAndWaitResponse(minecraftForm)
	}
}

// formatUserData ..
func (f *Function) formatUserData() (content string) {
	content = fmt.Sprintf(
		"基本用户信息\n  - 用户名: §r§e%s§r\n",
		f.userData.UserName,
	)

	switch f.userData.UserPermissionLevel {
	case define.UserPermissionSystem:
		content += "  - 用户权限: §r§a系统§r\n"
	case define.UserPermissionAdmin:
		content += "  - 用户权限: §r§a管理员§r\n"
	case define.UserPermissionManager:
		content += "  - 用户权限: §r§e租赁服管理员§r\n"
	case define.UserPermissionNormal:
		content += "  - 用户权限: §r§f普通用户§r\n"
	case define.UserPermissionNone:
		content += "  - 用户权限: §r§c访客§r\n"
	default:
		content += fmt.Sprintf("  - 用户权限: 未知 §r(§b%d§r)\n", f.userData.UserPermissionLevel)
	}

	currentTime := time.Now()
	if currentTime.Unix() < f.userData.UnbanUnixTime {
		totalSeconds := int64(time.Unix(f.userData.UnbanUnixTime, 0).Sub(currentTime).Seconds())
		days := totalSeconds / 86400
		hours := totalSeconds/3600 - days*24
		minutes := totalSeconds/60 - hours*60 - days*1440
		content += fmt.Sprintf(
			"  - 封禁状态: §r§c正被封禁 §f(还剩 §b%d §f天 §b%d §f时 §b%d §f秒)§r\n",
			days, hours, minutes,
		)
	} else {
		content += "  - 封禁状态: §r§a未封禁§r\n"
	}

	if len(f.userData.ProvidedPeAuthData) > 0 {
		content += "  - PE Auth: §r§a正在使用§r\n"
	} else {
		content += "  - PE-Auth: §r§7未使用§r\n"
	}

	if f.userData.DisableGlobalOpertorVerify {
		content += "  - 无权限进入任何租赁服: §r§a已授权§r\n"
	} else {
		content += "  - 无权限进入任何租赁服: §r§c未授权§r\n"
	}

	if f.userData.DisableGlobalOpertorVerify {
		content += "  - 可访问任意租赁服: §r§a已授权§r\n"
	} else {
		content += "  - 可访问任意租赁服: §r§c未授权§r\n"
	}

	if f.userData.CanGetGameSavesKeyCipher {
		content += "  - 可取得任意租赁服的存档密钥: §r§a已授权§r\n"
	} else {
		content += "  - 可取得任意租赁服的存档密钥: §r§c未授权§r\n"
	}

	if f.userData.CanGetHelperToken {
		content += "  - 可取得辅助用户令牌: §r§a已授权§r\n"
	} else {
		content += "  - 可取得辅助用户令牌: §r§c未授权§r\n"
	}

	return
}

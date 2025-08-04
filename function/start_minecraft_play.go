package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// SelectRentalServerAndStartPlay ..
func (f *Function) SelectRentalServerAndStartPlay() (
	backedMessage *MessageFromBacked,
	err error,
) {
	for {
		if len(f.userData.RentalServerConfig) != 0 {
			break
		}

		resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "询问",
			Content: "您目前没有添加任何租赁服, 要先添加租赁服吗？",
			Button1: "确定",
			Button2: "忽略",
		})
		if err != nil {
			return nil, fmt.Errorf("SelectRentalServerAndStartPlay: %v", err)
		}
		if isUserCancel {
			return nil, nil
		}
		if !resp.(bool) {
			break
		}

		err = f.AddNewRentalServerConfig()
		if err != nil {
			return nil, fmt.Errorf("SelectRentalServerAndStartPlay: %v", err)
		}
	}

	actionForm := form.ActionForm{
		Title: "选择租赁服进行游玩",
	}
	if len(f.userData.RentalServerConfig) == 0 {
		actionForm.Content = "您目前没有添加任何租赁服, 请§r§e回退§r到上一级菜单先进行添加。"
	} else {
		actionForm.Content = fmt.Sprintf(
			"您已添加 §r§b%d§r 个租赁服, 请§r§e选择§r一个进入游玩。",
			len(f.userData.RentalServerConfig),
		)
	}
	for _, value := range f.userData.RentalServerConfig {
		actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
			Text: value.ServerNumber,
			Icon: form.ActionFormIconNone{},
		})
	}

	resp, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
	if err != nil {
		return nil, fmt.Errorf("SelectRentalServerAndStartPlay: %v", err)
	}
	if isUserCancel {
		return nil, nil
	}

	backedMessage, err = f.startMinecraftPlay(f.userData.RentalServerConfig[resp])
	if err != nil {
		return nil, fmt.Errorf("SelectRentalServerAndStartPlay: %v", err)
	}

	return backedMessage, nil
}

// StartPlayByLastConfig ..
func (f *Function) StartPlayByLastConfig() (
	backedMessage *MessageFromBacked,
	err error,
) {
	if len(f.config.LastServerConfig.ServerNumber) == 0 {
		resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "通知",
			Content: "您目前没有进入过任何租赁服, 因此您无法按照上次的配置进服。要跳转到租赁服配置选择界面选择一个租赁服游玩吗？",
			Button1: "确定 (推荐)",
			Button2: "取消",
		})
		if err != nil {
			return nil, fmt.Errorf("StartPlayByLastConfig: %v", err)
		}
		if isUserCancel || !resp.(bool) {
			return nil, nil
		}

		backedMessage, err = f.SelectRentalServerAndStartPlay()
		if err != nil {
			return nil, fmt.Errorf("StartPlayByLastConfig: %v", err)
		}

		return backedMessage, nil
	}

	backedMessage, err = f.startMinecraftPlay(f.config.LastServerConfig)
	if err != nil {
		return nil, fmt.Errorf("StartPlayByLastConfig: %v", err)
	}

	return backedMessage, nil
}

// startMinecraftPlay ..
func (f *Function) startMinecraftPlay(config define.RentalServerConfig) (
	backedMessage *MessageFromBacked,
	err error,
) {
	var userIsSettingAccount bool

	if f.config.LastServerConfig != config {
		f.config.LastServerConfig = config
		err = define.WriteConfig(f.config)
		if err != nil {
			return nil, fmt.Errorf("startMinecraftPlay: %v", err)
		}
	}

	for {
		_, ok := f.userData.CurrentAuthServerAccount.Value()
		if ok && userIsSettingAccount {
			resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "通知",
				Content: "您已完成 MC 账号的设置, 您即将开始进入租赁服。",
				Button1: "我知道了",
				Button2: "取消",
			})
			if err != nil {
				return nil, fmt.Errorf("startMinecraftPlay: %v", err)
			}
			if isUserCancel || !resp.(bool) {
				return nil, nil
			}
		}
		if ok {
			break
		}
		userIsSettingAccount = true

		resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "询问",
			Content: "您目前没有正在使用的 MC 账户, 要先设置 MC 账户吗？",
			Button1: "确定",
			Button2: "取消",
		})
		if err != nil {
			return nil, fmt.Errorf("startMinecraftPlay: %v", err)
		}
		if isUserCancel || !resp.(bool) {
			_, _, err = f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "询问",
				Content: "您目前没有正在使用的 MC 账户, 因此, 在您完成配置前无法进服。",
				Button1: "我知道了",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return nil, fmt.Errorf("startMinecraftPlay: %v", err)
			}
			return nil, nil
		}

		resp, isUserCancel, err = f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "通知",
			Content: "您将跳转到 MC 账号的添加页面。",
			Button1: "继续",
			Button2: "取消",
		})
		if err != nil {
			return nil, fmt.Errorf("startMinecraftPlay: %v", err)
		}
		if isUserCancel || !resp.(bool) {
			continue
		}

		err = f.AddNewHelper()
		if err != nil {
			return nil, fmt.Errorf("startMinecraftPlay: %v", err)
		}
		if _, ok = f.userData.CurrentAuthServerAccount.Value(); ok {
			continue
		}
		if len(f.userData.MultipleAuthServerAccounts) == 0 {
			continue
		}

		resp, isUserCancel, err = f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "通知",
			Content: "您已添加了一些 MC 账号, 但没有设置位您目前所使用的账户, 是否要跳转到 MC 账户的设置页面？",
			Button1: "是",
			Button2: "否",
		})
		if err != nil {
			return nil, fmt.Errorf("startMinecraftPlay: %v", err)
		}
		if isUserCancel || !resp.(bool) {
			continue
		}

		err = f.ChangeCurrentHelper()
		if err != nil {
			return nil, fmt.Errorf("startMinecraftPlay: %v", err)
		}
	}

	aesCipher, err := f.BeforePlayPrepare(config.ServerNumber)
	if err != nil {
		return nil, fmt.Errorf("startMinecraftPlay: %v", err)
	}

	account, ok := f.userData.CurrentAuthServerAccount.Value()
	if !ok {
		return nil, fmt.Errorf("startMinecraftPlay: 已设置的 MC 账号未能找到")
	}

	frontedMsg := MessageFromFronted{
		RentalServerNumber:   config.ServerNumber,
		RentalServerPasscode: config.ServerPassCode,
		AuthServerAddress:    account.AuthServerAddress(),
		ProvidedPeAuthData:   f.userData.ProvidedPeAuthData,
		EulogistUniqueID:     f.userData.UserUniqueID,
		GameSavesAESCipher:   aesCipher,
	}
	if account.IsStdAccount() {
		frontedMsg.AuthServerToken = fmt.Sprintf("%s|%s", f.config.EulogistToken, account.AuthServerSecret())
	} else {
		frontedMsg.AuthServerToken = account.AuthServerSecret()
	}
	f.message.NotifyToBacked(frontedMsg)

	backedMessage, err = f.waitbackedMessage(config.ServerNumber)
	if err != nil {
		return nil, fmt.Errorf("startMinecraftPlay: %v", err)
	}

	return backedMessage, nil
}

// waitbackedMessage ..
func (f *Function) waitbackedMessage(rentalServerNumber string) (
	backedMessage *MessageFromBacked,
	err error,
) {
	backedMsgIncoming := make(chan struct{})

	go func() {
		for {
			ctx, closer := f.interact.SendFormOmitResponse(form.MessageForm{
				Title: "进行中...",
				Content: fmt.Sprintf(
					"您正在登陆到租赁服 §r§b%s§r, 此操作不可取消, 请耐心等待...",
					rentalServerNumber,
				),
				Button1: "确定",
				Button2: "继续",
			})
			select {
			case <-ctx.Done():
			case <-backedMsgIncoming:
				closer()
				return
			}
		}
	}()

	msg := <-f.message.BackedMessageChannel()
	backedMessage = &msg
	close(backedMsgIncoming)

	if backedMessage.CanTerminate {
		_ = f.interact.Server().MinecraftConn().WritePacket(&packet.Transfer{
			Address: msg.TransferAddress,
			Port:    msg.TransferPort,
		})
		_ = f.interact.Server().CloseServer()
	}

	return backedMessage, nil
}

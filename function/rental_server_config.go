package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

const (
	ActionTypeEditRentalServer uint8 = iota
	ActionTypeRemoveRentalServer
)

// RentalServerListRequest ..
type RentalServerListRequest struct {
	Token          string `json:"token,omitempty"`
	ActionType     uint8  `json:"action_type"`
	ServerNumber   string `json:"server_number,omitempty"`
	ServerPassCode string `json:"server_passcode,omitempty"`
}

// RentalServerListResponse ..
type RentalServerListResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// sendRentalServerConfigRequest ..
func (f *Function) sendRentalServerConfigRequest(
	actionType uint8,
	serverNumber string,
	serverPasscode string,
) error {
	rentalServerListResp, err := utils.SendAndGetHttpResponse[RentalServerListResponse](
		fmt.Sprintf("%s/rental_server_list", define.StdAuthServerAddress),
		RentalServerListRequest{
			Token:          f.config.EulogistToken,
			ActionType:     actionType,
			ServerNumber:   serverNumber,
			ServerPassCode: serverPasscode,
		},
	)
	if err != nil {
		return fmt.Errorf("sendRentalServerConfigRequest: %v", err)
	}

	if !rentalServerListResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: rentalServerListResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("sendRentalServerConfigRequest: %v", err)
		}
		return nil
	}

	if actionType != ActionTypeRemoveRentalServer {
		isRepeat := false
		for index, value := range f.userData.RentalServerConfig {
			if value.ServerNumber == serverNumber {
				f.userData.RentalServerConfig[index].ServerPassCode = serverPasscode
				isRepeat = true
				break
			}
		}
		if !isRepeat {
			f.userData.RentalServerConfig = append(f.userData.RentalServerConfig, define.RentalServerConfig{
				ServerNumber:   serverNumber,
				ServerPassCode: serverPasscode,
			})
		}
	} else {
		newList := make([]define.RentalServerConfig, 0)
		for _, value := range f.userData.RentalServerConfig {
			if value.ServerNumber != serverNumber {
				newList = append(newList, value)
			}
		}
		f.userData.RentalServerConfig = newList
	}

	return nil
}

// AddNewRentalServerConfig ..
func (f *Function) AddNewRentalServerConfig() error {
	minecraftForm := form.ModalForm{
		Title: "添加新租赁服配置",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "" +
					"您正在添加新的§r§e租赁服配置§r。\n" +
					"添加完成后, 您将回退到上一级菜单, \n" +
					"然后您可以通过§r§e选择租赁服进行游玩§r开始游玩",
			},
			form.ModalFormElementInput{
				Text:        "请输入要添加的租赁服号",
				Default:     "",
				PlaceHolder: "非空的租赁服号码",
			},
			form.ModalFormElementInput{
				Text:        "请输入被添加租赁服的密码",
				Default:     "",
				PlaceHolder: "如果租赁服没有密码, 请置空",
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("AddNewRentalServerConfig: %v", err)
	}
	if isUserCancel {
		return nil
	}

	respList := resp.([]any)
	serverNumber := respList[1].(string)
	serverPasscode := respList[2].(string)

	err = f.sendRentalServerConfigRequest(ActionTypeEditRentalServer, serverNumber, serverPasscode)
	if err != nil {
		return fmt.Errorf("AddNewRentalServerConfig: %v", err)
	}

	return nil
}

func (f *Function) ModifyRentalServerConfig() error {
	actionForm := form.ActionForm{
		Title: "修改租赁服配置",
	}
	if len(f.userData.RentalServerConfig) == 0 {
		actionForm.Content = "您目前没有添加任何租赁服, 请先§r§e回退到上一级菜单§r添加。"
	} else {
		actionForm.Content = fmt.Sprintf(
			"您已添加 §r§b%d§r 个租赁服, 请选择一个§r§e修改§r其配置信息。",
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
		return fmt.Errorf("ModifyRentalServerConfig: %v", err)
	}
	if isUserCancel {
		return nil
	}
	config := f.userData.RentalServerConfig[resp]

	minecraftForm := form.ModalForm{
		Title: fmt.Sprintf("修改 %s 的配置", config.ServerNumber),
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: fmt.Sprintf(
					"您正在修改租赁服 §r§b%s§r 的配置。\n请输入该租赁服的§r§e新密码§r。",
					config.ServerNumber,
				),
			},
			form.ModalFormElementInput{
				Text:        "新密码",
				Default:     config.ServerPassCode,
				PlaceHolder: "如果租赁服没有密码, 请置空",
			},
		},
	}

	answer, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("ModifyRentalServerConfig: %v", err)
	}
	if isUserCancel {
		return nil
	}

	answerList := answer.([]any)
	serverPasscode := answerList[1].(string)

	err = f.sendRentalServerConfigRequest(ActionTypeEditRentalServer, config.ServerNumber, serverPasscode)
	if err != nil {
		return fmt.Errorf("ModifyRentalServerConfig: %v", err)
	}

	if f.config.LastServerConfig.ServerNumber == config.ServerNumber {
		f.config.LastServerConfig.ServerPassCode = serverPasscode
		if err = define.WriteConfig(f.config); err != nil {
			return fmt.Errorf("ModifyRentalServerConfig: %v", err)
		}
	}

	return nil
}

// RemoveRentalServerConfig ..
func (f *Function) RemoveRentalServerConfig() error {
	for {
		actionForm := form.ActionForm{
			Title: "删除租赁服配置",
		}
		if len(f.userData.RentalServerConfig) == 0 {
			actionForm.Content = "您目前没有添加任何租赁服, 因此您无法§r§e删除§r任何租赁服配置。"
		} else {
			actionForm.Content = fmt.Sprintf(
				"您已添加 §r§b%d§r 个租赁服, 请选择一个将其§r§e删除§r。",
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
			return fmt.Errorf("RemoveRentalServerConfig: %v", err)
		}
		if isUserCancel {
			return nil
		}
		config := f.userData.RentalServerConfig[resp]

		ensure, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "二重确认",
			Content: "您确定要§r§e删除§r目标租赁服配置吗？",
			Button1: "确定",
			Button2: "取消",
		})
		if err != nil {
			return fmt.Errorf("RemoveRentalServerConfig: %v", err)
		}
		if isUserCancel || !ensure.(bool) {
			continue
		}

		err = f.sendRentalServerConfigRequest(ActionTypeRemoveRentalServer, config.ServerNumber, "")
		if err != nil {
			return fmt.Errorf("RemoveRentalServerConfig: %v", err)
		}

		if f.config.LastServerConfig.ServerNumber == config.ServerNumber {
			f.config.LastServerConfig = define.RentalServerConfig{}
			if err = define.WriteConfig(f.config); err != nil {
				return fmt.Errorf("RemoveRentalServerConfig: %v", err)
			}
		}

		return nil
	}
}

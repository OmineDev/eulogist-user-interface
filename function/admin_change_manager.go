package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

const (
	ActionTypeEditManageServer uint8 = iota
	ActionTypeRemoveManageServer
	ActionTypeRequestServerList
)

// ManagerChangeRequest ..
type ManagerChangeRequest struct {
	Token            string `json:"token,omitempty"`
	ActionType       uint8  `json:"action_type"`
	EulogistUserName string `json:"eulogist_user_name,omitempty"`
	ServerNumber     string `json:"server_number,omitempty"`
}

// ManagerChangeResponse ..
type ManagerChangeResponse struct {
	ErrorInfo       string   `json:"error_info"`
	Success         bool     `json:"success"`
	ServerCanManage []string `json:"server_can_manage"`
}

// AdminChangeManager ..
func (f *Function) AdminChangeManager(eulogistUserName string) error {
	for {
		actionForm := form.ActionForm{
			Title: fmt.Sprintf(
				"配置 %s 可管理的租赁服",
				eulogistUserName,
			),
			Content: "您将§r§e配置§r该用户可以管理的租赁服, 请选择您的操作。",
			Buttons: []form.ActionFormElement{
				{
					Text: "允诺其管理更多的租赁服",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "移除其正在管理的租赁服",
					Icon: form.ActionFormIconNone{},
				},
				{
					Text: "返回上一级菜单",
					Icon: form.ActionFormIconNone{},
				},
			},
		}

		resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(actionForm)
		if err != nil {
			return fmt.Errorf("AdminChangeManager: %v", err)
		}
		if isUserCancel || resp.(int32) == 2 {
			return nil
		}

		managerChangeResp, err := utils.SendAndGetHttpResponse[ManagerChangeResponse](
			fmt.Sprintf("%s/admin_change_manager", define.StdAuthServerAddress),
			ManagerChangeRequest{
				Token:            f.config.EulogistToken,
				ActionType:       ActionTypeRequestServerList,
				EulogistUserName: eulogistUserName,
			},
		)
		if err != nil {
			return fmt.Errorf("AdminChangeManager: %v", err)
		}

		if !managerChangeResp.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: managerChangeResp.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("AdminChangeManager: %v", err)
			}
			return nil
		}

		switch resp.(int32) {
		case 0:
			_ = f.adminChangeManagerEdit(eulogistUserName)
		case 1:
			_ = f.adminChangeManagerRemove(eulogistUserName, managerChangeResp.ServerCanManage)
		}

		if eulogistUserName == f.userData.UserName {
			if err = f.RequestUserInfo(true); err != nil {
				return fmt.Errorf("AdminChangeMainConfig: %v", err)
			}
		}
	}
}

// adminChangeManagerEdit ..
func (f *Function) adminChangeManagerEdit(eulogistUserName string) error {
	minecraftForm := form.ModalForm{
		Title: "下达神谕 (允诺)",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "" +
					"您将授权该用户§r§e管理§r一个新的租赁服, \n" +
					"这意味着他可以控制哪些人可以使用§r§e内置验证服务§r进入他的服务器。\n" +
					"除此外, 该用户还可以调整特定于该租赁服的§r§e无权限进服§r和§r§e存档密钥§r设置。",
			},
			form.ModalFormElementInput{
				Text:        "请输入租赁服号",
				Default:     "",
				PlaceHolder: "非空的租赁服号码",
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("adminChangeManagerEdit: %v", err)
	}
	if isUserCancel {
		return nil
	}

	managerChangeResp, err := utils.SendAndGetHttpResponse[ManagerChangeResponse](
		fmt.Sprintf("%s/admin_change_manager", define.StdAuthServerAddress),
		ManagerChangeRequest{
			Token:            f.config.EulogistToken,
			ActionType:       ActionTypeEditManageServer,
			EulogistUserName: eulogistUserName,
			ServerNumber:     resp.([]any)[1].(string),
		},
	)
	if err != nil {
		return fmt.Errorf("adminChangeManagerEdit: %v", err)
	}

	if !managerChangeResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: managerChangeResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("adminChangeManagerEdit: %v", err)
		}
		return nil
	}

	return nil
}

// adminChangeManagerRemove ..
func (f *Function) adminChangeManagerRemove(
	eulogistUserName string,
	serverCanManage []string,
) error {
	for {
		actionForm := form.ActionForm{
			Title:   "下达神谕 (剥夺)",
			Content: "您将移除该用户§r§e正在管理§r的一个租赁服, 请从下方选择一个租赁服。",
		}
		for _, value := range serverCanManage {
			actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
				Text: value,
				Icon: form.ActionFormIconNone{},
			})
		}

		answer, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
		if err != nil {
			return fmt.Errorf("adminChangeManagerRemove: %v", err)
		}
		if isUserCancel {
			return nil
		}

		ensure, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "二重确认",
			Content: "您确定要下达§r§c剥夺§r神谕吗？",
			Button1: "确定",
			Button2: "取消",
		})
		if err != nil {
			return fmt.Errorf("adminChangeManagerRemove: %v", err)
		}
		if isUserCancel || !ensure.(bool) {
			continue
		}

		managerChangeResp, err := utils.SendAndGetHttpResponse[ManagerChangeResponse](
			fmt.Sprintf("%s/admin_change_manager", define.StdAuthServerAddress),
			ManagerChangeRequest{
				Token:            f.config.EulogistToken,
				ActionType:       ActionTypeRemoveManageServer,
				EulogistUserName: eulogistUserName,
				ServerNumber:     serverCanManage[answer],
			},
		)
		if err != nil {
			return fmt.Errorf("adminChangeManagerRemove: %v", err)
		}

		if !managerChangeResp.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: managerChangeResp.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("adminChangeManagerRemove: %v", err)
			}
			return nil
		}

		return nil
	}
}

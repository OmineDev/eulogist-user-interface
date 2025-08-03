package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

// AllowListUpdateRequest ..
type AllowListUpdateRequest struct {
	Token                    string `json:"token,omitempty"`
	RentalServerNumber       string `json:"rental_server_number,omitempty"`
	EulogistUserName         string `json:"eulogist_user_name,omitempty"`
	DisableOpertorVerify     bool   `json:"disbale_operator_verify"`
	CanGetGameSavesKeyCipher bool   `json:"can_get_game_saves_key_cipher"`
}

// UserSearchResponse ..
type AllowListUpdateResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// rsmUpdatePermission ..
func (f *Function) rsmUpdatePermission(
	rentalServerNumber string,
	eulogistUserName string,
	isInitProgress bool,
	legacyOPSettings bool,
	legacyCipherSettings bool,
) error {
	minecraftForm := form.ModalForm{}

	if isInitProgress {
		minecraftForm.Title = fmt.Sprintf(
			"纳入并设置 %s 的权限",
			eulogistUserName,
		)
		minecraftForm.Contents = append(minecraftForm.Contents, form.ModalFormElementLabel{
			Text: "" +
				"您正在§r§e纳入§r一个的赞颂者用户到您的租赁服允许列表中。\n\n" +
				"这意味着其如果使用§r§e内置验证服务账户§r, 其将被允许登录到您的租赁服。\n\n" +
				"特别地, 使用第三方验证服务§r§e总是§r可以\n" + "登录到您的租赁服, 而无论您是否特别允许, \n" + "但下方设置依然有效。\n\n" +
				"下面列出了您可以§r§e配置§r的权限设置。",
		})
	} else {
		minecraftForm.Title = fmt.Sprintf(
			"修改 %s 的权限",
			eulogistUserName,
		)
		minecraftForm.Contents = append(minecraftForm.Contents, form.ModalFormElementLabel{
			Text: "" +
				"您正在§r§e修改§r一个的赞颂者用户在您租赁服上的权限。\n" +
				"应当注意的是, 这里指的§r§e不是§r操作员权限, 请注意区别。\n" +
				"下面列出了您可以§r§e配置§r的权限设置。",
		})
	}

	minecraftForm.Contents = append(minecraftForm.Contents, form.ModalFormElementToggle{
		Text:    "允许此用户无操作员权限进服",
		Default: legacyOPSettings,
	})
	minecraftForm.Contents = append(minecraftForm.Contents, form.ModalFormElementLabel{
		Text: "" +
			"下面的设置具有§r§e危险性§r, \n" +
			"这意味着你将授权此用户以某种形式取得您服务器的存档。\n\n" +
			"赞颂者在不久前引入了存档保存功能, \n" +
			"这意味着玩家所经之处将会保存下来。\n\n" +
			"存档总是会保存, §r§e而无论是否有您的授权§r。\n\n" +
			"在打开下方的选项后, 用户将可以取得密钥, \n" +
			"于是该用户将可以解密存档, 而因此得到存档！\n\n" +
			"另外, 对方§r§e一旦§r获得密钥, 您就再也§r§e无法反悔§r！",
	})
	minecraftForm.Contents = append(minecraftForm.Contents, form.ModalFormElementToggle{
		Text:    "允许此用户取得存档解密密钥",
		Default: legacyCipherSettings,
	})

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("rsmUpdatePermission: %v", err)
	}
	if isUserCancel {
		return nil
	}

	respList := resp.([]any)
	disableOpertorVerify := respList[1].(bool)
	canGetGameSavesKeyCipher := respList[3].(bool)

	allowListUpdateResp, err := utils.SendAndGetHttpResponse[AllowListUpdateResponse](
		fmt.Sprintf("%s/update_allow_list_config", define.StdAuthServerAddress),
		AllowListUpdateRequest{
			Token:                    f.config.EulogistToken,
			RentalServerNumber:       rentalServerNumber,
			EulogistUserName:         eulogistUserName,
			DisableOpertorVerify:     disableOpertorVerify,
			CanGetGameSavesKeyCipher: canGetGameSavesKeyCipher,
		},
	)
	if err != nil {
		return fmt.Errorf("rsmUpdatePermission: %v", err)
	}

	if !allowListUpdateResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: allowListUpdateResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("rsmUpdatePermission: %v", err)
		}
		return nil
	}

	return nil
}

// RentalServerManageAddNew ..
func (f *Function) RentalServerManageAddNew(rentalServerNumber string) error {
	eulogistUserName, isUserCancel, err := f.ShowEulogistSearch()
	if err != nil {
		return fmt.Errorf("RentalServerManageAddNew: %v", err)
	}
	if isUserCancel {
		return nil
	}
	err = f.rsmUpdatePermission(rentalServerNumber, eulogistUserName, true, false, false)
	if err != nil {
		return fmt.Errorf("RentalServerManageAddNew: %v", err)
	}
	return nil
}

// AllowListGetRequest ..
type AllowListGetRequest struct {
	Token              string `json:"token,omitempty"`
	RentalServerNumber string `json:"rental_server_number,omitempty"`
}

// AllowListGetResponse ..
type AllowListGetResponse struct {
	ErrorInfo                string   `json:"error_info"`
	Success                  bool     `json:"success"`
	UserNames                []string `json:"user_names"`
	DisableOpertorVerify     []bool   `json:"disbale_operator_verify"`
	CanGetGameSavesKeyCipher []bool   `json:"can_get_game_saves_key_cipher"`
}

// RentalServerManageModify ..
func (f *Function) RentalServerManageModify(rentalServerNumber string) error {
	allowListGetResp, err := utils.SendAndGetHttpResponse[AllowListGetResponse](
		fmt.Sprintf("%s/get_allow_list_config", define.StdAuthServerAddress),
		AllowListGetRequest{
			Token:              f.config.EulogistToken,
			RentalServerNumber: rentalServerNumber,
		},
	)
	if err != nil {
		return fmt.Errorf("RentalServerManageModify: %v", err)
	}

	if !allowListGetResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: allowListGetResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("RentalServerManageModify: %v", err)
		}
		return nil
	}

	actionForm := form.ActionForm{
		Title: "修改已有授权",
	}
	if len(allowListGetResp.UserNames) == 0 {
		actionForm.Content = "您目前还没有授权任何赞颂者用户, 或许您需要§r§e返回上一级菜单§r添加。"
	} else {
		actionForm.Content = fmt.Sprintf(
			"您目前已授权了 §r§b%d§r 个赞颂者用户。\n请从中选择一个以操作他的权限。",
			len(allowListGetResp.UserNames),
		)
	}
	for _, value := range allowListGetResp.UserNames {
		actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
			Text: value,
			Icon: form.ActionFormIconNone{},
		})
	}

	resp, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
	if err != nil {
		return fmt.Errorf("RentalServerManageModify: %v", err)
	}
	if isUserCancel {
		return nil
	}

	err = f.rsmUpdatePermission(
		rentalServerNumber,
		allowListGetResp.UserNames[resp],
		false,
		allowListGetResp.DisableOpertorVerify[resp],
		allowListGetResp.CanGetGameSavesKeyCipher[resp],
	)
	if err != nil {
		return fmt.Errorf("RentalServerManageModify: %v", err)
	}

	return nil
}

// AllowListDeleteRequest ..
type AllowListDeleteRequest struct {
	Token              string `json:"token,omitempty"`
	RentalServerNumber string `json:"rental_server_number,omitempty"`
	EulogistUserName   string `json:"eulogist_user_name,omitempty"`
}

// AllowListDeleteResponse ..
type AllowListDeleteResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// RentalServerManageDelete ..
func (f *Function) RentalServerManageDelete(rentalServerNumber string) error {
	for {
		allowListGetResp, err := utils.SendAndGetHttpResponse[AllowListGetResponse](
			fmt.Sprintf("%s/get_allow_list_config", define.StdAuthServerAddress),
			AllowListGetRequest{
				Token:              f.config.EulogistToken,
				RentalServerNumber: rentalServerNumber,
			},
		)
		if err != nil {
			return fmt.Errorf("RentalServerManageDelete: %v", err)
		}

		if !allowListGetResp.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: allowListGetResp.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("RentalServerManageDelete: %v", err)
			}
			return nil
		}

		actionForm := form.ActionForm{
			Title: "删除已有授权",
		}
		if len(allowListGetResp.UserNames) == 0 {
			actionForm.Content = "您目前还没有授权任何赞颂者用户, 因此您无法§r§e删除§r任何授权。"
		} else {
			actionForm.Content = fmt.Sprintf(
				"您目前已授权了 §r§b%d§r 个赞颂者用户。\n请从中选择一个, 然后将其§r§e删除§r。",
				len(allowListGetResp.UserNames),
			)
		}
		for _, value := range allowListGetResp.UserNames {
			actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
				Text: value,
				Icon: form.ActionFormIconNone{},
			})
		}

		resp, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
		if err != nil {
			return fmt.Errorf("RentalServerManageDelete: %v", err)
		}
		if isUserCancel {
			return nil
		}

		ensure, isUserCancel, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "二重确认",
			Content: "您确定要§r§c删除§r目标赞颂者用户的授权吗？",
			Button1: "确定",
			Button2: "取消",
		})
		if err != nil {
			return fmt.Errorf("RentalServerManageDelete: %v", err)
		}
		if isUserCancel || !ensure.(bool) {
			continue
		}

		allowListDeleteResp, err := utils.SendAndGetHttpResponse[AllowListDeleteResponse](
			fmt.Sprintf("%s/delete_allow_list_config", define.StdAuthServerAddress),
			AllowListDeleteRequest{
				Token:              f.config.EulogistToken,
				RentalServerNumber: rentalServerNumber,
				EulogistUserName:   allowListGetResp.UserNames[resp],
			},
		)
		if err != nil {
			return fmt.Errorf("RentalServerManageDelete: %v", err)
		}

		if !allowListDeleteResp.Success {
			_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "错误",
				Content: allowListDeleteResp.ErrorInfo,
				Button1: "确定",
				Button2: "返回上一级菜单",
			})
			if err != nil {
				return fmt.Errorf("RentalServerManageDelete: %v", err)
			}
			return nil
		}

		return nil
	}
}

// RentalServerManage ..
func (f *Function) RentalServerManage(
	afterSelectRentalServer func(rentalServerNumber string) error,
) (isUserCancel bool, err error) {
	for {
		rentalServerNumber, isUserCancel, err := f.ListServerCanManage()
		if err != nil {
			return false, fmt.Errorf("RentalServerManage: %v", err)
		}
		if isUserCancel {
			return true, nil
		}
		err = afterSelectRentalServer(rentalServerNumber)
		if err != nil {
			return false, fmt.Errorf("RentalServerManage: %v", err)
		}
	}
}

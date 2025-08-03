package function

import (
	"fmt"
	"strings"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// UserSearchRequest ..
type UserSearchRequest struct {
	Token        string `json:"token,omitempty"`
	FilterString string `json:"filter_string,omitempty"`
}

// UserSearchResponse ..
type UserSearchResponse struct {
	ErrorInfo   string   `json:"error_info"`
	Success     bool     `json:"success"`
	HitUserName []string `json:"hit_user_name"`
}

// ShowEulogistSearch ..
func (f *Function) ShowEulogistSearch() (eulogistUserName string, isUserCancel bool, err error) {
	minecraftForm := form.ModalForm{
		Title: "搜索用户",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "请在下方输入框输入您希望搜索的用户。\n提供的用户名可以是不完整的, 但不得为空。",
			},
			form.ModalFormElementInput{
				Text:        "用户名",
				Default:     "",
				PlaceHolder: "赞颂者用户的用户名 (或用户名的一部分)",
			},
		},
	}

	resp, isUserCancel, err := f.Interact().SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return "", false, fmt.Errorf("ShowEulogistSearch: %v", err)
	}
	if isUserCancel {
		return "", true, nil
	}
	filterString := resp.([]any)[1].(string)

	userSearchResponse, err := SendAndGetHttpResponse[UserSearchResponse](
		fmt.Sprintf("%s/search_eulogist_user", define.StdAuthServerAddress),
		UserSearchRequest{
			Token:        f.config.EulogistToken,
			FilterString: filterString,
		},
	)
	if err != nil {
		return "", false, fmt.Errorf("ShowEulogistSearch: %v", err)
	}

	if !userSearchResponse.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: userSearchResponse.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return "", false, fmt.Errorf("ShowEulogistSearch: %v", err)
		}
		return "", true, nil
	}

	replaceString := fmt.Sprintf("§r§c%s§r", filterString)
	actionForm := form.ActionForm{
		Title: "搜索结果",
	}
	if len(userSearchResponse.HitUserName) == 0 {
		actionForm.Content = "没有搜索到满足条件的任何赞颂者用户。您可能需要§r§e重试§r。"
	} else {
		actionForm.Content = fmt.Sprintf("已搜索到 §r§b%d§r 个用户。", len(userSearchResponse.HitUserName))
	}
	for _, value := range userSearchResponse.HitUserName {
		actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
			Text: strings.ReplaceAll(value, filterString, replaceString),
			Icon: form.ActionFormIconNone{},
		})
	}

	answer, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
	if err != nil {
		return "", false, fmt.Errorf("ShowEulogistSearch: %v", err)
	}
	if isUserCancel {
		return "", true, nil
	}

	return userSearchResponse.HitUserName[answer], false, nil
}

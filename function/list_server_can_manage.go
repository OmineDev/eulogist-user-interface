package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
)

// ListServerCanManage ..
func (f *Function) ListServerCanManage() (rentalServerNumber string, isUserCancel bool, err error) {
	actionForm := form.ActionForm{
		Title: "选择租赁服",
	}

	if len(f.userData.RentalServerCanManage) == 0 {
		actionForm.Content = "您旗下没有可管理的租赁服, 请§r§e返回上一级菜单§r。"
	} else {
		actionForm.Content = fmt.Sprintf(
			"您总共可以管理 §r§b%d§r 个租赁服。\n请从中选择一个进行管理。",
			len(f.userData.RentalServerCanManage),
		)
	}

	for _, value := range f.userData.RentalServerCanManage {
		actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
			Text: fmt.Sprintf("租赁服 %s", value),
			Icon: form.ActionFormIconNone{},
		})
	}

	resp, isUserCancel, err := f.Interact().SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
	if err != nil {
		return "", false, fmt.Errorf("ListServerCanManage: %v", err)
	}
	if isUserCancel {
		return "", true, nil
	}

	return f.userData.RentalServerCanManage[resp], false, nil
}

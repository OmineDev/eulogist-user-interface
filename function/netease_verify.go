package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/form"
)

// ShowAuthServerError 以表单的形式显示验证服务的错误回应。
// neteaseRequireVerify 指示是否是因为操作 MC 账号时触发网易验证，
// verifyURl 则是验证所用的链接。errorInfo 指示验证服务返回的错误信息
func (f *Function) ShowAuthServerError(
	neteaseRequireVerify bool,
	verifyURl string,
	errorInfo string,
) (isUserCancel bool, err error) {
	if !neteaseRequireVerify {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: errorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return false, fmt.Errorf("ShowAuthServerError: %v", err)
		}
		return false, nil
	}

	_, isUserCancel, err = f.interact.SendFormAndWaitResponse(form.ModalForm{
		Title: "需要验证",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "" +
					"操作 MC 账号时触发网易验证。\n" +
					"请§r§e复制§r下面输入框中的链接, \n" +
					"然后§r§e粘贴§r到浏览器打开以完成验证。\n" +
					"验证完成后, 请点击§r§e提交§r按钮以便能继续。\n" +
					"什么, 你不会复制粘贴？\n" +
					"  使用 §r§bCTRL+A§r 全选, 然后使用 §r§bCTRL+C§r 复制即可！",
			},
			form.ModalFormElementInput{
				Text:        "网易验证链接",
				Default:     verifyURl,
				PlaceHolder: "网易验证链接",
			},
		},
	})
	if err != nil {
		return false, fmt.Errorf("SendNeteaseVerifyInfo: %v", err)
	}

	return isUserCancel, nil
}

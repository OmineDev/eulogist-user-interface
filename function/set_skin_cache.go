package function

import (
	"encoding/hex"
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

// SetSkinCacheRequest ..
type SetSkinCacheRequest struct {
	Token           string `json:"token,omitempty"`
	SkinDownloadURL string `json:"skin_download_url,omitempty"`
}

// SetSkinCacheResponse ..
type SetSkinCacheResponse struct {
	ErrorInfo string `json:"error_info"`
	Success   bool   `json:"success"`
}

// sendSetSkinCacheRequest ..
func (f *Function) sendSetSkinCacheRequest(skinDownloadURL string) error {
	setSkinCacheResp, err := utils.SendAndGetHttpResponse[SetSkinCacheResponse](
		fmt.Sprintf("%s/set_skin_cache", define.StdAuthServerAddress),
		SetSkinCacheRequest{
			Token:           f.config.EulogistToken,
			SkinDownloadURL: skinDownloadURL,
		},
	)
	if err != nil {
		return fmt.Errorf("sendSetSkinCacheRequest: %v", err)
	}
	if !setSkinCacheResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: setSkinCacheResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("sendSetSkinCacheRequest: %v", err)
		}
		return nil
	}
	return nil
}

// UseBuiltInSkin ..
func (f *Function) UseBuiltInSkin() (isUserCancel bool, err error) {
	var skinDownloadURL string
	skinNameList := []string{
		"e79089e79283e4b99de5b0be",
		"e585abe788aae5a2a8e5a2a8",
		"e681b6e9ad94e7949ce5bf83",
		"e782bde784b0e78b82e683b3",
		"e6b3b3e8a385e5a48fe697a5202d20e591b1e591b1",
		"e6b3b3e8a385e5a48fe697a5202d20e595bee595be",
		"e6b3b3e8a385e5a48fe697a5202d20e597a1e597a1",
		"e6bdaee6b190e5a49ce69bb2",
		"e9a38ee4b98be5928fe58fb9",
		"e99d92e9b89f",
		"e5a587e5a699e789a9e8afade7b3bbe58897202d20e6aca7e68b89e4b98be6989f",
		"e5a587e5a699e789a9e8afade7b3bbe58897202d20e799bde993b6e68898e5a3ab",
		"e5a587e5a699e789a9e8afade7b3bbe58897202d20e7bbbfe889b2e6b0b4e699b6",
		"e5a587e5a699e789a9e8afade7b3bbe58897202d20e69ca8e5a4a7e69ca8e5a4a7",
		"e681b6e9ad94e6bdaee6b581",
		"e79280e792a8202d20e9a38ee4b98be5928fe58fb9",
		"e79280e792a8202d20e6bdaee6b190e5a49ce69bb2",
		"e9bb91e6b49ee4bba3e79086e4baba",
		"e88db7e9b281e696afe5ba87e68aa4",
		"e69c88e58589e78e8be5ad90",
		"e88b8de78e89e9ba92e9ba9f",
	}

	actionForm := form.ActionForm{
		Title:   "Use built in skin",
		Content: "§r§eChoose§r a built in skin so that you can use it in rental server.",
	}
	for _, value := range skinNameList {
		skinName, err := hex.DecodeString(value)
		if err != nil {
			return false, fmt.Errorf("UseBuiltInSkin: %v", err)
		}
		actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
			Text: string(skinName),
			Icon: form.ActionFormIconNone{},
		})
	}
	actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
		Text: "Back to main panel",
		Icon: form.ActionFormIconNone{},
	})

	resp, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
	if err != nil {
		return false, fmt.Errorf("UseBuiltInSkin: %v", err)
	}
	if isUserCancel {
		return true, nil
	}

	switch resp {
	case 0:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383130383733353837303630353234" +
			"38355f335f76335f387479646a306e6a2e7a6970"
	case 1:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383035333239373434313535363333" +
			"34305f335f76335f7a6b34396a7a726e2e7a6970"
	case 2:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383130383733353933393732373239" +
			"32395f345f76345f69637176797478722e7a6970"
	case 3:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383037383031313732383339313737" +
			"37385f325f76325f77386a77393776692e7a6970"
	case 4:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383130383733353933343439333939" +
			"39355f345f76345f6470666a793031772e7a6970"
	case 5:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383130383733353932313839313439" +
			"33375f345f76345f697078306e6632612e7a6970"
	case 6:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383130383733353932383230363635" +
			"34305f345f76345f74743167766677652e7a6970"
	case 7:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383035333239373333383237303134" +
			"30305f325f76325f326733726e7539392e7a6970"
	case 8:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383032343438363539343230393339" +
			"37325f335f76335f306b77336b7732792e7a6970"
	case 9:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383037373532373837363230373134" +
			"38355f325f76325f677571793677796b2e7a6970"
	case 10:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383135353131373831313438373333" +
			"32385f345f76345f6b626b3170656e682e7a6970"
	case 11:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383135353131373832363738393934" +
			"39315f335f76335f70697a79703275722e7a6970"
	case 12:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383135353131373831393334363232" +
			"35375f335f76335f736e78766e6f31382e7a6970"
	case 13:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383135353131373833353537393930" +
			"30325f335f76335f686d3931746471722e7a6970"
	case 14:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383035353531313730333434383737" +
			"31315f335f76335f306e6c68666c72352e7a6970"
	case 15:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383032343438363630343138313934" +
			"38325f345f76345f797a64676b77386a2e7a6970"
	case 16:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436383035333239373334363131313038" +
			"38395f335f76335f63793434616b37612e7a6970"
	case 17:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436373431363930333432383833343633" +
			"35355f31335f7631335f71706c72326830362e7a6970"
	case 18:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e6770682e6e6574656173652e63" +
			"6f6d2f6974656d5f3436373335303337313630323539373833" +
			"37335f335f76335f67313863646664772e7a6970"
	case 19:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e66702e70732e6e657465617365" +
			"2e636f6d2f66696c652f363233343362643838633337383433" +
			"33613238363635303971524e47466c5a303034"
	case 20:
		skinDownloadURL = "" +
			"68747470733a2f2f7831392e66702e70732e6e657465617365" +
			"2e636f6d2f66696c652f363362353664373163643032323638" +
			"613731653833316439357969384c43584e3034"
	case 21:
		return true, nil
	}

	realSkinDownloadURL, err := hex.DecodeString(skinDownloadURL)
	if err != nil {
		return false, fmt.Errorf("UseBuiltInSkin: %v", err)
	}
	err = f.sendSetSkinCacheRequest(string(realSkinDownloadURL))
	if err != nil {
		return false, fmt.Errorf("UseBuiltInSkin: %v", err)
	}

	return false, nil
}

// UseCustomSkin ..
func (f *Function) UseCustomSkin() (isUserCancel bool, err error) {
	minecraftForm := form.ModalForm{
		Title: "Use a custom skin",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "§r§eFill in§r the skin info that corresponding to the skin you want to use.",
			},
			form.ModalFormElementInput{
				Text:        "Skin Download URL",
				Default:     "",
				PlaceHolder: "A URL to download this skin",
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return false, fmt.Errorf("UseCustomSkin: %v", err)
	}
	if isUserCancel {
		return true, nil
	}

	err = f.sendSetSkinCacheRequest(resp.([]any)[1].(string))
	if err != nil {
		return false, fmt.Errorf("ChooseBuiltInSkin: %v", err)
	}

	return false, nil
}

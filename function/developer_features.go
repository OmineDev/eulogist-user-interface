package function

import (
	"bytes"
	"fmt"
	"time"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
	"github.com/ToolDelta-Basic/gophertunnel/minecraft/protocol"
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
		fmt.Sprintf("%s/set_skin_cache", define.AddressEulogistAPI),
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

// BuiltInSkinRequest ..
type BuiltInSkinRequest struct {
	Token string `json:"token,omitempty"`
}

// BuiltInSkinResponse ..
type BuiltInSkinResponse struct {
	ErrorInfo       string   `json:"error_info"`
	Success         bool     `json:"success"`
	SkinNameList    []string `json:"skin_name_list"`
	SkinDownloadURL []string `json:"skin_download_url"`
}

// UseBuiltInSkin ..
func (f *Function) UseBuiltInSkin() (success bool, err error) {
	builtInSkinResp, err := utils.SendAndGetHttpResponse[BuiltInSkinResponse](
		fmt.Sprintf("%s/get_built_in_skin", define.AddressEulogistAPI),
		BuiltInSkinRequest{
			Token: f.config.EulogistToken,
		},
	)
	if err != nil {
		return false, fmt.Errorf("UseBuiltInSkin: %v", err)
	}
	if !builtInSkinResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: builtInSkinResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return false, fmt.Errorf("UseBuiltInSkin: %v", err)
		}
		return false, nil
	}

	actionForm := form.ActionForm{
		Title:   "Use built in skin",
		Content: "§r§eChoose§r a built in skin so that you can use it in rental server.",
	}
	for _, value := range builtInSkinResp.SkinNameList {
		actionForm.Buttons = append(actionForm.Buttons, form.ActionFormElement{
			Text: value,
			Icon: form.ActionFormIconNone{},
		})
	}

	resp, isUserCancel, err := f.interact.SendLargeActionFormAndWaitResponse(actionForm, define.DefaultPageSize)
	if err != nil {
		return false, fmt.Errorf("UseBuiltInSkin: %v", err)
	}
	if isUserCancel {
		return false, nil
	}

	err = f.sendSetSkinCacheRequest(builtInSkinResp.SkinDownloadURL[resp])
	if err != nil {
		return false, fmt.Errorf("UseBuiltInSkin: %v", err)
	}

	f.customData.useCustomSkin = false
	f.customData.customSkinData = nil
	return true, nil
}

// UseCustomSkin ..
func (f *Function) UseCustomSkin() (success bool, err error) {
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
		return false, nil
	}

	err = f.sendSetSkinCacheRequest(resp.([]any)[1].(string))
	if err != nil {
		return false, fmt.Errorf("ChooseBuiltInSkin: %v", err)
	}

	f.customData.useCustomSkin = false
	f.customData.customSkinData = nil
	return true, nil
}

// UseSkinFromMinecraft ..
func (f *Function) UseSkinFromMinecraft() error {
	minecraftForm := form.MessageForm{
		Title: "Use a skin from standard Minecraft",
		Content: "" +
			"You will use a skin that from standard Minecraft.\n\n" +
			"" +
			"After you press §r§eOK§r, then you have §r§e30 seconds§r to change the skin §r§ein game§r.\n" +
			"Then, you will login to rental server with this skin.\n\n" +
			"" +
			"Note that we only accept the §r§efirst§r skin you select.\n" +
			"For §r§eany future skin§r we will not use.",
		Button1: "OK",
		Button2: "Cancel",
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("UseSkinFromMinecraft: %v", err)
	}
	if isUserCancel || !resp.(bool) {
		return nil
	}

	clientSkin, isTimeOut := f.interact.WaitClientUseSkin(time.Second * 30)
	if isTimeOut {
		minecraftForm = form.MessageForm{
			Title:   "Set skin time out",
			Content: "We failed to get your skin with in §r§e30 seconds§r!",
			Button1: "OK",
			Button2: "Back to last panel",
		}
		_, _, err := f.interact.SendFormAndWaitResponse(minecraftForm)
		if err != nil {
			return fmt.Errorf("UseSkinFromMinecraft: %v", err)
		}
		return nil
	}
	buf := bytes.NewBuffer(nil)
	writer := protocol.NewWriter(buf, 0)
	clientSkin.Marshal(writer)

	minecraftForm = form.MessageForm{
		Title:   "Success",
		Content: "We success to get and set skin that from standard Minecraft!",
		Button1: "OK",
		Button2: "Continue",
	}
	_, _, err = f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("UseSkinFromMinecraft: %v", err)
	}

	f.customData.useCustomSkin = true
	f.customData.customSkinData = buf.Bytes()
	return nil
}

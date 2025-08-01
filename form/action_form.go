package form

import "encoding/json"

type ActionFormIcon any
type ActionFormIconNone struct{}
type ActionFormIconPathImage struct {
	ImagePath string // e.g. `textures/ui/anvil_icon.png`
}
type ActionFormIconURLImage struct {
	ImageURL string // e.g. `https://avatars.githubusercontent.com/u/109064184`
}

// ActionFormElement ..
type ActionFormElement struct {
	Text string         // 按钮中的文字
	Icon ActionFormIcon // 按钮的图标，可以不使用
}

func (a ActionFormElement) asJsonStruct() any {
	switch icon := a.Icon.(type) {
	case ActionFormIconPathImage:
		var jsonStruct struct {
			Text  string `json:"text"`
			Image struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"image"`
		}
		jsonStruct.Text = a.Text
		jsonStruct.Image.Type = "path"
		jsonStruct.Image.Data = icon.ImagePath
		return jsonStruct
	case ActionFormIconURLImage:
		var jsonStruct struct {
			Text  string `json:"text"`
			Image struct {
				Type string `json:"type"`
				Data string `json:"data"`
			} `json:"image"`
		}
		jsonStruct.Text = a.Text
		jsonStruct.Image.Type = "url"
		jsonStruct.Image.Data = icon.ImageURL
		return jsonStruct
	default:
		var jsonStruct struct {
			Text string `json:"text"`
		}
		jsonStruct.Text = a.Text
		return jsonStruct
	}
}

// ActionForm ..
type ActionForm struct {
	Title   string
	Content string
	Buttons []ActionFormElement
}

func (a ActionForm) ID() uint8 {
	return FormTypeAction
}

func (a ActionForm) PackToJSON() string {
	var jsonStruct struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Content string `json:"content"`
		Buttons []any  `json:"buttons"`
	}

	jsonStruct.Type = "form"
	jsonStruct.Title = a.Title
	jsonStruct.Content = a.Content
	for _, button := range a.Buttons {
		jsonStruct.Buttons = append(jsonStruct.Buttons, button.asJsonStruct())
	}

	jsonBytes, _ := json.Marshal(jsonStruct)
	return string(jsonBytes)
}

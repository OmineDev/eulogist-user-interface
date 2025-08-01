package form

import "encoding/json"

// ModalFormElement 是模态表单中各种元素类别的总称
type ModalFormElement any

// 普通文本
type ModalFormElementLabel struct {
	Text string `json:"text"` // 对应的文本
}

// 输入框
type ModalFormElementInput struct {
	Text        string `json:"text"`        // 输入框的标题
	Default     string `json:"default"`     // 输入框内的默认内容
	PlaceHolder string `json:"placeholder"` // 输入框的提示内容
}

// 按钮
type ModalFormElementToggle struct {
	Text    string `json:"text"`    // 按钮的文字
	Default bool   `json:"default"` // 按钮的默认状态 (开/关)
}

// 多选
type ModalFormElementDropdown struct {
	Text    string   `json:"text"`    // 多选的标题
	Options []string `json:"options"` // 各个选项
	Default uint     `json:"default"` // 默认选项 (即上方 Options 的索引)
}

// 可拖拽的进度条
type ModalFormElementSlider struct {
	Text    string `json:"text"`    // 进度条的标题
	Min     int    `json:"min"`     // 进度条的最小值
	Max     int    `json:"max"`     // 进度条的最大值
	Step    int    `json:"step"`    // 用户调整进度条的最小单位 (或步进长度)
	Default int    `json:"default"` // 进度条的默认值
}

// 滑块
type ModalFormElementStepSlider struct {
	Text    string   `json:"text"`    // 滑块的标题
	Steps   []string `json:"steps"`   // 可以滑动选择的各个选项
	Default uint     `json:"default"` // 默认的选项 (即上方 Steps 的索引)
}

// packModalFormElementToStruct ..
func packModalFormElementToStruct(element ModalFormElement) any {
	switch e := element.(type) {
	case ModalFormElementLabel:
		var jsonStruct struct {
			Type string `json:"type"`
			ModalFormElementLabel
		}
		jsonStruct.Type = "label"
		jsonStruct.ModalFormElementLabel = e
		return jsonStruct
	case ModalFormElementInput:
		var jsonStruct struct {
			Type string `json:"type"`
			ModalFormElementInput
		}
		jsonStruct.Type = "input"
		jsonStruct.ModalFormElementInput = e
		return jsonStruct
	case ModalFormElementToggle:
		var jsonStruct struct {
			Type string `json:"type"`
			ModalFormElementToggle
		}
		jsonStruct.Type = "toggle"
		jsonStruct.ModalFormElementToggle = e
		return jsonStruct
	case ModalFormElementDropdown:
		var jsonStruct struct {
			Type string `json:"type"`
			ModalFormElementDropdown
		}
		jsonStruct.Type = "dropdown"
		jsonStruct.ModalFormElementDropdown = e
		return jsonStruct
	case ModalFormElementSlider:
		var jsonStruct struct {
			Type string `json:"type"`
			ModalFormElementSlider
		}
		jsonStruct.Type = "slider"
		jsonStruct.ModalFormElementSlider = e
		return jsonStruct
	case ModalFormElementStepSlider:
		var jsonStruct struct {
			Type string `json:"type"`
			ModalFormElementStepSlider
		}
		jsonStruct.Type = "step_slider"
		jsonStruct.ModalFormElementStepSlider = e
		return jsonStruct
	default:
		return map[string]any{}
	}
}

// ModalForm ..
type ModalForm struct {
	Title    string             // 模态表单的标题
	Contents []ModalFormElement // 模态表单的各个内容
}

func (m ModalForm) ID() uint8 {
	return FormTypeModal
}

func (m ModalForm) PackToJSON() string {
	var jsonStruct struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Content []any  `json:"content"`
	}

	jsonStruct.Type = "custom_form"
	jsonStruct.Title = m.Title
	for _, content := range m.Contents {
		jsonStruct.Content = append(jsonStruct.Content, packModalFormElementToStruct(content))
	}

	jsonBytes, _ := json.Marshal(jsonStruct)
	return string(jsonBytes)
}

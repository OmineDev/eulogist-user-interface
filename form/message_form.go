package form

import "encoding/json"

// MessageForm ..
type MessageForm struct {
	Title   string `json:"title"`   // 标题
	Content string `json:"content"` // 内容
	Button1 string `json:"button1"` // 可以表示 “确定/取消” 的 “确定” 按钮
	Button2 string `json:"button2"` // 可以表示 “确定/取消” 的 “取消” 按钮
}

func (m MessageForm) ID() uint8 {
	return FormTypeMessage
}

func (m MessageForm) PackToJSON() string {
	var jsonStruct struct {
		Type string `json:"type"`
		MessageForm
	}

	jsonStruct.Type = "modal"
	jsonStruct.MessageForm = m

	jsonBytes, _ := json.Marshal(jsonStruct)
	return string(jsonBytes)
}

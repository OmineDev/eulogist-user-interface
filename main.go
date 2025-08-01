package main

import (
	"fmt"
	"strings"

	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/server"
)

func main() {
	test()
}

func test() {
	s := server.NewServer()

	fmt.Println(s.RunServer("127.0.0.1:19132"))
	fmt.Println(s.WaitConnect())

	conn := s.MinecraftConn()
	interact := server.NewInteract(conn)

	for {
		resp, err := interact.SendFormAndWaitResponse(
			form.ModalForm{
				Title: "菜单标题",
				Contents: []form.ModalFormElement{
					form.ModalFormElementLabel{
						Text: "普通文本",
					},
					form.ModalFormElementInput{
						Text:        "输入框",
						Default:     "",
						PlaceHolder: strings.Repeat("233", 100),
					},
					form.ModalFormElementToggle{
						Text:    "开关按钮",
						Default: true,
					},
					form.ModalFormElementDropdown{
						Text:    "多选",
						Options: []string{"选项1", "选项2", "选项3"},
						Default: 2,
					},
					form.ModalFormElementSlider{
						Text:    "进度条",
						Min:     -10,
						Max:     100,
						Step:    30,
						Default: 7,
					},
					form.ModalFormElementStepSlider{
						Text:    "滑动选择",
						Steps:   []string{"选项1", "选项2", "选项3"},
						Default: 2,
					},
				},
			},
		)

		select {
		case <-conn.Context().Done():
			_ = s.CloseServer()
			return
		default:
		}

		fmt.Printf("%#v\n", resp)
		fmt.Println(err)
	}
}

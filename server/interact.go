package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// Interact 是客户端和赞颂者假服务器的表单交互实现
type Interact struct {
	mu         *sync.Mutex
	conn       *minecraft.Conn
	formID     uint32
	clientResp packet.ModalFormResponse
	waiter     context.Context
	downWaiter context.CancelFunc
}

// NewInteract 根据 server 创建并返回一个新的交互装置
func NewInteract(conn *minecraft.Conn) *Interact {
	interact := &Interact{
		mu:     new(sync.Mutex),
		conn:   conn,
		formID: 0,
	}
	go interact.handlePacket()
	return interact
}

// sendFormAndWaitResponse ..
func (i *Interact) sendFormAndWaitResponse(minecraftForm form.MinecraftForm) (resp any, isUserCanel bool, err error) {
	for {
		i.waiter, i.downWaiter = context.WithCancel(context.Background())

		err = i.conn.WritePacket(&packet.ModalFormRequest{
			FormID:   i.formID,
			FormData: []byte(minecraftForm.PackToJSON()),
		})
		if err != nil {
			return nil, false, fmt.Errorf("SendFormAndWaitResponse: %v", err)
		}
		i.formID++

		select {
		case <-i.waiter.Done():
		case <-i.conn.Context().Done():
			return nil, false, fmt.Errorf("SendFormAndWaitResponse: Minecraft connection has been closed")
		}

		if i.clientResp.FormID != i.formID-1 {
			return nil, false, fmt.Errorf(
				"SendFormAndWaitResponse: Form ID not match (server = %d, client = %d)",
				i.formID-1, i.clientResp.FormID,
			)
		}

		canelReason, ok := i.clientResp.CancelReason.Value()
		if ok {
			if canelReason == packet.ModalFormCancelReasonUserClosed {
				return nil, true, nil
			}
			time.Sleep(time.Second / 20)
			continue
		}

		resp, ok := i.clientResp.ResponseData.Value()
		if !ok {
			return nil, false, fmt.Errorf("SendFormAndWaitResponse: Response data is not exist")
		}

		switch minecraftForm.ID() {
		case form.FormTypeMessage:
			if strings.TrimSuffix(string(resp), "\n") == "true" {
				return true, false, nil
			}
			return false, false, nil
		case form.FormTypeAction:
			result, err := strconv.ParseInt(strings.TrimSuffix(string(resp), "\n"), 10, 32)
			if err != nil {
				return nil, false, fmt.Errorf("SendFormAndWaitResponse: %v", err)
			}
			return int32(result), false, nil
		case form.FormTypeModal:
			var respList []any

			err = json.Unmarshal(resp, &respList)
			if err != nil {
				return nil, false, fmt.Errorf("SendFormAndWaitResponse: %v", err)
			}

			result := make([]any, 0)
			f := minecraftForm.(form.ModalForm)

			for index, value := range respList {
				switch f.Contents[index].(type) {
				case form.ModalFormElementLabel:
					result = append(result, nil)
				case form.ModalFormElementInput:
					result = append(result, value.(string))
				case form.ModalFormElementToggle:
					result = append(result, value.(bool))
				case form.ModalFormElementDropdown:
					result = append(result, int32(value.(float64)))
				case form.ModalFormElementSlider:
					result = append(result, int32(value.(float64)))
				case form.ModalFormElementStepSlider:
					result = append(result, int32(value.(float64)))
				default:
					return nil, false, fmt.Errorf(
						"SendFormAndWaitResponse: Unknown modal form element type %T; f.Contents[%d] = %#v",
						f.Contents[index], index, f.Contents[index],
					)
				}
			}

			return result, false, nil
		default:
			return nil, false, fmt.Errorf(
				"SendFormAndWaitResponse: Unknown minecraft form type %T; minecraftForm = %#v",
				minecraftForm, minecraftForm,
			)
		}
	}
}

// sendLargeActionFormAndWaitResponse ..
func (i *Interact) sendLargeActionFormAndWaitResponse(
	actionForm form.ActionForm,
	pageSize int,
) (resp int32, isUserCanel bool, err error) {
	pageSize = max(1, pageSize)
	maxPage := max(1, (len(actionForm.Buttons)+pageSize-1)/pageSize)
	currentPage := 1

	for {
		var lastPageIndex int32
		var nextPageIndex int32
		var jumpPageIndex int32
		var exitIndex int32

		startIndexInclude := int32((currentPage - 1) * pageSize)
		endIndexNotInclude := int32(min(currentPage*pageSize, len(actionForm.Buttons)))

		newForm := form.ActionForm{
			Title:   actionForm.Title,
			Content: actionForm.Content,
		}
		newForm.Content += fmt.Sprintf(
			"\n\n§r§f当前第 §b%d §f页, 总计 §b%d §f页",
			currentPage, maxPage,
		)

		// Append normal entry
		for i := startIndexInclude; i < endIndexNotInclude; i++ {
			newForm.Buttons = append(newForm.Buttons, actionForm.Buttons[i])
		}
		// Last page button
		if currentPage > 1 {
			lastPageIndex = int32(len(newForm.Buttons))
			newForm.Buttons = append(newForm.Buttons, form.ActionFormElement{
				Text: "§r§l§2上一页",
				Icon: form.ActionFormIconNone{},
			})
		}
		// Next page button
		if currentPage*pageSize < len(actionForm.Buttons) {
			nextPageIndex = int32(len(newForm.Buttons))
			newForm.Buttons = append(newForm.Buttons, form.ActionFormElement{
				Text: "§r§l§2下一页",
				Icon: form.ActionFormIconNone{},
			})
		}
		// Jump to button
		jumpPageIndex = int32(len(newForm.Buttons))
		newForm.Buttons = append(newForm.Buttons, form.ActionFormElement{
			Text: "§r§l§2跳转到",
			Icon: form.ActionFormIconNone{},
		})
		// Exit button
		exitIndex = int32(len(newForm.Buttons))
		newForm.Buttons = append(newForm.Buttons, form.ActionFormElement{
			Text: "§r§l§c返回上一级菜单",
			Icon: form.ActionFormIconNone{},
		})

		anyResp, isUserCanel, err := i.sendFormAndWaitResponse(newForm)
		if err != nil {
			return 0, false, fmt.Errorf("SendLargeActionFormAndWaitResponse: %v", err)
		}
		if isUserCanel {
			return 0, true, nil
		}

		idx := anyResp.(int32)
		realIndex := startIndexInclude + idx
		if startIndexInclude <= realIndex && realIndex < endIndexNotInclude {
			return realIndex, false, nil
		}

		switch idx {
		case lastPageIndex:
			currentPage--
		case nextPageIndex:
			currentPage++
		case jumpPageIndex:
			anyResp, isUserCanel, err := i.sendFormAndWaitResponse(form.ModalForm{
				Title: "跳转",
				Contents: []form.ModalFormElement{
					form.ModalFormElementInput{
						Text:        "跳转到",
						Default:     "",
						PlaceHolder: fmt.Sprintf("页数 (当前第 %d 页 | 最多 %d 页)", currentPage, maxPage),
					},
				},
			})
			if err != nil {
				return 0, false, fmt.Errorf("SendLargeActionFormAndWaitResponse: %v", err)
			}
			if !isUserCanel {
				jumpTo, err := strconv.ParseInt(anyResp.([]any)[0].(string), 10, 32)
				if err != nil {
					jumpTo = int64(currentPage)
				}
				currentPage = min(max(int(jumpTo), 1), maxPage)
			}
		case exitIndex:
			return 0, true, nil
		default:
			panic("SendLargeActionFormAndWaitResponse: Should nerver happened")
		}
	}
}

// SendFormAndWaitResponse 发送 minecraftForm 所指示的表单给客户端并等待回应。
//
// resp 是客户端的回应，只可能为：
//   - minecraftForm.ID() 为 [from.FormTypeMessage] 时：bool
//   - minecraftForm.ID() 为 [from.FormTypeAction] 时：int32
//   - minecraftForm.ID() 为 [from.FormTypeModal] 时：[]any
//
// 如果回应是 []any，则其中的元素只可能是：
//   - [form.ModalFormElementLabel] -> nil
//   - [form.ModalFormElementInput] -> string
//   - [form.ModalFormElementToggle] -> bool
//   - [form.ModalFormElementDropdown] -> int32
//   - [form.ModalFormElementSlider] -> int32
//   - [form.ModalFormElementStepSlider] -> int32
//
// isUserCanel 指示表单是否是由用户通过叉号 (×) 关闭的
func (i *Interact) SendFormAndWaitResponse(minecraftForm form.MinecraftForm) (resp any, isUserCanel bool, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.sendFormAndWaitResponse(minecraftForm)
}

// SendLargeActionFormAndWaitResponse 向客户端发送大型的 ActionForm，
// 这意味着 actionForm.Buttons 具有很多项目，需要按 pageSize 分页拆分。
// isUserCanel 指示表单是否是由用户通过叉号 (×) 关闭的
func (i *Interact) SendLargeActionFormAndWaitResponse(
	actionForm form.ActionForm,
	pageSize int,
) (resp int32, isUserCanel bool, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.sendLargeActionFormAndWaitResponse(actionForm, pageSize)
}

// handlePacket 不断地读取数据包，
// 并期望下一个抵达的数据包是客户端对表单的响应
func (i *Interact) handlePacket() {
	for {
		pk, err := i.conn.ReadPacket()
		if err != nil {
			return
		}

		p, ok := pk.(*packet.ModalFormResponse)
		if !ok {
			continue
		}

		i.clientResp = *p
		i.downWaiter()
	}
}

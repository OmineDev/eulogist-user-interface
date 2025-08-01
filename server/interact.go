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
func (i *Interact) sendFormAndWaitResponse(minecraftForm form.MinecraftForm) (resp any, err error) {
	for {
		i.waiter, i.downWaiter = context.WithCancel(context.Background())

		err = i.conn.WritePacket(&packet.ModalFormRequest{
			FormID:   i.formID,
			FormData: []byte(minecraftForm.PackToJSON()),
		})
		if err != nil {
			return nil, fmt.Errorf("SendFormAndWaitResponse: %v", err)
		}
		i.formID++

		select {
		case <-i.waiter.Done():
		case <-i.conn.Context().Done():
			return nil, fmt.Errorf("SendFormAndWaitResponse: Minecraft connection has been closed")
		}

		if i.clientResp.FormID != i.formID-1 {
			return nil, fmt.Errorf(
				"SendFormAndWaitResponse: Form ID not match (server = %d, client = %d)",
				i.formID-1, i.clientResp.FormID,
			)
		}

		_, ok := i.clientResp.CancelReason.Value()
		if ok {
			time.Sleep(time.Second / 20)
			continue
		}

		resp, ok := i.clientResp.ResponseData.Value()
		if !ok {
			return nil, fmt.Errorf("SendFormAndWaitResponse: Response data is not exist")
		}

		switch minecraftForm.ID() {
		case form.FormTypeMessage:
			if strings.Contains(strings.ToLower(string(resp)), "true") {
				return true, nil
			}
			return false, nil
		case form.FormTypeAction:
			result, err := strconv.ParseInt(string(resp), 10, 32)
			if err != nil {
				return nil, fmt.Errorf("SendFormAndWaitResponse: %v", err)
			}
			return int32(result), nil
		case form.FormTypeModal:
			var respList []any

			err = json.Unmarshal(resp, &respList)
			if err != nil {
				return nil, fmt.Errorf("SendFormAndWaitResponse: %v", err)
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
					return nil, fmt.Errorf(
						"SendFormAndWaitResponse: Unknown modal form element type %T; f.Contents[%d] = %#v",
						f.Contents[index], index, f.Contents[index],
					)
				}
			}

			return result, nil
		default:
			return nil, fmt.Errorf(
				"SendFormAndWaitResponse: Unknown minecraft form type %T; minecraftForm = %#v",
				minecraftForm, minecraftForm,
			)
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
func (i *Interact) SendFormAndWaitResponse(minecraftForm form.MinecraftForm) (resp any, err error) {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.sendFormAndWaitResponse(minecraftForm)
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

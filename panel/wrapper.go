package panel_wrapper

import (
	"context"
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/function"
	"github.com/OmineDev/eulogist-user-interface/panel/panel"
	"github.com/OmineDev/eulogist-user-interface/server"
	"github.com/YingLunTown-DreamLand/gophertunnel/minecraft/protocol/packet"
	"github.com/pterm/pterm"
)

// PanelWrapper 是赞颂者前端的总包装者
type PanelWrapper struct {
	server   *server.Server
	interact *server.Interact
	function *function.Function
	panel    *panel.Panel
}

// MessageChannel 用于提供前端和后端的消息传递
func (p *PanelWrapper) MessageChannel() *function.MessageChannel {
	return p.panel.MessageChannel()
}

// MainEntry 在 realAddress 上打开一个假赞颂者服务器以供客户端连接。
// 当且仅当客户端连接到假赞颂者服务器后，此函数将返回值。
//
// 在 MainEntry 返回后，其调用者应当尽快调用 [MessageChannel] 并接
// 收来自前端的数据。如果 ctx 已经关闭，则意味着用户没有进入租赁服，
// 并且退出了前端所指示的赞颂者假服务器。此时 [MessageChannel] 将不会收到任何值。
//
// showAddress 指示显示给客户端的假服务器地址，
// bedrockVersion 是用户应该使用的基岩版版本
func MainEntry(realAddress string, showAddress string, bedrockVersion string) (
	wrapper *PanelWrapper,
	ctx context.Context,
	err error,
) {
	var cancelFunc context.CancelFunc

	wrapper = &PanelWrapper{server: server.NewServer()}
	ctx, cancelFunc = context.WithCancel(context.Background())

	err = wrapper.server.RunServer(realAddress)
	if err != nil {
		cancelFunc()
		return nil, nil, fmt.Errorf("MainEntry: %v", err)
	}
	pterm.Success.Printf(
		"赞颂者已经成功启动。\n请使用版本为 %s 的国际基岩版连接到该服务器: %s\n",
		bedrockVersion, showAddress,
	)

	err = wrapper.server.WaitConnect()
	if err != nil {
		cancelFunc()
		return nil, nil, fmt.Errorf("MainEntry: %v", err)
	}

	wrapper.interact = server.NewInteract(wrapper.server)
	wrapper.function, err = function.NewFunction(wrapper.interact)
	if err != nil {
		cancelFunc()
		return nil, nil, fmt.Errorf("MainEntry: %v", err)
	}
	wrapper.panel = panel.NewPanel(wrapper.function)

	go func() {
		// Cancel context
		defer cancelFunc()

		// Register or login
		for {
			if len(wrapper.function.EulogistConfig().EulogistToken) != 0 {
				break
			}

			exitGame, err := wrapper.function.RegisterOrLogin()
			select {
			case <-wrapper.server.MinecraftConn().Context().Done():
				return
			default:
			}

			if err != nil {
				_, _, _ = wrapper.interact.SendFormAndWaitResponse(form.MessageForm{
					Title:   "注册或登录失败",
					Content: fmt.Sprintf("%v", err),
					Button1: "确定",
					Button2: "继续",
				})
				continue
			}

			if exitGame {
				_ = wrapper.server.MinecraftConn().WritePacket(&packet.Disconnect{
					Message: "感谢您的使用, 期待下次与您相会！",
				})
				_ = wrapper.server.CloseServer()
				return
			}
		}

		// Request user info
		for {
			err = wrapper.function.RequestUserInfo(false)
			if err == nil {
				break
			}

			select {
			case <-wrapper.server.MinecraftConn().Context().Done():
				return
			default:
			}

			_, _, _ = wrapper.interact.SendFormAndWaitResponse(form.MessageForm{
				Title:   "请求赞颂者账户信息失败",
				Content: fmt.Sprintf("%v", err),
				Button1: "确定",
				Button2: "继续",
			})
		}

		// Show main panel
		for {
			exitGame, err := wrapper.panel.MainPanel()

			select {
			case <-wrapper.server.MinecraftConn().Context().Done():
				return
			default:
			}

			if err != nil {
				_, _, _ = wrapper.interact.SendFormAndWaitResponse(form.MessageForm{
					Title:   "错误",
					Content: fmt.Sprintf("%v", err),
					Button1: "确定",
					Button2: "继续",
				})
				continue
			}

			if exitGame {
				_ = wrapper.server.MinecraftConn().WritePacket(&packet.Disconnect{
					Message: "感谢您的使用, 期待下次与您相会！",
				})
				_ = wrapper.server.CloseServer()
				return
			}
		}
	}()

	return wrapper, ctx, nil
}

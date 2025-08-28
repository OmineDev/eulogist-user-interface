package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/server"
)

// CustomData 指示额外的自定义数据
type CustomData struct {
	useCustomSkin  bool
	customSkinData []byte
}

// NewCustomData 创建并返回一个新的 CustomData
func NewCustomData() *CustomData {
	return new(CustomData)
}

// Function 是基于 [server.Interact] 实现的赞颂者功能
type Function struct {
	interact   *server.Interact
	config     *define.LocalConfig
	userData   *define.EulogistUser
	customData *CustomData
	message    *MessageChannel
}

// NewFunction 根据 interact 创建并返回一个新的 Function
func NewFunction(interact *server.Interact) (result *Function, err error) {
	cfg, err := define.ReadConfig()
	if err != nil {
		return nil, fmt.Errorf("NewFunction: %v", err)
	}
	if cfg == nil {
		cfg = define.DefaultLocalConfig()
	}
	return &Function{
		interact:   interact,
		config:     cfg,
		message:    NewMessageChannel(),
		customData: NewCustomData(),
	}, nil
}

// Interact 返回表单交互实现
func (f *Function) Interact() *server.Interact {
	return f.interact
}

// EulogistConfig 返回赞颂者配置文件
func (f *Function) EulogistConfig() *define.LocalConfig {
	return f.config
}

// EulogistUserData 返回赞颂者用户数据
func (f *Function) EulogistUserData() *define.EulogistUser {
	return f.userData
}

// CustomData 返回其他额外的自定义数据
func (f *Function) CustomData() *CustomData {
	return f.customData
}

// MessageChannel 返回 [*MessageChannel]
func (f *Function) MessageChannel() *MessageChannel {
	return f.message
}

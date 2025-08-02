package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/server"
)

// Function 是基于 [server.Interact] 实现的赞颂者功能
type Function struct {
	interact *server.Interact
	config   *define.LocalConfig
	userData *define.EulogistUser
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
		interact: interact,
		config:   cfg,
	}, nil
}

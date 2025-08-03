package panel

import "github.com/OmineDev/eulogist-user-interface/function"

// Panel 是基于 [function.Function] 实现的赞颂者主要控制界面
type Panel struct {
	f *function.Function
}

// NewPanel 基于 function 创建并返回一个新的 Panel
func NewPanel(function *function.Function) *Panel {
	return &Panel{
		f: function,
	}
}

// Function 返回 [Panel] 底层中的 [*function.Function]
func (p *Panel) Function() *function.Function {
	return p.f
}

package function

import (
	"fmt"

	"github.com/OmineDev/eulogist-user-interface/define"
)

// Logout 清除底层记录的数据，从而使用户登出
func (f *Function) Logout() error {
	f.userData = nil
	f.config = new(define.LocalConfig)

	err := define.WriteConfig(f.config)
	if err != nil {
		return fmt.Errorf("Logout: %v", err)
	}

	return nil
}

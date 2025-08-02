package define

import (
	"encoding/json"
	"fmt"
	"os"
)

// LocalConfig 是赞颂者本地储存的配置数据
type LocalConfig struct {
	EulogistToken    string             `json:"eulogist_token"`
	LastServerConfig RentalServerConfig `json:"last_server_config"`
}

// DefaultLocalConfig 创建并返回一个非 nil 的默认赞颂者配置
func DefaultLocalConfig() *LocalConfig {
	return new(LocalConfig)
}

// ReadConfig 从赞颂者配置文件阅读对应的配置。
// 如果配置文件不存在，则返回的 config 为 nil
func ReadConfig() (config *LocalConfig, err error) {
	var result LocalConfig

	file, _ := os.ReadFile(EulogistConfigFileName)
	if len(file) == 0 {
		return nil, nil
	}

	err = json.Unmarshal(file, &result)
	if err != nil {
		return nil, fmt.Errorf("ReadConfig: %v", err)
	}

	return &result, nil
}

// WriteConfig 将 config 指示的赞颂者配置写入到对应的配置文件。
// 如果给出的 config 为 nil，则相应的配置文件将被删除
func WriteConfig(config *LocalConfig) error {
	if config == nil {
		if err := os.Remove(EulogistConfigFileName); err != nil {
			return fmt.Errorf("WriteConfig: %v", err)
		}
	}

	jsonBytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return fmt.Errorf("WriteConfig: %v", err)
	}

	err = os.WriteFile(EulogistConfigFileName, jsonBytes, 0600)
	if err != nil {
		return fmt.Errorf("WriteConfig: %v", err)
	}

	return nil
}

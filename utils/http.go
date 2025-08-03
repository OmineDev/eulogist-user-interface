package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SendAndGetHttpResponse 向 url 发送 request 的 POST 请求，
// 并解析服务端的 JSON 响应为 result
func SendAndGetHttpResponse[T any](url string, request any) (result T, err error) {
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		err = fmt.Errorf("SendAndGetHttpResponse: %v", err)
		return
	}

	buf := bytes.NewBuffer(jsonBytes)
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		err = fmt.Errorf("SendAndGetHttpResponse: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("SendAndGetHttpResponse: Status code (%d) is not 200", resp.StatusCode)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("SendAndGetHttpResponse: %v", err)
		return
	}

	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		err = fmt.Errorf("SendAndGetHttpResponse: %v", err)
		return
	}

	return
}

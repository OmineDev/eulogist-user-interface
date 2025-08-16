package function

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/OmineDev/eulogist-user-interface/define"
	"github.com/OmineDev/eulogist-user-interface/form"
	"github.com/OmineDev/eulogist-user-interface/utils"
)

// GameSavesKeyRequest ..
type GameSavesKeyRequest struct {
	Token              string `json:"token,omitempty"`
	RentalServerNumber string `json:"rental_server_number,omitempty"`
}

// GameSavesKeyResponse ..
type GameSavesKeyResponse struct {
	ErrorInfo              string `json:"error_info"`
	Success                bool   `json:"success"`
	RentelServerNumber     string `json:"rental_server_number"`
	GameSavesAESCipher     []byte `json:"game_saves_aes_cipher"`
	DisableOpertorVerify   bool   `json:"disable_operator_verify"`
	ResponseExpireUnixTime int64  `json:"response_expire_unix_time"`
}

// GetGameSavesKey ..
func (f *Function) GetGameSavesKey() error {
	minecraftForm := form.ModalForm{
		Title: "获取存档加密密钥",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: "" +
					"您将获取存档的§r§e加密密钥§r。\n" +
					"您必须得到租赁服§r§e有关人员§r的授权。\n" +
					"请输入§r§e租赁服号§r以获取对应租赁服的存档加密密钥。",
			},
			form.ModalFormElementInput{
				Text:        "租赁服号",
				Default:     "",
				PlaceHolder: "要获取存档密钥的租赁服号",
			},
		},
	}

	resp, isUserCancel, err := f.interact.SendFormAndWaitResponse(minecraftForm)
	if err != nil {
		return fmt.Errorf("GetGameSavesKey: %v", err)
	}
	if isUserCancel {
		return nil
	}
	rentalServerNumber := resp.([]any)[1].(string)

	gameSavesKeyResp, err := utils.SendAndGetHttpResponse[GameSavesKeyResponse](
		fmt.Sprintf("%s/get_game_saves_key", define.StdAuthServerAddress),
		GameSavesKeyRequest{
			Token:              f.config.EulogistToken,
			RentalServerNumber: rentalServerNumber,
		},
	)
	if err != nil {
		return fmt.Errorf("GetGameSavesKey: %v", err)
	}
	if !gameSavesKeyResp.Success {
		_, _, err := f.interact.SendFormAndWaitResponse(form.MessageForm{
			Title:   "错误",
			Content: gameSavesKeyResp.ErrorInfo,
			Button1: "确定",
			Button2: "返回上一级菜单",
		})
		if err != nil {
			return fmt.Errorf("GetGameSavesKey: %v", err)
		}
		return nil
	}

	_, _, err = f.interact.SendFormAndWaitResponse(form.ModalForm{
		Title: "存档密钥",
		Contents: []form.ModalFormElement{
			form.ModalFormElementLabel{
				Text: fmt.Sprintf(
					""+
						"您已成功获得租赁服 §r§b%s§r 的§r§e存档解密密钥§r, \n"+
						"它只对您§r§e当前§r赞颂者账户有效。\n"+
						"请务必§r§e妥善保管§r, 不要遗失！",
					rentalServerNumber,
				),
			},
			form.ModalFormElementInput{
				Text:        "AES 密钥 (16位, Hex 字符串)",
				Default:     hex.EncodeToString(gameSavesKeyResp.GameSavesAESCipher),
				PlaceHolder: "AES Cipher (Hex string)",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("GetGameSavesKey: %v", err)
	}

	return nil
}

// BeforePlayPrepare ..
func (f *Function) BeforePlayPrepare(rentalServerNumber string) (
	providedPeAuthData string,
	aesCipher []byte,
	disableOpertorVerify bool,
	err error,
) {
	var gameSavesKeyResp GameSavesKeyResponse
	request := GameSavesKeyRequest{
		Token:              f.config.EulogistToken,
		RentalServerNumber: rentalServerNumber,
	}

	jsonBytes, err := json.Marshal(request)
	if err != nil {
		return "", nil, disableOpertorVerify, fmt.Errorf("BeforePlayPrepare: %v", err)
	}

	encrypted, err := utils.EncryptPKCS1v15(&define.GameSavesEncryptKey.PublicKey, jsonBytes)
	if err != nil {
		return "", nil, false, fmt.Errorf("BeforePlayPrepare: %v", err)
	}

	buf := bytes.NewBuffer(encrypted)
	resp, err := http.Post(
		fmt.Sprintf("%s/get_game_saves_key", define.StdAuthServerAddress),
		"application/json",
		buf,
	)
	if err != nil {
		err = fmt.Errorf("BeforePlayPrepare: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, false, fmt.Errorf("BeforePlayPrepare: Status code (%d) is not 200", resp.StatusCode)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, false, fmt.Errorf("BeforePlayPrepare: %v", err)
	}

	decrypted, err := utils.DecryptPKCS1v15(define.GameSavesEncryptKey, bodyBytes)
	if err == nil {
		bodyBytes = decrypted
	}

	err = json.Unmarshal(bodyBytes, &gameSavesKeyResp)
	if err != nil {
		return "", nil, false, fmt.Errorf("BeforePlayPrepare: %v", err)
	}

	if !gameSavesKeyResp.Success {
		return "", nil, false, fmt.Errorf("BeforePlayPrepare: Failed to get game saves key due to %v", gameSavesKeyResp.ErrorInfo)
	}
	if time.Now().Unix() >= gameSavesKeyResp.ResponseExpireUnixTime {
		return "", nil, false, fmt.Errorf("BeforePlayPrepare: Unsuccessful hacking attempt (mark 0)")
	}
	if gameSavesKeyResp.RentelServerNumber != rentalServerNumber {
		return "", nil, false, fmt.Errorf("BeforePlayPrepare: Unsuccessful hacking attempt (mark 1)")
	}

	providedPeAuthData = f.userData.ProvidedPeAuthData
	f.userData.ProvidedPeAuthData = ""
	return providedPeAuthData, gameSavesKeyResp.GameSavesAESCipher, gameSavesKeyResp.DisableOpertorVerify, nil
}

package define

import (
	"crypto/rsa"
	"crypto/x509"
	_ "embed"
	"encoding/pem"
)

const (
	EulogistConfigFileName = "eulogist_config.json"
	UserPasswordSlat       = "YoRHa"
	DefaultPageSize        = 5
)

const (
	StdAuthServerPhoenixAPI = "http://127.0.0.1:8080"
	StdAuthServerAddress    = "http://127.0.0.1:8080/eulogist_api"
)

const (
	AuthServerAccountTypeStd uint8 = iota
	AuthServerAccountTypeCustom
)

const (
	UserPermissionSystem = iota
	UserPermissionAdmin
	UserPermissionManager
	UserPermissionNormal
	UserPermissionNone
	UserPermissionDefault = UserPermissionNormal
)

//go:embed game_saves_encrypt.key
var keyBytes []byte
var GameSavesEncryptKey *rsa.PrivateKey

func init() {
	var err error
	keyBlock, _ := pem.Decode(keyBytes)
	GameSavesEncryptKey, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		panic(err)
	}
}

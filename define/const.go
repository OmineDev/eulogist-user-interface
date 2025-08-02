package define

const EulogistConfigFileName = "eulogist_config.json"

const (
	StdAuthServerAddress = "http://127.0.0.1:8080/eulogist_api"
)

const UserPasswordSlat = "YoRHa"

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

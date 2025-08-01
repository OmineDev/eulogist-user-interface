package form

const (
	FormTypeMessage uint8 = iota
	FormTypeAction
	FormTypeModal
)

// MinecraftForm 是各种表单的类型的总称
type MinecraftForm interface {
	ID() uint8
	PackToJSON() string
}

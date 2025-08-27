package function

// MessageFromFronted ..
type MessageFromFronted struct {
	EulogistUniqueID string

	RentalServerNumber   string
	RentalServerPasscode string

	AuthServerAddress  string
	AuthServerToken    string
	ProvidedPeAuthData string

	GameSavesAESCipher   []byte
	DisableOpertorVerify bool

	UseCustomSkin  bool
	CustomSkinData []byte

	HaveSkinCacheData bool
	SkinDownloadURL   string
}

// MessageFromBacked ..
type MessageFromBacked struct {
	CanTerminate         bool
	LoginServerMeetError bool
	LoginServrErrorInfo  string
	TransferAddress      string
	TransferPort         uint16
}

// MessageChannel ..
type MessageChannel struct {
	msgFromFronted chan MessageFromFronted
	msgFromBacked  chan MessageFromBacked
}

// NewMessageChannel ..
func NewMessageChannel() *MessageChannel {
	return &MessageChannel{
		msgFromFronted: make(chan MessageFromFronted, 1),
		msgFromBacked:  make(chan MessageFromBacked, 1),
	}
}

// NotifyToFronted ..
func (m *MessageChannel) NotifyToFronted(message MessageFromBacked) {
	m.msgFromBacked <- message
}

// NotifyToBacked ..
func (m *MessageChannel) NotifyToBacked(message MessageFromFronted) {
	m.msgFromFronted <- message
}

// ReceiveFrontedMessage ..
func (m *MessageChannel) FrontedMessageChannel() <-chan MessageFromFronted {
	return m.msgFromFronted
}

// ReceiveBackedMessage ..
func (m *MessageChannel) BackedMessageChannel() <-chan MessageFromBacked {
	return m.msgFromBacked
}

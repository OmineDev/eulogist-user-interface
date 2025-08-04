package define

import "github.com/ToolDelta-Basic/gophertunnel/minecraft/protocol"

type RentalServerConfig struct {
	ServerNumber   string `json:"server_number"`
	ServerPassCode string `json:"server_passcode"`
}

func (r *RentalServerConfig) Marshal(io protocol.IO) {
	io.String(&r.ServerNumber)
	io.String(&r.ServerPassCode)
}

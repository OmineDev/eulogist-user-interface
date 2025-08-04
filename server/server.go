package server

import (
	"bytes"
	_ "embed"
	"fmt"
	"sync"

	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/resource"
)

//go:embed depends.mcpack
var dependsResourcesPack []byte

// Server 简单的实现了一个 MC 服务器，
// 以用于运行一个赞颂者前置交互服务
type Server struct {
	mu       *sync.Mutex
	closed   bool
	listener *minecraft.Listener
	conn     *minecraft.Conn
}

// NewServer 创建并返回一个新的 Server
func NewServer() *Server {
	return &Server{
		mu:       new(sync.Mutex),
		closed:   false,
		listener: nil,
		conn:     nil,
	}
}

// RunServer 在 address 所指示的地址上运行服务器
func (s *Server) RunServer(address string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return fmt.Errorf("RunServer: Server has been closed")
	}

	pack, err := resource.Read(bytes.NewBuffer(dependsResourcesPack))
	if err != nil {
		return fmt.Errorf("RunServer: %v", err)
	}

	config := minecraft.ListenConfig{
		AllowUnknownPackets: true,
		StatusProvider: minecraft.NewStatusProvider(
			"Eulogist", "Eulogist",
		),
		ResourcePacks:        []*resource.Pack{pack},
		TexturePacksRequired: true,
	}

	listener, err := config.Listen("raknet", address)
	if err != nil {
		return fmt.Errorf("RunServer: %v", err)
	}

	s.listener, s.conn = listener, nil
	return nil
}

// WaitConnect 等待一个 MC 客户端创建连接
func (s *Server) WaitConnect() error {
	s.mu.Lock()
	if s.listener == nil {
		s.mu.Unlock()
		return fmt.Errorf("RunServer: Server is not initialized")
	}
	if s.closed {
		s.mu.Unlock()
		return fmt.Errorf("RunServer: Server has been closed")
	}
	s.mu.Unlock()

	netConn, err := s.listener.Accept()
	if err != nil {
		return fmt.Errorf("WaitConnect: %v", err)
	}
	s.conn = netConn.(*minecraft.Conn)

	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				return
			}
			_ = conn.Close()
		}
	}()

	err = s.conn.StartGame(
		minecraft.GameData{
			WorldName:           "Eulogist-User-Interface",
			WorldSeed:           0,
			Difficulty:          0, // Peaceful
			EntityUniqueID:      0,
			EntityRuntimeID:     0,
			PlayerGameMode:      1, // Creative
			PersonaDisabled:     false,
			CustomSkinsDisabled: false,
			EmoteChatMuted:      true,
			BaseGameVersion:     "*",
			PlayerPosition:      [3]float32{0.5, 1.5, 0.5},
			Pitch:               0,
			Yaw:                 0,
			Dimension:           0, // Overworld
			WorldSpawn:          [3]int32{0, 0, 0},
			EditorWorldType:     packet.EditorWorldTypeNotEditor,
			CreatedInEditor:     false,
			WorldGameMode:       1, // Creative
			Hardcore:            false,
			GameRules: []protocol.GameRule{
				{
					Name:  "doDayLightCycle",
					Value: false,
				},
			},
			Time:                     0,
			ServerBlockStateChecksum: 0,
			CustomBlocks:             nil,
			Items:                    nil,
			PlayerMovementSettings: protocol.PlayerMovementSettings{
				RewindHistorySize:                0,
				ServerAuthoritativeBlockBreaking: false,
			},
			ServerAuthoritativeInventory: true,
			Experiments:                  nil,
			PlayerPermissions:            1, // Member
			ChunkRadius:                  4,
			ClientSideGeneration:         false,
			ChatRestrictionLevel:         packet.ChatRestrictionLevelDisabled,
			DisablePlayerInteractions:    true,
			UseBlockNetworkIDHashes:      false,
		},
	)
	if err != nil {
		return fmt.Errorf("WaitConnect: %v", err)
	}

	return nil
}

// MinecraftConn 返回已创建的连接
func (s *Server) MinecraftConn() *minecraft.Conn {
	return s.conn
}

// CloseServer 关闭正在运行的服务器。
// CloseServer 调用后不应再次调用 RunServer
func (s *Server) CloseServer() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return fmt.Errorf("CloseServer: %v", err)
		}
	}
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("CloseServer: %v", err)
		}
	}

	s.closed = true
	return nil
}

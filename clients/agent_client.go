package clients

import (
	"github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/client"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/config"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
)

type AgentServerClient struct {
	PacketChannel chan network.Packet
	Connected     bool
	Config        config.Config
	*client.BackendClient
}

type LoginTokenRequest struct {
	Username  string
	Password  string
	Token     uint32
	AccountID uint32
	ShardID   uint16
}

func NewAgentServerClient(config config.Config) *AgentServerClient {
	c := client.NewBackendClient(
		config.AgentServer.IP,
		config.AgentServer.Port,
		config.GatewayServer.ModuleID,
		config.GatewayServer.Secret)
	logrus.Infoln("connecting to agent server")
	c.Connect()
	ac := AgentServerClient{Config: config, PacketChannel: c.IncomingPacketChannel, BackendClient: c}
	return &ac
}

func (ac *AgentServerClient) SendLoginToken(request LoginTokenRequest) {
	p := network.EmptyPacket()
	p.MessageID = opcode.GatewayLoginTokenRequest
	p.WriteUInt32(request.AccountID)
	p.WriteString(request.Username)
	p.WriteString(request.Password)
	p.WriteUInt32(request.Token)
	p.WriteUInt16(request.ShardID)
	ac.OutgoingPacketChannel <- p
}

package main

import (
	"bufio"
	"github.com/ferdoran/go-sro-framework/logging"
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/ferdoran/go-sro-gateway-server/clients"
	"github.com/ferdoran/go-sro-gateway-server/config"
	"github.com/ferdoran/go-sro-gateway-server/handler"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"math/rand"
	"os"
	"time"
)

type GatewayServer struct {
	Server            server.Server
	failedLogins      map[string]int
	AgentServerClient *clients.AgentServerClient
}

func NewGatewayServer() GatewayServer {
	s := server.NewEngine(
		viper.GetString(config.GatewayHost),
		viper.GetInt(config.GatewayPort),
		network.EncodingOptions{
			None:         false,
			Disabled:     false,
			Encryption:   true,
			EDC:          true,
			KeyExchange:  true,
			KeyChallenge: false,
		},
	)
	s.ModuleID = viper.GetString(config.GatewayModuleId)

	agentClient := clients.NewAgentServerClient()
	return GatewayServer{s, make(map[string]int), agentClient}
}

func (g *GatewayServer) Start() {
	go g.Server.Start()
	g.handlePackets()
}

func (g *GatewayServer) handlePackets() {

	handler.NewPatchRequestHandler()
	handler.NewShardlistRequestHandler()
	handler.NewShardlistPingHandler()
	handler.NewLoginRequestHandler(g.failedLogins, g.AgentServerClient)
	handler.NewNoticeRequestHandler()
	for {
		data := <-g.Server.PacketChannel
		switch data.MessageID {
		default:
			log.Debugf("Unhandled packet %v\n", data.Packet)
		}
	}
}

func main() {
	config.Initialize()
	logging.Init()

	log.Info("Agent Public IP: %s", viper.GetString(config.AgentPublicIp))
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	log.Infoln("Loading Config")
	log.Infoln("Starting server...")
	gw := NewGatewayServer()

	gw.Start()
	log.Println("Press Enter to exit...")
	reader.ReadString('\n')
}

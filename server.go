package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/config"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/logging"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
	"gitlab.ferdoran.de/game-dev/go-sro/gateway-server/clients"
	"gitlab.ferdoran.de/game-dev/go-sro/gateway-server/handler"
	"math/rand"
	"os"
	"time"
)

type GatewayServer struct {
	Server            server.Server
	Config            config.Config
	failedLogins      map[string]int
	AgentServerClient *clients.AgentServerClient
}

func NewGatewayServer(config config.Config) GatewayServer {
	server := server.NewEngine(
		config.GatewayServer.IP,
		config.GatewayServer.Port,
		network.EncodingOptions{
			None:         false,
			Disabled:     false,
			Encryption:   true,
			EDC:          true,
			KeyExchange:  true,
			KeyChallenge: false,
		},
		config)
	server.ModuleID = config.GatewayServer.ModuleID
	agentClient := clients.NewAgentServerClient(config)
	return GatewayServer{server, config, make(map[string]int), agentClient}
}

func (g *GatewayServer) Start() {
	go g.Server.Start()
	g.handlePackets()
}

func (g *GatewayServer) handlePackets() {

	handler.NewPatchRequestHandler()
	handler.NewShardlistRequestHandler()
	handler.NewShardlistPingHandler()
	handler.NewLoginRequestHandler(g.Config, g.failedLogins, g.AgentServerClient)
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
	logging.Init()
	rand.Seed(time.Now().UnixNano())
	reader := bufio.NewReader(os.Stdin)

	log.Infoln("Loading Config")
	config.LoadConfig("config.json")
	log.Infoln("Starting server...")
	gw := NewGatewayServer(config.GlobalConfig)

	gw.Start()
	log.Println("Press Enter to exit...")
	reader.ReadString('\n')
}

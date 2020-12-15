package handler

import (
	log "github.com/sirupsen/logrus"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
	"gitlab.ferdoran.de/game-dev/go-sro/gateway-server/db"
)

type ShardlistRequestHandler struct {
}

func NewShardlistRequestHandler() server.PacketHandler {
	handler := ShardlistRequestHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.ShardlistRequest, handler)
	return handler
}

func (h ShardlistRequestHandler) Handle(packet server.PacketChannelData) {
	log.Println("CLIENT_GATEWAY_SHARD_LIST_REQUEST")
	shards := db.GetShards()
	// TODO make this configurable
	// Farm
	p := network.EmptyPacket()
	p.MessageID = opcode.ShardlistResponse

	p.WriteByte(1)                   // hasEntries
	p.WriteByte(1)                   // Farm.ID
	p.WriteString("SRO_DEUTSCHLAND") //Farm.Name

	p.WriteByte(0) // Divider / End of Farms

	// Shards
	for _, shard := range shards {
		p.WriteByte(1)                             // hasEntries
		p.WriteUInt16(uint16(shard.Id))            // Shard.ID
		p.WriteString(shard.Name)                  // Shard.Name
		p.WriteUInt16(uint16(shard.OnlinePlayers)) // Shard.OnlinePlayers
		p.WriteUInt16(uint16(shard.Capacity))      // Shard.PlayerCapacity
		p.WriteByte(byte(shard.Status))            // Shard.IsOperating
		p.WriteByte(1)                             // Shard.FarmID
	}
	p.WriteByte(0) // Divider / End of Shards

	packet.Session.Conn.Write(p.ToBytes())
}

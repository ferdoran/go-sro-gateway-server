package handler

import (
	"github.com/ferdoran/go-sro-framework/network"
	"github.com/ferdoran/go-sro-framework/network/opcode"
	"github.com/ferdoran/go-sro-framework/server"
	"github.com/ferdoran/go-sro-framework/utils"
)

type ShardlistPingHandler struct {
}

func NewShardlistPingHandler() server.PacketHandler {
	handler := ShardlistPingHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.ShardlistPing, handler)
	return handler
}

func (h ShardlistPingHandler) Handle(packet server.PacketChannelData) {
	p := network.EmptyPacket()
	p.MessageID = opcode.ShardlistPong
	// TODO can also return an error code
	p.WriteByte(1) // result = 1 = Successful, else error
	p.WriteByte(1) // result Farm.ID
	// TODO how do we want to pass the Farm.IP here?
	p.WriteUInt32(utils.ByteArrayToUint32([]byte{127, 0, 0, 1})) // Farm.IP

	packet.Session.Conn.Write(p.ToBytes())
}

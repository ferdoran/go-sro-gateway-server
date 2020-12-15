package handler

import (
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
)

type PatchRequestHandler struct{}

func NewPatchRequestHandler() server.PacketHandler {
	handler := PatchRequestHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.PatchRequest, handler)
	return handler
}

func (h PatchRequestHandler) Handle(packet server.PacketChannelData) {
	// TODO an actual check for patches
	// TODO add patch error codes
	//  InvalidVersion = 1,
	//  UPDATE = 2,
	//  NotInService = 3,
	//  AbnormalModule = 4,
	//  PatchDisabled = 5,
	p := network.EmptyPacket()
	p.MessageID = 0x600D
	p.WriteByte(1)
	p.WriteUInt16(1)
	p.WriteUInt16(opcode.PatchResponse)

	packet.Session.Conn.Write(p.ToBytes())
	p = network.EmptyPacket()
	p.MessageID = 0x600D
	p.WriteByte(0)
	p.WriteUInt16(1)

	packet.Session.Conn.Write(p.ToBytes())
}

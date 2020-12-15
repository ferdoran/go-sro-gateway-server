package handler

import (
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/network/opcode"
	"gitlab.ferdoran.de/game-dev/go-sro/framework/server"
	"time"

	"gitlab.ferdoran.de/game-dev/go-sro/gateway-server/db"
)

type NoticeRequestHandler struct {
}

func NewNoticeRequestHandler() server.PacketHandler {
	handler := NoticeRequestHandler{}
	server.PacketManagerInstance.RegisterHandler(opcode.NoticeRequest, handler)
	return handler
}

func (h NoticeRequestHandler) Handle(packet server.PacketChannelData) {
	p := network.EmptyPacket()
	p.MessageID = 0x600D
	p.WriteByte(1)
	p.WriteUInt16(1)
	p.WriteUInt16(opcode.NoticeResponse)
	packet.Session.Conn.Write(p.ToBytes())

	notices := db.GetNotices()
	p = network.EmptyPacket()
	p.MessageID = 0x600D
	p.WriteByte(0)
	p.WriteByte(byte(len(notices)))
	for _, notice := range notices {
		p.WriteString(notice.Subject)
		p.WriteString(notice.Article)
		p.WriteUInt16(uint16(notice.Ctime.Year()))
		p.WriteUInt16(uint16(notice.Ctime.Month()))
		p.WriteUInt16(uint16(notice.Ctime.Day()))
		p.WriteUInt16(uint16(notice.Ctime.Hour()))
		p.WriteUInt16(uint16(notice.Ctime.Minute()))
		p.WriteUInt16(uint16(notice.Ctime.Second()))
		p.WriteUInt32(uint32((notice.Ctime.Nanosecond()) / int(time.Millisecond)))
	}

	packet.Session.Conn.Write(p.ToBytes())
}

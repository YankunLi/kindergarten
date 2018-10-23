package fly

import (
	"log"
	"net"
)

type MsgHeadT struct {
	Tag     uint8
	MType   uint8
	Length  uint32
	SrcName [256]byte
	DesName [256]byte
	//  CRC32 uint32
}

type MsgACKT struct {
	Tag    uint8
	MType  uint8
	Code   uint8
	Length uint32
}

type ConfType struct {
	SrcFilePath *string
	DesFileName *string
	ConfigFile  *string
	LogFile     *string
	RemoteSer   *string
	CMD         string

	Logger *log.Logger
}

type FlyContext struct {
	Config *ConfType
	Client *NetClient
}

type NetClient struct {
	Conn net.Conn
}

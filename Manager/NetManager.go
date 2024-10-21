package manager

import (
	"encoding/binary"
	"log"

	pb "testServer/Messages"

	"google.golang.org/protobuf/proto"
)

type NetManager struct {
}

var netManager *NetManager

func GetNetManager() *NetManager {
	if netManager == nil {
		netManager = &NetManager{}
	}

	return netManager
}

func (nm *NetManager) MakePacket(msg *pb.GameMessage) []byte {
	response, err := proto.Marshal(msg)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return nil
	}

	lengthBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(lengthBuf, uint32(len(response)))
	lengthBuf = append(lengthBuf, response...)

	return lengthBuf
}

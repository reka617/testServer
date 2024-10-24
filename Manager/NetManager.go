package manager

import (
	"encoding/binary"
	"log"

	pb "testServer/Messages"

	"google.golang.org/protobuf/proto"
)

type NetworkManager struct {
}

var netManager *NetworkManager

func GetNetManager() *NetworkManager {
	if netManager == nil {
		netManager = &NetworkManager{}
	}

	return netManager
}

func (nm *NetworkManager) MakePacket(msg *pb.GameMessage) []byte {
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

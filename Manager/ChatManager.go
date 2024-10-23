package manager

import (
	"encoding/binary"
	"log"

	pb "testServer/Messages"
	behave "testServer/behavior"

	"google.golang.org/protobuf/proto"
)

var chatManager *ChatManager

// PlayerManager manages a list of playersH
type ChatManager struct {
	players map[int]behave.Player
	nextID  int
}

// NewPlayerManager creates a new PlayerManager
func GetChatManager() *ChatManager {
	if chatManager == nil {
		chatManager = &ChatManager{
			players: make(map[int]behave.Player),
			nextID:  1,
		}
	}
	return chatManager
}

// AddPlayer adds a new player to the manager
func (pm *ChatManager) Broadcast(name string, content string) {
	for _, player := range GetPlayerManager().ListPlayers() {
		gameMessage := &pb.GameMessage{
			Message: &pb.GameMessage_Chat{
				Chat: &pb.ChatMessage{
					Sender:  name,
					Content: content,
				},
			},
		}
		response, err := proto.Marshal(gameMessage)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			continue
		}

		lengthBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(lengthBuf, uint32(len(response)))
		lengthBuf = append(lengthBuf, response...)
		(*player.Conn).Write(lengthBuf)
	}
}

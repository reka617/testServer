package manager

import (
	"encoding/binary"
	"errors"
	"log"
	"net"

	pb "testServer/Messages"
	behave "testServer/behavior"

	"google.golang.org/protobuf/proto"
)

var playerManager *PlayerManager

// PlayerManager manages a list of players
type PlayerManager struct {
	players map[string]*behave.Player
	nextID  int
}

// NewPlayerManager creates a new PlayerManager
func GetPlayerManager() *PlayerManager {
	if playerManager == nil {
		playerManager = &PlayerManager{
			players: make(map[string]*behave.Player),
			nextID:  1,
		}
	}

	return playerManager
}

// AddPlayer adds a new player to the manager
func (pm *PlayerManager) AddPlayer(name string, age int, conn *net.Conn) *behave.Player {
	player := behave.Player{
		ID:        pm.nextID,
		Name:      name,
		Age:       age,
		Conn:      conn,
		X:         0,
		Y:         0,
		Z:         0,
		RotationY: 0,
	}

	pm.players[name] = &player
	pm.nextID++

	// 내가 로그인 되었음을 나한테 알려준다.
	myPlayerSapwn := &pb.GameMessage{
		Message: &pb.GameMessage_SpawnMyPlayer{
			SpawnMyPlayer: &pb.SpawnMyPlayer{
				X:         player.X,
				Y:         player.Y,
				Z:         player.Z,
				RotationY: player.RotationY,
			},
		},
	}

	pathTest := &pb.GameMessage{
		Message: &pb.GameMessage_PathTest{
			PathTest: &pb.PathTest{},
		},
	}

	path, err := GetNavMeshManager().PathFinding(-230, 0, -291, 235, 0, 180)
	if err == nil {
		for _, path := range path.PathList {
			pathTest.GetPathTest().Paths = append(pathTest.GetPathTest().Paths, &pb.NavV3{X: float32(path.X), Y: float32(path.Y), Z: float32(path.Z)})
		}

		response := GetNetManager().MakePacket(pathTest)
		(*player.Conn).Write(response)
	}

	response := GetNetManager().MakePacket(myPlayerSapwn)
	(*player.Conn).Write(response)

	otherPlayerSpawnPacket := &pb.GameMessage{
		Message: &pb.GameMessage_SpawnOtherPlayer{
			SpawnOtherPlayer: &pb.SpawnOtherPlayer{
				PlayerId:  name,
				X:         player.X,
				Y:         player.Y,
				Z:         player.Z,
				RotationY: player.RotationY,
			},
		},
	}

	// 이 코드를 들어온 유저를 제외한 플레이어들에게 스폰시켜달라고 한다.
	for _, p := range pm.players {
		if p.Name == name {
			continue
		}

		response = GetNetManager().MakePacket(otherPlayerSpawnPacket)

		(*p.Conn).Write(response)
	}

	// 다른 플레이어의 위치정보를 접속한 인원에게 보낸다.
	for _, p := range pm.players {
		if p.Name == name {
			continue
		}

		otherPlayerSpawnPacket := &pb.GameMessage{
			Message: &pb.GameMessage_SpawnOtherPlayer{
				SpawnOtherPlayer: &pb.SpawnOtherPlayer{
					PlayerId:  p.Name,
					X:         p.X,
					Y:         p.Y,
					Z:         p.Z,
					RotationY: player.RotationY,
				},
			},
		}

		response = GetNetManager().MakePacket(otherPlayerSpawnPacket)

		(*player.Conn).Write(response)
	}

	return &player
}

func (pm *PlayerManager) MovePlayer(p *pb.GameMessage_PlayerPosition) {

	pm.players[p.PlayerPosition.PlayerId].X = p.PlayerPosition.X
	pm.players[p.PlayerPosition.PlayerId].Y = p.PlayerPosition.Y
	pm.players[p.PlayerPosition.PlayerId].Z = p.PlayerPosition.Z
	pm.players[p.PlayerPosition.PlayerId].RotationY = p.PlayerPosition.RotationY

	response, err := proto.Marshal(&pb.GameMessage{
		Message: p,
	})

	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return
	}

	for _, player := range pm.players {
		if player.Name == p.PlayerPosition.PlayerId {
			continue
		}

		lengthBuf := make([]byte, 4)
		binary.LittleEndian.PutUint32(lengthBuf, uint32(len(response)))
		lengthBuf = append(lengthBuf, response...)
		(*player.Conn).Write(lengthBuf)
	}
}

// GetPlayer retrieves a player by ID
func (pm *PlayerManager) GetPlayer(id string) (*behave.Player, error) {
	player, exists := pm.players[id]
	if !exists {
		return nil, errors.New("player not found")
	}
	return player, nil
}

// RemovePlayer removes a player by ID
func (pm *PlayerManager) RemovePlayer(id string) error {
	if _, exists := pm.players[id]; !exists {
		return errors.New("player not found")
	}
	delete(pm.players, id)

	logoutPacket := &pb.GameMessage{
		Message: &pb.GameMessage_Logout{
			Logout: &pb.LogoutMessage{
				PlayerId: id,
			},
		},
	}

	response := GetNetManager().MakePacket(logoutPacket)

	// 이 코드를 들어온 유저를 제외한 플레이어들에게 스폰시켜달라고 한다.
	for _, p := range pm.players {
		(*p.Conn).Write(response)
	}

	return nil
}

// ListPlayers returns all players in the manager
func (pm *PlayerManager) ListPlayers() []*behave.Player {
	playerList := []*behave.Player{}
	for _, player := range pm.players {
		playerList = append(playerList, player)
	}
	return playerList
}

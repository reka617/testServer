package manager

import (
	"testServer/behavior"

	pb "testServer/Messages"
)

var monsterManager *MonsterManager

// PlayerManager manages a list of players
type MonsterManager struct {
	monsters map[string]*behavior.Monster
	nextID   int32
}

// NewPlayerManager creates a new PlayerManager
func GetMonsterManager() *MonsterManager {
	if monsterManager == nil {
		monsterManager = &MonsterManager{
			monsters: make(map[string]*behavior.Monster),
			nextID:   1,
		}
	}

	return monsterManager
}

// AddPlayer adds a new player to the manager
func (mm *MonsterManager) AddMonster(id int32) *behavior.Monster {
	monster := behavior.Monster{
		MonsterId: mm.nextID,
		X:         0,
		Z:         0,
	}

	mm.monsters[id] = &monster
	mm.nextID++

	// 내가 로그인 되었음을 나한테 알려준다.
	MonsterSapwn := &pb.GameMessage{
		Message: &pb.GameMessage_SpawnMonster{
			SpawnMonster: &pb.SpawnMonster{
				X:         monster.X,
				Z:         monster.Z,
				MonsterId: monster.MonsterId,
			},
		},
	}

	// 이 코드를 들어온 유저를 제외한 플레이어들에게 스폰시켜달라고 한다.
	for _, p := range GetPlayerManager().players {

		response := GetNetManager().MakePacket(MonsterSapwn)
		(*p.Conn).Write(response)
	}

	return &monster
}

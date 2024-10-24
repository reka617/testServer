package manager

import (
	pb "testServer/Messages"
	"testServer/common"
	"time"
)

var monsterManager *MonsterManager

// PlayerManager manages a list of players
type MonsterManager struct {
	monsters map[int32]*Monster
	nextID   int32
}

// NewPlayerManager creates a new PlayerManager
func GetMonsterManager() *MonsterManager {
	if monsterManager == nil {
		monsterManager = &MonsterManager{
			monsters: make(map[int32]*Monster),
			nextID:   1,
		}
	}

	return monsterManager
}

func (mm *MonsterManager) UpdateMonster() {

	ticker := time.NewTicker(16 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			for _, m := range mm.monsters {
				m.AI.Execute()
			}
		}
	}
}

// AddPlayer adds a new player to the manager
func (mm *MonsterManager) AddMonster(id int32) *Monster {
	path := make([]common.Point, 0)

	monster := NewMonster(0, 0, 0, 100, path)

	mm.monsters[id] = monster
	mm.nextID++

	// 내가 로그인 되었음을 나한테 알려준다.
	MonsterSapwn := &pb.GameMessage{
		Message: &pb.GameMessage_SpawnMonster{
			SpawnMonster: &pb.SpawnMonster{
				X:         monster.X,
				Z:         monster.Z,
				MonsterId: int32(monster.ID),
			},
		},
	}

	// 이 코드를 들어온 유저를 제외한 플레이어들에게 스폰시켜달라고 한다.
	for _, p := range GetPlayerManager().players {

		response := GetNetManager().MakePacket(MonsterSapwn)
		(*p.Conn).Write(response)
	}

	return monster
}

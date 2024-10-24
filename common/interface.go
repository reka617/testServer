// common/types.go
package common

import pb "testServer/Messages"

type Point struct {
	X, Z float32
}

type IMonster interface {
	GetPosition() Point
	SetPosition(x, y float32)
	GetTarget() *Point
	SetTarget(point *Point)
	GetPath() []Point
	GetPathIndex() int
	SetPathIndex(idx int)
	GetID() int
	GetHealth() int
	SetHealth(health int)
	GetMaxHealth() int
	IsDead() bool
}

type INetworkManager interface {
	MakePacket(sg *pb.GameMessage) []byte
}

type IPlayerManager interface {
	ListPoints() []*Point
	Broadcast(sg *pb.GameMessage)
}

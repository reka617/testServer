// manager/monster.go
package manager

import (
	"testServer/behavior"
	"testServer/common"
	"time"
)

type Monster struct {
	ID        int
	X, Z      float32
	Health    int
	MaxHealth int
	Target    *common.Point
	Path      []common.Point
	PathIndex int
	AI        behavior.Node
}

func NewMonster(id int, x, y float32, maxHealth int, path []common.Point) *Monster {
	m := &Monster{
		ID:        id,
		X:         x,
		Z:         y,
		Health:    maxHealth,
		MaxHealth: maxHealth,
		Path:      path,
		PathIndex: 0,
	}
	m.AI = CreateMonsterBehaviorTree(m)
	return m
}

// IMonster 인터페이스 구현
func (m *Monster) GetPosition() common.Point {
	return common.Point{X: m.X, Z: m.Z}
}

func (m *Monster) SetPosition(x, y float32) {
	m.X = x
	m.Z = y
}

func (m *Monster) SetTarget(target *common.Point) {
	m.Target = target
}

func (m *Monster) GetTarget() *common.Point {
	return m.Target
}

func (m *Monster) GetPath() []common.Point {
	return m.Path
}

func (m *Monster) GetPathIndex() int {
	return m.PathIndex
}

func (m *Monster) SetPathIndex(idx int) {
	m.PathIndex = idx
}

func (m *Monster) GetID() int {
	return m.ID
}

func (m *Monster) GetHealth() int {
	return m.Health
}

func (m *Monster) SetHealth(health int) {
	m.Health = health
	if m.Health < 0 {
		m.Health = 0
	}
	if m.Health > m.MaxHealth {
		m.Health = m.MaxHealth
	}
}

func (m *Monster) GetMaxHealth() int {
	return m.MaxHealth
}

func (m *Monster) IsDead() bool {
	return m.Health <= 0
}

func (m *Monster) Update() {
	if !m.IsDead() {
		m.AI.Execute()
	}
}

// CreateMonsterBehaviorTree creates the AI behavior tree for the monster
func CreateMonsterBehaviorTree(monster *Monster) behavior.Node {
	return behavior.NewSelector(
		// Combat sequence
		behavior.NewSequence(
			behavior.NewDetectPlayer(monster, 10.0, GetPlayerManager(), GetNetManager()), // Detection range
			behavior.NewSelector(
				// Attack sequence
				behavior.NewSequence(
					behavior.NewDetectPlayer(monster, 2.0, GetPlayerManager(), GetNetManager()), // Attack range
					behavior.NewAttack(monster, 2.0, 10, time.Second),                           // Attack damage and cooldown
				),
				// Chase sequence
				behavior.NewChase(monster, 3.0, GetPlayerManager(), GetNetManager()), // Chase speed
			),
		),
		// Patrol behavior
		behavior.NewPatrol(monster, 2.0), // Patrol speed
	)
}

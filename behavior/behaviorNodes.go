package behavior

import (
	"math"
	"time"
	manager "testServer/Manager"
)

// 행동 트리의 상태를 나타내는 상수
type Status int

const (
	Success Status = iota
	Failure
	Running
)

// Node 인터페이스 - 모든 노드가 구현해야 함
type Node interface {
	Execute() Status
}

// 몬스터 정보를 담는 구조체
type Monster struct {
	X, Z float32
	HP   int
	Target  *manager.Player
	Path    []Point
	PathIdx int
	MonsterId int32
}

// 위치 정보를 담는 구조체
type Point struct {
	X, Y float32
}

// Sequence 노드 - 자식 노드들을 순차적으로 실행
type Sequence struct {
	children []Node
}

func NewSequence(children ...Node) *Sequence {
	return &Sequence{children: children}
}

func (s *Sequence) Execute() Status {
	for _, child := range s.children {
		switch child.Execute() {
		case Failure:
			return Failure
		case Running:
			return Running
		}
	}
	return Success
}

// Selector 노드 - 자식 노드들 중 하나라도 성공할 때까지 실행
type Selector struct {
	children []Node
}

func NewSelector(children ...Node) *Selector {
	return &Selector{children: children}
}

func (s *Selector) Execute() Status {
	for _, child := range s.children {
		switch child.Execute() {
		case Success:
			return Success
		case Running:
			return Running
		}
	}
	return Failure
}

// 순찰 행동을 담당하는 노드
type Patrol struct {
	monster *Monster
}

func NewPatrol(monster *Monster) *Patrol {
	return &Patrol{monster: monster}
}

func (p *Patrol) Execute() Status {
	// 현재 목표 지점까지의 거리 계산
	currentPoint := p.monster.Path[p.monster.PathIdx]
	dist := distance(p.monster.X, p.monster.Z, currentPoint.X, currentPoint.Y)

	// 목표 지점에 도달했으면 다음 지점으로
	if dist < 1.0 {
		p.monster.PathIdx = (p.monster.PathIdx + 1) % len(p.monster.Path)
		return Success
	}

	// 목표 지점을 향해 이동
	speed := float32(2.0)
	dx := currentPoint.X - p.monster.X
	dy := currentPoint.Y - p.monster.Z
	norm := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	p.monster.X += (dx / norm) * speed
	p.monster.Z += (dy / norm) * speed

	return Running
}

// 플레이어 감지를 담당하는 노드
type DetectPlayer struct {
	monster *Monster
	_range   float32
}

func NewDetectPlayer(monster *Monster, detectRange float32) *DetectPlayer {
	return &DetectPlayer{monster: monster, _range: detectRange}
}

func (d *DetectPlayer) Execute() Status {
	if d.monster.Target == nil {
		return Failure
	}

	dist := distance(d.monster.X, d.monster.Z, d.monster.Target.X, d.monster.Target.Y)
	if dist <= d.range {
		return Success
	}
	return Failure
}

// 공격 행동을 담당하는 노드
type Attack struct {
	monster     *Monster
	attackRange float32
	damage      int
	lastAttack  time.Time
	cooldown    time.Duration
}

func NewAttack(monster *Monster, attackRange float32, damage int, cooldown time.Duration) *Attack {
	return &Attack{
		monster:     monster,
		attackRange: attackRange,
		damage:      damage,
		cooldown:    cooldown,
	}
}

func (a *Attack) Execute() Status {
	if a.monster.Target == nil {
		return Failure
	}

	// 공격 범위 확인
	dist := distance(a.monster.X, a.monster.Z, a.monster.Target.X, a.monster.Target.Y)
	if dist > a.attackRange {
		return Failure
	}

	// 쿨다운 확인
	now := time.Now()
	if now.Sub(a.lastAttack) < a.cooldown {
		return Running
	}

	// 공격 실행
	//a.monster.Target.HP -= a.damage
	a.lastAttack = now
	return Success
}

// 추적 행동을 담당하는 노드
type Chase struct {
	monster *Monster
	speed   float32
}

func NewChase(monster *Monster, speed float32) *Chase {
	return &Chase{monster: monster, speed: speed}
}

func (c *Chase) Execute() Status {
	if c.monster.Target == nil {
		return Failure
	}

	// 목표를 향해 이동
	dx := c.monster.Target.X - c.monster.X
	dy := c.monster.Target.Y - c.monster.Z
	norm := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	// 이미 충분히 가까우면 성공
	if norm < 1.0 {
		return Success
	}

	c.monster.X += (dx / norm) * c.speed
	c.monster.Z += (dy / norm) * c.speed
	return Running
}

// 거리 계산 유틸리티 함수
func distance(x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

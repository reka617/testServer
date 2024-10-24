package behavior

import (
	"math"
	pb "testServer/Messages"
	"testServer/common"
	"time"
)

type Status int

const (
	Success Status = iota
	Failure
	Running
)

type Node interface {
	Execute() Status
}

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

type DetectPlayer struct {
	monster common.IMonster
	range_  float32
	p       common.IPlayerManager
	n       common.INetworkManager
}

func (d *DetectPlayer) FindTarget() *common.Point {
	if len(d.p.ListPoints()) > 0 {
		return d.p.ListPoints()[0]
	}
	return nil
}

func NewDetectPlayer(monster common.IMonster, range_ float32, p common.IPlayerManager, n common.INetworkManager) *DetectPlayer {
	return &DetectPlayer{
		monster: monster,
		range_:  range_,
		p:       p,
		n:       n,
	}
}

func (d *DetectPlayer) Execute() Status {
	target := d.monster.GetTarget()
	if target == nil {
		target = d.FindTarget()
		if target == nil {
			return Failure
		}

		d.monster.SetTarget(target)
	}

	pos := d.monster.GetPosition()
	dist := distance(pos.X, pos.Z, target.X, target.Z)

	if dist <= d.range_ {
		return Success
	}
	return Failure
}

type Attack struct {
	monster     common.IMonster
	attackRange float32
	damage      int
	lastAttack  time.Time
	cooldown    time.Duration
}

func NewAttack(monster common.IMonster, attackRange float32, damage int, cooldown time.Duration) *Attack {
	return &Attack{
		monster:     monster,
		attackRange: attackRange,
		damage:      damage,
		cooldown:    cooldown,
	}
}

func (a *Attack) Execute() Status {
	if a.monster.IsDead() {
		return Failure
	}

	target := a.monster.GetTarget()
	if target == nil {
		return Failure
	}

	pos := a.monster.GetPosition()
	dist := distance(pos.X, pos.Z, target.X, target.Z)
	if dist > a.attackRange {
		return Failure
	}

	now := time.Now()
	if now.Sub(a.lastAttack) < a.cooldown {
		return Running
	}

	// 실제 공격 로직 구현
	// 여기서는 예시로 데미지만 처리
	currentHealth := a.monster.GetHealth()
	a.monster.SetHealth(currentHealth - a.damage)
	a.lastAttack = now

	return Success
}

type Chase struct {
	monster common.IMonster
	speed   float32

	p common.IPlayerManager
	n common.INetworkManager
}

func NewChase(monster common.IMonster, speed float32, p common.IPlayerManager, n common.INetworkManager) *Chase {
	return &Chase{monster: monster, speed: speed, p: p, n: n}
}

func (c *Chase) Execute() Status {
	target := c.monster.GetTarget()
	if target == nil {
		return Failure
	}

	println(target.X, target.Z)

	pos := c.monster.GetPosition()
	dx := target.X - pos.X
	dy := target.Z - pos.Z
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if dist < 0.1 {
		return Success
	}

	// 정규화된 방향으로 이동
	newX := pos.X + (dx/dist)*0.01
	newZ := pos.Z + (dy/dist)*0.01
	c.monster.SetPosition(newX, newZ)

	// 내가 로그인 되었음을 나한테 알려준다.
	MonsterMove := &pb.GameMessage{
		Message: &pb.GameMessage_MoveMonster{
			MoveMonster: &pb.MoveMonster{
				X:         newX,
				Z:         newZ,
				MonsterId: int32(c.monster.GetID()),
			},
		},
	}

	// 이 코드를 들어온 유저를 제외한 플레이어들에게 스폰시켜달라고 한다.
	c.p.Broadcast(MonsterMove)

	return Running
}

type Patrol struct {
	monster common.IMonster
	speed   float32
}

func NewPatrol(monster common.IMonster, speed float32) *Patrol {
	return &Patrol{
		monster: monster,
		speed:   speed,
	}
}

func (p *Patrol) Execute() Status {
	pos := p.monster.GetPosition()
	path := p.monster.GetPath()
	if len(path) == 0 {
		return Failure
	}

	currentIdx := p.monster.GetPathIndex()
	target := path[currentIdx]

	dx := target.X - pos.X
	dy := target.Z - pos.Z
	dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))

	if dist < 0.1 {
		// 다음 경로 포인트로 이동
		nextIdx := (currentIdx + 1) % len(path)
		p.monster.SetPathIndex(nextIdx)
		return Success
	}

	// 현재 목표 지점을 향해 이동
	newX := pos.X + (dx/dist)*p.speed
	newY := pos.Z + (dy/dist)*p.speed
	p.monster.SetPosition(newX, newY)

	return Running
}

func distance(x1, y1, x2, y2 float32) float32 {
	dx := x2 - x1
	dy := y2 - y1
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

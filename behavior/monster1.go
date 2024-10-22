package behavior

import "time"

// 몬스터의 행동 트리 생성 함수
func CreateMonsterBehaviorTree(monster *Monster) Node {
	return NewSelector(
		// 전투 시퀀스
		NewSequence(
			NewDetectPlayer(monster, 10.0), // 감지 범위 10
			NewSelector(
				// 공격 시퀀스
				NewSequence(
					NewDetectPlayer(monster, 2.0),            // 공격 범위 2
					NewAttack(monster, 2.0, 10, time.Second), // 데미지 10, 쿨다운 1초
				),
				// 추적 시퀀스
				NewChase(monster, 3.0), // 이동 속도 3
			),
		),
		// 순찰 행동
		NewPatrol(monster),
	)
}

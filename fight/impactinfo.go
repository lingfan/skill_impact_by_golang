package fight

type ImpactInfo struct {
	SkillID     int
	ImpactID    int
	TargetList  [MaxMatrixCellCount]int
	TargetCount int
	ConAttTimes int
	Hurts       [MaxMatrixCellCount][MaxConAttackTimes]int //本次技能带给的伤害变化 >0 减血 < 0 加血
	Mp          [MaxMatrixCellCount]int                    //本次技能带给的魔法量变化 >0 减蓝 < 0 加蓝
}

func (ii *ImpactInfo) GetTargetIndex(uid int) int {
	for i := 0; i < ii.TargetCount; i++ {
		if ii.TargetList[i] == uid {
			return i
		}
	}
	return InvalidId
}

func (ii *ImpactInfo) CleanUp() {
	ii.SkillID = InvalidId
	ii.ImpactID = InvalidId
}

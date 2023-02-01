package fight

type SkillAttack struct {
	SkillID     int                              //技能ID
	SkillTarget int                              //技能目标
	CostMp      int                              //消耗魔法量
	Impact      [MaxSkillImpactCount]*ImpactInfo //impact列表
	ImpactCount int                              //impact个数
}

func (sa *SkillAttack) CleanUp() {
	sa.SkillID = InvalidId
	for i := 0; i < MaxSkillImpactCount; i++ {
		sa.Impact[i].CleanUp()
	}
}

func (sa *SkillAttack) GetImpactInfo(impactId int) *ImpactInfo {
	for i := 0; i < MaxSkillImpactCount; i++ {
		if sa.Impact[i] != nil && sa.Impact[i].ImpactID == impactId {
			return sa.Impact[i]
		}
	}
	return &ImpactInfo{}

}

func (sa *SkillAttack) AddImpactInfo(info *ImpactInfo) {
	sa.Impact[sa.ImpactCount] = info
	sa.ImpactCount++
}
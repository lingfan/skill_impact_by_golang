package fight

type AttackInfo struct {
	CastUid int //对象guid

	BSkilled bool //魔法攻击还是普通攻击
	//魔法攻击
	SkillAttack      [6]*SkillAttack //如果是普通攻击先进行[装备]技能攻击
	SkillAttackCount int

	//普通攻击
	SkillTarget    int  //普通攻击目标
	BHit           bool //是否命中
	BStrike        bool //是否暴击
	Hurt           int  //伤害值
	BBackAttack    bool //是否有反击
	BackAttackHurt int
}

func (ai *AttackInfo) CleanUp() {
	ai.CastUid = 0
	ai.SkillTarget = 0
	ai.SkillAttackCount = 0

}

func (ai *AttackInfo) IsValid() bool {
	if ai.CastUid > 0 {
		return true
	}
	return false
}

func (ai *AttackInfo) NewSkillAttack() *SkillAttack {
	//fmt.Printf("AttackInfo AllocSkillAttack ai.SkillAttack %#v\n",ai.SkillAttack)
	index := ai.SkillAttackCount
	ai.SkillAttackCount++

	ai.SkillAttack[index] = &SkillAttack{}
	return ai.SkillAttack[index]

}

func (ai *AttackInfo) GetSkillAttack(nSkillID int) *SkillAttack {
	for i := 0; i < ai.SkillAttackCount; i++ {
		if ai.SkillAttack[i].SkillID == nSkillID {
			return ai.SkillAttack[i]
		}
	}
	return &SkillAttack{}
}



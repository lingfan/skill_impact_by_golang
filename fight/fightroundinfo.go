package fight

type FightRoundInfo struct {
	AttackObjectInfo  [MaxMatrixCellCount]*FObjectInfo //本回合开始前状态数据
	AttackObjectCount int
	DefendObjectInfo  [MaxMatrixCellCount]*FObjectInfo
	DefendObjectCount int
	AttackInfo        [MaxMatrixCellCount * 2]*AttackInfo //出手数据
	AttInfoCount      int
}

func (fri *FightRoundInfo) CleanUp() {
	for i := 0; i < MaxMatrixCellCount; i++ {
		fri.AttackObjectInfo[i].CleanUp()
		fri.DefendObjectInfo[i].CleanUp()
		fri.AttackInfo[i].CleanUp()
		fri.AttackInfo[i+MaxMatrixCellCount].CleanUp()
	}
	fri.AttackObjectCount = 0
	fri.DefendObjectCount = 0
	fri.AttInfoCount = 0

}

func (fri *FightRoundInfo) AddAttackInfo(ai *AttackInfo) {

	fri.AttInfoCount += 1
	fri.AttackInfo[fri.AttInfoCount] = ai
}

func (fri *FightRoundInfo) GetFObjectInfoByGuid(uid int) *FObjectInfo {
	for i := 0; i < fri.AttackObjectCount; i++ {
		if fri.AttackObjectInfo[i].Uid == uid {
			return fri.AttackObjectInfo[i]
		}
	}

	for i := 0; i < fri.DefendObjectCount; i++ {
		if fri.DefendObjectInfo[i].Uid == uid {
			return fri.DefendObjectInfo[i]
		}
	}
	return &FObjectInfo{}
}

func (fri *FightRoundInfo) AddAttackObjectInfo(objInfo *FObjectInfo) {

	//fmt.Printf("FightRoundInfo AddAttackObjectInfo fri.AttackObjectCount %#v, %#v\n",fri.AttackObjectCount, fri.AttackObjectInfo)

	fri.AttackObjectInfo[fri.AttackObjectCount] = objInfo
	fri.AttackObjectCount++
}
func (fri *FightRoundInfo) AddDefendObjectInfo(objInfo *FObjectInfo) {
	//fmt.Printf("FightRoundInfo AddDefendObjectInfo fri.DefendObjectCount %#v, %#v\n",fri.DefendObjectCount, fri.DefendObjectInfo)

	fri.DefendObjectInfo[fri.DefendObjectCount] = objInfo
	fri.DefendObjectCount++
}

type FObjectInfo struct {
	Uid           int
	HP            int //血量
	MaxHP         int //最大血量
	MP            int //魔法值
	MaxMP         int //最大魔法
	FightDistance int //战斗条长度
	AttackSpeed   int //速度
	EndDistance   int //最后位置
	ImpactCount   int
	ImpactList    [MaxImpactNumber]int //身上impact
	ImpactHurt    [MaxImpactNumber]int //持续impact伤害 + 掉血 - 加血
	ImpactMP      [MaxImpactNumber]int //持续impact 蓝  + 掉蓝 - 加蓝
}

func (fri *FObjectInfo) CleanUp() {

}

func (fri *FObjectInfo) AddImpact(impactId int, hurt int, mp int) {

	fri.ImpactList[fri.ImpactCount] = impactId
	fri.ImpactHurt[fri.ImpactCount] = hurt
	fri.ImpactMP[fri.ImpactCount] = mp
	fri.ImpactCount++
}

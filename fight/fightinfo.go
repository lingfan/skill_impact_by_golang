package fight

import "fmt"

type FightInfo struct {
	AttackObjectData  [MaxMatrixCellCount]*FObjectData
	AttackObjectCount int
	DefendObjectData  [MaxMatrixCellCount]*FObjectData
	DefendObjectCount int
	DefendType        int
	RoundInfo         [MaxFightRound]*FightRoundInfo
	Rounds            int  //总回合
	MaxFightDistance  int  //战斗条长度
	BWin              bool //挑战者是否胜利
}

func (fi *FightInfo) CleanUp() {
	for i := 0; i < MaxMatrixCellCount; i++ {
		fi.AttackObjectData[i].CleanUp()
		fi.DefendObjectData[i].CleanUp()
	}

	for i := 0; i < MaxFightRound; i++ {
		fi.RoundInfo[i].CleanUp()
	}

	fi.Rounds = 0
	fi.AttackObjectCount = 0
	fi.DefendObjectCount = 0
	fi.MaxFightDistance = 0
	fi.BWin = false
	fi.DefendType = 0

}

func (fi *FightInfo) SetWin(b bool) {
	fmt.Printf("SetWin %#v\n", b)

	fi.BWin = b
}

func (fi *FightInfo) AddRoundInfo(roundInfo *FightRoundInfo) {
	//fmt.Printf("FightObjList AddRoundInfo %#v\n",roundInfo)

	fi.RoundInfo[fi.Rounds] = roundInfo
	fi.Rounds++
}

func (fi *FightInfo) AddAttackObjectData(objData *FObjectData) {
	fi.AttackObjectData[fi.AttackObjectCount] = objData
	fi.AttackObjectCount++
}

func (fi *FightInfo) AddDefendObjectData(objData *FObjectData) {
	fi.DefendObjectData[fi.DefendObjectCount] = objData
	fi.DefendObjectCount++
}

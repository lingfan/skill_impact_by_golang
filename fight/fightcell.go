package fight

import "fmt"

//FightCell --
type FightCell struct {
	AttackerList *FightObjList //攻击方
	DefenderList *FightObjList //防御方
	FightInfo    *FightInfo    //战斗信息提取
	RoundInfo    *FightRoundInfo
	FightType    int
	PlusAtk      int //攻加成
	Log []Log
	FLog []FLog
}

func NewFightCell() *FightCell {
	fc := &FightCell{}
	fc.AttackerList = &FightObjList{}
	fc.DefenderList = &FightObjList{}
	fc.FightInfo = &FightInfo{}
	fc.RoundInfo = &FightRoundInfo{}

	return fc
}

func (fc *FightCell) CleanUp() {
	fc.AttackerList.CleanUp()
	fc.DefenderList.CleanUp()
	fc.FightInfo.CleanUp()
	fc.RoundInfo.CleanUp()
	fc.FightType = EM_TYPE_FIGHT_NORMAL
	fc.PlusAtk = 0
}

func (fc *FightCell) InitAttackerList() {
	//fmt.Printf("InitAttackerList AttackerList %#v\n", fc.AttackerList)

	//最大阵型位置id
	for i := 0; i < MaxMatrixCellCount; i++ {
		fo := fc.AttackerList.GetFightObject(i)
		fmt.Printf("InitAttackerList i %#v, fo %#v\n", i, fo.FightDBData)
		if fo.FightDBData.GUid > 0 {
			fo.BAttacker = true
			fo.FightCell = fc
			fo.InitSkill()
		}
	}
}

func (fc *FightCell) InitDefenderList(defendType int) {
	//最大阵型位置id
	for i := 0; i < MaxMatrixCellCount; i++ {
		fo := fc.DefenderList.GetFightObject(i)
		if fo.FightDBData.GUid > 0 {
			fo.BAttacker = false
			fo.FightCell = fc
			fo.InitSkill()
		}
	}
	fc.FightInfo.DefendType = defendType
}

func (fc *FightCell) IsWin() bool {
	return fc.FightInfo.BWin
}

func (fc *FightCell) IsOver() bool {
	if fc.AttackerList.GetActiveCount() <= 0 {
		return true
	}
	if fc.DefenderList.GetActiveCount() <= 0 {
		return true
	}
	return false
}

func (fc *FightCell) Fight() bool {
	fc.initFightInfo()
	for nRound := 1; nRound <= MaxFightRound; nRound++ {
		//fmt.Printf("Fight nRound %#v\n", nRound)

		if fc.IsOver() {

			fmt.Printf("Fight IsOver AttackerList %#v\n", fc.AttackerList.GetActiveCount())
			fmt.Printf("Fight IsOver DefenderList %#v\n", fc.DefenderList.GetActiveCount())

			if fc.AttackerList.GetActiveCount() > 0 {
				fc.FightInfo.SetWin(true)
			} else {
				fc.FightInfo.SetWin(false)
			}

			return true
		}
		fc.initRoundInfo()

		fc.AttackerList.ImpactHeartBeat(nRound)
		fc.DefenderList.ImpactHeartBeat(nRound)

		fc.AttackerList.HeartBeat(nRound)
		fc.DefenderList.HeartBeat(nRound)

		fc.FightInfo.AddRoundInfo(fc.RoundInfo)
	}

	fc.FightInfo.SetWin(false)
	return true
}

func (fc *FightCell) initFightInfo() {
	fc.FightInfo.MaxFightDistance = FightDistance

	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := fc.AttackerList.GetFightObject(i)
		if pFightObj.IsValid() {
			objData := &FObjectData{}
			objData.Uid = pFightObj.GetGUid()
			objData.MatrixID = pFightObj.GetMatrixID()
			objData.Level = pFightObj.GetLevel()

			fc.FightInfo.AddAttackObjectData(objData)
		}

		pFightObj = fc.AttackerList.GetFightObject(i)
		if pFightObj.IsValid() {
			objData := &FObjectData{}
			objData.Uid = pFightObj.GetGUid()
			objData.MatrixID = pFightObj.GetMatrixID()
			objData.Level = pFightObj.GetLevel()

			fc.FightInfo.AddDefendObjectData(objData)
		}
	}

}

func (fc *FightCell) initRoundInfo() {
	fc.RoundInfo = &FightRoundInfo{}

	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := fc.AttackerList.GetFightObject(i)
		if pFightObj.IsValid() {
			objInfo := &FObjectInfo{}
			objInfo.Uid = pFightObj.GetGUid()

			//fmt.Printf("AttackerList %v 开始战斗计算\n",objInfo.Uid)

			objInfo.HP = pFightObj.GetHP()
			objInfo.MaxHP = pFightObj.GetMaxHP()
			objInfo.MP = pFightObj.GetMP()
			objInfo.MaxMP = pFightObj.GetMaxMP()
			objInfo.FightDistance = pFightObj.FightDistance
			objInfo.AttackSpeed = pFightObj.GetAttackSpeed()

			pImpactList := pFightObj.GetImpactList()
			for idx := 0; idx < MaxImpactNumber; idx++ {
				if pImpactList[idx] != nil && pImpactList[idx].IsValid() {

					//fmt.Printf("%v 开始技能计算 pImpactList[idx].ImpactID %v\n",objInfo.Uid,pImpactList[idx].ImpactID)

					objInfo.AddImpact(pImpactList[idx].ImpactID, 0, 0)
				}
			}
			fc.RoundInfo.AddAttackObjectInfo(objInfo)
		}

		pFightObj = fc.DefenderList.GetFightObject(i)
		if pFightObj.IsValid() {
			objInfo := &FObjectInfo{}
			objInfo.Uid = pFightObj.GetGUid()
			//fmt.Printf("DefenderList %v 开始战斗计算\n",objInfo.Uid)

			objInfo.HP = pFightObj.GetHP()
			objInfo.MaxHP = pFightObj.GetMaxHP()
			objInfo.MP = pFightObj.GetMP()
			objInfo.MaxMP = pFightObj.GetMaxMP()
			objInfo.FightDistance = pFightObj.FightDistance
			objInfo.AttackSpeed = pFightObj.GetAttackSpeed()

			pImpactList := pFightObj.GetImpactList()
			for idx := 0; idx < MaxImpactNumber; idx++ {
				if pImpactList[idx] != nil && pImpactList[idx].IsValid() {
					//fmt.Printf("%v 开始技能计算 pImpactList[idx].ImpactID %v\n",objInfo.Uid,pImpactList[idx].ImpactID)

					objInfo.AddImpact(pImpactList[idx].ImpactID, 0, 0)
				}
			}
			fc.RoundInfo.AddDefendObjectInfo(objInfo)
		}

	}

}

func (fc *FightCell) GetDefenceList() *FightObjList {
	return fc.DefenderList
}

func (fc *FightCell) GetAttackList() *FightObjList {
	return fc.AttackerList
}

func (fc *FightCell) GetRoundInfo() *FightRoundInfo {
	return fc.RoundInfo
}
func (fc *FightCell) GetFightInfo() *FightInfo {
	return fc.FightInfo
}

func (fc *FightCell) GetFightType() int {
	return fc.FightType
}
func (fc *FightCell) GetPlusAtt() int {
	return fc.PlusAtk
}


func (fc *FightCell) GetLog()[]Log{
	return fc.Log
}

func (fc *FightCell) AddLog(log Log) {
	fc.Log = append(fc.Log, log)
}


func (fc *FightCell) AddFLog(log FLog) {
	fc.FLog = append(fc.FLog, log)
}

func (fc *FightCell) GetFLog()[]FLog{
	return fc.FLog
}


func (fc *FightCell) GetLogList() []LFLog {
	var ret []LFLog
	var xx []string
	vv  := map[string][]FLog{}

	for i, v := range fc.GetFLog() {
		v.No = i
		newKey := fmt.Sprintf("%d_%d_%d", v.Round, v.TeamNo, v.Attacker)
		if len(vv[newKey]) == 0 {
			xx = append(xx, newKey)
		}
		vv[newKey] = append(vv[newKey], v)
	}

	for _,k := range xx {
		lfl := LFLog{}
		if len(vv[k]) == 0 {
			continue
		}
		lfl.Round = vv[k][0].Round
		lfl.TeamNo = vv[k][0].TeamNo
		lfl.Attacker = vv[k][0].Attacker
		lfl.List = vv[k]
		lfl.AttackerName = vv[k][0].AttackerName
		lfl.SkillName = vv[k][0].SkillName
		lfl.SkillID = vv[k][0].SkillID

		ret = append(ret,lfl)
	}
	return ret

}


type LFLog struct {
	Round    int
	TeamNo   int
	Attacker int
	AttackerName string
	SkillName string
	SkillID int
	List     []FLog
}
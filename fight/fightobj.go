package fight

import (
	"cyhd/library/tools"
	"fmt"
)

//FightObj --
type FightObj struct {
	FightDBData   *FightDBData
	FightCell     *FightCell
	BAttacker     bool
	SkillList     [MaxSkillNum]*Skill
	CommonKill    *Skill
	ImpactList    [MaxImpactNumber]*Impact
	ImpactEffect  [EmAttributeNumber]int
	FightDistance int
	AttackInfo    *AttackInfo
}

//InitFightDBData --
func (fo *FightObj) InitFightDBData(fightDBData *FightDBData) {

	fmt.Printf("InitFightDBData fightDBData %#v\n", fightDBData)

	fo.FightDBData = fightDBData
	fo.AttackInfo = &AttackInfo{}

}

func (fo *FightObj) GetFightCell() *FightCell {
	return fo.FightCell
}

func (fo *FightObj) InitSkill() {

	//fmt.Printf("FightObj InitSkill\n")

	for i, id := range fo.FightDBData.Skill {
		if id > 0 {
			fo.SkillList[i] = NewSkill(id, fo)
		}
	}

	fo.CommonKill = NewSkill(fo.FightDBData.Profession, fo)

}

func (fo *FightObj) GetGUid() int {
	return fo.FightDBData.GUid
}

func (fo *FightObj) GetTplId() int {
	return fo.FightDBData.TableID
}


func (fo *FightObj) IsAttacker() bool {
	return fo.BAttacker
}

func (fo *FightObj) CleanUp() {

	for i := 0; i < MaxSkillNum; i++ {
		fo.SkillList[i].CleanUp()
	}
	for i := 0; i < MaxImpactNumber; i++ {
		fo.ImpactList[i].CleanUp()
	}

	for i := 0; i < EmAttributeNumber; i++ {
		fo.ImpactEffect[i] = 0
	}

	fo.CommonKill.CleanUp()

	fo.FightDistance = 0

}

func (fo *FightObj) GetProfession() int {
	return fo.FightDBData.Profession
}

func (fo *FightObj) GetMaxHP() int {
	nEffectValue := float64(fo.ImpactEffect[EmAttributeMaxHP])

	nEffectValue += float64(fo.FightDBData.MaxHP*fo.ImpactEffect[EmAttributePercentMaxHP]) * 0.01

	nEndValue := fo.FightDBData.MaxHP + int(nEffectValue)

	return nEndValue
}

func (fo *FightObj) GetMaxMP() int {
	nEffectValue := float64(fo.ImpactEffect[EmAttributeMaxMP])

	nEffectValue += float64(fo.FightDBData.MaxMP*fo.ImpactEffect[EmAttributeMaxMP]) * 0.01

	nEndValue := fo.FightDBData.MaxMP + int(nEffectValue)

	return nEndValue
}

func (fo *FightObj) HeartBeat(uTime int) bool {

	//fmt.Printf("FightObj HeartBeat %#v, uid %#v, TableID %#v\n",uTime, fo.GetGUid(), fo.FightDBData.TableID)

	fo.AttackInfo.CleanUp()

	pRoundInfo := fo.FightCell.GetRoundInfo()
	pFObjectInfo := pRoundInfo.GetFObjectInfoByGuid(fo.FightDBData.GUid)

	pFObjectInfo.AttackSpeed = fo.GetAttackSpeed()

	if !fo.IsActive() {
		return true
	}

	if uTime == 1 {
		fo.CastPassiveSkill(uTime)
	}

	pEnemyList := fo.GetEnemyList()

	//fmt.Printf("FightObj HeartBeat pEnemyList %#v\n",pEnemyList)

	if pEnemyList.GetActiveCount() == 0 {
		return true
	}

	fo.FightDistance += fo.GetAttackSpeed()

	//fmt.Printf("FightObj HeartBeat fo.FightDistance %#v, uid %#v, TableID %#v\n",fo.FightDistance, fo.GetGUid(), fo.FightDBData.TableID)

	if fo.FightDistance >= FightDistance {
		fo.FightDistance = 0
		bRet := fo.SkillHeartBeat(uTime)
		if !bRet {
			//fmt.Printf("FightObj HeartBeat fo.CastCommonSkill uTime %#v, guid %#v\n",uTime, fo.GetGUid())
			fo.CastCommonSkill(uTime)
		}
	}

	if pFObjectInfo != nil {
		pFObjectInfo.EndDistance = fo.FightDistance
		if pFObjectInfo.MaxHP < fo.GetMaxHP() {
			pFObjectInfo.MaxHP = fo.GetMaxHP()
		}
		if pFObjectInfo.MaxMP < fo.GetMaxMP() {
			pFObjectInfo.MaxMP = fo.GetMaxMP()
		}
	}

	return true
}

func (fo *FightObj) GetEnemyList() *FightObjList {
	if fo.BAttacker {
		return fo.FightCell.GetDefenceList()
	}

	return fo.FightCell.GetAttackList()
}

func (fo *FightObj) GetOwnerList() *FightObjList {
	if fo.BAttacker {
		return fo.FightCell.GetAttackList()
	}

	return fo.FightCell.GetDefenceList()
}

func (fo *FightObj) GetLevel() int {
	return fo.FightDBData.Level
}

func (fo *FightObj) GetMatrixID() int {
	return fo.FightDBData.MatrixID
}
func (fo *FightObj) SetMatrixID(index int) {
	fo.FightDBData.MatrixID = index
}

func (fo *FightObj) IsValid() bool {
	if fo.FightDBData.GUid == 0 {
		return false
	}
	return true
}

func (fo *FightObj) IsActive() bool {
	//fmt.Printf("FightObj IsActive %#v\n",fo.FightDBData)

	if fo.GetHP() > 0 {
		return true
	}
	return false
}
func (fo *FightObj) GetAttackSpeed() int {
	nEffectValue := float64(fo.ImpactEffect[EmAttributeAttackSpeed])
	nEffectValue += float64(fo.FightDBData.AttackSpeed * fo.ImpactEffect[EmAttributePercentAttackSpeed])*0.01
	nEndValue := fo.FightDBData.AttackSpeed + int(nEffectValue)
	return nEndValue
}

func (fo *FightObj) GetPhysicAttack() int {
	//fmt.Printf("%v GetPhysicAttack %v, %v, %v\n",fo.GetGUid(), fo.FightDBData.PhysicAttack, fo.ImpactEffect[EmAttributePhysicAttack],fo.ImpactEffect[EmAttributePercentPhysicAttack])
	nEffectValue := float64(fo.ImpactEffect[EmAttributePhysicAttack])
	nEffectValue += float64(fo.FightDBData.PhysicAttack * fo.ImpactEffect[EmAttributePercentPhysicAttack])*0.01
	nEndValue := fo.FightDBData.PhysicAttack + int(nEffectValue)
	return nEndValue
}

func (fo *FightObj) GetMagicAttack() int {
	fmt.Printf("%v GetMagicAttack %v, %v, %v\n",fo.GetGUid(), fo.FightDBData.MagicAttack, fo.ImpactEffect[EmAttributeMagicAttack],fo.ImpactEffect[EmAttributePercentMagicAttack])
	nEffectValue := float64(fo.ImpactEffect[EmAttributeMagicAttack])
	nEffectValue += float64(fo.FightDBData.MagicAttack * fo.ImpactEffect[EmAttributePercentMagicAttack])*0.01
	nEndValue := fo.FightDBData.MagicAttack + int(nEffectValue)
	return nEndValue
}

func (fo *FightObj) GetPhysicDefend() int {
	//fmt.Printf("%v GetPhysicDefend %v, %v, %v\n",fo.GetGUid(), fo.FightDBData.PhysicDefence, fo.ImpactEffect[EmAttributePhysicDefence],fo.ImpactEffect[EmAttributePercentPhysicDefence])
	nEffectValue := float64(fo.ImpactEffect[EmAttributePhysicDefence])
	nEffectValue += float64(fo.FightDBData.PhysicDefence * fo.ImpactEffect[EmAttributePercentPhysicDefence])*0.01
	nEndValue := fo.FightDBData.PhysicDefence + int(nEffectValue)
	return nEndValue
}

func (fo *FightObj) GetMagicDefend() int {
	nEffectValue := float64(fo.ImpactEffect[EmAttributeMagicDefence])
	nEffectValue += float64(fo.FightDBData.MagicDefence*fo.ImpactEffect[EmAttributePercentMagicDefence]) * 0.01
	nEndValue := fo.FightDBData.MagicDefence + int(nEffectValue)
	return nEndValue
}

//GetPhysicHurtDecay 物理减免
func (fo *FightObj) GetPhysicHurtDecay() int {
	nEffectValue := fo.ImpactEffect[EmAttributePhysicHurtDecay]
	nEndValue := fo.FightDBData.PhysicHurtDecay + nEffectValue
	return nEndValue
}

//GetMagicHurtDecay 魔法减免
func (fo *FightObj) GetMagicHurtDecay() int {
	nEffectValue := fo.ImpactEffect[EmAttributeMagicHurtDecay]
	nEndValue := fo.FightDBData.MagicHurtDecay + nEffectValue
	return nEndValue
}

//GetHit
func (fo *FightObj) GetHit() int {
	nEffectValue := fo.ImpactEffect[EmAttributeHit]
	nEndValue := fo.FightDBData.Hit + nEffectValue
	return nEndValue
}

//GetDodge
func (fo *FightObj) GetDodge() int {
	nEffectValue := fo.ImpactEffect[EmAttributeDodge]
	nEndValue := fo.FightDBData.Dodge + nEffectValue
	return nEndValue
}

//GetFloatingHurt 伤害浮动
func (fo *FightObj) GetFloatingHurt() int {
	nFloatingHurt := 0
	if fo.FightDBData.Type == EM_HERO_HUMAN {
		nFloatingHurt = fo.FightDBData.FloatingHurt
	}

	if nFloatingHurt > 0 {
		nFloatingHurt = tools.RandNumRange(1,nFloatingHurt)
	}
	return nFloatingHurt
}

func (fo *FightObj) ClearImpactEffect() {
	//fmt.Printf("FightObj # %v, ClearImpactEffect  \n",fo.GetGUid())

	for i := 0; i < EmAttributeNumber; i++ {
		fo.ImpactEffect[i] = 0
	}

}


func (fo *FightObj) ImpactHeartBeat(uTime int) {

	for i := 0; i < MaxImpactNumber; i++ {
		if fo.ImpactList[i] != nil && fo.ImpactList[i].IsValid() {
			//fmt.Printf("FightObj ImpactHeartBeat # %v, i %v, uTime %v, fo.ImpactList %#v\n",fo.GetGUid(),i, uTime,fo.ImpactList[i])
			fo.ImpactList[i].HeartBeat(uTime)
		}
	}

}

//GetContinuous 连击
func (fo *FightObj) GetContinuous() int {
	nEffectValue := fo.ImpactEffect[EmAttributeContinuous]
	nEndValue := fo.FightDBData.Continuous + nEffectValue
	return nEndValue
}

//GetConAttTimes 连击次数
func (fo *FightObj) GetConAttTimes() int {
	nEffectValue := fo.ImpactEffect[EmAttributeContinuousTimes]
	nEndValue := fo.FightDBData.ConAttTimes + nEffectValue
	return nEndValue
}

//GetConAttHurt 连击伤害
func (fo *FightObj) GetConAttHurt() int {
	nEffectValue := fo.ImpactEffect[EmAttributeHurtContinuous]
	nEndValue := fo.FightDBData.ConAttHurt + nEffectValue
	return nEndValue
}

//GetAttackBack 反击
func (fo *FightObj) GetAttackBack() int {
	nEffectValue := fo.ImpactEffect[EmAttributeBackAttack]
	nEndValue := fo.FightDBData.BackAttack + nEffectValue
	return nEndValue
}

//GetBackAttHurt 反击伤害
func (fo *FightObj) GetBackAttHurt() int {
	nEffectValue := fo.ImpactEffect[EmAttributeHurtBackAttack]
	nEndValue := fo.FightDBData.BackAttHurt + nEffectValue
	return nEndValue
}

//GetStrikeHurt 暴击伤害
func (fo *FightObj) GetStrikeHurt() int {
	nEffectValue := fo.ImpactEffect[EmAttributeHurtStrike]
	nEndValue := fo.FightDBData.StrikeHurt + nEffectValue
	return nEndValue
}

//GetStrike 暴击
func (fo *FightObj) GetStrike() int {
	nEffectValue := fo.ImpactEffect[EmAttributeStrike]
	nEndValue := fo.FightDBData.Strike + nEffectValue
	return nEndValue
}

func (fo *FightObj) SetHP(nHP int) {
	nMaxHP := fo.FightDBData.MaxHP
	if nHP > nMaxHP {
		nHP = nMaxHP
	}

	if nHP < 0 {
		nHP = 0
	}

	fo.FightDBData.HP = nHP

	if nHP == 0 {
		fo.ClearImpact()
	}

}

//GetHP 获得hp
func (fo *FightObj) GetHP() int {

	nEffectValue := fo.ImpactEffect[EmAttributeHp]
	nEndValue := fo.FightDBData.HP + nEffectValue

	if nEndValue > fo.FightDBData.MaxHP {
		nEndValue = fo.FightDBData.MaxHP
	}

	fo.FightDBData.HP = nEndValue

	fo.ImpactEffect[EmAttributeHp] = 0

	//fmt.Printf("FightObj GetHP %#v %#v %#v %#v\n",fo.GetGUid(), fo.FightDBData.HP , nEffectValue, nEndValue)


	return nEndValue
}

func (fo *FightObj) SetMP(nMP int) {
	nMaxMP := fo.FightDBData.MaxMP
	if nMP > nMaxMP {
		nMP = nMaxMP
	}

	if nMP < 0 {
		nMP = 0
	}

	fo.FightDBData.MP = nMP

}

//GetHP --
func (fo *FightObj) GetMP() int {

	nEffectValue := fo.ImpactEffect[EmAttributeMp]
	nEndValue := fo.FightDBData.MP + nEffectValue

	if nEndValue > fo.FightDBData.MaxMP {
		nEndValue = fo.FightDBData.MaxMP
	}

	fo.FightDBData.MP = nEndValue

	fo.ImpactEffect[EmAttributeMp] = 0

	return nEndValue
}

func (fo *FightObj) GetImpactLogicType(nImpactID int) int {
	pRowImpact := GetImpactRow(nImpactID)

	ret := EM_TYPE_IMPACTLOGIC_INVALID
	switch pRowImpact.LogicID {
	case EmImpactLogic0:
		fallthrough
	case EmImpactLogic1:
		fallthrough
	case EmImpactLogic6:
		ret = EM_TYPE_IMPACTLOGIC_SINGLE
		break
	case EmImpactLogic2:
		fallthrough
	case EmImpactLogic3:
		fallthrough
	case EmImpactLogic5:
		ret = EM_TYPE_IMPACTLOGIC_DEBUFF
		break
	case EmImpactLogic4:
		ret = EM_TYPE_IMPACTLOGIC_BUFF
		break
	}

	//fmt.Printf("FightObj AddImpact nImpactID %#v, pRowImpact.LogicID %#v, ret %#v\n",nImpactID, pRowImpact.LogicID, ret)
	return ret
}

//AddImpact --
func (fo *FightObj) AddImpact(nImpactID int, conAttTimes int, nRound int, pCaster *FightObj, nSkillID int) int {
	pRowImpactNew := GetImpactRow(nImpactID)
	logicType := fo.GetImpactLogicType(nImpactID)

	//fmt.Printf("FightObj AddImpact guid %v, nImpactID %#v, logicType %#v, ImpactMutexID %v\n", fo.GetGUid(), nImpactID, logicType,pRowImpactNew.ImpactMutexID)

	if logicType == EM_TYPE_IMPACTLOGIC_SINGLE {
		newImpact := NewImpact(nImpactID, conAttTimes, nRound, pCaster, fo, nSkillID)
		newImpact.HeartBeat(nRound)
		return EM_IMPACT_RESULT_NORMAL
	} else {

		if pRowImpactNew.ImpactMutexID >= 0 {
			for i := 0; i < MaxImpactNumber; i++ {
				if fo.ImpactList[i] != nil && fo.ImpactList[i].IsValid() {
					pRowImpact := GetImpactRow(fo.ImpactList[i].ImpactID)

					if pRowImpactNew.ImpactMutexID == pRowImpact.ImpactMutexID {
						fo.ImpactList[i] = NewImpact(nImpactID, conAttTimes, nRound, pCaster, fo, nSkillID)
						fo.ImpactList[i].HeartBeat(nRound)
						return EM_IMPACT_RESULT_NORMAL
					} else {
						fmt.Printf("FightObj AddImpact pRowImpactNew.ImpactMutexID %#v   \n", []interface{}{nImpactID, pRowImpactNew.ImpactMutexID, fo.ImpactList[i].ImpactID, pRowImpact.ImpactMutexID})

						return EM_IMPACT_RESULT_FAIL
					}
				}
			}
		}

		if logicType == EM_TYPE_IMPACTLOGIC_BUFF {
			for i := 0; i < MaxBuffNumber; i++ {
				if fo.ImpactList[i] == nil || !fo.ImpactList[i].IsValid() {
					fo.ImpactList[i] = NewImpact(nImpactID, conAttTimes, nRound, pCaster, fo, nSkillID)
					fo.ImpactList[i].HeartBeat(nRound)

					//fmt.Printf("FightObj AddImpact guid %v EM_TYPE_IMPACTLOGIC_BUFF %v\n", fo.GetGUid(), i)
					return EM_IMPACT_RESULT_NORMAL
				}
			}
		}

		if logicType == EM_TYPE_IMPACTLOGIC_DEBUFF {
			for i := MaxBuffNumber; i < MaxImpactNumber; i++ {

				if fo.ImpactList[i] == nil || !fo.ImpactList[i].IsValid() {
					fo.ImpactList[i] = NewImpact(nImpactID, conAttTimes, nRound, pCaster, fo, nSkillID)
					fo.ImpactList[i].HeartBeat(nRound)

					//fmt.Printf("FightObj AddImpact guid %v EM_TYPE_IMPACTLOGIC_DEBUFF\n", fo.GetGUid())

					return EM_IMPACT_RESULT_NORMAL
				}
			}
		}
	}
	return EM_IMPACT_RESULT_FAIL

}

func (fo *FightObj) GetImpactList() [MaxImpactNumber]*Impact {
	//fmt.Printf("FightObj GetImpactList %#v\n",fo.ImpactList)
	return fo.ImpactList
}

//ClearImpact --
func (fo *FightObj) ClearImpact() {
	pOwnList := fo.GetOwnerList()
	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := pOwnList.GetFightObject(i)
		if pFightObj.IsActive() {
			pImpactList := pFightObj.GetImpactList()

			for j := 0; j < MaxImpactNumber; j++ {
				if pImpactList[j] != nil && pImpactList[j].IsValid() && pImpactList[j].GetCaster() == fo {
					pImpactList[j].CleanUp()
				}
			}
		}
	}

	pOwnList = fo.GetEnemyList()
	for i := 0; i < MaxMatrixCellCount; i++ {
		pFightObj := pOwnList.GetFightObject(i)
		if pFightObj.IsActive() {
			pImpactList := pFightObj.GetImpactList()

			for j := 0; j < MaxImpactNumber; j++ {
				if pImpactList[j] != nil && pImpactList[j].IsValid() && pImpactList[j].GetCaster() == fo {
					pImpactList[j].CleanUp()
				}
			}
		}
	}
}

func (fo *FightObj) ChangeEffect(nAttrType int, nValue int) {
	fo.ImpactEffect[nAttrType] += nValue
}

//SkillHeartBeat -
func (fo *FightObj) SkillHeartBeat(uTime int) bool {
	for i := MaxSkillNum - 1; i >= 0; i-- {
		if fo.SkillList[i] != nil {
			bLogic := fo.SkillList[i].SkillLogic(uTime)
			if bLogic {
				return true
			}
		}

	}
	return false
}

//CastCommonSkill -
func (fo *FightObj) CastCommonSkill(uTime int) {
	fo.CommonKill.CommonSkillLogic(uTime)
}

//CastPassiveSkill --
func (fo *FightObj) CastPassiveSkill(uTime int) {
	for i := 0; i < MaxSkillNum; i++ {
		if fo.SkillList[i] == nil {
			continue
		}
		fo.SkillList[i].PassiveSkillLogic(uTime)
	}
	// EquipSkillList[i].PassiveSkillLogic(uTime);
}

func (fo *FightObj) GetAttackInfo() *AttackInfo {
	return fo.AttackInfo
}

func (fo *FightObj) TeamNo() int {
	n := 2
	if fo.BAttacker {
		n = 1
	}
	return n
}

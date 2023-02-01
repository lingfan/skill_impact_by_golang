package fight

import (
	"cyhd/app/cy/config"
	"cyhd/library/tools"
	"fmt"
)

//Skill 技能数据
type Skill struct {
	SkillID      int       // 技能ID
	SkillType    int       // 技能类型
	SkillTime    int       // 技能时间
	CoolDownTime int       // 冷却时间
	PCaster      *FightObj // 释放者
	PTarget      *FightObj // 技能的目标
}

//Init --
func NewSkill(nSkillID int, pCaster *FightObj) *Skill {
	s := &Skill{}

	s.SkillID = nSkillID
	s.PCaster = pCaster

	pSkillRow := s.GetSkillRow()

	s.SkillType = pSkillRow.SkillType
	s.SkillTime = pSkillRow.StartRound
	s.CoolDownTime = pSkillRow.CoolDownTime

	//fmt.Printf("Skill Init %#v\n",s)

	return s
}

//CleanUp --
func (s *Skill) CleanUp() {
	s.SkillID = InvalidId
	s.SkillType = InvalidId
}

//IsValid --
func (s *Skill) IsValid() bool {
	if s.SkillID > 0 {
		return true
	}
	return false
}

//GetSkillRow 技能静态数据
func (s *Skill) GetSkillRow() config.TableSkill {
	return config.GetSkillRow(s.SkillID)
}

//GetSkillCastRate 获得技能释放几率
func (s *Skill) GetSkillCastRate() int {

	return 0
}

//SetSkillTarget --
func (s *Skill) SetSkillTarget(pTarget *FightObj) {

	s.PTarget = pTarget
}

//SkillLogic 魔法技能-主动
func (s *Skill) SkillLogic(nRound int) bool {
	if !s.IsValid() {
		return false
	}

	//fmt.Printf("SkillLogic IsValid %#v,%#v,%#v\n",nRound,s.SkillID,s.PCaster.GetGUid())

	if !s.CheckCondition(nRound) {
		return false
	}

	//fmt.Printf("SkillLogic CheckCondition %#v,%#v,%#v\n",nRound,s.SkillID,s.PCaster.GetGUid())

	if !s.SelectTarget() {
		return false
	}

	//fmt.Printf("SkillLogic SelectTarget %#v,%#v,%#v\n",nRound,s.SkillID,s.PCaster.GetGUid())

	if !s.CastSkill(nRound) {
		return false
	}
	//fmt.Printf("SkillLogic CastSkill %#v,%#v,%#v\n",nRound,s.SkillID,s.PCaster.GetGUid())

	return true
}

//PassiveSkillLogic 魔法技能-被动
func (s *Skill) PassiveSkillLogic(nRound int) bool {

	if !s.IsValid() {
		return false
	}

	if nRound != 1 {
		return false
	}

	if s.SkillType != EM_SKILL_TYPE_HERO_PASSIVE {
		return false
	}
	fmt.Printf("Skill PassiveSkillLogic uid %v, nRound %#v, SkillID %#v, SkillType %#v\n",s.PCaster.GetGUid(), nRound, s.SkillID,s.SkillType)

	if !s.SelectTarget() {
		return false
	}
	//fmt.Printf("Skill PassiveSkillLogic1 uid %v, nRound %#v, SkillID %#v, SkillType %#v\n",s.PCaster.GetGUid(), nRound, s.SkillID,s.SkillType)

	if !s.CastSkill(nRound) {
		return false
	}
	//fmt.Printf("Skill PassiveSkillLogic1 uid %v, nRound %#v, SkillID %#v, SkillType %#v\n",s.PCaster.GetGUid(), nRound, s.SkillID,s.SkillType)
	return true

}

//CheckCondition --
func (s *Skill) CheckCondition(nRound int) bool {


	if s.SkillType == EM_SKILL_TYPE_HERO_PASSIVE {
		return false
	}
	//fmt.Printf("Skill CheckCondition uid %v, nRound %#v, SkillID %#v, SkillType %#v\n",s.PCaster.GetGUid(), nRound, s.SkillID,s.SkillType)
	if nRound < s.SkillTime {
		return false
	}

	pRow := s.GetSkillRow()

	//fmt.Printf("Skill CheckCondition GetMP uid %v %#v %#v \n",s.PCaster.GetGUid(), s.PCaster.GetMP() , pRow.NeedMP)

	if s.PCaster.GetMP() < pRow.NeedMP {
		return false
	}

	nRand := tools.RandNumRange(0, 10000)
	if nRand > pRow.SkillRate {
		return false
	}

	return true
}

//SelectTarget --
//
//	3 4 5		3 0 | 0 3
//	0 1 2		4 1 | 1 4
//	------		5 2 | 2 5
//	0 1 2
//	3 4 5
//
func (s *Skill) SelectTarget() bool {
	pRow := s.GetSkillRow()
	pEnemyList := s.PCaster.GetEnemyList()

	//fmt.Printf("Skill SelectTarget SelectTargetOpt %#v\n",pRow.SelectTargetOpt)

	switch pRow.SelectTargetOpt {
	case EM_SKILL_TARGET_OPT_AUTO: //自动选择
		s.PTarget = s.PCaster
		break
	case EM_SKILL_TARGET_OPT_ORDER: //顺序选择
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			frontIndex := (s.PCaster.GetMatrixID() + i) % (MaxMatrixCellCount / 2)
			backIndex := frontIndex + MaxMatrixCellCount/2
			pFightObject := pEnemyList.GetFightObject(frontIndex)

			if pFightObject.IsActive() {
				s.PTarget = pFightObject
				break
			}

			pFightObject = pEnemyList.GetFightObject(backIndex)

			if pFightObject.IsActive() {
				s.PTarget = pFightObject
				break
			}

		}

		break
	case EM_SKILL_TARGET_OPT_RAND: //随机选择

		nCount := 0
		var nIndexList []int
		for i := 0; i < MaxMatrixCellCount; i++ {
			pFightObject := pEnemyList.GetFightObject(i)

			if pFightObject.IsActive() {
				nIndexList[nCount] = i
				nCount++
			}
		}

		if nCount <= 0 {
			return false
		}

		randIndex := tools.RandNumRange(0, nCount-1)

		s.PTarget = pEnemyList.GetFightObject(nIndexList[randIndex])

		break
	case EM_SKILL_TARGET_OPT_SLOW: //后排优先

		frontIndex := s.PCaster.GetMatrixID() % (MaxMatrixCellCount / 2)
		//后排 对线
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			index := (frontIndex+i)%(MaxMatrixCellCount/2) + MaxMatrixCellCount/2
			pFightObject := pEnemyList.GetFightObject(index)
			if pFightObject.IsActive() {
				s.PTarget = pFightObject
				return true
			}
		}
		//前排 对线
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			index := (frontIndex + i) % (MaxMatrixCellCount / 2)
			pFightObject := pEnemyList.GetFightObject(index)
			if pFightObject.IsActive() {
				s.PTarget = pFightObject
				return true
			}
		}

		break
	}

	//fmt.Printf("Skill SelectTarget PTarget.GetGUid %#v\n",s.PTarget.GetGUid())

	if s.PTarget.IsActive() {
		return true
	}

	return false
}

func (s *Skill) CastSkill(nRound int) bool {

	pRow := s.GetSkillRow()
	//fmt.Printf("=======================\nSkill CastSkill GetSkillRow uid %v, %#v, %#v, %#v\n",s.PCaster.GetGUid(), pRow.Name,pRow.SkillID,pRow.ImpactID)

	s.PCaster.SetMP(s.PCaster.FightDBData.MP - pRow.NeedMP)
	s.SkillTime = nRound + pRow.CoolDownTime
	s.PCaster.AttackInfo.CastUid = s.PCaster.GetGUid()
	s.PCaster.AttackInfo.BSkilled = true

	skillAttack := s.PCaster.AttackInfo.NewSkillAttack()
	if skillAttack == nil {
		return false
	}

	skillAttack.SkillID = s.SkillID
	skillAttack.SkillTarget = s.PTarget.GetGUid()
	skillAttack.CostMp = pRow.NeedMP
	//连击
	nConAttTimes := s.GetConAttTimes()
	for i := 0; i < MaxSkillImpactCount; i++ {
		if pRow.ImpactID[i] <= 0 {
			continue
		}

		nRand := tools.RandNumRange(1,10000)

		//fmt.Printf("Skill CastSkill GetSkillRow %v, nRand %v, pRow.ImpactRate[i] %#v, pRow.ImpactID[i] %v\n",i, nRand, pRow.ImpactRate[i],pRow.ImpactID[i])

		if nRand <= pRow.ImpactRate[i]{
			pRowImpact := GetImpactRow(pRow.ImpactID[i])

			//fmt.Printf("Skill CastSkill pRowImpact %#v\n",pRowImpact)

			if pRowImpact.LogicID == 0 || pRowImpact.LogicID == 1 {
				nConAttTimes = 0
			}

			impactInfo := &ImpactInfo{}
			impactInfo.SkillID = s.SkillID
			impactInfo.ImpactID = pRow.ImpactID[i]
			impactInfo.ConAttTimes = nConAttTimes

			tl := &TargetList{}

			s.GetTargetList(pRow.ImpactTargetType[i], tl)

			targetList := tl

			//fmt.Printf("----Skill CastSkill targetList %#v %#v %#v\n",pRow.ImpactTargetType[i], targetList.GetCount(), targetList)
			//fmt.Printf("----Skill CastSkill tl %#v %#v %v\n", tl.GetCount(),targetList.GetCount(), tl)
			if targetList.GetCount() == 0 {
				continue
			}

			for index := 0; index < targetList.GetCount(); index++ {
				s.PTarget = targetList.GetFightObject(index)
				if s.PTarget.IsActive() {
					//fmt.Printf("----impactInfo.TargetList[impactInfo.TargetCount] %#v\n", s.PTarget.GetGUid())

					impactInfo.TargetList[impactInfo.TargetCount] = s.PTarget.GetGUid()
					impactInfo.TargetCount++
				}
			}

			skillAttack.AddImpactInfo(impactInfo)

			for index := 0; index < targetList.GetCount(); index++ {
				s.PTarget = targetList.GetFightObject(index)
				if s.PTarget.IsActive() {
					//fmt.Printf("Skill CastSkill AddImpact %#v\n",[]interface{}{pRow.ImpactID[i], nConAttTimes, nRound, s.PCaster.GetGUid(), s.SkillID})
					//fmt.Printf("----s.PTarget.AddImpact %#v\n", []interface{}{s.PTarget.GetGUid(), pRow.ImpactID[i], nConAttTimes, nRound, s.PCaster.GetGUid(), s.SkillID})
					//impact结果
					nResult := s.PTarget.AddImpact(pRow.ImpactID[i], nConAttTimes, nRound, s.PCaster, s.SkillID)
					if nResult > 0 {
						fmt.Printf("----Skill CastSkill AddImpact Fail %#v\n", []interface{}{nResult, pRow.ImpactID[i], nConAttTimes, nRound, s.PCaster.GetGUid(), s.SkillID})
					}
				}
			}
		}
	}

	return true
}
func (s *Skill) GetTargetList(nType int, targetList *TargetList) bool {

	switch nType {
	case EmImpactTargetOptSelf: //0=自身
		targetList.Add(s.PCaster)
		break
	case EmImpactTargetOwnerSingle: //1=己方个体
		pOwnerList := s.PCaster.GetOwnerList()
		var nIndexList []int
		nCount := 0
		for i := 0; i < MaxMatrixCellCount; i++ {
			pFightObject := pOwnerList.GetFightObject(i)
			if pFightObject.IsActive() {
				nIndexList[nCount] = i
				nCount++
			}
		}
		if nCount <= 0 {
			return false
		}

		randIndex := tools.RandNumRange(0, nCount-1)
		targetList.Add(pOwnerList.GetFightObject(nIndexList[randIndex]))

		break
	case EmImpactTargetOwnerAll: //2=己方全体
		pOwnerList := s.PCaster.GetOwnerList()
		for i := 0; i < MaxMatrixCellCount; i++ {
			pFightObject := pOwnerList.GetFightObject(i)
			if pFightObject.IsActive() {
				targetList.Add(pFightObject)
			}
		}

		break
	case EmImpactTargetEnemySingle: //3=敌方个体
		targetList.Add(s.PTarget)
		break
	case EmImpactTargetEnemyFront: //4=敌方横排
		matrixIndex := s.PTarget.GetMatrixID()
		pEnemyList := s.PCaster.GetEnemyList()
		if pEnemyList == nil {
			break
		}

		//fmt.Printf("----Skill GetTargetList EmImpactTargetEnemyFront %#v\n",matrixIndex)

		var frontTargetList TargetList
		var backTargetList TargetList
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			//fmt.Printf("----Skill GetTargetList frontTargetList %#v\n",i)

			pFightObject := pEnemyList.GetFightObject(i)
			if pFightObject.IsActive() {
				frontTargetList.Add(pFightObject)
				//fmt.Printf("---------Skill frontTargetList IsActive uid %#v %#v\n",pFightObject.GetGUid(),frontTargetList.GetCount())

			}
		}
		for i := MaxMatrixCellCount / 2; i < MaxMatrixCellCount; i++ {
			//fmt.Printf("----Skill GetTargetList backTargetList %#v\n",i)
			pFightObject := pEnemyList.GetFightObject(i)
			if pFightObject.IsActive() {

				backTargetList.Add(pFightObject)
				//fmt.Printf("---------Skill backTargetList IsActive uid %#v %#v\n",pFightObject.GetGUid(),backTargetList.GetCount())
			}
		}

		targetList = &frontTargetList
		if matrixIndex < MaxMatrixCellCount/2 {
			//fmt.Printf("------Skill frontTargetList.GetCount %#v \n",frontTargetList.GetCount())

			if frontTargetList.GetCount() <= 0 {
				targetList = &backTargetList
			}
		} else {
			//fmt.Printf("------Skill backTargetList.GetCount %#v \n",backTargetList.GetCount())

			if backTargetList.GetCount() > 0 {
				targetList = &backTargetList
			}

		}

		//fmt.Printf("=====Skill targetList.GetCount %#v \n",targetList.GetCount())

		break
	case EmImpactTargetEnemyBehind: //5=敌方后排；
		pEnemyList := s.PCaster.GetEnemyList()
		for i := MaxMatrixCellCount / 2; i < MaxMatrixCellCount; i++ {
			pFightObject := pEnemyList.GetFightObject(i)
			if pFightObject.IsActive() {
				targetList.Add(pFightObject)
			}
		}

		if targetList.GetCount() == 0 {
			for i := 0; i < MaxMatrixCellCount/2; i++ {
				pFightObject := pEnemyList.GetFightObject(i)
				if pFightObject.IsActive() {
					targetList.Add(pFightObject)
				}
			}
		}

		break
	case EmImpactTargetEnemyAll: //6=敌方全体
		pEnemyList := s.PCaster.GetEnemyList()
		for i := 0; i < MaxMatrixCellCount; i++ {
			pFightObject := pEnemyList.GetFightObject(i)
			if pFightObject.IsActive() {
				targetList.Add(pFightObject)
			}
		}
		break
	case EmImpactTargetEnemyLine: //7=敌方目标竖排
		matrixIndex := s.PTarget.GetMatrixID()
		lineIndex := (matrixIndex + MaxMatrixCellCount/2) % MaxMatrixCellCount
		targetList.Add(s.PTarget)

		pEnemyList := s.PCaster.GetEnemyList()

		pFightObject := pEnemyList.GetFightObject(lineIndex)
		if pFightObject.IsActive() {
			targetList.Add(pFightObject)
		}

		break
	case EmImpactTargetEnemyAround: //8=敌方目标及周围

		matrixIndex := s.PTarget.GetMatrixID()
		lineIndex := (matrixIndex + MaxMatrixCellCount/2) % MaxMatrixCellCount
		targetList.Add(s.PTarget)

		pEnemyList := s.PCaster.GetEnemyList()

		pFightObject := pEnemyList.GetFightObject(lineIndex)
		if pFightObject.IsActive() {
			targetList.Add(pFightObject)
		}

		if (((matrixIndex + 1) * 2) % MaxMatrixCellCount) > 0 {
			pFightObject := pEnemyList.GetFightObject(matrixIndex + 1)
			if pFightObject.IsActive() {
				targetList.Add(pFightObject)
			}
		}

		if (matrixIndex-1) > 0 && ((matrixIndex*2)%MaxMatrixCellCount) > 0 {
			pFightObject := pEnemyList.GetFightObject(matrixIndex - 1)
			if pFightObject.IsActive() {
				targetList.Add(pFightObject)
			}
		}

		break
	case EmImpactTargetEnemyBehindone: //后排个体
		pEnemyList := s.PCaster.GetEnemyList()

		frontIndex := s.PTarget.GetMatrixID() % (MaxMatrixCellCount / 2)
		//后排 对线
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			index := (frontIndex+i)%(MaxMatrixCellCount/2) + MaxMatrixCellCount/2
			pFightObject := pEnemyList.GetFightObject(index)
			if pFightObject.IsActive() {
				targetList.Add(pFightObject)
				return true
			}

		}
		//前排 对线
		for i := 0; i < MaxMatrixCellCount/2; i++ {
			index := (frontIndex + i) % (MaxMatrixCellCount / 2)
			pFightObject := pEnemyList.GetFightObject(index)
			if pFightObject.IsActive() {
				targetList.Add(pFightObject)
				return true
			}

		}

		break
	case EmImpactTargetOwnerMinHP:
		nMinHP := InvalidValue
		index := InvalidId

		pOwnerList := s.PCaster.GetOwnerList()
		for i := 0; i < MaxMatrixCellCount; i++ {
			pFightObject := pOwnerList.GetFightObject(i)
			if pFightObject.IsActive() {
				if nMinHP < 0 {
					nMinHP = pFightObject.GetHP()
					index = i
				}

				if nMinHP > pFightObject.GetHP() {
					nMinHP = pFightObject.GetHP()
					index = i
				}
			}
		}

		if index > 0 {
			targetList.Add(pOwnerList.GetFightObject(index))

		}

		break
	case EmImpactTargetOwnerMinMP:
		nMinMP := InvalidValue
		index := InvalidId

		pOwnerList := s.PCaster.GetOwnerList()
		for i := 0; i < MaxMatrixCellCount; i++ {
			pFightObject := pOwnerList.GetFightObject(i)
			if pFightObject.IsActive() {
				if nMinMP < 0 {
					nMinMP = pFightObject.GetMP()
					index = i
				}

				if nMinMP > pFightObject.GetMP() {
					nMinMP = pFightObject.GetMP()
					index = i
				}
			}
		}

		if index > 0 {
			targetList.Add(pOwnerList.GetFightObject(index))

		}
		break
	}

	//fmt.Printf("=====Skill targetList.GetCount %#v \n", targetList.GetCount())

	return true
}

//GetConAttTimes 本次技能连击次数
func (s *Skill) GetConAttTimes() int {
	nConAttTimes := 0
	nContinuous := s.PCaster.GetContinuous()
	if nContinuous > 100 {
		nConAttTimes = s.PCaster.GetConAttTimes()
		return nConAttTimes
	}

	nRand := tools.RandNumRange(1, 100)
	if nRand < nContinuous {
		nConAttTimes = s.PCaster.GetConAttTimes()
	}

	return nConAttTimes
}

//CommonSkillLogic 普通攻击技能
func (s *Skill) CommonSkillLogic(nRound int) bool {

	//fmt.Printf("Skill CommonSkillLogic %#v\n",nRound)
	msg := fmt.Sprintf("%d 发起普通谈判请求",s.PCaster.GetGUid())
	log := NewLog(s.PCaster.TeamNo(),s.PCaster.GetGUid(),LogCommonAtk,0,msg)





	pAttackInfo := s.PCaster.AttackInfo
	if !s.SelectTarget() {
		return false
	}

	pAttackInfo.CastUid = s.PCaster.GetGUid()

	//命中
	bRet := s.CanHit()
	pAttackInfo.BHit = bRet
	if !bRet {
		return true
	}


	nPhysicAttack := float64(s.PCaster.GetPhysicAttack())
	bStrike := s.CanStrike()

	pAttackInfo.BStrike = bStrike
	pAttackInfo.SkillTarget = s.PTarget.GetGUid()
	if bStrike {
		//暴击攻击力=当前攻击力*（1+自身暴击伤害/100）
		nPhysicAttack = nPhysicAttack * (1 + float64(s.PCaster.GetStrikeHurt())*0.01)
	}

	pFightCell := s.PCaster.GetFightCell()

	//单排
	if pFightCell.GetFightType() == EM_TYPE_FIGHT_STAIR && s.PCaster.IsAttacker() {
		nPhysicAttack += nPhysicAttack * float64(pFightCell.GetPlusAtt()) * 0.01
	}

	nDefend := s.PTarget.GetPhysicDefend()
	nDecay := s.PTarget.GetPhysicHurtDecay()

	//fmt.Printf("Skill CommonSkillLogic nPhysicAttack  %#v,nDefend  %#v,nDecay %#v\n", nPhysicAttack, nDefend, nDecay)

	//本次物理攻击伤害=(自身当前经过暴击计算后的物理攻击-目标物理防御)*(1-目标物理伤害减免/100)*(1+物理攻击伤害浮动)
	nDamage := CALCDAMAGE(float64(nPhysicAttack), float64(nDefend), float64(nDecay))
	nDamage = nDamage * (1 + float64(s.PCaster.GetFloatingHurt())*0.01)
	s.PTarget.SetHP(s.PTarget.FightDBData.HP - int(nDamage))

	pAttackInfo.SkillTarget = s.PTarget.GetGUid()

	pAttackInfo.Hurt = int(nDamage)

	msg = fmt.Sprintf("%d 士气下降了 %d 当前士气 %d",s.PTarget.GetTplId(),int(nDamage),s.PTarget.GetHP())
	log.AddAct(s.PTarget.TeamNo(),int(nDamage),s.PTarget.GetTplId(),s.PTarget.GetHP(),LogCommonAtk,msg)


	//fmt.Printf("Skill CommonSkillLogic PCasteruid %#v, PTargetuid %#v, MatrixID %#v, nDamage %#v,PTargetHP %#v\n", s.PCaster.GetGUid(), s.PTarget.GetGUid(), s.PTarget.GetMatrixID(), int(nDamage), s.PTarget.FightDBData.HP)



	_flog := FLog{}
	_flog.Round = nRound
	_flog.TeamNo = s.PCaster.TeamNo()
	_flog.Attacker = s.PCaster.GetGUid()
	_flog.DefendTeamNo = s.PTarget.TeamNo()
	_flog.Defender = s.PTarget.GetGUid()
	_flog.LeftHP = s.PTarget.GetHP()
	_flog.Damage = int(nDamage)
	_flog.Strike = bStrike
	_flog.SkillName = "普通策略"
	_flog.AttackerName = s.PCaster.FightDBData.Name
	_flog.DefenderName = s.PTarget.FightDBData.Name
	pFightCell.AddFLog(_flog)



	//反击
	if s.PTarget.IsActive() && s.CanBackAttack() {
		backAttack := float64(s.PTarget.GetPhysicAttack()) * float64(s.PTarget.GetBackAttHurt()) * 0.01

		s.PCaster.SetHP(s.PCaster.FightDBData.HP - int(backAttack))

		pAttackInfo.BBackAttack = true
		pAttackInfo.BackAttackHurt = int(backAttack)


		_flog := FLog{}
		_flog.Round = nRound
		_flog.TeamNo = s.PTarget.TeamNo()
		_flog.Attacker = s.PTarget.GetGUid()
		_flog.DefendTeamNo = s.PTarget.TeamNo()
		_flog.Defender = s.PCaster.GetGUid()
		_flog.LeftHP = s.PCaster.GetHP()
		_flog.Damage = int(backAttack)
		_flog.BackAttack = pAttackInfo.BBackAttack
		_flog.SkillName = "反问策略"
		_flog.AttackerName = s.PTarget.FightDBData.Name
		_flog.DefenderName = s.PCaster.FightDBData.Name
		pFightCell.AddFLog(_flog)

	}

	pFightCell.AddLog(*log)

	return true
}

//CanHit 是否命中
func (s *Skill) CanHit() bool {
	nHit := s.PCaster.GetHit() - s.PTarget.GetDodge()

	if nHit >= 100 {
		return true
	}

	if nHit <= 0 {
		return false
	}

	nRand := tools.RandNumRange(1,100)

	fmt.Printf("Skill CanHit nRand %#v,nHit %#v\n", nRand,nHit)

	if nRand <= nHit {
		return true
	}

	return false

}

//CanStrike 是否暴击
func (s *Skill) CanStrike() bool {
	nStrike := s.PCaster.GetStrike()
	if nStrike >= 100 {
		return true
	}

	if nStrike <= 0 {
		return false
	}

	nRand := tools.RandNumRange(1,100)

	fmt.Printf("Skill CanStrike nRand %#v,nStrike %#v\n", nRand,nStrike)
	if nRand <= nStrike {
		return true
	}

	return false
}

//CanBackAttack 是否反击
func (s *Skill) CanBackAttack() bool {
	nAttackBack := s.PTarget.GetAttackBack()
	if nAttackBack >= 100 {
		return true
	}

	if nAttackBack <= 0 {
		return false
	}
	nRand := tools.RandNumRange(1,100)
	if nRand <= nAttackBack {

		fmt.Printf("Skill CanBackAttack nRand %#v,nAttackBack %#v\n", nRand, nAttackBack)

		return true
	}

	return false
}

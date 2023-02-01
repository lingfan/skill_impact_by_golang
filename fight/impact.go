package fight

import (
	"cyhd/app/cy/config"
	"fmt"
)

type Impact struct {
	ImpactID    int       //ID
	PHolder     *FightObj //拥有者
	PCaster     *FightObj //施放者
	SkillID     int       //
	StartTime   int       //开始时间
	Life        int       //生命周期
	LogicTime   int       //开始作用时间
	Term        int       //开始间隔
	ConAttTimes int       //连击次数
}

var handler = map[int]func(*Impact){
	0: ImpactLogic0,
	1: ImpactLogic1,
	2: ImpactLogic2,
	3: ImpactLogic3,
	4: ImpactLogic4,
	5: ImpactLogic5,
	6: ImpactLogic6,
}

func NewImpact(nImpactID int, conAttTimes int, nRound int, pCaster, pHolder *FightObj, nSkillID int) *Impact {
	im := &Impact{}
	im.ImpactID = nImpactID
	im.PHolder = pHolder
	im.PCaster = pCaster
	im.SkillID = nSkillID
	im.StartTime = nRound
	im.Life = nRound
	im.LogicTime = nRound
	im.Term = 1
	im.ConAttTimes = 0
	pRow := GetImpactRow(im.ImpactID)
	switch pRow.LogicID {
	case EmImpactLogic0:
		fallthrough
	case EmImpactLogic1:
		fallthrough
	case EmImpactLogic6:
		if conAttTimes > 0 {
			//连击次数
			im.ConAttTimes = conAttTimes
		}
		break
	case EmImpactLogic2:
		fallthrough
	case EmImpactLogic3:
		//2=物理持续攻击
		//逻辑参数3：持续时间，单位回合（10）
		//逻辑参数4：生效间隔，单位回合（2）
		im.Life += pRow.Param[2] - pRow.Param[3]
		im.Term = pRow.Param[3]
		break
	case EmImpactLogic4:
		fallthrough
	case EmImpactLogic5:
		//逻辑参数3：持续时间，单位回合（10）
		im.Life += pRow.Param[2] - 1

		break
	default:
		break

	}
	return im
}

func (im *Impact) GetCaster() *FightObj {
	return im.PCaster
}

func (im *Impact) GetHolder() *FightObj {
	return im.PHolder
}

func (im *Impact) CleanUp() {
	im.ImpactID = InvalidId
	im.PHolder = nil
	im.PCaster = nil
	im.StartTime = 0
	im.Life = 0
	im.LogicTime = 0
	im.Term = 0
	im.ConAttTimes = 0
}

func (im *Impact) IsValid() bool {
	if im.ImpactID > 0 {
		return true
	}
	return false
}

func (im *Impact) HeartBeat(uTime int) bool {

	if uTime == im.LogicTime {
		pRow := GetImpactRow(im.ImpactID)
		//fmt.Printf("Impact HeartBeat PCaster.GetGUid() %v, ImpactID %v, uTime %#v,im.LogicTime %#v, pRow.LogicID %v, skillID %v\n", im.PCaster.GetGUid(), im.ImpactID, uTime, im.LogicTime, pRow.LogicID, im.SkillID)

		if pRow.LogicID >= 0 && pRow.LogicID < EmImpactLogicCount {
			//fmt.Printf("Impact HeartBeat uTime %#v,im.Life %#v\n", uTime, im.Life)
			//fmt.Printf("Impact HeartBeat LogicID %#v\n", pRow.LogicID)
			h := handler[pRow.LogicID]
			h(im)
		}

	}

	if uTime >= im.Life {
		im.CleanUp()
	}
	return true
}

//ImpactLogic0 单次物理
//	0=单次物理攻击；
//	逻辑参数1：从英雄物理攻击中取得的倍率（例：150）
//	逻辑参数2：额外增加的物理伤害（例：50）
//	如英雄物理攻击为100，则最终的技能物理攻击=英雄物理攻击（100）*逻辑参数1（150/100）+逻辑参数2（50）

func ImpactLogic0(im *Impact) {
	//fmt.Printf("Impact ImpactLogic0\n")

	pAttackInfo := im.PCaster.GetAttackInfo()

	pSkillAttack := pAttackInfo.GetSkillAttack(im.SkillID)
	pImpactInfo := pSkillAttack.GetImpactInfo(im.ImpactID)
	index := pImpactInfo.GetTargetIndex(im.PHolder.GetGUid())

	pFightCell := im.PHolder.GetFightCell()

	_skill := config.GetSkillRow(im.SkillID)
	msg := fmt.Sprintf("%d 谈判请求 采用了直接策略 %s(%d)", im.PCaster.GetTplId(), _skill.Name, im.SkillID)

	log := NewLog(im.PCaster.TeamNo(), im.PCaster.GetTplId(), LogSkillAtk, im.SkillID, msg)

	nPhysicAttack := im.PCaster.GetPhysicAttack()
	pRow := GetImpactRow(im.ImpactID)

	//fmt.Printf("Impact ImpactLogic0 pRow %#v\n", pRow)

	nAttack := int(float64(nPhysicAttack*pRow.Param[0])*0.01) + pRow.Param[1]

	if pFightCell.GetFightType() == EM_TYPE_FIGHT_STAIR && im.PCaster.IsAttacker() {
		nAttack += nAttack * pFightCell.GetPlusAtt() / 100
	}

	nDefend := im.PHolder.GetPhysicDefend()
	nDecay := im.PHolder.GetPhysicHurtDecay()
	//本次物理攻击伤害=(自身当前经过连击计算后的物理攻击-目标物理防御)*(1-目标物理伤害减免/100)
	nDamage := int(CALCDAMAGE(float64(nAttack), float64(nDefend), float64(nDecay)))
	im.PHolder.SetHP(im.PHolder.GetHP() - nDamage)

	_flog := FLog{}
	_flog.Round = im.LogicTime
	_flog.TeamNo = im.PCaster.TeamNo()
	_flog.Attacker = im.PCaster.GetGUid()
	_flog.DefendTeamNo = im.PHolder.TeamNo()
	_flog.Defender = im.PHolder.GetGUid()
	_flog.LeftHP = im.PHolder.GetHP()
	_flog.Damage = int(nDamage)
	_flog.SkillID = im.SkillID
	_flog.ImpactID = im.ImpactID
	_flog.LogicID = pRow.LogicID
	_flog.SkillName = _skill.Name
	_flog.AttackerName = im.PCaster.FightDBData.Name
	_flog.DefenderName = im.PHolder.FightDBData.Name
	pFightCell.AddFLog(_flog)

	msg = fmt.Sprintf("%d 士气下降了 %d 当前士气 %d", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())

	log.AddAct(im.PHolder.TeamNo(), nDamage, im.PHolder.GetTplId(), im.PHolder.GetHP(), LogSkillAtk, msg)

	if index >= 0 {
		pImpactInfo.Hurts[index][0] = nDamage
	}

	fConAttHurt := float64(im.PCaster.GetConAttHurt()) * 0.01
	for i := 1; i <= im.ConAttTimes; i++ {
		nAttack := float64(nAttack) * fConAttHurt
		nDamage := int(CALCDAMAGE(float64(nAttack), float64(nDefend), float64(nDecay)))
		im.PHolder.SetHP(im.PHolder.GetHP() - nDamage)
		if index >= 0 {
			pImpactInfo.Hurts[index][i] = nDamage

			msg = fmt.Sprintf("%d 士气持续下降了 %d 当前士气 %d", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())

			log.AddAct(im.PHolder.TeamNo(), nDamage, im.PHolder.GetTplId(), im.PHolder.GetHP(), LogSkillAtk, msg)

			_flog := FLog{}
			_flog.Round = im.LogicTime
			_flog.TeamNo = im.PCaster.TeamNo()
			_flog.Attacker = im.PCaster.GetGUid()
			_flog.DefendTeamNo = im.PHolder.TeamNo()
			_flog.Defender = im.PHolder.GetGUid()
			_flog.LeftHP = im.PHolder.GetHP()
			_flog.Damage = int(nDamage)
			_flog.SkillID = im.SkillID
			_flog.ImpactID = im.ImpactID
			_flog.LogicID = pRow.LogicID
			_flog.ConAttTimes = i
			_flog.SkillName = _skill.Name
			_flog.AttackerName = im.PCaster.FightDBData.Name
			_flog.DefenderName = im.PHolder.FightDBData.Name
			pFightCell.AddFLog(_flog)
		}
	}
	pFightCell.AddLog(*log)
}

//ImpactLogic1 单次魔法
//1=单次魔法攻击
//	逻辑参数1：从英雄物理攻击（因为英雄默认没有魔法攻击）中取得的倍率（例：150）
//	逻辑参数2：额外增加的魔法伤害（例：50）
//	如英雄物理攻击为100，则最终的技能魔法攻击=英雄物理攻击（100）*逻辑参数1（150/100）+逻辑参数2（50）

func ImpactLogic1(im *Impact) {
	//fmt.Printf("Impact ImpactLogic1\n")
	pAttackInfo := im.PCaster.GetAttackInfo()
	pSkillAttack := pAttackInfo.GetSkillAttack(im.SkillID)
	pImpactInfo := pSkillAttack.GetImpactInfo(im.ImpactID)
	index := pImpactInfo.GetTargetIndex(im.PHolder.GetGUid())

	pFightCell := im.PHolder.GetFightCell()

	_skill := config.GetSkillRow(im.SkillID)
	msg := fmt.Sprintf("%d 谈判请求 采用了特殊策略 %s(%d)", im.PCaster.GetTplId(), _skill.Name, im.SkillID)
	log := NewLog(im.PCaster.TeamNo(), im.PCaster.GetTplId(), LogSkillAtk, im.SkillID, msg)

	nPhysicAttack := im.PCaster.GetPhysicAttack()
	pRow := GetImpactRow(im.ImpactID)

	nAttack := int(float64(nPhysicAttack*pRow.Param[0])*0.01) + pRow.Param[1]

	if pFightCell.GetFightType() == EM_TYPE_FIGHT_STAIR && im.PCaster.IsAttacker() {
		nAttack += nAttack * pFightCell.GetPlusAtt() / 100
	}

	nDefend := im.PHolder.GetMagicDefend()
	nDecay := im.PHolder.GetMagicHurtDecay()
	//本次魔法攻击伤害=(自身当前经过连击计算后的魔法攻击-目标魔法防御)*(1-目标魔法伤害减免/100)
	nDamage := int(CALCDAMAGE(float64(nAttack), float64(nDefend), float64(nDecay)))
	im.PHolder.SetHP(im.PHolder.GetHP() - nDamage)

	//fmt.Printf("Impact ImpactLogic1 nAttack %#v,nDefend %#v,nDecay %#v,nDamage %#v, hp %#v\n", nAttack, nDefend, nDecay, nDamage, im.PHolder.GetHP())

	_flog := FLog{}
	_flog.Round = im.LogicTime
	_flog.TeamNo = im.PCaster.TeamNo()
	_flog.Attacker = im.PCaster.GetGUid()
	_flog.DefendTeamNo = im.PHolder.TeamNo()
	_flog.Defender = im.PHolder.GetGUid()
	_flog.LeftHP = im.PHolder.GetHP()
	_flog.Damage = int(nDamage)
	_flog.SkillID = im.SkillID
	_flog.ImpactID = im.ImpactID
	_flog.LogicID = pRow.LogicID
	_flog.SkillName = _skill.Name
	_flog.AttackerName = im.PCaster.FightDBData.Name
	_flog.DefenderName = im.PHolder.FightDBData.Name
	pFightCell.AddFLog(_flog)

	if index >= 0 {
		pImpactInfo.Hurts[index][0] = nDamage
		msg = fmt.Sprintf("%d 士气下降了 %d 当前士气 %d", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())
		log.AddAct(im.PHolder.TeamNo(), nDamage, im.PHolder.GetTplId(), im.PHolder.GetHP(), LogSkillAtk, msg)
	}

	fConAttHurt := float64(im.PCaster.GetConAttHurt()) * 0.01
	for i := 1; i <= im.ConAttTimes; i++ {
		nAttack := float64(nAttack) * fConAttHurt
		nDamage := int(CALCDAMAGE(float64(nAttack), float64(nDefend), float64(nDecay)))
		im.PHolder.SetHP(im.PHolder.GetHP() - nDamage)
		if index >= 0 {
			pImpactInfo.Hurts[index][i] = nDamage
			msg = fmt.Sprintf("%d 士气持续下降了 %d 当前士气 %d", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())
			log.AddAct(im.PHolder.TeamNo(), nDamage, im.PHolder.GetTplId(), im.PHolder.GetHP(), LogSkillAtk, msg)

			_flog := FLog{}
			_flog.Round = im.LogicTime
			_flog.TeamNo = im.PCaster.TeamNo()
			_flog.Attacker = im.PCaster.GetGUid()
			_flog.DefendTeamNo = im.PHolder.TeamNo()
			_flog.Defender = im.PHolder.GetGUid()
			_flog.LeftHP = im.PHolder.GetHP()
			_flog.Damage = int(nDamage)
			_flog.SkillID = im.SkillID
			_flog.ImpactID = im.ImpactID
			_flog.LogicID = pRow.LogicID
			_flog.ConAttTimes = i
			_flog.SkillName = _skill.Name
			_flog.AttackerName = im.PCaster.FightDBData.Name
			_flog.DefenderName = im.PHolder.FightDBData.Name
			pFightCell.AddFLog(_flog)
		}
	}
	pFightCell.AddLog(*log)
	//fmt.Printf("Impact ImpactLogic1 %#v\n",pImpactInfo)

}

//ImpactLogic2 持续物理
//	2=物理持续攻击
//	逻辑参数1：从英雄物理攻击中取得的倍率（例：150）
//	逻辑参数2：额外增加的物理伤害（例：50）
//	逻辑参数3：持续时间，单位回合（10）
//	逻辑参数4：生效间隔，单位回合（2）
//	此效果的最终效果为：每2回合对目标造成一次物理伤害，持续10回合（生效5次），每次造成的物理攻击具体数值=英雄物理攻击（100）*逻辑参数1（150/100）+逻辑参数2（50）
func ImpactLogic2(im *Impact) {
	fmt.Printf("Impact ImpactLogic2\n")

	pAttackInfo := im.PCaster.GetAttackInfo()
	pSkillAttack := pAttackInfo.GetSkillAttack(im.SkillID)
	pImpactInfo := pSkillAttack.GetImpactInfo(im.ImpactID)
	index := pImpactInfo.GetTargetIndex(im.PHolder.GetGUid())

	pFightCell := im.PHolder.GetFightCell()

	if im.LogicTime <= im.Life {
		_skill := config.GetSkillRow(im.SkillID)
		msg := fmt.Sprintf("%d 谈判请求 采用了持续直接策略 %s(%d)", im.PCaster.GetTplId(), _skill.Name, im.SkillID)

		log := NewLog(im.PCaster.TeamNo(), im.PCaster.GetTplId(), LogSkillAtk, im.SkillID, msg)

		nPhysicAttack := im.PCaster.GetPhysicAttack()
		pRow := GetImpactRow(im.ImpactID)
		nAttack := int(float64(nPhysicAttack*pRow.Param[0])*0.01) + pRow.Param[1]
		if pFightCell.GetFightType() == EM_TYPE_FIGHT_STAIR && im.PCaster.IsAttacker() {
			nAttack += nAttack * pFightCell.GetPlusAtt() / 100
		}
		nDefend := im.PHolder.GetPhysicDefend()
		nDecay := im.PHolder.GetPhysicHurtDecay()
		//本次物理攻击伤害=(自身当前经过连击计算后的物理攻击-目标物理防御)*(1-目标物理伤害减免/100)
		nDamage := int(CALCDAMAGE(float64(nAttack), float64(nDefend), float64(nDecay)))
		im.PHolder.SetHP(im.PHolder.GetHP() - nDamage)

		_flog := FLog{}
		_flog.Round = im.LogicTime
		_flog.TeamNo = im.PCaster.TeamNo()
		_flog.Attacker = im.PCaster.GetGUid()
		_flog.DefendTeamNo = im.PHolder.TeamNo()
		_flog.Defender = im.PHolder.GetGUid()
		_flog.LeftHP = im.PHolder.GetHP()
		_flog.Damage = int(nDamage)
		_flog.SkillID = im.SkillID
		_flog.ImpactID = im.ImpactID
		_flog.LogicID = pRow.LogicID
		_flog.SkillName = _skill.Name
		_flog.AttackerName = im.PCaster.FightDBData.Name
		_flog.DefenderName = im.PHolder.FightDBData.Name
		pFightCell.AddFLog(_flog)

		fmt.Printf("Impact ImpactLogic2 nAttack %#v,nDefend %#v,nDecay %#v,nDamage %#v, hp %#v\n", nAttack, nDefend, nDecay, nDamage, im.PHolder.GetHP())

		if im.StartTime == im.LogicTime {

			if index >= 0 {
				pImpactInfo.Hurts[index][0] = nDamage
				fmt.Printf("%d 士气下降了 %d 当前士气 %d\n", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())

				//log.AddAct(im.PHolder.TeamNo(),nDamage,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
			}
		} else if im.StartTime < im.LogicTime {
			pRoundInfo := pFightCell.GetRoundInfo()
			pFObjectInfo := pRoundInfo.GetFObjectInfoByGuid(im.PHolder.GetGUid())
			for i := 0; i < pFObjectInfo.ImpactCount; i++ {
				if pFObjectInfo.ImpactList[i] == im.ImpactID {
					pFObjectInfo.ImpactHurt[i] = nDamage
					fmt.Printf("%d 士气下降了 %d 当前士气 %d\n", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())
					//log.AddAct(im.PHolder.TeamNo(),nDamage,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
				}
			}
		}
		pFightCell.AddLog(*log)
	}
	im.LogicTime += im.Term

}

//ImpactLogic3 持续魔法
//	3=持续魔法攻击
//	逻辑参数1：从英雄物理攻击（因为英雄默认没有魔法攻击）中取得的倍率（例：150）
//	逻辑参数2：额外增加的魔法伤害（例：50）
//	逻辑参数3：持续时间，单位回合（10）
//	逻辑参数4：生效间隔，单位回合（2）
//	此效果的最终效果为：每2回合对目标造成一次魔法伤害，持续10回合（生效5次），每次造成的魔法攻击具体数值=英雄物理攻击（100）*逻辑参数1（150/100）+逻辑参数2（50）
func ImpactLogic3(im *Impact) {
	fmt.Printf("Impact ImpactLogic3\n")

	pAttackInfo := im.PCaster.GetAttackInfo()
	pSkillAttack := pAttackInfo.GetSkillAttack(im.SkillID)
	pImpactInfo := pSkillAttack.GetImpactInfo(im.ImpactID)
	index := pImpactInfo.GetTargetIndex(im.PHolder.GetGUid())

	pFightCell := im.PHolder.GetFightCell()

	if im.LogicTime <= im.Life {
		_skill := config.GetSkillRow(im.SkillID)
		msg := fmt.Sprintf("%d 谈判请求 采用了持续特殊策略 %s(%d)", im.PCaster.GetTplId(), _skill.Name, im.SkillID)

		log := NewLog(im.PCaster.TeamNo(), im.PCaster.GetTplId(), LogSkillAtk, im.SkillID, msg)

		nPhysicAttack := im.PCaster.GetPhysicAttack()
		pRow := GetImpactRow(im.ImpactID)
		nAttack := int(float64(nPhysicAttack*pRow.Param[0])*0.01) + pRow.Param[1]
		if pFightCell.GetFightType() == EM_TYPE_FIGHT_STAIR && im.PCaster.IsAttacker() {
			nAttack += nAttack * pFightCell.GetPlusAtt() / 100
		}
		nDefend := im.PHolder.GetMagicDefend()
		nDecay := im.PHolder.GetMagicHurtDecay()
		//本次物理攻击伤害=(自身当前经过连击计算后的物理攻击-目标物理防御)*(1-目标物理伤害减免/100)
		nDamage := int(CALCDAMAGE(float64(nAttack), float64(nDefend), float64(nDecay)))
		im.PHolder.SetHP(im.PHolder.GetHP() - nDamage)

		fmt.Printf("Impact ImpactLogic3 nAttack %#v,nDefend %#v,nDecay %#v,nDamage %#v, hp %#v\n", nAttack, nDefend, nDecay, nDamage, im.PHolder.GetHP())
		_flog := FLog{}
		_flog.Round = im.LogicTime
		_flog.TeamNo = im.PCaster.TeamNo()
		_flog.Attacker = im.PCaster.GetGUid()
		_flog.DefendTeamNo = im.PHolder.TeamNo()
		_flog.Defender = im.PHolder.GetGUid()
		_flog.LeftHP = im.PHolder.GetHP()
		_flog.Damage = int(nDamage)
		_flog.SkillID = im.SkillID
		_flog.ImpactID = im.ImpactID
		_flog.LogicID = pRow.LogicID
		_flog.SkillName = _skill.Name
		_flog.AttackerName = im.PCaster.FightDBData.Name
		_flog.DefenderName = im.PHolder.FightDBData.Name
		pFightCell.AddFLog(_flog)

		if im.StartTime == im.LogicTime {

			if index >= 0 {
				pImpactInfo.Hurts[index][0] = nDamage
				fmt.Printf("%d 士气下降了 %d 当前士气 %d\n", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())
				//log.AddAct(im.PHolder.TeamNo(),nDamage,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
			}
		} else if im.StartTime < im.LogicTime {
			pRoundInfo := pFightCell.GetRoundInfo()
			pFObjectInfo := pRoundInfo.GetFObjectInfoByGuid(im.PHolder.GetGUid())
			for i := 0; i < pFObjectInfo.ImpactCount; i++ {
				if pFObjectInfo.ImpactList[i] == im.ImpactID {
					pFObjectInfo.ImpactHurt[i] = nDamage
					fmt.Printf("%d 士气下降了 %d 当前士气 %d\n", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())
					//log.AddAct(im.PHolder.TeamNo(),nDamage,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
				}
			}
		}
		pFightCell.AddLog(*log)
	}
	im.LogicTime += im.Term
}

//ImpactLogic4 强化Buff
//	4=buff强化类
//	逻辑参数1：改变的英雄属性id，读取AttributeData.tab表。
//	逻辑参数2：改变的具体数值
//	逻辑参数3：持续时间，单位回合（10）
//	最终可实现的效果如英雄攻击增加X点持续10回合。
func ImpactLogic4(im *Impact) {
	//fmt.Printf("Impact ImpactLogic4\n")
	pAttackInfo := im.PCaster.GetAttackInfo()
	pSkillAttack := pAttackInfo.GetSkillAttack(im.SkillID)
	pImpactInfo := pSkillAttack.GetImpactInfo(im.ImpactID)
	index := pImpactInfo.GetTargetIndex(im.PHolder.GetGUid())

	if im.LogicTime <= im.Life {
		pFightCell := im.PHolder.GetFightCell()

		_skill := config.GetSkillRow(im.SkillID)
		msg := fmt.Sprintf("%d 谈判请求 采用了加强策略 %s(%d)", im.PCaster.GetTplId(), _skill.Name, im.SkillID)
		log := NewLog(im.PCaster.TeamNo(), im.PCaster.GetTplId(), LogSkillAtk, im.SkillID, msg)

		pRow := GetImpactRow(im.ImpactID)
		im.PHolder.ChangeEffect(pRow.Param[0], pRow.Param[1])
		msg = fmt.Sprintf("%d 属性 %d 加强 %d", im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])



		if im.StartTime == im.LogicTime {
			_flog := FLog{}
			_flog.Round = im.LogicTime
			_flog.TeamNo = im.PCaster.TeamNo()
			_flog.Attacker = im.PCaster.GetGUid()
			_flog.DefendTeamNo = im.PHolder.TeamNo()
			_flog.Defender = im.PHolder.GetGUid()
			_flog.SkillID = im.SkillID
			_flog.ImpactID = im.ImpactID
			_flog.LogicID = pRow.LogicID
			_flog.EffectIdx = pRow.Param[0]
			_flog.EffectVal = pRow.Param[1]
			_flog.EffectName = LangEmAttribute[pRow.Param[0]]

			_flog.SkillName = _skill.Name
			_flog.AttackerName = im.PCaster.FightDBData.Name
			_flog.DefenderName = im.PHolder.FightDBData.Name
			pFightCell.AddFLog(_flog)

			if pRow.Param[0] == EmAttributeHp {
				im.PHolder.SetHP(im.PHolder.GetHP())
				pImpactInfo.Hurts[index][0] = pRow.Param[1] * (-1)

				//fmt.Printf("im.StartTime %d , im.LogicTime %d , %d 属性 %d 加强了 %d\n", im.StartTime, im.LogicTime, im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])

			}
			if pRow.Param[0] == EmAttributeMp {
				//nMP := im.PHolder.GetMP()
				pImpactInfo.Mp[index] = pRow.Param[1] * (-1)

				//fmt.Printf("im.StartTime %d , im.LogicTime %d , %d 属性 %d 加强了 %d\n", im.StartTime, im.LogicTime, im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])
			}
		} else if im.StartTime < im.LogicTime {
			//fmt.Printf("4 im.StartTime %d < im.LogicTime %d , %d\n", im.StartTime, im.LogicTime, im.PHolder.GetGUid())

			pFightCell := im.PHolder.GetFightCell()
			pRoundInfo := pFightCell.GetRoundInfo()
			pFObjectInfo := pRoundInfo.GetFObjectInfoByGuid(im.PHolder.GetGUid())

			for i := 0; i < pFObjectInfo.ImpactCount; i++ {
				if pFObjectInfo.ImpactList[i] == im.ImpactID {
					if pRow.Param[0] == EmAttributeHp {
						pFObjectInfo.ImpactHurt[i] = pRow.Param[1] * -1
						fmt.Printf("im.StartTime %d , im.LogicTime %d , %d 属性 %d 加强了 %d\n", im.StartTime, im.LogicTime, im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])
					}
					if pRow.Param[0] == EmAttributeMp {
						pFObjectInfo.ImpactMP[i] = pRow.Param[1] * -1
						fmt.Printf("im.StartTime %d , im.LogicTime %d , %d 属性 %d 加强了 %d\n", im.StartTime, im.LogicTime, im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])
					}
				}
			}
		}
		pFightCell.AddLog(*log)

	}

	im.LogicTime++

}

//ImpactLogic5 削弱Debuff
//	5=debuff削弱类
//	逻辑参数1：改变的英雄属性id，读取AttributeData.tab表。
//	逻辑参数2：改变的具体数值
//	逻辑参数3：持续时间，单位回合（10）
//	最终可实现的效果如敌人攻击减少X点持续10回合。
func ImpactLogic5(im *Impact) {
	//fmt.Printf("Impact ImpactLogic5\n")
	pAttackInfo := im.PCaster.GetAttackInfo()
	pSkillAttack := pAttackInfo.GetSkillAttack(im.SkillID)
	pImpactInfo := pSkillAttack.GetImpactInfo(im.ImpactID)
	index := pImpactInfo.GetTargetIndex(im.PHolder.GetGUid())

	if im.LogicTime <= im.Life {
		pFightCell := im.PHolder.GetFightCell()

		_skill := config.GetSkillRow(im.SkillID)
		msg := fmt.Sprintf("%d 谈判请求 采用削弱策略 %s(%d)", im.PCaster.GetTplId(), _skill.Name, im.SkillID)
		log := NewLog(im.PCaster.TeamNo(), im.PCaster.GetTplId(), LogSkillAtk, im.SkillID, msg)

		pRow := GetImpactRow(im.ImpactID)
		im.PHolder.ChangeEffect(pRow.Param[0], pRow.Param[1]*-1)
		msg = fmt.Sprintf("%d 属性 %d 削弱 %d", im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])
		log.AddAct(im.PHolder.TeamNo(), 0, im.PHolder.GetTplId(), im.PHolder.GetHP(), LogSkillAtk, msg)


		if im.StartTime == im.LogicTime {
			_flog := FLog{}
			_flog.Round = im.LogicTime
			_flog.TeamNo = im.PCaster.TeamNo()
			_flog.Attacker = im.PCaster.GetGUid()
			_flog.DefendTeamNo = im.PHolder.TeamNo()
			_flog.Defender = im.PHolder.GetGUid()
			_flog.SkillID = im.SkillID
			_flog.ImpactID = im.ImpactID
			_flog.LogicID = pRow.LogicID
			_flog.EffectIdx = pRow.Param[0]
			_flog.EffectVal = pRow.Param[1]
			_flog.EffectName = LangEmAttribute[pRow.Param[0]]

			_flog.SkillName = _skill.Name
			_flog.AttackerName = im.PCaster.FightDBData.Name
			_flog.DefenderName = im.PHolder.FightDBData.Name
			pFightCell.AddFLog(_flog)

			if pRow.Param[0] == EmAttributeHp {
				im.PHolder.SetHP(im.PHolder.GetHP())
				pImpactInfo.Hurts[index][0] = pRow.Param[1]
				fmt.Printf("im.StartTime %d , im.LogicTime %d , %d 属性 %d 削弱了 %d\n", im.StartTime, im.LogicTime, im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])
				//log.AddAct(im.PHolder.TeamNo(),0,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
			}
			if pRow.Param[0] == EmAttributeMp {
				im.PHolder.SetMP(im.PHolder.GetMP())
				pImpactInfo.Mp[index] = pRow.Param[1]
				fmt.Printf("im.StartTime %d , im.LogicTime %d , %d 属性 %d 削弱了 %d\n", im.StartTime, im.LogicTime, im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])
				//log.AddAct(im.PHolder.TeamNo(),0,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
			}
		} else if im.StartTime < im.LogicTime {
			//fmt.Printf("5 im.StartTime %d < im.LogicTime %d , %d\n", im.StartTime, im.LogicTime, im.PHolder.GetGUid())

			pFightCell := im.PHolder.GetFightCell()
			pRoundInfo := pFightCell.GetRoundInfo()
			pFObjectInfo := pRoundInfo.GetFObjectInfoByGuid(im.PHolder.GetGUid())

			for i := 0; i < pFObjectInfo.ImpactCount; i++ {
				if pFObjectInfo.ImpactList[i] == im.ImpactID {
					if pRow.Param[0] == EmAttributeHp {
						pFObjectInfo.ImpactHurt[i] = pRow.Param[1]
						fmt.Printf("im.StartTime %d , im.LogicTime %d , %d 属性 %d 削弱了 %d\n", im.StartTime, im.LogicTime, im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])
						//log.AddAct(im.PHolder.TeamNo(),0,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
					}
					if pRow.Param[0] == EmAttributeMp {
						pFObjectInfo.ImpactMP[i] = pRow.Param[1]
						fmt.Printf("im.StartTime %d , im.LogicTime %d , %d 属性 %d 削弱了 %d\n", im.StartTime, im.LogicTime, im.PHolder.GetTplId(), pRow.Param[0], pRow.Param[1])
						//log.AddAct(im.PHolder.TeamNo(),0,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
					}
				}
			}
		}
		pFightCell.AddLog(*log)

	}

	im.LogicTime++

}

//ImpactLogic6 单次魔法攻击
//	6=单次魔法攻击
//	逻辑参数1：从英雄魔法攻击中取得的倍率（例：150）
//	逻辑参数2：额外增加的魔法伤害（例：50）
//	如英雄魔法攻击为100，则最终的技能魔法攻击=英雄魔法攻击（100）*逻辑参数1（150/100）+逻辑参数2（50）
func ImpactLogic6(im *Impact) {
	fmt.Printf("Impact ImpactLogic6\n")

	pAttackInfo := im.PCaster.GetAttackInfo()
	pSkillAttack := pAttackInfo.GetSkillAttack(im.SkillID)
	pImpactInfo := pSkillAttack.GetImpactInfo(im.ImpactID)
	index := pImpactInfo.GetTargetIndex(im.PHolder.GetGUid())

	pFightCell := im.PHolder.GetFightCell()

	_skill := config.GetSkillRow(im.SkillID)
	msg := fmt.Sprintf("%d 谈判请求 使用单次特殊策略 %s(%d)", im.PCaster.GetTplId(), _skill.Name, im.SkillID)
	log := NewLog(im.PCaster.TeamNo(), im.PCaster.GetTplId(), LogSkillAtk, im.SkillID, msg)

	nMagicAttack := im.PCaster.GetMagicAttack()
	pRow := GetImpactRow(im.ImpactID)
	nAttack := int(float64(nMagicAttack*pRow.Param[0])*0.01) + pRow.Param[1]
	if pFightCell.GetFightType() == EM_TYPE_FIGHT_STAIR && im.PCaster.IsAttacker() {
		nAttack += nAttack * pFightCell.GetPlusAtt() / 100
	}
	nDefend := im.PHolder.GetMagicDefend()
	nDecay := im.PHolder.GetMagicHurtDecay()
	//本次魔法攻击伤害=(自身当前经过连击计算后的魔法攻击-目标魔法防御)*(1-目标魔法伤害减免/100)
	nDamage := int(CALCDAMAGE(float64(nAttack), float64(nDefend), float64(nDecay)))
	im.PHolder.SetHP(im.PHolder.GetHP() - nDamage)

	fmt.Printf("Impact ImpactLogic6 nAttack %#v,nDefend %#v,nDecay %#v,nDamage %#v, hp %#v\n", nAttack, nDefend, nDecay, nDamage, im.PHolder.GetHP())
	_flog := FLog{}
	_flog.Round = im.LogicTime
	_flog.TeamNo = im.PCaster.TeamNo()
	_flog.Attacker = im.PCaster.GetGUid()
	_flog.DefendTeamNo = im.PHolder.TeamNo()
	_flog.Defender = im.PHolder.GetGUid()
	_flog.LeftHP = im.PHolder.GetHP()
	_flog.Damage = int(nDamage)
	_flog.SkillID = im.SkillID
	_flog.ImpactID = im.ImpactID
	_flog.LogicID = pRow.LogicID
	_flog.SkillName = _skill.Name
	_flog.AttackerName = im.PCaster.FightDBData.Name
	_flog.DefenderName = im.PHolder.FightDBData.Name
	pFightCell.AddFLog(_flog)

	if index >= 0 {
		pImpactInfo.Hurts[index][0] = nDamage
		fmt.Printf("%d 士气下降了 %d 当前士气 %d\n", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())

		//log.AddAct(im.PHolder.TeamNo(),nDamage,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
	}
	fConAttHurt := float64(im.PCaster.GetConAttHurt()) * 0.01
	for i := 1; i <= im.ConAttTimes; i++ {
		nAttack := float64(nAttack) * fConAttHurt
		nDamage := int(CALCDAMAGE(float64(nAttack), float64(nDefend), float64(nDecay)))
		im.PHolder.SetHP(im.PHolder.GetHP() - nDamage)
		if index >= 0 {
			pImpactInfo.Hurts[index][i] = nDamage
			fmt.Printf("%d 士气下降了 %d 当前士气 %d\n", im.PHolder.GetTplId(), int(nDamage), im.PHolder.GetHP())
			//log.AddAct(im.PHolder.TeamNo(),nDamage,im.PHolder.GetTplId(),im.PHolder.GetHP(),LogSkillAtk,msg)
		}
	}
	pFightCell.AddLog(*log)

}

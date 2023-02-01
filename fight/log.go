package fight

type Log struct {
	AtkTeamNo int
	AtkHeroId int
	AtkType   int
	ActLog    []ActLog
	SkillId   int
	Msg       string
}

type ActLog struct {
	HeroId  int
	Damage  int
	LeftHp  int
	IsCrit  int
	ActType int //1发起普通攻击 2发起技能攻击
	TeamNo  int //1 進攻方 2 防守方
	Msg     string
}

const (
	LogCommonAtk = iota + 1
	LogSkillAtk
)

func NewLog(atkTeamNo int, atkHeroId int, atkType int, skillId int, msg string) *Log {
	log := &Log{}
	log.AtkTeamNo = atkTeamNo
	log.AtkType = atkType
	log.AtkHeroId = atkHeroId
	log.ActLog = []ActLog{}
	log.SkillId = skillId
	log.Msg = msg

	return log
}

func (log *Log) AddAct(teamNo, damage, heroId, leftHp, actType int, msg string) {

	actLog := ActLog{}
	actLog.TeamNo = teamNo
	actLog.Damage = damage
	actLog.HeroId = heroId
	actLog.LeftHp = leftHp
	actLog.ActType = actType
	actLog.Msg = msg

	log.ActLog = append(log.ActLog, actLog)
}

type FLog struct {
	No int
	Round int
	TeamNo int
	Attacker int
	DefendTeamNo int
	Defender int
	LogicID int
	SkillID int
	ImpactID int
	Damage int
	LeftHP int
	EffectIdx int
	EffectVal int
	Strike bool
	BackAttack bool
	ConAttTimes int

	SkillName string
	EffectName string
	AttackerName string
	DefenderName string
}

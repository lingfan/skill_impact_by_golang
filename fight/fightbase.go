package fight

import (
	"cyhd/app/cy/config"
)

const (
	MaxMatrixCellCount  = 6   //技能最大选择目标
	MaxFightRound       = 128 //最大战斗回合
	MaxImpactNumber     = 20  //impact数量
	MaxSkillNum         = 2   //技能数量
	MaxSkillImpactCount = 4   //最大技能所带impact
	MaxConAttackTimes   = 4   //最大连击次数+1
	MaxBuffNumber       = 10  //最大buff数量
	InvalidId           = -1
	InvalidValue        = -1
	FightDistance       = 100
)

const (
	EM_HERO_HUMAN   = 0
	EM_HERO_MONSTER = 1
)

//技能选择目标方式
const (
	EM_SKILL_TARGET_OPT_AUTO  = -1 //-1=自动选择
	EM_SKILL_TARGET_OPT_ORDER = 0  //0=顺序选择
	EM_SKILL_TARGET_OPT_RAND  = 1  //1=随机选择
	EM_SKILL_TARGET_OPT_SLOW  = 2  //2=后排优先
)

const (
	EM_SKILL_TYPE_INVALID       = -1
	EM_SKILL_TYPE_HERO_ACTIVE   = 0 //0=英雄主动技能
	EM_SKILL_TYPE_HERO_PASSIVE  = 1 //1=英雄被动技能
	EM_SKILL_TYPE_EQUIP_ACTIVE  = 2 //2=装备主动技能
	EM_SKILL_TYPE_EQUIP_PASSIVE = 3 //3=装备被动技能

)

//IMPACT选择目标方式
const (
	EmImpactTargetOptAuto        = -1
	EmImpactTargetOptSelf        = 0  //0=自身；
	EmImpactTargetOwnerSingle    = 1  //1=己方个体；
	EmImpactTargetOwnerAll       = 2  //2=己方全体；
	EmImpactTargetEnemySingle    = 3  //3=敌方个体；
	EmImpactTargetEnemyFront     = 4  //4=敌方横排；
	EmImpactTargetEnemyBehind    = 5  //5=敌方后排优先；
	EmImpactTargetEnemyAll       = 6  //6=敌方全体；
	EmImpactTargetEnemyLine      = 7  //7=敌方目标竖排；
	EmImpactTargetEnemyAround    = 8  //8=敌方目标及周围；
	EmImpactTargetEnemyBehindone = 9  //9=敌方后排个体
	EmImpactTargetOwnerMinHP     = 10 //10=己方血最少
	EmImpactTargetOwnerMinMP     = 11 //11=己方蓝最少
	EmImpactTargetOptNumber      = 12
)

const (
	EmImpactLogic0 = iota
	EmImpactLogic1
	EmImpactLogic2  //持续物伤
	EmImpactLogic3  //持续法伤
	EmImpactLogic4
	EmImpactLogic5
	EmImpactLogic6
	EmImpactLogicCount
)

//impact结果
const (
	EM_IMPACT_RESULT_NORMAL    = iota //正常加目标身上
	EM_IMPACT_RESULT_FAIL             //不能加在目标身上
	EM_IMPACT_RESULT_DISSAPEAR        //抵消
)

const (
	EM_TYPE_IMPACTLOGIC_INVALID = -1
	EM_TYPE_IMPACTLOGIC_SINGLE  = iota //单次生效
	EM_TYPE_IMPACTLOGIC_BUFF           //强化
	EM_TYPE_IMPACTLOGIC_DEBUFF         //削弱

	EM_TYPE_IMPACTLOGIC_COUNT
)

const (
	EM_TYPE_FIGHT_NORMAL = iota
	EM_TYPE_FIGHT_STAIR   //单排

	EM_TYPE_FIGHT_COUNT
)

type FightDBData struct {
	GUid            int              //对象guid
	Name 	string
	TableID         int              //表格id
	Type            int              //英雄出处
	Quality         int              //品质
	MatrixID        int              //阵法Id
	Profession      int              //职业 近战0 远程1
	Level           int              //等级
	HP              int              //血量
	MP              int              //魔法值
	SkillCount      int              //技能数量
	Skill           [MaxSkillNum]int //技能列表
	EquipSkillCount int
	EquipSkill      []int //装备技能列表

	//战斗属性
	PhysicAttack    int //物理攻击
	MagicAttack     int //魔法攻击，
	PhysicDefence   int //物理防御，
	MagicDefence    int //魔法防御，
	MaxHP           int //最大生命
	MaxMP           int //最大魔法值
	Hit             int //命中值，
	Dodge           int //闪避值，
	Strike          int //暴击值，
	StrikeHurt      int //暴击伤害
	Continuous      int //连击值
	ConAttHurt      int //连击伤害
	ConAttTimes     int //连击次数
	BackAttack      int //反击值
	BackAttHurt     int //反击伤害
	AttackSpeed     int //攻击速度，
	PhysicHurtDecay int //物理减免
	MagicHurtDecay  int //魔法减免
	FloatingHurt int//浮动伤害

	//				DamageReduce int				//总伤害减免（百分比）
	//				nReduceCriticalDamage int	//减免暴击伤害
	//				nExtraHurt int				//附加伤害
	///新增属性
	Exp       int   //经验
	GrowRate  int   //成长
	BearPoint int   //负重
	Equip     []int //装备ID
	Color     int   //颜色
}

type FIGHTDB struct {
	Uid         int
	FightDBData []*FightDBData
}

var LangEmAttribute = map[int]string{

	EmAttributeMaxHP                :"最大生命",
	EmAttributeMoveSpeed            :"移动速度",
	EmAttributeAttackSpeed          :"攻击速度",
	EmAttributePhysicAttack         :"物理攻击",
	EmAttributePhysicDefence        :"物理防御",
	EmAttributeHit                  :"命中点数",
	EmAttributeDodge                :"闪避点数",
	EmAttributeStrike               :"暴击",
	EmAttributeContinuous           :"连击",
	EmAttributeBackAttack           :"反击",
	EmAttributeContinuousTimes      :"连击次数",
	EmAttributeHurtContinuous       :"连击伤害",
	EmAttributeHurtBackAttack       :"反击伤害",
	EmAttributeHurtStrike           :"暴击伤害",
	EmAttributePhysicHurtDecay      :"物理伤害减免",
	EmAttributeStrikelhurtDecay     :"暴击伤害减免",
	EmAttributeHurtExtral           :"附加伤害",
	EmAttributeHurtPhysic           :"普通攻击伤害",
	EmAttributeMagicAttack          :"魔法攻击",
	EmAttributeMagicDefence         :"魔法防御",
	EmAttributeMagicHurtDecay       :"魔法伤害减免",
	EmAttributeMaxMP                :"魔法值上限",
	EmAttributePercentAttackSpeed   :"攻击速度",//百分比
	EmAttributePercentPhysicAttack  :"物理攻击",//百分比
	EmAttributePercentMagicAttack   :"魔法攻击",//百分比
	EmAttributePercentPhysicDefence :"物理防御",//百分比
	EmAttributePercentMagicDefence  :"魔法防御",//百分比
	EmAttributePercentMaxHP         :"最大生命值",//百分比
	EmAttributePercentMaxMP         :"最大魔法值",//百分比
	EmAttributeLevel                :"等级",
	EmAttributeHp                   :"血量",//血量
	EmAttributeMp                   :"魔法值",
	EmAttributeCurrentexp           :"当前经验",
	EmAttributeAction               :"行动力",
}
const (
	EmAttributeInvalid              = 0
	EmAttributeMaxHP                = 1  //最大生命
	EmAttributeMoveSpeed            = 2  //移动速度
	EmAttributeAttackSpeed          = 3  //攻击速度
	EmAttributePhysicAttack         = 4  //物理攻击
	EmAttributePhysicDefence        = 5  //物理防御
	EmAttributeHit                  = 6  //命中点数
	EmAttributeDodge                = 7  //闪避点数
	EmAttributeStrike               = 8  //暴击
	EmAttributeContinuous           = 9  //连击
	EmAttributeBackAttack           = 10 //反击
	EmAttributeContinuousTimes      = 11 //连击次数
	EmAttributeHurtContinuous       = 12 //连击伤害
	EmAttributeHurtBackAttack       = 13 //反击伤害
	EmAttributeHurtStrike           = 14 //暴击伤害
	EmAttributePhysicHurtDecay      = 15 //物理伤害减免
	EmAttributeStrikelhurtDecay     = 16 //暴击伤害减免
	EmAttributeHurtExtral           = 17 //附加伤害
	EmAttributeHurtPhysic           = 18 //普通攻击伤害
	EmAttributeMagicAttack          = 19 //魔法攻击
	EmAttributeMagicDefence         = 20 //魔法防御
	EmAttributeMagicHurtDecay       = 21 //魔法伤害减免
	EmAttributeMaxMP                = 22 //最大魔法值
	EmAttributePercentAttackSpeed   = 23 //攻击速度百分比
	EmAttributePercentPhysicAttack  = 24 //物理攻击百分比
	EmAttributePercentMagicAttack   = 25 //魔法攻击百分比
	EmAttributePercentPhysicDefence = 26 //物理防御百分比
	EmAttributePercentMagicDefence  = 27 //魔法防御百分比
	EmAttributePercentMaxHP         = 28 //最大生命值百分比
	EmAttributePercentMaxMP         = 29 //最大魔法值百分比
	EmAttributeLevel                = 30 //等级
	EmAttributeHp                   = 31 //血量
	EmAttributeMp                   = 32 //魔法值
	EmAttributeCurrentexp           = 33 //当前经验
	EmAttributeAction               = 34 //行动力
	EmAttributeNumber               = 35
)

//TableRowHeroAttr --
type TableRowHeroAttr struct {
	ID                        int    //英雄ID
	SpiritID                  int    //英雄魂魄ID
	Name                      string //英雄名称
	Profession                int    //英雄职业
	InitQuality               int    //初始品质
	Initlevel                 int    //初始等级
	LevelLimit                int    //等级上限
	LevelCrossRole            int    //英雄比人物高出等级上限
	InitExp                   int    //初始经验
	TakeLevel                 int    //可携带等级
	EffectAttackByGrowRate    int    //成长值对英雄攻击点数的影响系数
	EffectDefendByGrowRate    int    //成长值对英雄防御点数的影响系数
	EffectHPByGrowRate        int    //成长值对英雄生命上限点数的影响系数
	EffectMPByGrowRate        int    //成长值对英雄魔法上限点数的影响系数
	EffectPhysicAttackbyLevel int    //英雄物理攻击随等级的增长系数
	EffectPhysicDefendbyLevel int    //英雄物理防御随等级的增长系数
	EffectHpbyLevel           int    //英雄生命上限随等级的增长系数
	EffectMpbyLevel           int    //英雄魔法上限随等级的增长系数
	InitAttackSpeed           int    //初始速度
	InitPhysicAttack          int    //初始物理攻击
	InitMagicAttack           int    //初始魔法攻击
	InitPhysicDefence         int    //初始物理防御
	InitMagicDefence          int    //初始魔法防御
	InitHP                    int    //初始生命值
	InitMP                    int    //初始魔法值
	InitHit                   int    //初始命中值
	InitDodge                 int    //初始闪避值
	InitStrike                int    //初始暴击值
	InitContinuous            int    //初始连击值	计算发生概率
	InitBackAttack            int    //初始反击值
	InitStrikeHurt            int    //初始暴击伤害
	InitConAttTimes           int    //初始连击次数
	InitConAttHurt            int    //初始连击伤害	计算伤害倍率
	InitBackAttHurt           int    //初始反击伤害
	InitPhysicHurtDecay       int    //初始物理减免
	InitMagicHurtDecay        int    //初始魔法减免
	FloatingHurt              int    //英雄伤害浮动
	PhysicSkillID             int    //普通攻击技能
	MagicSkillID1             int    //魔法技能1
	LearnedLevel1             int    //学会等级1
	MagicSkillID2             int    //魔法技能2
	LearnedLevel2             int    //学会等级2
	MaxGrowPoint              int    //成长值上限
	InitBearPoint             int    //初始承载力
	BearParam                 int    //承载力随等级的增长系数
	Icon                      string //图标
	Sound                     string //出场声音
	RequireHumanLevel         int    //需要召唤师等级
}

type TableRowImpact struct {
	ImpactID      int    //ImpactID
	Description   string //策划描述
	LogicID       int    //Impact逻辑id
	Param         []int  //逻辑参数
	Icon          string //Buff图标
	ImpactMutexID int    //IMPACT互斥id
	ReplaceLevel  int    //顶替优先级
	DeadDisapeer  int    //死亡后是否消失
	OfflineTimeGO int    //下线是否计时	INT
	ScriptID      int    //脚本ID
	szEffect      string //特效
	szName        string //名称
	szDesc        string //描述
	szSkillEffect string //技能效果
}

//GetImpactRow --
func GetImpactRow(id int) config.TableImpact {
	_list := config.GetTableImpactList()
	for _, v := range _list {
		if v.ImpactID == id {
			return v
		}
	}

	return config.TableImpact{}
}

type FObjectData struct {
	Uid        int
	Level      int //等级
	MatrixID   int //位置
	Profession int //职业
	TableId    int //表格ID
	Quality    int //品质
}

func (fod *FObjectData) CleanUp() {
	fod.Uid = 0
}

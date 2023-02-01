package fight

import (
	"cyhd/app/cy/config"
	"cyhd/app/cy/models/userbo"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

type User struct {
	Uid      int
	HeroList [MaxMatrixCellCount]*FightDBData
}

func TestFightCell_Fight(t *testing.T) {
	TFightCell(1)

}

func TFightCell(lv int) {


	user1 := User{}
	user1.Uid = 1
	user1.HeroList = [MaxMatrixCellCount]*FightDBData{}



	hid1 := userbo.DrawHero(1,6)
	hid2 := userbo.DrawHero(1,6)

	for i := 0; i < MaxMatrixCellCount; i++ {
		fd := &FightDBData{}

		tplId := hid1[i]
		fd.GUid = (i+1)*10000+tplId
		if tplId > 0 {
			hero := userbo.NewHero(tplId)
			hero.CurQuality = 1
			hero.Lv = lv
			hero.CalcHeroAttr()
			_tabHero := config.GetHeroData(tplId)
			fd.Name = _tabHero.Name

			fd.TableID = tplId
			fd.Quality = hero.CurQuality
			fd.Level = hero.Lv
			fd.HP = hero.Hp
			fd.MaxHP = hero.Hp
			fd.MP = hero.Mp
			fd.MaxMP = hero.Hp

			fd.PhysicAttack = hero.Atk
			fd.PhysicDefence = hero.Def
			fd.MagicAttack = hero.MAtk
			fd.MagicDefence = hero.MDef

			fd.Hit = hero.Hit
			fd.Dodge = hero.Dodge
			fd.Strike = hero.Strike

			fd.AttackSpeed = hero.AtkSpeed



			for i := 0; i < MaxSkillNum; i++ {
				fd.Skill[i] = 0
				if hero.FightSkills[i].Lv > 0 {
					fd.Skill[i] = hero.FightSkills[i].ID
				}
			}

		}
		user1.HeroList[i] = fd
	}

	user2 := User{}
	user2.Uid = 2
	for i := 0; i < MaxMatrixCellCount; i++ {
		fd := &FightDBData{}
		tplId := hid2[i]
		fd.GUid = (i+1)*10000+tplId
		if tplId > 0 {
			hero := userbo.NewHero(tplId)
			hero.CurQuality = 1
			hero.Lv = lv
			hero.CalcHeroAttr()
			_tabHero := config.GetHeroData(tplId)
			fd.Name = _tabHero.Name

			fd.TableID = tplId
			fd.Quality = hero.CurQuality
			fd.Level = hero.Lv
			fd.HP = hero.Hp
			fd.MaxHP = hero.Hp
			fd.MP = hero.Mp
			fd.MaxMP = hero.Hp

			fd.PhysicAttack = hero.Atk
			fd.PhysicDefence = hero.Def
			fd.MagicAttack = hero.MAtk
			fd.MagicDefence = hero.MDef

			fd.Hit = hero.Hit
			fd.Dodge = hero.Dodge
			fd.Strike = hero.Strike

			fd.AttackSpeed = hero.AtkSpeed

			for i := 0; i < MaxSkillNum; i++ {
				fd.Skill[i] = 0
				if hero.FightSkills[i].Lv > 0 {
					fd.Skill[i] = hero.FightSkills[i].ID
				}
			}

		}

		user2.HeroList[i] = fd
	}

	fc := NewFightCell()

	//fmt.Printf("FightCell %#v\n", fc)

	FillFightObjList(user1, fc.GetAttackList())
	fc.InitAttackerList()
	FillFightObjList(user2, fc.GetDefenceList())
	fc.InitDefenderList(0)
	fc.Fight()

	//fmt.Printf("FightCell %#v\n", fc)

	//_log1, _ := jsoniter.MarshalToString(fc.FightInfo)
	//fmt.Printf("FightCell FightInfo %s\n", _log1)
	//_log2, _ := jsoniter.MarshalToString(fc.RoundInfo)
	//fmt.Printf("FightCell RoundInfo %s\n", _log2)

	//_log3, _ := jsoniter.MarshalToString(fc.GetLog())
	//fmt.Printf("FightCell GetLog %s\n", _log3)
	//_log4, _ := jsoniter.MarshalToString(fc.GetFLog())
	//fmt.Printf("FightCell GetFLog %s\n", _log4)

	ret := fc.GetLogList()
	fmt.Printf("FightCell LFLog %v\n", len(ret))
	_log5, _ := jsoniter.MarshalToString(ret)
	fmt.Printf("FightCell LFLog %s\n", _log5)

}


func FillFightObjList(u User, fightObjList *FightObjList) {
	//fmt.Printf("FillFightObjList User %#v\n", u)
	//fmt.Printf("FillFightObjList HeroList %#v\n", len(u.HeroList))
	fightObjList.Owner = u.Uid
	for i := 0; i < MaxMatrixCellCount; i++ {
		fo := &FightObj{}
		fo.InitFightDBData(u.HeroList[i])
		fightObjList.FillObject(i, fo)
	}

	//fmt.Printf("FillFightObjList FightObjList %#v\n", fightObjList)
}

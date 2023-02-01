package fight

//FightObjList --
type FightObjList struct {
	Owner   int
	ObjList [MaxMatrixCellCount]*FightObj
}

func (fol *FightObjList) CleanUp() {
	fol.Owner = 0
	for i := 0; i < MaxMatrixCellCount; i++ {
		fol.ObjList[i].CleanUp()
	}
}

func (fol *FightObjList) GetActiveCount() int {
	nCount := 0
	for i := 0; i < MaxMatrixCellCount; i++ {
		if fol.ObjList[i].IsActive() {
			nCount++
		}
	}
	return nCount
}

func (fol *FightObjList) GetInactiveCount() int {
	nCount := 0
	for i := 0; i < MaxMatrixCellCount; i++ {
		if fol.ObjList[i].IsValid() && !fol.ObjList[i].IsActive() {
			nCount++
		}
	}
	return nCount
}

func (fol *FightObjList) GetFightObject(idx int) *FightObj {
	//fmt.Printf("GetFightObject idx %#v\n",idx)

	if idx >= 0 && idx < MaxMatrixCellCount {
		return fol.ObjList[idx]
	}
	return &FightObj{}
}

func (fol *FightObjList) FillObject(idx int, obj *FightObj) {
	obj.SetMatrixID(idx)
	fol.ObjList[idx] = obj

	//fmt.Printf("FightObjList FillObject idx %#v, obj %#v\n",idx,obj)
}

func (fol *FightObjList) HeartBeat(uTime int) {
	//fmt.Printf("FightObjList HeartBeat %#v\n",uTime)

	for i := 0; i < MaxMatrixCellCount; i++ {
		if fol.ObjList[i].IsActive() {
			fol.ObjList[i].HeartBeat(uTime)
			pAttackInfo := fol.ObjList[i].GetAttackInfo()
			if pAttackInfo != nil && pAttackInfo.IsValid() {
				pFightCell := fol.ObjList[i].GetFightCell()
				pRoundInfo := pFightCell.GetRoundInfo()
				pRoundInfo.AddAttackInfo(pAttackInfo)

			}
		}
	}
}

func (fol *FightObjList) ClearImpactEffect() {
	for i := 0; i < MaxMatrixCellCount; i++ {
		if fol.ObjList[i] != nil && fol.ObjList[i].IsActive() {
			fol.ObjList[i].ClearImpactEffect()

		}
	}
}

func (fol *FightObjList) ImpactHeartBeat(uTime int) {
	//fmt.Printf("FightObjList ImpactHeartBeat %#v\n",uTime)

	for i := 0; i < MaxMatrixCellCount; i++ {
		if fol.ObjList[i].IsActive() {
			//fmt.Printf("FightObjList ImpactHeartBeat uTime %#v, i %#v, uid %#v\n",uTime,i,fol.ObjList[i].GetGUid())
			fol.ObjList[i].ClearImpactEffect()
			fol.ObjList[i].ImpactHeartBeat(uTime)
		}
	}
}

func (fol *FightObjList) Init() {

}

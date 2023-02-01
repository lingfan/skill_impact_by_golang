package fight

type TargetList struct {
	PObjectList [MaxMatrixCellCount]*FightObj
	NCount      int
}

func (tl *TargetList) GetCount() int {
	return tl.NCount
}

func (tl *TargetList) GetFightObject(index int) *FightObj {
	if index >= 0 && index < tl.NCount {
		return tl.PObjectList[index]
	}
	return &FightObj{}
}

func (tl *TargetList) Add(pObj *FightObj) bool {
	if tl.NCount >= MaxMatrixCellCount {
		return false
	}

	tl.PObjectList[tl.NCount] = pObj
	tl.NCount++
	return true

}

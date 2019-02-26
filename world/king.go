package world

type King struct {
	world      *World
	Color      int32
	Units      []*Unit
	UnitCount  int32
	Formations []*Formation
}

func (m *World) NewKing(color int32) *King {
	k := new(King)
	k.world = m
	k.Color = color
	k.Units = make([]*Unit, 50)
	k.Formations = make([]*Formation, 11)
	for i := 0; i < len(k.Formations); i++ {
		k.Formations[i] = NewFormation()
	}
	return k
}

func (k *King) AddUnit(u *Unit) {
	if k.UnitCount == int32(len(k.Units)) {
		array := make([]*Unit, k.UnitCount+5)
		copy(array, k.Units)
		k.Units = array
	}
	k.Units[k.UnitCount] = u
	k.UnitCount++
}

func (k *King) RemoveUnit(u *Unit) {
	for i := int32(0); i < k.UnitCount; i++ {
		if k.Units[i] == u {
			for j := i; j < k.UnitCount-1; j++ {
				k.Units[j] = k.Units[j+1]
			}
			k.UnitCount--
			break
		}
	}
}

func (k *King) MoveOrder(radian float64, x, y int32) {
	k.Formations[0].Order(STATUS_MOVE, radian, x, y)
}

func (k *King) DoodadOrder(d *Doodad) {
	f := k.Formations[0]
	for i := int32(0); i < f.UnitCount; i++ {
		u := f.Units[i]
		u.TargetDoodad = d
		u.Command = STATUS_DOODAD
	}
}

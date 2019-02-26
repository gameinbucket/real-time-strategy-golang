package world

type Tile struct {
	world         *World
	Closed        bool
	SpriteID      int32
	Units         []*Unit
	UnitCount     int32
	Physical      []*Unit
	PhysicalCount int32
	Doodads       []*Doodad
	DoodadCount   int32
	Nav           *NavBox
	IsEdge        bool
}

func NewTile(world *World) *Tile {
	t := new(Tile)
	t.world = world
	t.Units = make([]*Unit, 5)
	t.Physical = make([]*Unit, 10)
	t.Doodads = make([]*Doodad, 0)
	t.Closed = false
	t.SpriteID = 0
	return t
}

func (t *Tile) AddUnit(u *Unit) {
	if t.UnitCount == int32(len(t.Units)) {
		array := make([]*Unit, t.UnitCount+5)
		copy(array, t.Units)
		t.Units = array
	}
	t.Units[t.UnitCount] = u
	t.UnitCount++
}

func (t *Tile) RemoveUnit(u *Unit) {
	for i := int32(0); i < t.UnitCount; i++ {
		if t.Units[i] == u {
			for j := i; j < t.UnitCount-1; j++ {
				t.Units[j] = t.Units[j+1]
			}
			t.UnitCount--
			break
		}
	}
}

func (t *Tile) AddPhysical(u *Unit) {
	if t.PhysicalCount == int32(len(t.Physical)) {
		array := make([]*Unit, t.PhysicalCount+5)
		copy(array, t.Physical)
		t.Physical = array
	}
	t.Physical[t.PhysicalCount] = u
	t.PhysicalCount++

	if t.PhysicalCount == 2 {
		t.world.AddTileCache(t)
	}
}

func (t *Tile) RemovePhysical(u *Unit) {
	for i := int32(0); i < t.PhysicalCount; i++ {
		if t.Physical[i] == u {
			for j := i; j < t.PhysicalCount-1; j++ {
				t.Physical[j] = t.Physical[j+1]
			}
			t.PhysicalCount--
			if t.PhysicalCount == 1 {
				t.world.RemoveTileCache(t)
			}
			break
		}
	}
}

func (t *Tile) AddDoodad(d *Doodad) {
	if t.DoodadCount == int32(len(t.Doodads)) {
		array := make([]*Doodad, t.DoodadCount+1)
		copy(array, t.Doodads)
		t.Doodads = array
	}
	t.Doodads[t.DoodadCount] = d
	t.DoodadCount++
}

func (t *Tile) RemoveDoodad(d *Doodad) {
	for i := int32(0); i < t.DoodadCount; i++ {
		if t.Doodads[i] == d {
			for j := i; j < t.DoodadCount-1; j++ {
				t.Doodads[j] = t.Doodads[j+1]
			}
			t.DoodadCount--
			break
		}
	}
}

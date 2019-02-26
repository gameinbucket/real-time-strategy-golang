package world

import (
	"fmt"
	"strconv"
	"sync"

	"../fixed"
)

type World struct {
	MapSize               int32
	Shift                 uint32
	TilePixels            int32
	FixedTilePixels       int32
	Pixels                int32
	PixelsFixed           int32
	HalfPixelsFixed       int32
	Tiles                 []*Tile
	Kings                 []*King
	DeadUnits             []*Unit
	DeadUnitCount         int32
	TilePhysicsCache      []*Tile
	TilePhysicsCacheCount int32
	CollisionSet          map[string]bool
}

func NewWorld(size int32, shift uint32) *World {
	m := new(World)
	m.MapSize = size
	m.Shift = shift
	m.TilePixels = int32(1 << uint32(shift))
	m.FixedTilePixels = fixed.Whole(m.TilePixels)
	m.Pixels = size * m.TilePixels
	m.PixelsFixed = fixed.Whole(m.Pixels)
	m.HalfPixelsFixed = fixed.Whole(m.Pixels >> 1)
	m.Tiles = make([]*Tile, size*size)
	m.DeadUnits = make([]*Unit, 10)
	m.TilePhysicsCache = make([]*Tile, size)
	for i := int32(0); i < size*size; i++ {
		m.Tiles[i] = NewTile(m)
	}
	m.CollisionSet = make(map[string]bool)
	return m
}

func (m *World) GetTile(x, y int32) *Tile {
	for x < 0 {
		x += m.MapSize
	}
	for y < 0 {
		y += m.MapSize
	}
	for x >= m.MapSize {
		x -= m.MapSize
	}
	for y >= m.MapSize {
		y -= m.MapSize
	}
	return m.Tiles[x+y*m.MapSize]
}

func (m *World) AddTileCache(t *Tile) {
	if m.TilePhysicsCacheCount == int32(len(m.TilePhysicsCache)) {
		array := make([]*Tile, m.TilePhysicsCacheCount+9)
		copy(array, m.TilePhysicsCache)
		m.TilePhysicsCache = array
	}
	m.TilePhysicsCache[m.TilePhysicsCacheCount] = t
	m.TilePhysicsCacheCount++
}

func (m *World) RemoveTileCache(t *Tile) {
	for i := int32(0); i < m.TilePhysicsCacheCount; i++ {
		if m.TilePhysicsCache[i] == t {
			for j := i; j < m.TilePhysicsCacheCount-1; j++ {
				m.TilePhysicsCache[j] = m.TilePhysicsCache[j+1]
			}
			m.TilePhysicsCacheCount--
			break
		}
	}
}

func (m *World) AddDeadUnit(u *Unit) {
	if m.DeadUnitCount == int32(len(m.DeadUnits)) {
		array := make([]*Unit, m.DeadUnitCount+9)
		copy(array, m.DeadUnits)
		m.DeadUnits = array
	}
	m.DeadUnits[m.DeadUnitCount] = u
	m.DeadUnitCount++
}

func (m *World) Integrate() {
	var wg sync.WaitGroup
	for i := 0; i < len(m.Kings); i++ {
		k := m.Kings[i]
		wg.Add(int(k.UnitCount))
		for j := int32(0); j < k.UnitCount; j++ {
			go k.Units[j].Vision(&wg, m)
		}
	}
	for k := range m.CollisionSet {
		delete(m.CollisionSet, k)
	}
	for i := int32(0); i < m.TilePhysicsCacheCount; i++ {
		t := m.TilePhysicsCache[i]
		for j := int32(0); j < t.PhysicalCount; j++ {
			a := t.Physical[j]
			for k := int32(j + 1); k < t.PhysicalCount; k++ {
				b := t.Physical[k]
				var minX int32
				var minY int32
				var maxX int32
				var maxY int32
				if a.X > b.X {
					minX = b.X
					maxX = a.X
				} else {
					minX = a.X
					maxX = b.X
				}
				if a.Y > b.Y {
					minY = b.Y
					maxY = a.Y
				} else {
					minY = a.Y
					maxY = b.Y
				}
				key := strconv.Itoa(int(minX)) + "-" + strconv.Itoa(int(minY)) + "-" + strconv.Itoa(int(maxX)) + "-" + strconv.Itoa(int(maxY))
				if _, ok := m.CollisionSet[key]; !ok {
					m.UnitOverlap(a, b)
					m.CollisionSet[key] = true
				}
			}
		}
	}
	wg.Wait()
	for i := 0; i < len(m.Kings); i++ {
		k := m.Kings[i]
		for j := int32(0); j < k.UnitCount; j++ {
			k.Units[j].Integrate(m)
		}
	}
	for i := int32(0); i < m.DeadUnitCount; i++ {
		u := m.DeadUnits[i]
		m.GetTile(u.GX, u.GY).RemoveUnit(u)
		for gx := u.MinGX; gx <= u.MaxGX; gx++ {
			for gy := u.MinGY; gy <= u.MaxGY; gy++ {
				m.GetTile(gx, gy).RemovePhysical(u)
			}
		}
		u.KingFor.RemoveUnit(u)
		if u.Form != nil {
			u.Form.RemoveUnit(u)
		}
	}
}

func (m *World) UnitOverlap(a, b *Unit) {
	dxx := a.X - b.X
	dyy := a.Y - b.Y
	if dxx < 0 {
		if -dxx > m.HalfPixelsFixed {
			dxx += m.PixelsFixed
		}
	} else {
		if dxx > m.HalfPixelsFixed {
			dxx -= m.PixelsFixed
		}
	}
	if dyy < 0 {
		if -dyy > m.HalfPixelsFixed {
			dyy += m.PixelsFixed
		}
	} else {
		if dyy > m.HalfPixelsFixed {
			dyy -= m.PixelsFixed
		}
	}
	repel := a.Radius + b.Radius
	dist := fixed.Mul(dxx, dxx) + fixed.Mul(dyy, dyy)
	if dist < fixed.Mul(repel, repel)*2 {
		UnitInteraction(a, b)
	}
	if dist > fixed.Mul(repel, repel) {
		return
	}
	dist = fixed.Sqrt(dist)
	if dist == 0 {
		dist = 1
	}
	mult := fixed.Div(repel, dist)
	fx := fixed.Mul(dxx, mult) - dxx
	fy := fixed.Mul(dyy, mult) - dyy
	a.DX += fx
	a.DY += fy
	b.DX -= fx
	b.DY -= fy
}

func UnitInteraction(a, b *Unit) {
	/*
		if !FormationUpdate(a, b) {
			if a.Direction != b.Direction || a.Mirror != b.Mirror {
				if a.Position != nil && b.Position != nil {
					if !a.Position.Filled && b.Position.Filled {
						b.Status = STATUS_STEP_ASIDE
						b.MoveY -= fixed.Whole(32)
						fmt.Println("step aside b")
					} else if a.Position.Filled && !b.Position.Filled {
						a.Status = STATUS_STEP_ASIDE
						a.MoveY -= fixed.Whole(32)
						fmt.Println("step aside a")
					}
				}
			}
		}
	*/
	/*
		rules for unit movement:
			if both moving in same direction:
				ignore each other
			else if both moving different directions:
				if moving to same coordinate:
					unit farther from coordinate waits at closest different tile
				else:
					each unit moves to closest side of other unit (if open space, else pathfind)
			else only one is moving:
				if moving unit heading towards idle unit:
					idle unit steps aside by temporarily moving perpendicular to moving units direction (if open space, else pathfind)

	*/
	if a.Status == STATUS_MOVE {
		if b.Status == STATUS_MOVE {
			UnitInteractBothMove(a, b)
		} else {
			UnitInteractOneMove(a, b)
		}
	} else if b.Status == STATUS_MOVE {
		UnitInteractOneMove(b, a)
	}
}

func UnitInteractBothMove(a, b *Unit) {
	if a.Direction == b.Direction && a.Mirror == b.Mirror {
		return
	}
	dx := a.X - b.X
	//dy := a.Y - b.Y
	if a.Direction == b.Direction {
		if dx > 0 {
			a.MoveY = a.Y + a.Radius + b.Radius
			b.MoveY = b.Y - b.Radius - a.Radius
		}
	}
}

func UnitInteractOneMove(mov, idl *Unit) {
	dx := mov.X - idl.X
	dy := mov.Y - idl.Y
	switch mov.Direction {
	case 0: // up
		if dy > 0 {
			if dx > 0 {
				idl.MoveX = idl.X - mov.Radius
			} else {
				idl.MoveX = idl.X + mov.Radius
			}
		}
	case 1: // up right or left
		adx := dx
		if adx < 0 {
			adx = -adx
		}
		ady := dy
		if ady < 0 {
			ady = -ady
		}
		if mov.Mirror {
			if adx > ady {
				if dx > 0 {
					idl.MoveY = idl.Y + mov.Radius
					idl.MoveX = idl.X - mov.Radius
				}
			} else {
				if dy > 0 {
					idl.MoveY = idl.Y - mov.Radius
					idl.MoveX = idl.X + mov.Radius
				}
			}
		} else {
			if adx > ady {
				if dx < 0 {
					idl.MoveY = idl.Y + mov.Radius
					idl.MoveX = idl.X + mov.Radius
				}
			} else {
				if dy > 0 {
					idl.MoveY = idl.Y - mov.Radius
					idl.MoveX = idl.X - mov.Radius
				}
			}
		}
	case 2: // right or left
		if mov.Mirror {
			if dx > 0 {
				if dy > 0 {
					idl.MoveY = idl.Y - mov.Radius
				} else {
					idl.MoveY = idl.Y + mov.Radius
				}
			}
		} else {
			if dx < 0 {
				if dy > 0 {
					idl.MoveY = idl.Y - mov.Radius
				} else {
					idl.MoveY = idl.Y + mov.Radius
				}
			}
		}
	case 3: // down right or left
		adx := dx
		if adx < 0 {
			adx = -adx
		}
		ady := dy
		if ady < 0 {
			ady = -ady
		}
		if mov.Mirror {
			if adx > ady {
				if dx > 0 {
					idl.MoveY = idl.Y - mov.Radius
					idl.MoveX = idl.X - mov.Radius
				}
			} else {
				if dy < 0 {
					idl.MoveY = idl.Y + mov.Radius
					idl.MoveX = idl.X + mov.Radius
				}
			}
		} else {
			if adx > ady {
				if dx < 0 {
					idl.MoveY = idl.Y - mov.Radius
					idl.MoveX = idl.X + mov.Radius
				}
			} else {
				if dy < 0 {
					idl.MoveY = idl.Y + mov.Radius
					idl.MoveX = idl.X - mov.Radius
				}
			}
		}
	case 4: // down
		if dy < 0 {
			if dx > 0 {
				idl.MoveX = idl.X - mov.Radius
			} else {
				idl.MoveX = idl.X + mov.Radius
			}
		}
	}
}

func (m *World) Print() {
	fmt.Println(m.Shift)
}

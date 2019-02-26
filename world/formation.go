package world

import (
	"math"
	"math/rand"

	"../fixed"
)

type Position struct {
	Form   *Formation
	X      int32
	Y      int32
	Left   *Position
	Right  *Position
	Front  *Position
	Back   *Position
	Filled bool
}

type Formation struct {
	Units     []*Unit
	UnitCount int32
	Positions []Position
}

func NewFormation() *Formation {
	f := new(Formation)
	f.Units = make([]*Unit, 1)
	return f
}

func (f *Formation) Clear() {
	for i := int32(0); i < f.UnitCount; i++ {
		f.Units[i].Form = nil
		f.Units[i].Position = nil
	}
	f.UnitCount = 0
}

func (f *Formation) AddUnit(u *Unit) {
	if f.UnitCount == int32(len(f.Units)) {
		array := make([]*Unit, f.UnitCount+5)
		copy(array, f.Units)
		f.Units = array
	}
	u.Form = f
	f.Units[f.UnitCount] = u
	f.UnitCount++
}

func (f *Formation) RemoveUnit(u *Unit) {
	for i := int32(0); i < f.UnitCount; i++ {
		if f.Units[i] == u {
			for j := i; j < f.UnitCount-1; j++ {
				f.Units[j] = f.Units[j+1]
			}
			f.UnitCount--
			break
		}
	}
	u.Form = nil
	u.Position = nil
}

func (f *Formation) Order(command int32, radian float64, x, y int32) {
	if f.UnitCount == 0 {
		return
	}
	f.Positions = make([]Position, f.UnitCount)
	space := 20
	columns := int(float64(f.UnitCount)/3.0 + 0.5)
	rows := int(float64(f.UnitCount)/float64(columns) - 0.5)
	ox := (columns - 1) / 2 * space
	oy := rows * space
	sin := math.Sin(radian)
	cos := math.Cos(radian)
	col := 0
	row := 0
	for i := 0; i < int(f.UnitCount); i++ {
		px := ox - col*space
		py := oy - row*space
		xx := -(float64(px)*cos - float64(py)*sin)
		yy := float64(px)*sin + float64(py)*cos

		p := &f.Positions[i]
		p.Form = f
		p.X = x + fixed.Whole(int32(xx))
		p.Y = y + fixed.Whole(int32(yy))
		p.Filled = false
		if col > 0 {
			p.Left = &f.Positions[i-1]
			f.Positions[i-1].Right = p
		}
		if row > 0 {
			p.Front = &f.Positions[i-col]
			f.Positions[i-col].Back = p
		}
		u := f.Units[i]
		u.AnimMod = rand.Intn(8)
		u.Command = command
		u.MoveFinalX = p.X
		u.MoveFinalY = p.Y
		u.Form = f
		u.Position = p

		col++
		if col == columns {
			col = 0
			row++
			if int(f.UnitCount)-i-1 < columns {
				col = (columns - int(f.UnitCount) + i + 1) / 2
			}
		}
	}
}

func SwapPositions(a, b *Unit) {
	a.Position.Filled = false
	b.Position.Filled = false
	temp := a.Position
	a.Position = b.Position
	a.MoveX = a.Position.X
	a.MoveY = a.Position.Y
	b.Position = temp
	b.MoveX = b.Position.X
	b.MoveY = b.Position.Y
}

func FormationUpdate(a, b *Unit) bool {
	if a.Position == nil || b.Position == nil || a.Position.Form != b.Position.Form {
		return false
	}
	if a.Path != nil || b.Path != nil {
		return false
	}
	if a.Position.Filled && b.Position.Filled {
		return false
	}
	abspx := a.Position.X - b.Position.X
	if abspx < 0 {
		abspx = -abspx
	}
	abspy := a.Position.Y - b.Position.Y
	if abspy < 0 {
		abspy = -abspy
	}
	if abspx > abspy {
		if a.Position.X < b.Position.X {
			if a.X > b.X {
				SwapPositions(a, b)
			}
		} else {
			if a.X < b.X {
				SwapPositions(a, b)
			}
		}
	} else {
		if a.Position.Y < b.Position.Y {
			if a.Y > b.Y {
				SwapPositions(a, b)
			}
		} else {
			if a.Y < b.Y {
				SwapPositions(a, b)
			}
		}
	}
	return true
}

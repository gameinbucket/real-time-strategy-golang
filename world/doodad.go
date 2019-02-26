package world

import (
	"../fixed"
	"../graphics"
)

type Doodad struct {
	Mirror    bool
	Anim      []*graphics.Sprite
	AnimMod   int
	AnimFrame int
	X         int32
	Y         int32
	RenderX   float32
	RenderY   float32
	GX        int32
	GY        int32
}

func (m *World) NewDoodad(anim []*graphics.Sprite, x, y int32) *Doodad {
	d := new(Doodad)

	d.Anim = anim
	d.Mirror = false
	d.AnimMod = 0
	d.AnimFrame = 0

	for x < 0 {
		x += m.PixelsFixed
	}

	for y < 0 {
		y += m.PixelsFixed
	}

	for x >= m.PixelsFixed {
		x -= m.PixelsFixed
	}

	for y >= m.PixelsFixed {
		y -= m.PixelsFixed
	}

	d.X = x
	d.Y = y

	d.GX = fixed.Integer(x) >> m.Shift
	d.GY = fixed.Integer(y) >> m.Shift

	m.GetTile(d.GX, d.GY).AddDoodad(d)

	return d
}

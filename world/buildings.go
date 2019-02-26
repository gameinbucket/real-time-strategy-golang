package world

import (
	"../res"
)

func (m *World) MakeHouse(x, y int32) {
	w := int32(9)
	h := int32(9)
	for xx := int32(0); xx < w; xx++ {
		for yy := int32(0); yy < h; yy++ {
			t := m.GetTile(x+xx, y+yy)
			if yy == h-2 {
				t.SpriteID = dat.CavernWallEdge
				t.Closed = true
			} else if yy == h-1 {
				t.SpriteID = dat.CavernWallCorner
				t.Closed = true
			} else if xx == 0 || xx == w-1 {
				t.SpriteID = dat.CavernWall
				t.Closed = true
			} else if yy == 0 {
				if xx > 0 && xx < w-1 {
					t.SpriteID = dat.CavernWallEdge
					t.Closed = true
				} else {
					t.SpriteID = dat.CavernWall
					t.Closed = true
				}
			} else if yy == 1 && xx > 0 && xx < w-1 {
				t.SpriteID = dat.CavernWallCorner
				t.Closed = true
			} else {
				t.SpriteID = dat.CavernStoneFloor
			}
		}
	}
	t := m.GetTile(x+w/2, y+h-1)
	t.SpriteID = dat.CavernStoneFloor
	t.Closed = false
	t = m.GetTile(x+w/2, y+h-2)
	t.SpriteID = dat.CavernStoneFloor
	t.Closed = false
}

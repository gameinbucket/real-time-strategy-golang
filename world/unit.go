package world

import (
	"sync"

	"../fixed"
	"../graphics"
)

const (
	ANIMATION_RATE = 8

	ANIMATION_NOT_DONE    = 0
	ANIMATION_DONE        = 1
	ANIMATION_ALMOST_DONE = 2

	STATUS_IDLE        = 0
	STATUS_CHASE       = 1
	STATUS_MELEE       = 2
	STATUS_MISSILE     = 3
	STATUS_DEAD        = 4
	STATUS_MOVE        = 5
	STATUS_ATTACK_MOVE = 6
	STATUS_DOODAD      = 7
	STATUS_STEP_ASIDE  = 8
)

type Unit struct {
	KingFor        *King
	Command        int32
	Status         int32
	Radius         int32
	Speed          int32
	TargetUnit     *Unit
	TargetDoodad   *Doodad
	AttackTime     int32
	AttackCooldown int32
	Range          int32
	Sight          int32
	Health         int32
	Form           *Formation
	Position       *Position
	Holding        int32
	HoldingCount   int32
	Mirror         bool
	Direction      int32
	SpriteID       uint32
	AnimWalk       [][]*graphics.Sprite
	AnimAttack     [][]*graphics.Sprite
	AnimDeath      [][]*graphics.Sprite
	Anim           [][]*graphics.Sprite
	AnimMod        int
	AnimFrame      int
	X              int32
	Y              int32
	RenderX        float32
	RenderY        float32
	GX             int32
	GY             int32
	DX             int32
	DY             int32
	MinGX          int32
	MinGY          int32
	MaxGX          int32
	MaxGY          int32
	MoveX          int32
	MoveY          int32
	MoveFinalX     int32
	MoveFinalY     int32
	Path           *PathList
}

func (m *World) NewUnit(spriteID uint32, walk, attack, death [][]*graphics.Sprite, k *King, x, y, radius int32) *Unit {
	u := new(Unit)

	u.KingFor = k
	k.AddUnit(u)

	u.Status = STATUS_IDLE
	u.SpriteID = spriteID
	u.AnimWalk = walk
	u.AnimAttack = attack
	u.AnimDeath = death
	u.Mirror = false
	u.Anim = walk
	u.AnimMod = 0
	u.AnimFrame = 0
	u.AttackTime = 32

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

	u.X = x
	u.Y = y

	u.MoveX = x
	u.MoveY = y

	u.GX = fixed.Integer(x) >> m.Shift
	u.GY = fixed.Integer(y) >> m.Shift

	m.GetTile(u.GX, u.GY).AddUnit(u)

	u.Health = 2
	u.Radius = radius
	u.Speed = fixed.Whole(1)
	u.Range = fixed.Whole(8)
	u.Sight = 10

	u.MinGX = fixed.Integer(x-radius) >> m.Shift
	u.MinGY = fixed.Integer(y-radius) >> m.Shift

	u.MaxGX = fixed.Integer(x+radius) >> m.Shift
	u.MaxGY = fixed.Integer(y+radius) >> m.Shift

	for gx := u.MinGX; gx <= u.MaxGX; gx++ {
		for gy := u.MinGY; gy <= u.MaxGY; gy++ {
			m.GetTile(gx, gy).AddPhysical(u)
		}
	}

	return u
}

func (u *Unit) Missile() {

}

func (u *Unit) Melee(m *World) {
	u.TargetUnit.Damage(m, 1)
}

func (u *Unit) Damage(m *World, d int32) {
	u.Health -= d
	if u.Status != STATUS_DEAD {
		if u.Health < 1 {
			u.Status = STATUS_DEAD
			u.AnimFrame = 0
			u.AnimMod = 0
			u.Anim = u.AnimDeath
		}
	}
}

func (u *Unit) Animation() int32 {
	u.AnimMod++
	if u.AnimMod == ANIMATION_RATE {
		u.AnimMod = 0
		u.AnimFrame++
		if u.AnimFrame == len(u.Anim[0])-1 {
			return ANIMATION_ALMOST_DONE
		}
		if u.AnimFrame == len(u.Anim[0]) {
			return ANIMATION_DONE
		}
	}
	return ANIMATION_NOT_DONE
}

func (u *Unit) Vision(wg *sync.WaitGroup, m *World) {
	var nearest *Unit
	nearX := m.MapSize
	nearY := m.MapSize
	for x := -u.Sight; x <= u.Sight; x++ {
		for y := -u.Sight; y <= u.Sight; y++ {
			ax := x
			if ax < 0 {
				ax = -ax
			}
			ay := y
			if ay < 0 {
				ay = -ay
			}
			pos := ax + ay
			if pos < nearX+nearY {
				t := m.GetTile(u.GX+x, u.GY+y)
				for i := int32(0); i < t.UnitCount; i++ {
					o := t.Units[i]
					if o.KingFor != u.KingFor {
						nearest = o
						nearX = ax
						nearY = ay
					}
				}
			}
		}
	}
	u.TargetUnit = nearest
	wg.Done()
}

func delta(x, y int32, m *World) (int32, int32) {
	if x < 0 {
		if -x > m.HalfPixelsFixed {
			x += m.PixelsFixed
		}
	} else {
		if x > m.HalfPixelsFixed {
			x -= m.PixelsFixed
		}
	}
	if y < 0 {
		if -y > m.HalfPixelsFixed {
			y += m.PixelsFixed
		}
	} else {
		if y > m.HalfPixelsFixed {
			y -= m.PixelsFixed
		}
	}
	return x, y
}

func Approx(x, y int32) int32 {
	if x < 0 {
		if y < 0 {
			if -x > -y {
				return (-x - y - (-y >> 1))
			} else {
				return (-x - y - (-x >> 1))
			}
		} else {
			if -x > y {
				return (-x + y - (y >> 1))
			} else {
				return (-x + y - (-x >> 1))
			}
		}
	} else {
		if y < 0 {
			if x > -y {
				return (x - y - (-y >> 1))
			} else {
				return (x - y - (x >> 1))
			}
		} else {
			if x > y {
				return (x + y - (y >> 1))
			} else {
				return (x + y - (x >> 1))
			}
		}
	}
}

func (u *Unit) Integrate(m *World) {
	if u.Status == STATUS_DEAD {
		if u.AnimFrame < len(u.Anim[0])-1 {
			u.Animation()
		} else {
			m.AddDeadUnit(u)
		}
		return
	}
	u.X += u.DX
	u.Y += u.DY
	u.DX = 0
	u.DY = 0
	if u.Status == STATUS_MELEE {
		anim := u.Animation()
		if anim == ANIMATION_ALMOST_DONE {
			if u.TargetUnit != nil {
				u.Melee(m)
			}
		} else if anim == ANIMATION_DONE {
			u.AnimFrame = 0
			u.Anim = u.AnimWalk
			u.Status = STATUS_CHASE
		}
		return
	}
	if u.AttackCooldown > 0 {
		u.AttackCooldown--
	}
	if u.Status == STATUS_IDLE {
		if u.Command == STATUS_MOVE {
			u.Status = STATUS_MOVE
		} else if u.Command == STATUS_DOODAD {
			u.Status = STATUS_DOODAD
		} else if u.TargetUnit != nil {
			u.Status = STATUS_CHASE
		}
	}
	if u.Status == STATUS_CHASE {
		if u.TargetUnit == nil || u.TargetUnit.Health < 1 {
			u.Status = STATUS_IDLE
			u.AnimFrame = 0
			u.AnimMod = 0
		} else {
			dx, dy := delta(u.TargetUnit.X-u.X, u.TargetUnit.Y-u.Y, m)
			if u.AttackCooldown == 0 && Approx(dx, dy) <= u.Radius+u.Range+u.TargetUnit.Radius {
				u.Status = STATUS_MELEE
				u.AttackCooldown = u.AttackTime
				u.Anim = u.AnimAttack
				u.AnimFrame = 0
				u.AnimMod = 0
			} else {
				u.MoveTowards(m, dx, dy)
			}
		}
	} else if u.Status == STATUS_DOODAD {
		dx, dy := delta(u.TargetDoodad.X-u.X, u.TargetDoodad.Y-u.Y, m)
		if Approx(dx, dy) <= u.Radius {
			u.Status = STATUS_IDLE
			u.Command = STATUS_IDLE
			u.AnimFrame = 0
			u.AnimMod = 0
			u.Holding = 1
		} else {
			u.MoveTowards(m, dx, dy)
		}
	} else if u.Status == STATUS_STEP_ASIDE {
		dx, dy := delta(u.MoveX-u.X, u.MoveY-u.Y, m)
		if Approx(dx, dy) <= (u.Radius >> 1) {
			u.Status = STATUS_IDLE
			u.Command = STATUS_IDLE
			u.AnimFrame = 0
			u.AnimMod = 0
		} else {
			u.MoveTowards(m, dx, dy)
		}
	} else {
		dx, dy := delta(u.MoveX-u.X, u.MoveY-u.Y, m)
		if Approx(dx, dy) <= (u.Radius >> 1) {
			if u.Path != nil {
				u.Path = u.Path.Next
				if u.Path != nil {
					u.MoveX = fixed.Whole((u.Path.X << m.Shift) + (m.TilePixels >> 1))
					u.MoveY = fixed.Whole((u.Path.Y << m.Shift) + (m.TilePixels >> 1))
				} else if u.MoveX != u.MoveFinalX || u.MoveY != u.MoveFinalY {
					u.MoveX = u.MoveFinalX
					u.MoveY = u.MoveFinalY
				} else {
					if u.Position != nil {
						u.Position.Filled = true
					}
					u.Status = STATUS_IDLE
					u.Command = STATUS_IDLE
					u.AnimFrame = 0
					u.AnimMod = 0
				}
			} else {
				if u.Position != nil {
					u.Position.Filled = true
				}
				u.Status = STATUS_IDLE
				u.Command = STATUS_IDLE
				u.AnimFrame = 0
				u.AnimMod = 0
			}
		} else {
			u.MoveTowards(m, dx, dy)
		}
	}
}

func (u *Unit) MoveTowards(m *World, dx, dy int32) {
	if u.Animation() == ANIMATION_DONE {
		u.AnimFrame = 0
	}
	radian := fixed.Atan2(dy, dx)
	degree := fixed.Floating(radian) * 57.2958
	if degree < 0 {
		degree += 360
	}
	if degree > 337.5 {
		u.Direction = 2
		u.Mirror = false
	} else if degree > 292.5 {
		u.Direction = 1
		u.Mirror = false
	} else if degree > 247.5 {
		u.Direction = 0
		u.Mirror = false
	} else if degree > 202.5 {
		u.Direction = 1
		u.Mirror = true
	} else if degree > 157.5 {
		u.Direction = 2
		u.Mirror = true
	} else if degree > 112.5 {
		u.Direction = 3
		u.Mirror = true
	} else if degree > 67.5 {
		u.Direction = 4
		u.Mirror = false
	} else if degree > 22.5 {
		u.Direction = 3
		u.Mirror = false
	} else {
		u.Direction = 2
		u.Mirror = false
	}
	u.Move(m, fixed.Mul(fixed.Cos(radian), u.Speed), fixed.Mul(fixed.Sin(radian), u.Speed))
}

func (u *Unit) Move(m *World, dx, dy int32) {
	u.X += dx
	u.Y += dy

	gxMin := fixed.Integer(u.X-u.Radius) >> m.Shift
	gyMin := fixed.Integer(u.Y-u.Radius) >> m.Shift
	gxMax := fixed.Integer(u.X+u.Radius) >> m.Shift
	gyMax := fixed.Integer(u.Y+u.Radius) >> m.Shift
	for x := gxMin; x <= gxMax; x++ {
		for y := gyMin; y <= gyMax; y++ {
			if m.GetTile(x, y).Closed {
				xx := fixed.Whole(x << m.Shift)
				yy := fixed.Whole(y << m.Shift)
				closeX := u.X
				if closeX < xx {
					closeX = xx
				} else if closeX > xx+m.FixedTilePixels {
					closeX = xx + m.FixedTilePixels
				}
				closeY := u.Y
				if closeY < yy {
					closeY = yy
				} else if closeY > yy+m.FixedTilePixels {
					closeY = yy + m.FixedTilePixels
				}
				dxx := u.X - closeX
				dyy := u.Y - closeY
				dist := fixed.Mul(dxx, dxx) + fixed.Mul(dyy, dyy)
				if dist > fixed.Mul(u.Radius, u.Radius) {
					continue
				}
				dist = fixed.Sqrt(dist)
				if dist == 0 {
					dist = 1
				}
				mult := fixed.Div(u.Radius, dist)
				u.X += fixed.Mul(dxx, mult) - dxx
				u.Y += fixed.Mul(dyy, mult) - dyy
			}
		}
	}

	for u.X < 0 {
		u.X += m.PixelsFixed
	}
	for u.Y < 0 {
		u.Y += m.PixelsFixed
	}
	for u.X >= m.PixelsFixed {
		u.X -= m.PixelsFixed
	}
	for u.Y >= m.PixelsFixed {
		u.Y -= m.PixelsFixed
	}
	gx := fixed.Integer(u.X) >> m.Shift
	gy := fixed.Integer(u.Y) >> m.Shift
	if gx != u.GX || gy != u.GY {
		old := m.GetTile(u.GX, u.GY)
		old.RemoveUnit(u)
		next := m.GetTile(gx, gy)
		next.AddUnit(u)
		u.GX = gx
		u.GY = gy
	}

	gxMin = fixed.Integer(u.X-u.Radius) >> m.Shift
	gyMin = fixed.Integer(u.Y-u.Radius) >> m.Shift
	gxMax = fixed.Integer(u.X+u.Radius) >> m.Shift
	gyMax = fixed.Integer(u.Y+u.Radius) >> m.Shift
	if gxMin != u.MinGX || gyMin != u.MinGY || gxMax != u.MaxGX || gyMax != u.MaxGY {
		for gx = u.MinGX; gx <= u.MaxGX; gx++ {
			for gy = u.MinGY; gy <= u.MaxGY; gy++ {
				m.GetTile(gx, gy).RemovePhysical(u)
			}
		}
		for gx = gxMin; gx <= gxMax; gx++ {
			for gy = gyMin; gy <= gyMax; gy++ {
				m.GetTile(gx, gy).AddPhysical(u)
			}
		}
		u.MinGX = gxMin
		u.MinGY = gyMin
		u.MaxGX = gxMax
		u.MaxGY = gyMax
	}
}

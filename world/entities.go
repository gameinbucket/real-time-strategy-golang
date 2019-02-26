package world

type Entity struct {
	Type    int32
	Radius  int32
	X       int32
	Y       int32
	RenderX float32
	RenderY float32
	GX      int32
	GY      int32
	MinGX   int32
	MinGY   int32
	MaxGX   int32
	MaxGY   int32
}

/*
type Entity interface {
	area() float64
	}

type Unit struct {
	Entity
	KingFor        *King
	Command        int32
	Status         int32
	Speed          int32
	TargetUnit     *Unit
	TargetDoodad   *Doodad
	AttackTime     int32
	AttackCooldown int32
	Range          int32
	Sight          int32
	Health         int32
	Form           *Formation
	Holding        []*Item
	Mirror         bool
	Direction      int32
	SpriteID       uint32
	AnimWalk       [][]*graphics.Sprite
	AnimAttack     [][]*graphics.Sprite
	AnimDeath      [][]*graphics.Sprite
	Anim           [][]*graphics.Sprite
	AnimMod        int
	AnimFrame      int
	MoveX          int32
	MoveY          int32
}

type Doodad struct {
	Entity
	Sprite *graphics.Sprite
}

type Item struct {
	Entity
	Sprite *graphics.Sprite
}
*/

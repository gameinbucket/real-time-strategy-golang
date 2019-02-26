package main

import (
	"runtime"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"./fixed"
	"./graphics"
	"./res"
	"./world"
)

const (
	EventNone     = 0
	EventWorld    = 1
	EventMiniMap  = 2
	EventBuild    = 3
	DEBUG_TIMING  = false
	DEBUG_PATHING = false
)

var (
	canvasWidth       = 11 * 75
	canvasHeight      = 7 * 75
	window            *glfw.Window
	leftClick         = false
	leftUnset         = false
	rightClick        = false
	rightUnset        = false
	programTriangle   uint32
	programPreColor   uint32
	programPreTexture uint32
	myOrthographic    = make([]float32, 16)
	unitBuffer        map[uint32]*graphics.Buffer
	mapBuffer         *graphics.Buffer
	colorBuffer       *graphics.Buffer
	lineBuffer        *graphics.Buffer
	imgMan            uint32
	imgFootman        uint32
	imgCavern         uint32
	manWalk           [][]*graphics.Sprite
	footmanWalk       [][]*graphics.Sprite
	footmanAttack     [][]*graphics.Sprite
	footmanDeath      [][]*graphics.Sprite
	doodadPot         []*graphics.Sprite
	spriteCavern      []*graphics.Sprite
	viewX             = int32(0)
	viewY             = int32(0)
	viewGX            = int32(0)
	viewGY            = int32(0)
	viewGW            = int32(0)
	viewGH            = int32(0)
	viewPadding       = int32(2)
	unitsSelected     = false
	originSelectionX  int32
	originSelectionY  int32
	selectionLeft     int32
	selectionRight    int32
	selectionTop      int32
	selectionBottom   int32
	orderingMove      bool
	orderingMoveX     int32
	orderingMoveY     int32
	orderingMoveX2    int32
	orderingMoveY2    int32
	EventState        = EventNone
	mWorld            *world.World
	player            *world.King
)

func Resize(w *glfw.Window, width, height int) {
	canvasWidth = width
	canvasHeight = height

	graphics.Orthographic(myOrthographic, 0.0, 0.0, float32(canvasWidth), float32(canvasHeight), 0, 1)

	viewGW = int32(canvasWidth) >> mWorld.Shift
	viewGH = int32(canvasHeight) >> mWorld.Shift
}

func main() {
	runtime.LockOSThread()

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	defer glfw.Terminate()

	/* window */

	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Visible, glfw.False)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	windowLocal, err := glfw.CreateWindow(canvasWidth, canvasHeight, "Grand Campaign", nil, nil)
	if err != nil {
		panic(err)
	}

	window = windowLocal

	videoMode := glfw.GetPrimaryMonitor().GetVideoMode()
	window.SetPos((videoMode.Width-canvasWidth)/2, (videoMode.Height-canvasHeight)/2)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	window.SetSizeCallback(Resize)

	/* gl */

	if err := gl.Init(); err != nil {
		panic(err)
	}

	graphics.SetClearColor(0.0, 0.0, 0.0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	/* graphics */

	programTriangle = graphics.MakeProgram("tri", 'r')
	programPreColor = graphics.MakeProgram("pre-color", 'r')
	programPreTexture = graphics.MakeProgram("pre-texture", 'r')

	colorBuffer = graphics.NewBuffer(2, 4, 0, 0, 50000, 50000)
	lineBuffer = graphics.NewBuffer(2, 4, 0, 0, 1000, 1000)
	mapBuffer = graphics.NewBuffer(2, 0, 2, 0, 50000, 50000)

	imgMan = graphics.MakeTexture("man.png", true, false)
	imgFootman = graphics.MakeTexture("footman.png", true, false)
	imgCavern = graphics.MakeTexture("caverns.png", true, false)

	unitBuffer = make(map[uint32]*graphics.Buffer)
	unitBuffer[imgMan] = graphics.NewBuffer(2, 0, 2, 0, 200, 200)
	unitBuffer[imgFootman] = graphics.NewBuffer(2, 0, 2, 0, 4000, 4000)

	var w float32
	var h float32
	var s float32
	var d float32
	var z float32
	var t float32

	w = 1.0 / 512.0
	h = 1.0 / 512.0
	d = 32.0
	z = 0.0
	t = 16.0
	manWalk = make([][]*graphics.Sprite, 5)
	manWalk[0] = make([]*graphics.Sprite, 8)
	manWalk[1] = make([]*graphics.Sprite, 8)
	manWalk[2] = make([]*graphics.Sprite, 8)
	manWalk[3] = make([]*graphics.Sprite, 8)
	manWalk[4] = make([]*graphics.Sprite, 8)
	for i := 0; i < 5; i++ {
		manWalk[i][0] = graphics.NewSprite(float32(i)*d, 0, d, d, w, h, z, t)
		manWalk[i][1] = graphics.NewSprite(float32(i)*d, 4.0*d, d, d, w, h, z, t)
		manWalk[i][2] = graphics.NewSprite(float32(i)*d, 3.0*d, d, d, w, h, z, t)
		manWalk[i][3] = manWalk[i][1]
		manWalk[i][4] = manWalk[i][0]
		manWalk[i][5] = graphics.NewSprite(float32(i)*d, 0.0*d, d, d, w, h, z, t)
		manWalk[i][6] = graphics.NewSprite(float32(i)*d, 1.0*d, d, d, w, h, z, t)
		manWalk[i][7] = manWalk[i][5]
	}

	w = 1.0 / 1024.0
	h = 1.0 / 256.0
	d = 48.0
	z = 0.0
	t = 24.0
	footmanWalk = make([][]*graphics.Sprite, 5)
	footmanWalk[0] = make([]*graphics.Sprite, 8)
	footmanWalk[1] = make([]*graphics.Sprite, 8)
	footmanWalk[2] = make([]*graphics.Sprite, 8)
	footmanWalk[3] = make([]*graphics.Sprite, 8)
	footmanWalk[4] = make([]*graphics.Sprite, 8)
	for i := 0; i < 5; i++ {
		footmanWalk[i][0] = graphics.NewSprite(float32(i)*d, 0, d, d, w, h, z, t)
		footmanWalk[i][1] = graphics.NewSprite(float32(i)*d, 4.0*d, d, d, w, h, z, t)
		footmanWalk[i][2] = graphics.NewSprite(float32(i)*d, 3.0*d, d, d, w, h, z, t)
		footmanWalk[i][3] = footmanWalk[i][1]
		footmanWalk[i][4] = footmanWalk[i][0]
		footmanWalk[i][5] = graphics.NewSprite(float32(i)*d, 0.0*d, d, d, w, h, z, t)
		footmanWalk[i][6] = graphics.NewSprite(float32(i)*d, 1.0*d, d, d, w, h, z, t)
		footmanWalk[i][7] = footmanWalk[i][5]
	}
	footmanAttack = make([][]*graphics.Sprite, 5)
	footmanAttack[0] = make([]*graphics.Sprite, 5)
	footmanAttack[1] = make([]*graphics.Sprite, 5)
	footmanAttack[2] = make([]*graphics.Sprite, 5)
	footmanAttack[3] = make([]*graphics.Sprite, 5)
	footmanAttack[4] = make([]*graphics.Sprite, 5)
	for i := 0; i < 5; i++ {
		j := i + 5
		footmanAttack[i][0] = graphics.NewSprite(float32(j)*d, 0, d, d, w, h, z, t)
		footmanAttack[i][1] = graphics.NewSprite(float32(j)*d, 4.0*d, d, d, w, h, z, t)
		footmanAttack[i][2] = graphics.NewSprite(float32(j)*d, 3.0*d, d, d, w, h, z, t)
		footmanAttack[i][3] = graphics.NewSprite(float32(j)*d, 0.0*d, d, d, w, h, z, t)
		footmanAttack[i][4] = graphics.NewSprite(float32(j)*d, 1.0*d, d, d, w, h, z, t)
	}
	footmanDeath = make([][]*graphics.Sprite, 5)
	footmanDeath[0] = make([]*graphics.Sprite, 3)
	footmanDeath[1] = make([]*graphics.Sprite, 3)
	footmanDeath[2] = make([]*graphics.Sprite, 3)
	footmanDeath[3] = make([]*graphics.Sprite, 3)
	footmanDeath[4] = make([]*graphics.Sprite, 3)
	for i := 0; i < 5; i++ {
		j := i + 10
		footmanDeath[i][0] = graphics.NewSprite(float32(j)*d, 0, d, d, w, h, z, t)
		footmanDeath[i][1] = graphics.NewSprite(float32(j)*d, 1.0*d, d, d, w, h, z, t)
		footmanDeath[i][2] = graphics.NewSprite(float32(j)*d, 2.0*d, d, d, w, h, z, t)
	}

	s = 16.0
	w = 1.0 / 256.0
	h = 1.0 / 128.0
	spriteCavern = make([]*graphics.Sprite, 7)
	spriteCavern[dat.CavernDirt] = graphics.NewSprite(1+17*0, 1+17*0, s, s, w, h, 0, 0)
	spriteCavern[dat.CavernDirtLight] = graphics.NewSprite(1+17*0, 1+17*1, s, s, w, h, 0, 0)
	spriteCavern[dat.CavernDirtLightest] = graphics.NewSprite(1+17*0, 1+17*2, s, s, w, h, 0, 0)
	spriteCavern[dat.CavernWall] = graphics.NewSprite(1+17*1, 1+17*0, s, s, w, h, 0, 0)
	spriteCavern[dat.CavernWallEdge] = graphics.NewSprite(1+17*1, 1+17*1, s, s, w, h, 0, 0)
	spriteCavern[dat.CavernWallCorner] = graphics.NewSprite(1+17*1, 1+17*2, s, s, w, h, 0, 0)
	spriteCavern[dat.CavernStoneFloor] = graphics.NewSprite(1+17*1, 1+17*3, s, s, w, h, 0, 0)

	doodadPot = make([]*graphics.Sprite, 1)
	doodadPot[0] = graphics.NewSprite(103, 34, 14, 15, w, h, 0, 0)

	/* game */

	mWorld = world.NewWorld(64, 4)
	mWorld.MakeHouse(0, 0)
	mWorld.MakeHouse(32, 16)
	mWorld.GetTile(21, 4).SpriteID = dat.CavernDirtLight
	mWorld.GetTile(21, 5).SpriteID = dat.CavernDirtLight
	mWorld.GetTile(21, 6).SpriteID = dat.CavernDirtLight
	mWorld.GetTile(22, 4).SpriteID = dat.CavernDirtLight
	mWorld.GetTile(22, 5).SpriteID = dat.CavernDirtLightest
	mWorld.GetTile(22, 6).SpriteID = dat.CavernDirtLight
	mWorld.GetTile(23, 4).SpriteID = dat.CavernDirtLight
	mWorld.GetTile(23, 5).SpriteID = dat.CavernDirtLight
	mWorld.GetTile(23, 6).SpriteID = dat.CavernDirtLight

	mWorld.ComputePathMesh()

	mWorld.Kings = make([]*world.King, 2)
	mWorld.Kings[0] = mWorld.NewKing(0xff0000)
	mWorld.Kings[1] = mWorld.NewKing(0x00ff00)
	player = mWorld.Kings[0]
	for x := int32(20); x < 32; x++ {
		for y := int32(1); y < 2; y++ {
			player.Formations[0].AddUnit(mWorld.NewUnit(imgFootman, footmanWalk, footmanAttack, footmanDeath, player, fixed.Whole(16*x), fixed.Whole(16*y), fixed.Whole(8)))
		}
	}
	for x := int32(20); x < 20; x++ {
		player.Formations[0].AddUnit(mWorld.NewUnit(imgMan, manWalk, manWalk, manWalk, player, fixed.Whole(16*13), fixed.Whole(16*0), fixed.Whole(8)))
	}
	enemy := mWorld.Kings[1]
	for x := int32(20); x < 20; x++ {
		for y := int32(12); y < 12; y++ {
			mWorld.NewUnit(imgFootman, footmanWalk, footmanAttack, footmanDeath, enemy, fixed.Whole(16*x), fixed.Whole(16*y), fixed.Whole(8))
		}
	}

	mWorld.NewDoodad(doodadPot, fixed.Whole(16*3), fixed.Whole(16*3))

	/* run */

	Resize(window, canvasWidth, canvasHeight)

	window.Show()

	for !window.ShouldClose() {
		glfw.PollEvents()

		if window.GetKey(glfw.KeyEscape) == glfw.Press {
			break
		}

		Event()
		Turn()
		Draw()

		window.SwapBuffers()
	}
}

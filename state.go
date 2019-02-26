package main

import (
	"fmt"
	"math"

	"./fixed"
	"./graphics"
	"./world"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func Event() {
	if !unitsSelected {
		dx := int32(0)
		dy := int32(0)
		const viewSpeed = int32(4)
		if window.GetKey(glfw.KeyW) == glfw.Press || window.GetKey(glfw.KeyUp) == glfw.Press {
			dy -= viewSpeed
		}
		if window.GetKey(glfw.KeyS) == glfw.Press || window.GetKey(glfw.KeyDown) == glfw.Press {
			dy += viewSpeed
		}
		if window.GetKey(glfw.KeyA) == glfw.Press || window.GetKey(glfw.KeyLeft) == glfw.Press {
			dx -= viewSpeed
		}
		if window.GetKey(glfw.KeyD) == glfw.Press || window.GetKey(glfw.KeyRight) == glfw.Press {
			dx += viewSpeed
		}

		if dx != 0 || dy != 0 {
			viewX += dx
			viewY += dy

			for viewX < 0 {
				viewX += mWorld.Pixels
			}

			for viewX > mWorld.Pixels {
				viewX -= mWorld.Pixels
			}

			for viewY < 0 {
				viewY += mWorld.Pixels
			}

			for viewY > mWorld.Pixels {
				viewY -= mWorld.Pixels
			}

			viewGX = viewX >> mWorld.Shift
			viewGY = viewY >> mWorld.Shift
		}
	}

	FcurX, FcurY := window.GetCursorPos()
	curX := int32(FcurX)
	curY := int32(FcurY)

	leftClick = false
	if window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		if leftUnset {
			leftClick = true
			leftUnset = false
		}
	} else {
		leftUnset = true
	}

	rightClick = false
	if window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press {
		if rightUnset {
			rightClick = true
			rightUnset = false
		}
	} else {
		rightUnset = true
	}

	if window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
		x := viewX + curX
		y := viewY + curY
		if unitsSelected {
			if x < originSelectionX {
				selectionLeft = x
				selectionRight = originSelectionX
			} else {
				selectionLeft = originSelectionX
				selectionRight = x
			}
			if y < originSelectionY {
				selectionTop = y
				selectionBottom = originSelectionY
			} else {
				selectionTop = originSelectionY
				selectionBottom = y
			}

			x = selectionLeft >> mWorld.Shift
			y = selectionTop >> mWorld.Shift

			if x >= mWorld.MapSize {
				x -= mWorld.MapSize
			}
			if y >= mWorld.MapSize {
				y -= mWorld.MapSize
			}

			w := (selectionRight - selectionLeft) >> mWorld.Shift
			h := (selectionBottom - selectionTop) >> mWorld.Shift

			form := player.Formations[0]
			form.Clear()

			for gx := x - 1; gx <= x+w+1; gx++ {
				col := gx
				for col < 0 {
					col += mWorld.MapSize
				}
				for col >= mWorld.MapSize {
					col -= mWorld.MapSize
				}
				for gy := y - 1; gy <= y+h+1; gy++ {
					row := gy
					for row < 0 {
						row += mWorld.MapSize
					}
					for row >= mWorld.MapSize {
						row -= mWorld.MapSize
					}
					tile := mWorld.Tiles[col+row*mWorld.MapSize]
					for i := int32(0); i < tile.UnitCount; i++ {
						u := tile.Units[i]
						if player != u.KingFor {
							continue
						}
						r := fixed.Floating(u.Radius)
						ux := u.RenderX
						uy := u.RenderY
						if float32(viewX) > ux {
							ux += float32(mWorld.Pixels)
						}
						if float32(viewY) > uy {
							uy += float32(mWorld.Pixels)
						}
						if selectionLeft <= int32(ux+r) &&
							selectionRight >= int32(ux-r) &&
							selectionTop <= int32(uy+r) &&
							selectionBottom >= int32(uy-r) {
							form.AddUnit(u)
						}
					}
				}
			}
		} else {
			unitsSelected = true
			selectionLeft = x
			selectionRight = x
			selectionTop = y
			selectionBottom = y
			originSelectionX = x
			originSelectionY = y
		}
	} else if unitsSelected {
		unitsSelected = false
	}

	if orderingMove {
		if window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press {
			x := viewX + curX
			y := viewY + curY
			orderingMoveX2 = x
			orderingMoveY2 = y
		} else {
			xm := orderingMoveX
			ym := orderingMoveY
			if xm >= mWorld.Pixels {
				xm -= mWorld.Pixels
			}
			if ym >= mWorld.Pixels {
				ym -= mWorld.Pixels
			}

			player.MoveOrder(math.Atan2(float64(orderingMoveX2-orderingMoveX), float64(orderingMoveY2-orderingMoveY)), fixed.Whole(xm), fixed.Whole(ym))
			orderingMove = false

			for i := int32(0); i < player.Formations[0].UnitCount; i++ {
				u := player.Formations[0].Units[i]
				u.Path = mWorld.FindPath(u.GX, u.GY, fixed.Integer(u.MoveFinalX)>>mWorld.Shift, fixed.Integer(u.MoveFinalY)>>mWorld.Shift)
				if u.Path != nil {
					u.MoveX = fixed.Whole((u.Path.X << mWorld.Shift) + (mWorld.TilePixels >> 1))
					u.MoveY = fixed.Whole((u.Path.Y << mWorld.Shift) + (mWorld.TilePixels >> 1))
				} else {
					u.MoveX = u.MoveFinalX
					u.MoveY = u.MoveFinalY
				}
			}
		}
	}

	if rightClick {
		x := viewX + curX
		y := viewY + curY
		xm := x
		ym := y
		if xm >= mWorld.Pixels {
			xm -= mWorld.Pixels
		}
		if ym >= mWorld.Pixels {
			ym -= mWorld.Pixels
		}
		t := mWorld.GetTile(xm>>mWorld.Shift, ym>>mWorld.Shift)
		if t.DoodadCount > 0 {
			player.DoodadOrder(t.Doodads[0])
		} else if player.Formations[0].UnitCount > 0 {
			orderingMove = true
			orderingMoveX = x
			orderingMoveY = y
			orderingMoveX2 = x
			orderingMoveY2 = y
		}
	}
}

func Turn() {
	mWorld.Integrate()
}

func Draw() {
	graphics.SetFramebuffer(0)
	graphics.SetView(0, 0, int32(canvasWidth), int32(canvasHeight))

	graphics.SetProgram(programPreTexture)
	graphics.SetOrthographic(myOrthographic, float32(-viewX), float32(-viewY))
	graphics.SetMVP()
	mapBuffer.Zero()
	for _, v := range unitBuffer {
		v.Zero()
	}

	timeStart := glfw.GetTime()

	for y := viewGY - viewPadding; y <= viewGY+viewGH+viewPadding; y++ {
		for x := viewGX - viewPadding; x <= viewGX+viewGW+viewPadding; x++ {
			t := mWorld.GetTile(x, y)
			mapBuffer.RenderSprite(spriteCavern[t.SpriteID], float32(x<<mWorld.Shift), float32(y<<mWorld.Shift))
			for i := int32(0); i < t.DoodadCount; i++ {
				d := t.Doodads[i]
				sp := d.Anim[d.AnimFrame]
				d.RenderX = fixed.Floating(d.X)
				d.RenderY = fixed.Floating(d.Y)
				ux := d.RenderX - sp.W/2 + sp.OX
				uy := d.RenderY + sp.OY
				if float32(viewX) > ux {
					ux += float32(mWorld.Pixels)
				}
				if float32(viewY) > uy {
					uy += float32(mWorld.Pixels)
				}
				if d.Mirror {
					mapBuffer.RenderSpriteMirror(sp, ux, uy)
				} else {
					mapBuffer.RenderSprite(sp, ux, uy)
				}
			}
			for i := int32(0); i < t.UnitCount; i++ {
				u := t.Units[i]
				sp := u.Anim[u.Direction][u.AnimFrame]
				u.RenderX = fixed.Floating(u.X)
				u.RenderY = fixed.Floating(u.Y)
				ux := u.RenderX - sp.W/2 + sp.OX
				uy := u.RenderY + sp.OY
				if float32(viewX) > ux {
					ux += float32(mWorld.Pixels)
				}
				if float32(viewY) > uy {
					uy += float32(mWorld.Pixels)
				}
				uy -= sp.H
				if u.Mirror {
					unitBuffer[u.SpriteID].RenderSpriteMirror(sp, ux, uy)
				} else {
					unitBuffer[u.SpriteID].RenderSprite(sp, ux, uy)
				}
			}
		}
	}

	timeEnd := glfw.GetTime()
	if DEBUG_TIMING {
		fmt.Println("Time", (timeEnd-timeStart)*1000)
	}

	graphics.SetTexture0(imgCavern)
	mapBuffer.DrawElements()
	for k, v := range unitBuffer {
		graphics.SetTexture0(k)
		v.DrawElements()
	}

	graphics.SetProgram(programPreColor)
	graphics.SetMVP()

	if DEBUG_PATHING {
		colorBuffer.Zero()
		for y := viewGY - viewPadding; y <= viewGY+viewGH+viewPadding; y++ {
			for x := viewGX - viewPadding; x <= viewGX+viewGW+viewPadding; x++ {
				t := mWorld.GetTile(x, y)
				if t.Nav == nil {
					continue
				}
				if t.IsEdge {
					colorBuffer.RenderRectangle(
						float32(x<<mWorld.Shift), float32(y<<mWorld.Shift),
						float32(mWorld.TilePixels), float32(mWorld.TilePixels),
						1.0, 0.0, 0.0, 0.4)
				} else {
					colorBuffer.RenderRectangle(
						float32(x<<mWorld.Shift), float32(y<<mWorld.Shift),
						float32(mWorld.TilePixels), float32(mWorld.TilePixels),
						t.Nav.Red, t.Nav.Green, t.Nav.Blue, 0.6)
				}
			}
		}
		colorBuffer.RenderRectangle(
			float32(world.DEBUG_FROM_X<<mWorld.Shift), float32(world.DEBUG_FROM_Y<<mWorld.Shift),
			float32(mWorld.TilePixels), float32(mWorld.TilePixels),
			0.0, 1.0, 0.0, 0.8)
		colorBuffer.RenderRectangle(
			float32(world.DEBUG_TO_X<<mWorld.Shift), float32(world.DEBUG_TO_Y<<mWorld.Shift),
			float32(mWorld.TilePixels), float32(mWorld.TilePixels),
			1.0, 0.0, 0.0, 0.8)
		path := world.DEBUG_PATH_LIST
		for path != nil {
			colorBuffer.RenderRectangle(
				float32(path.X<<mWorld.Shift), float32(path.Y<<mWorld.Shift),
				float32(mWorld.TilePixels), float32(mWorld.TilePixels),
				0.0, 0.0, 1.0, 0.8)
			path = path.Next
		}
		path = world.DEBUG_RAY_TRACE
		for path != nil {
			colorBuffer.RenderRectangle(
				float32(path.X<<mWorld.Shift), float32(path.Y<<mWorld.Shift),
				float32(mWorld.TilePixels), float32(mWorld.TilePixels),
				1.0, 0.0, 1.0, 0.8)
			path = path.Next
		}
		colorBuffer.DrawElements()
	}

	lineBuffer.Zero()
	colorBuffer.Zero()
	form := player.Formations[0]
	for i := int32(0); i < form.UnitCount; i++ {
		u := form.Units[i]
		ux := u.RenderX
		uy := u.RenderY - fixed.Floating(u.Radius)
		if float32(viewX) > ux {
			ux += float32(mWorld.Pixels)
		}
		if float32(viewY) > uy {
			uy += float32(mWorld.Pixels)
		}
		lineBuffer.RenderLine(float32(ux-4), float32(uy-2), float32(ux+4), float32(uy-2), 0.0, 1.0, 0.0)
		lineBuffer.RenderLine(float32(ux-4), float32(uy-1), float32(ux+4), float32(uy-1), 0.0, 0.0, 0.0)
		if u.Status == 5 {
			ux = fixed.Floating(u.MoveFinalX)
			uy = fixed.Floating(u.MoveFinalY)
			if float32(viewX) > ux {
				ux += float32(mWorld.Pixels)
			}
			if float32(viewY) > uy {
				uy += float32(mWorld.Pixels)
			}
			lineBuffer.RenderLine(float32(ux-2), float32(uy-1), float32(ux+3), float32(uy+3), 0.0, 0.0, 0.0)
			lineBuffer.RenderLine(float32(ux-2), float32(uy+3), float32(ux+3), float32(uy-1), 0.0, 0.0, 0.0)
			lineBuffer.RenderLine(float32(ux-2), float32(uy-2), float32(ux+2), float32(uy+2), 0.0, 1.0, 0.0)
			lineBuffer.RenderLine(float32(ux-2), float32(uy+2), float32(ux+2), float32(uy-2), 0.0, 1.0, 0.0)
			/*colorBuffer.RenderRectangle(float32(ux-1), float32(uy-1), 4.0, 4.0, 0.0, 0.0, 0.0, 1.0)
			colorBuffer.RenderRectangle(float32(ux-2), float32(uy-2), 4.0, 4.0, 0.0, 1.0, 0.0, 1.0)*/
		}
	}
	if unitsSelected {
		lineBuffer.RenderLine(float32(selectionLeft+1), float32(selectionTop+1), float32(selectionRight+1), float32(selectionTop+1), 0.0, 0.0, 0.0)
		lineBuffer.RenderLine(float32(selectionLeft+1), float32(selectionBottom+1), float32(selectionRight+1), float32(selectionBottom+1), 0.0, 0.0, 0.0)
		lineBuffer.RenderLine(float32(selectionLeft+1), float32(selectionTop+1), float32(selectionLeft+1), float32(selectionBottom+1), 0.0, 0.0, 0.0)
		lineBuffer.RenderLine(float32(selectionRight+1), float32(selectionTop+1), float32(selectionRight+1), float32(selectionBottom+1), 0.0, 0.0, 0.0)

		lineBuffer.RenderLine(float32(selectionLeft), float32(selectionTop), float32(selectionRight), float32(selectionTop), 0.0, 1.0, 0.0)
		lineBuffer.RenderLine(float32(selectionLeft), float32(selectionBottom), float32(selectionRight), float32(selectionBottom), 0.0, 1.0, 0.0)
		lineBuffer.RenderLine(float32(selectionLeft), float32(selectionTop), float32(selectionLeft), float32(selectionBottom), 0.0, 1.0, 0.0)
		lineBuffer.RenderLine(float32(selectionRight), float32(selectionTop), float32(selectionRight), float32(selectionBottom), 0.0, 1.0, 0.0)
	} else if orderingMove {
		lineBuffer.RenderLine(float32(orderingMoveX+1), float32(orderingMoveY), float32(orderingMoveX2+1), float32(orderingMoveY2), 0.0, 0.0, 0.0)
		lineBuffer.RenderLine(float32(orderingMoveX), float32(orderingMoveY+1), float32(orderingMoveX2), float32(orderingMoveY2+1), 0.0, 0.0, 0.0)
		lineBuffer.RenderLine(float32(orderingMoveX), float32(orderingMoveY), float32(orderingMoveX2), float32(orderingMoveY2), 0.0, 1.0, 0.0)
	}
	lineBuffer.DrawLineElements()
	colorBuffer.DrawElements()
}

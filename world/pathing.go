package world

import (
	"math/rand"
)

const (
	LEFT   = 0
	RIGHT  = 1
	TOP    = 2
	BOTTOM = 3
)

var (
	DEBUG_FROM_X    int32
	DEBUG_FROM_Y    int32
	DEBUG_TO_X      int32
	DEBUG_TO_Y      int32
	DEBUG_PATH_LIST *PathList
	DEBUG_RAY_TRACE *PathList
)

type NavBox struct {
	Red        float32
	Green      float32
	Blue       float32
	Top        int32
	Left       int32
	Width      int32
	Height     int32
	LeftList   *NavList
	RightList  *NavList
	TopList    *NavList
	BottomList *NavList
}

type NavPortal struct {
	Left           int32
	Top            int32
	Width          int32
	Height         int32
	InBox          *NavBox
	NeighborPortal *NavPortal
	Edges          *EdgeList
}

type PathNode struct {
	X      int32
	Y      int32
	Dir    int32
	F      int32
	G      int32
	From   *PathNode
	Portal *NavPortal
}

type NavList struct {
	Portal *NavPortal
	Next   *NavList
}

type PathList struct {
	X    int32
	Y    int32
	Next *PathList
}

type EdgeList struct {
	X    int32
	Y    int32
	Next *EdgeList
}

func (m *World) ComputePathMesh() {
	mesh := make([]*NavBox, m.MapSize*m.MapSize)
	count := 0
	for r := int32(0); r < m.MapSize; r++ {
		for c := int32(0); c < m.MapSize; c++ {
			// fill the map with boxes
			box := m.ComputePathBox(c, r)
			if box == nil {
				continue
			}
			mesh[count] = box
			count++
			c += box.Width
		}
	}
	red := 0.2 + rand.Float32()*0.8
	green := 0.2 + rand.Float32()*0.8
	blue := 0.2 + rand.Float32()*0.8
	for i := 0; i < count; i++ {
		box := mesh[i]

		box.Red = red
		box.Green = green
		box.Blue = blue
		blue = green
		green = red
		red = 0.2 + rand.Float32()*0.8

		x := box.Left
		y := box.Top

		var top *NavPortal
		var bottom *NavPortal
		for w := int32(0); w < box.Width; w++ {
			t := m.GetTile(x+w, y-1)
			if !t.Closed {
				neighbor := t.Nav
				if top == nil || neighbor != bottom.InBox {
					top = new(NavPortal)
					bottom = new(NavPortal)

					top.FlagEdge(m, x+w, y)
					top.Left = x + w
					top.Top = y
					top.Width = 1
					top.Height = 1
					top.InBox = box
					top.NeighborPortal = bottom

					bottom.FlagEdge(m, x+w, y-1)
					bottom.Left = x + w
					bottom.Top = y - 1
					bottom.Width = 1
					bottom.Height = 1
					bottom.InBox = neighbor
					bottom.NeighborPortal = top

					prev := box.TopList
					box.TopList = new(NavList)
					box.TopList.Portal = top
					box.TopList.Next = prev

					prev = neighbor.BottomList
					neighbor.BottomList = new(NavList)
					neighbor.BottomList.Portal = bottom
					neighbor.BottomList.Next = prev
				} else {
					top.Width++
					bottom.Width++
					top.FlagEdge(m, x+w, y)
					bottom.FlagEdge(m, x+w, y-1)
				}
			} else if top != nil {
				top = nil
				bottom = nil
			}
		}
		var left *NavPortal
		var right *NavPortal
		for h := int32(0); h < box.Height; h++ {
			t := m.GetTile(x-1, y+h)
			if !t.Closed {
				neighbor := t.Nav
				if left == nil || neighbor != right.InBox {
					left = new(NavPortal)
					right = new(NavPortal)

					left.FlagEdge(m, x, y+h)
					left.Left = x
					left.Top = y + h
					left.Width = 1
					left.Height = 1
					left.InBox = box
					left.NeighborPortal = right

					right.FlagEdge(m, x-1, y+h)
					right.Left = x - 1
					right.Top = y + h
					right.Width = 1
					right.Height = 1
					right.InBox = neighbor
					right.NeighborPortal = left

					prev := box.LeftList
					box.LeftList = new(NavList)
					box.LeftList.Portal = left
					box.LeftList.Next = prev

					prev = neighbor.RightList
					neighbor.RightList = new(NavList)
					neighbor.RightList.Portal = right
					neighbor.RightList.Next = prev
				} else {
					left.Height++
					right.Height++
					left.FlagEdge(m, x, y+h)
					right.FlagEdge(m, x-1, y+h)
				}
			} else if top != nil {
				left = nil
				right = nil
			}
		}
	}
}

func (m *World) ComputePathBox(c, r int32) *NavBox {
	t := m.Tiles[c+r*m.MapSize]
	if t.Closed || t.Nav != nil {
		return nil
	}
	// go as far right as possible, then as far down
	box := new(NavBox)
	w := int32(1)
	h := int32(1)
	for w < m.MapSize {
		t = m.GetTile(c+w, r)
		if t.Closed || t.Nav != nil {
			break
		}
		w++
	}
rectangle:
	for h < m.MapSize {
		for i := int32(0); i < w; i++ {
			t = m.GetTile(c+i, r+h)
			if t.Closed || t.Nav != nil {
				break rectangle
			}
		}
		h++
	}
	box.Left = c
	box.Top = r
	box.Width = w
	box.Height = h
	for i := int32(0); i < w; i++ {
		for j := int32(0); j < h; j++ {
			m.GetTile(c+i, r+j).Nav = box
		}
	}
	return box
}

func (p *NavPortal) FlagEdge(m *World, c, r int32) {
	for c < 0 {
		c += m.MapSize
	}
	for r < 0 {
		r += m.MapSize
	}
	for c >= m.MapSize {
		c -= m.MapSize
	}
	for r >= m.MapSize {
		r -= m.MapSize
	}
	if m.GetTile(c-1, r).Closed || m.GetTile(c+1, r).Closed || m.GetTile(c, r-1).Closed || m.GetTile(c, r+1).Closed {
		return
	}
	if m.GetTile(c-1, r-1).Closed || m.GetTile(c+1, r-1).Closed || m.GetTile(c-1, r+1).Closed || m.GetTile(c+1, r+1).Closed {
		m.Tiles[c+r*m.MapSize].IsEdge = true
		prev := p.Edges
		p.Edges = new(EdgeList)
		p.Edges.X = c
		p.Edges.Y = r
		p.Edges.Next = prev
	}
}

func (m *World) FindPath(fromX, fromY, toX, toY int32) *PathList {
	start := m.GetTile(fromX, fromY).Nav
	goal := m.GetTile(toX, toY).Nav
	DEBUG_FROM_X = fromX
	DEBUG_FROM_Y = fromY
	DEBUG_TO_X = toX
	DEBUG_TO_Y = toY
	DEBUG_PATH_LIST = nil
	if start == goal {
		return nil
	}
	open := make([]*PathNode, 100)
	openCount := 0
	closed := make([]*PathNode, 100)
	closedCount := 0
	// add portals in starting box to open list
	iter := start.LeftList
	for iter != nil {
		PushPortal(open, &openCount, fromX, fromY, toX, toY, LEFT, 0, nil, iter.Portal)
		iter = iter.Next
	}
	iter = start.RightList
	for iter != nil {
		PushPortal(open, &openCount, fromX, fromY, toX, toY, RIGHT, 0, nil, iter.Portal)
		iter = iter.Next
	}
	iter = start.TopList
	for iter != nil {
		PushPortal(open, &openCount, fromX, fromY, toX, toY, TOP, 0, nil, iter.Portal)
		iter = iter.Next
	}
	iter = start.BottomList
	for iter != nil {
		PushPortal(open, &openCount, fromX, fromY, toX, toY, BOTTOM, 0, nil, iter.Portal)
		iter = iter.Next
	}
	for openCount > 0 {
		index := 0
		current := open[0]
		for i := 1; i < openCount; i++ {
			if open[i].F < current.F {
				index = i
				current = open[i]
			}
		}
		if current.Portal.InBox == goal {
			path := new(PathList)
			path.Next = nil
			path.X = toX
			path.Y = toY
			end := new(PathList)
			end.Next = path
			end.X = current.X
			end.Y = current.Y
			for current.From != nil {
				current = current.From
				prev := path
				path = new(PathList)
				path.Next = prev
				path.X = current.X
				path.Y = current.Y
			}
			prev := path
			path = new(PathList)
			path.Next = prev
			path.X = fromX
			path.Y = fromY
			path.Refine(m)
			end.Next = nil
			DEBUG_PATH_LIST = path
			return path.Next
		}
		// remove current from open list and add to closed list
		for i := index; i < openCount-1; i++ {
			open[i] = open[i+1]
		}
		openCount--
		closed[closedCount] = current
		closedCount++
		// if we have not yet visited the box of the neighbor portal, then add its portals
		box := current.Portal.NeighborPortal.InBox
		if !box.Found(open, openCount, closed, closedCount) {
			iter := box.LeftList
			for iter != nil {
				PushPortal(open, &openCount, current.X, current.Y, toX, toY, LEFT, current.G, current, iter.Portal)
				iter = iter.Next
			}
			iter = box.RightList
			for iter != nil {
				PushPortal(open, &openCount, current.X, current.Y, toX, toY, RIGHT, current.G, current, iter.Portal)
				iter = iter.Next
			}
			iter = box.TopList
			for iter != nil {
				PushPortal(open, &openCount, current.X, current.Y, toX, toY, TOP, current.G, current, iter.Portal)
				iter = iter.Next
			}
			iter = box.BottomList
			for iter != nil {
				PushPortal(open, &openCount, current.X, current.Y, toX, toY, BOTTOM, current.G, current, iter.Portal)
				iter = iter.Next
			}
		}
	}
	return nil
}

func PushPortal(open []*PathNode, openCount *int, fromX, fromY, toX, toY, dir, g int32, from *PathNode, portal *NavPortal) {
	var x int32
	var y int32
	if portal.Edges != nil {
		x = portal.Edges.X
		y = portal.Edges.Y
	} else if portal.NeighborPortal.Edges != nil {
		x = portal.NeighborPortal.Edges.X
		y = portal.NeighborPortal.Edges.Y
	} else if dir == LEFT || dir == RIGHT {
		x = portal.Left
		y = toY
		if y < portal.Top {
			y = portal.Top
		} else if y > portal.Top+portal.Height-1 {
			y = portal.Top + portal.Height - 1
		}
	} else {
		x = toX
		y = portal.Top
		if x < portal.Left {
			x = portal.Left
		} else if x > portal.Left+portal.Width-1 {
			x = portal.Left + portal.Width - 1
		}
	}
	g += Approx(x-fromX, y-fromY)
	node := new(PathNode)
	node.From = from
	node.Portal = portal
	node.X = x
	node.Y = y
	node.Dir = dir
	node.F = g + Approx(x-toX, y-toY)
	node.G = g
	open[*openCount] = node
	*openCount++
}

func (find *NavBox) Found(open []*PathNode, openCount int, closed []*PathNode, closedCount int) bool {
	for i := 0; i < openCount; i++ {
		if open[i].Portal.InBox == find {
			return true
		}
	}
	for i := 0; i < closedCount; i++ {
		if closed[i].Portal.InBox == find {
			return true
		}
	}
	return false
}

func direction(dir int32) (int32, int32) {
	switch dir {
	case LEFT:
		return -1, 0
	case RIGHT:
		return 1, 0
	case TOP:
		return 0, -1
	case BOTTOM:
		return 0, 1
	}
	return 0, 0
}

func (p *PathList) Refine(m *World) {
	iter := p
	for iter.Next != nil {
		if iter.X == iter.Next.X && iter.Y == iter.Next.Y {
			iter.Next = iter.Next.Next
		} else {
			iter = iter.Next
		}
	}
	iter = p
	for iter.Next != nil && iter.Next.Next != nil {
		if m.GetTile(iter.X, iter.Y).Nav == m.GetTile(iter.Next.Next.X, iter.Next.Next.Y).Nav ||
			m.RayCast(iter.X, iter.Y, iter.Next.Next.X, iter.Next.Next.Y) {
			iter.Next = iter.Next.Next
		} else {
			iter = iter.Next
		}
	}
}

func (m *World) RayCastPoopy(fromX, fromY, toX, toY int32) bool {
	dx := toX - fromX
	if dx < 0 {
		dx = -dx
	}
	dy := toY - fromY
	if dy < 0 {
		dy = -dy
	}
	x := fromX
	y := fromY
	n := 1 + dx + dy
	var incX int32
	var incY int32
	if toX > fromX {
		incX = 1
	} else {
		incX = -1
	}
	if toY > fromY {
		incY = 1
	} else {
		incY = -1
	}
	er := dx - dy
	dx <<= 1
	dy <<= 1
	goal := m.GetTile(toX, toY).Nav
	DEBUG_RAY_TRACE = new(PathList)
	DEBUG_RAY_TRACE.X = x
	DEBUG_RAY_TRACE.Y = y
	for n > 0 {
		t := m.GetTile(x, y)
		if t.Closed {
			return false
		}
		prev := DEBUG_RAY_TRACE
		DEBUG_RAY_TRACE = new(PathList)
		DEBUG_RAY_TRACE.X = x
		DEBUG_RAY_TRACE.Y = y
		DEBUG_RAY_TRACE.Next = prev

		box := t.Nav
		if box == goal {
			return true
		}

		var bx int32
		var by int32
		if incX == 1 {
			bx = box.Left + box.Width
		} else {
			bx = box.Left - 1
		}
		if incY == 1 {
			by = box.Top + box.Height
		} else {
			by = box.Top - 1
		}
		bdx := bx - x
		bdy := by - y
		abdx := bdx
		abdy := bdy
		if abdx < 0 {
			abdx = -abdx
		}
		if abdy < 0 {
			abdy = -abdy
		}
		if abdx < abdy {
			x += bdx
			er -= dy * abdx
			for dx != 0 && er < 0 {
				y += incY
				er += dx
			}
		} else {
			y += bdy
			er += dx * abdy
			for dy != 0 && er > 0 {
				x += incX
				er -= dy
			}
		}

		n--
	}
	return true
}

func (m *World) RayCast(fromX, fromY, toX, toY int32) bool {
	dx := toX - fromX
	if dx < 0 {
		dx = -dx
	}
	dy := toY - fromY
	if dy < 0 {
		dy = -dy
	}
	x := fromX
	y := fromY
	n := 1 + dx + dy
	var incX int32
	var incY int32
	if toX > fromX {
		incX = 1
	} else {
		incX = -1
	}
	if toY > fromY {
		incY = 1
	} else {
		incY = -1
	}
	er := dx - dy
	dx <<= 1
	dy <<= 1
	DEBUG_RAY_TRACE = new(PathList)
	DEBUG_RAY_TRACE.X = x
	DEBUG_RAY_TRACE.Y = y
	for n > 0 {
		if m.GetTile(x, y).Closed {
			return false
		}
		prev := DEBUG_RAY_TRACE
		DEBUG_RAY_TRACE = new(PathList)
		DEBUG_RAY_TRACE.X = x
		DEBUG_RAY_TRACE.Y = y
		DEBUG_RAY_TRACE.Next = prev
		if er > 0 {
			x += incX
			er -= dy
		} else {
			y += incY
			er += dx
		}
		n--
	}
	return true
}

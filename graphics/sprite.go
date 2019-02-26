package graphics

type Sprite struct {
	OX float32
	OY float32

	W float32
	H float32

	U float32
	V float32
	S float32
	T float32
}

func NewSprite(left, top, width, height, sheetWidth, sheetHeight, offsetX, offsetY float32) *Sprite {
	sp := new(Sprite)

	sp.W = width
	sp.H = height

	sp.U = left * sheetWidth
	sp.V = top * sheetHeight
	sp.S = (left + width) * sheetWidth
	sp.T = (top + height) * sheetHeight

	sp.OX = offsetX
	sp.OY = offsetY

	return sp
}

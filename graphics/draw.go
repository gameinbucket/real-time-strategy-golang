package graphics

import (
	"unsafe"
)

func (b *Buffer) quad() {
	*(*uint32)(b.iPos) = 0 + b.index_offset
	b.iPos = unsafe.Pointer(uintptr(b.iPos) + uintptr(4))
	*(*uint32)(b.iPos) = 1 + b.index_offset
	b.iPos = unsafe.Pointer(uintptr(b.iPos) + uintptr(4))
	*(*uint32)(b.iPos) = 2 + b.index_offset
	b.iPos = unsafe.Pointer(uintptr(b.iPos) + uintptr(4))
	*(*uint32)(b.iPos) = 2 + b.index_offset
	b.iPos = unsafe.Pointer(uintptr(b.iPos) + uintptr(4))
	*(*uint32)(b.iPos) = 3 + b.index_offset
	b.iPos = unsafe.Pointer(uintptr(b.iPos) + uintptr(4))
	*(*uint32)(b.iPos) = 0 + b.index_offset
	b.iPos = unsafe.Pointer(uintptr(b.iPos) + uintptr(4))

	b.index_offset += 4
}

func (b *Buffer) RenderLine(xA, yA, xB, yB, red, green, blue float32) {
	*(*float32)(b.vPos) = xA
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = yA
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = red
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = green
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = blue
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = 1.0
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = xB
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = yB
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = red
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = green
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = blue
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = 1.0
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*uint32)(b.iPos) = 0 + b.index_offset
	b.iPos = unsafe.Pointer(uintptr(b.iPos) + uintptr(4))
	*(*uint32)(b.iPos) = 1 + b.index_offset
	b.iPos = unsafe.Pointer(uintptr(b.iPos) + uintptr(4))

	b.index_offset += 2
}

func (b *Buffer) RenderRectangle(x, y, w, h, red, green, blue, alpha float32) {
	*(*float32)(b.vPos) = x
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = red
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = green
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = blue
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = alpha
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y + h
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = red
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = green
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = blue
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = alpha
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x + w
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y + h
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = red
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = green
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = blue
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = alpha
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x + w
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = red
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = green
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = blue
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = alpha
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	b.quad()
}

func (b *Buffer) RenderImage(x, y, w, h, u, v, s, t float32) {
	*(*float32)(b.vPos) = x
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = u
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = v
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y + h
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = u
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = t
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x + w
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y + h
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = s
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = t
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x + w
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = s
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = v
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	b.quad()
}

func (b *Buffer) RenderSprite(sp *Sprite, x, y float32) {
	*(*float32)(b.vPos) = x
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.U
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.V
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y + sp.H
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.U
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.T
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x + sp.W
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y + sp.H
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.S
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.T
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x + sp.W
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.S
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.V
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	b.quad()
}

func (b *Buffer) RenderSpriteMirror(sp *Sprite, x, y float32) {
	*(*float32)(b.vPos) = x
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.S
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.V
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y + sp.H
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.S
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.T
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x + sp.W
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y + sp.H
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.U
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.T
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	*(*float32)(b.vPos) = x + sp.W
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = y
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.U
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))
	*(*float32)(b.vPos) = sp.V
	b.vPos = unsafe.Pointer(uintptr(b.vPos) + uintptr(4))

	b.quad()
}

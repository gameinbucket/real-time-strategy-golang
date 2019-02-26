package graphics

import (
	"unsafe"
)

type Buffer struct {
	vao uint32
	vbo uint32
	ebo uint32

	vertices unsafe.Pointer
	indices  unsafe.Pointer
	vPos     unsafe.Pointer
	iPos     unsafe.Pointer

	index_offset uint32

	vertex_limit int
	index_limit  int
}

func NewBuffer(pos, col, tex, norm, vertex_limit, index_limit int) *Buffer {
	b := new(Buffer)

	b.vertex_limit = (pos + col + tex + norm) * vertex_limit
	b.index_limit = index_limit

	b.MakeVAO(pos, col, tex, norm)

	return b
}

func (b *Buffer) Zero() {
	b.vPos = b.vertices
	b.iPos = b.indices
	b.index_offset = 0
}

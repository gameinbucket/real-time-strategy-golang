package graphics

import (
	"fmt"
	"os"
)

type Frame struct {
	fbo int

	textures                []int
	texture_internal_format []int
	texture_format          []int
	texture_type            []int

	width  int
	height int

	linear bool
	depth  bool

	draw_buffers []int

	depth_texture int
}

func NewFrame(width, height int, internal_format, format, givenType []int, linear, depth bool) *Frame {
	f := new(Frame)

	f.width = width
	f.height = height
	f.linear = linear
	f.depth = depth

	if len(internal_format) != len(format) || len(internal_format) != len(givenType) {
		fmt.Println("frame format lengths differ")
		os.Exit(1)
	}

	f.textures = make([]int, len(format))

	f.texture_internal_format = internal_format
	f.texture_format = format
	f.texture_type = givenType

	f.Framebuffer()

	return f
}

func (f *Frame) Resize(width, height int) {
	f.width = width
	f.height = height
	f.UpdateFramebuffer()
}

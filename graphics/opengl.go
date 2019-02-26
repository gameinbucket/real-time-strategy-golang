package graphics

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

const size32 = 4

var (
	res = "./res/"

	mv  = make([]float32, 16)
	mvp = make([]float32, 16)

	program uint32
)

func SetClearColor(red, green, blue float32) {
	gl.ClearColor(red, green, blue, 1.0)
}

func ClearColor() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func SetTexture0(id uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, id)
}

func SetProgram(prog uint32) {
	program = prog
	gl.UseProgram(program)
}

func (f *Frame) Framebuffer() {

}

func SetFramebuffer(id uint32) {
	gl.BindFramebuffer(gl.FRAMEBUFFER, id)
}

func (f *Frame) UpdateFramebuffer() {

}

func SetView(x, y, width, height int32) {
	gl.Viewport(x, y, width, height)
	gl.Scissor(x, y, width, height)
}

func (b *Buffer) DrawElements() {
	gl.BindVertexArray(b.vao)
	gl.DrawElements(gl.TRIANGLES, int32((uintptr(b.iPos)-uintptr(b.indices))>>2), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func (b *Buffer) DrawLineElements() {
	gl.BindVertexArray(b.vao)
	gl.DrawElements(gl.LINES, int32((uintptr(b.iPos)-uintptr(b.indices))>>2), gl.UNSIGNED_INT, nil)
	gl.BindVertexArray(0)
}

func SetOrthographic(orthographic []float32, x, y float32) {
	Identity(mv)
	Translate(mv, x, y, 0)
	Multiply(mvp, orthographic, mv)
}

func SetMVP() {
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("u_mvp\x00")), 1, false, &mvp[0])
}

func (b *Buffer) MakeVAO(pos, col, tex, norm int) {
	stride := pos + col + tex + norm

	gl.GenVertexArrays(1, &b.vao)
	gl.BindVertexArray(b.vao)

	gl.GenBuffers(1, &b.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vbo)
	gl.BufferStorage(gl.ARRAY_BUFFER, b.vertex_limit*size32, nil, gl.MAP_WRITE_BIT|gl.MAP_PERSISTENT_BIT|gl.MAP_COHERENT_BIT)
	b.vertices = gl.MapBufferRange(gl.ARRAY_BUFFER, 0, b.vertex_limit*size32, gl.MAP_WRITE_BIT|gl.MAP_PERSISTENT_BIT|gl.MAP_COHERENT_BIT)

	gl.GenBuffers(1, &b.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.BufferStorage(gl.ELEMENT_ARRAY_BUFFER, b.index_limit*size32, nil, gl.MAP_WRITE_BIT|gl.MAP_PERSISTENT_BIT|gl.MAP_COHERENT_BIT)
	b.indices = gl.MapBufferRange(gl.ELEMENT_ARRAY_BUFFER, 0, b.index_limit*size32, gl.MAP_WRITE_BIT|gl.MAP_PERSISTENT_BIT|gl.MAP_COHERENT_BIT)

	index := uint32(0)

	gl.VertexAttribPointer(index, int32(pos), gl.FLOAT, false, int32(stride*size32), gl.PtrOffset(0))
	gl.EnableVertexAttribArray(index)
	index++

	if col > 0 {
		gl.VertexAttribPointer(index, int32(col), gl.FLOAT, false, int32(stride*size32), gl.PtrOffset(pos*size32))
		gl.EnableVertexAttribArray(index)
		index++
	}

	if tex > 0 {
		gl.VertexAttribPointer(index, int32(tex), gl.FLOAT, false, int32(stride*size32), gl.PtrOffset((pos+col)*size32))
		gl.EnableVertexAttribArray(index)
		index++
	}

	if norm > 0 {
		gl.VertexAttribPointer(index, int32(norm), gl.FLOAT, false, int32(stride*size32), gl.PtrOffset((pos+col+tex)*size32))
		gl.EnableVertexAttribArray(index)
		index++
	}

	gl.BindVertexArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

func makeShader(file string, shaderType uint32) uint32 {
	b, err := ioutil.ReadFile(file)

	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	shader := gl.CreateShader(shaderType)
	source, free := gl.Strs(string(b) + "\x00")
	gl.ShaderSource(shader, 1, source, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		fmt.Println("shader error " + log)
		os.Exit(1)
	}

	return shader
}

func MakeProgram(name string, v rune) uint32 {
	var vertex uint32

	if v == 's' {
		vertex = makeShader(res+"screen_space.v", gl.VERTEX_SHADER)
	} else {
		vertex = makeShader(res+name+"-v.txt", gl.VERTEX_SHADER)
	}

	fragment := makeShader(res+name+"-f.txt", gl.FRAGMENT_SHADER)

	program := gl.CreateProgram()

	gl.AttachShader(program, vertex)
	gl.AttachShader(program, fragment)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
		fmt.Println("program error " + log)
		os.Exit(1)
	}

	gl.DeleteShader(vertex)
	gl.DeleteShader(fragment)

	gl.UseProgram(program)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("u_texture0\x00")), 0)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("u_texture1\x00")), 1)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("u_texture2\x00")), 2)
	gl.Uniform1i(gl.GetUniformLocation(program, gl.Str("u_texture3\x00")), 3)
	gl.UseProgram(0)

	return program
}

func MakeTexture(name string, clamp, linear bool) uint32 {
	file, err := os.Open(res + name)
	if err != nil {
		fmt.Println("texture", file, "not found", err)
		os.Exit(1)
	}

	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("texture", file, "error", err)
		os.Exit(1)
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		fmt.Println("texture", file, "unsupported stride")
		os.Exit(1)
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32

	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)

	if clamp {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	}

	if linear {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	}

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	return texture
}

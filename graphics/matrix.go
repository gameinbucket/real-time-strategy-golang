package graphics

import (
	"math"
)

var (
	temp   = make([]float32, 16)
	copied = make([]float32, 16)
	tmp    = make([]float32, 16)
	src    = make([]float32, 16)
	dst    = make([]float32, 16)

	vectorX float32
	vectorY float32
	vectorZ float32
)

func Identity(matrix []float32) {
	matrix[0] = 1.0
	matrix[1] = 0.0
	matrix[2] = 0.0
	matrix[3] = 0.0

	matrix[4] = 0.0
	matrix[5] = 1.0
	matrix[6] = 0.0
	matrix[7] = 0.0

	matrix[8] = 0.0
	matrix[9] = 0.0
	matrix[10] = 1.0
	matrix[11] = 0.0

	matrix[12] = 0.0
	matrix[13] = 0.0
	matrix[14] = 0.0
	matrix[15] = 1.0
}

func Orthographic(matrix []float32, left, top, right, bottom, near, far float32) {
	matrix[0] = 2.0 / (right - left)
	matrix[1] = 0.0
	matrix[2] = 0.0
	matrix[3] = 0.0

	matrix[4] = 0.0
	matrix[5] = 2.0 / (top - bottom)
	matrix[6] = 0.0
	matrix[7] = 0.0

	matrix[8] = 0.0
	matrix[9] = 0.0
	matrix[10] = -2.0 / (far - near)
	matrix[11] = 0.0

	matrix[12] = -(right + left) / (right - left)
	matrix[13] = -(top + bottom) / (top - bottom)
	matrix[14] = -(far + near) / (far - near)
	matrix[15] = 1.0
}

func Perspective(matrix []float32, fov, near, far, aspect float32) {
	top := near * float32(math.Tan(float64(fov)*math.Pi/360.0))
	bottom := -top
	left := bottom * aspect
	right := top * aspect

	Frustum(matrix, left, right, bottom, top, near, far)
}

func Frustum(matrix []float32, left, right, bottom, top, near, far float32) {
	matrix[0] = (2.0 * near) / (right - left)
	matrix[1] = 0.0
	matrix[2] = 0.0
	matrix[3] = 0.0

	matrix[4] = 0.0
	matrix[5] = (2.0 * near) / (top - bottom)
	matrix[6] = 0.0
	matrix[7] = 0.0

	matrix[8] = (right + left) / (right - left)
	matrix[9] = (top + bottom) / (top - bottom)
	matrix[10] = -(far + near) / (far - near)
	matrix[11] = -1.0

	matrix[12] = 0.0
	matrix[13] = 0.0
	matrix[14] = -(2.0 * far * near) / (far - near)
	matrix[15] = 0.0
}

func Translate(matrix []float32, x, y, z float32) {
	matrix[12] = x*matrix[0] + y*matrix[4] + z*matrix[8] + matrix[12]
	matrix[13] = x*matrix[1] + y*matrix[5] + z*matrix[9] + matrix[13]
	matrix[14] = x*matrix[2] + y*matrix[6] + z*matrix[10] + matrix[14]
	matrix[15] = x*matrix[3] + y*matrix[7] + z*matrix[11] + matrix[15]
}

func RotateX(matrix []float32, sin, cos float32) {
	temp[0] = 1.0
	temp[1] = 0.0
	temp[2] = 0.0
	temp[3] = 0.0

	temp[4] = 0.0
	temp[5] = cos
	temp[6] = sin
	temp[7] = 0.0

	temp[8] = 0.0
	temp[9] = -sin
	temp[10] = cos
	temp[11] = 0.0

	temp[12] = 0.0
	temp[13] = 0.0
	temp[14] = 0.0
	temp[15] = 1.0

	for i := 0; i < 16; i++ {
		copied[i] = matrix[i]
	}

	Multiply(matrix, copied, temp)
}

func RotateY(matrix []float32, sin, cos float32) {
	temp[0] = cos
	temp[1] = 0.0
	temp[2] = -sin
	temp[3] = 0.0

	temp[4] = 0.0
	temp[5] = 1.0
	temp[6] = 0.0
	temp[7] = 0.0

	temp[8] = sin
	temp[9] = 0.0
	temp[10] = cos
	temp[11] = 0.0

	temp[12] = 0.0
	temp[13] = 0.0
	temp[14] = 0.0
	temp[15] = 1.0

	for i := 0; i < 16; i++ {
		copied[i] = matrix[i]
	}

	Multiply(matrix, copied, temp)
}

func RotateZ(matrix []float32, sin, cos float32) {
	temp[0] = cos
	temp[1] = sin
	temp[2] = 0.0
	temp[3] = 0.0

	temp[4] = -sin
	temp[5] = cos
	temp[6] = 0.0
	temp[7] = 0.0

	temp[8] = 0.0
	temp[9] = 0.0
	temp[10] = 1.0
	temp[11] = 0.0

	temp[12] = 0.0
	temp[13] = 0.0
	temp[14] = 0.0
	temp[15] = 1.0

	for i := 0; i < 16; i++ {
		copied[i] = matrix[i]
	}

	Multiply(matrix, copied, temp)
}

func Multiply(matrix, b, c []float32) {
	matrix[0] = b[0]*c[0] + b[4]*c[1] + b[8]*c[2] + b[12]*c[3]
	matrix[1] = b[1]*c[0] + b[5]*c[1] + b[9]*c[2] + b[13]*c[3]
	matrix[2] = b[2]*c[0] + b[6]*c[1] + b[10]*c[2] + b[14]*c[3]
	matrix[3] = b[3]*c[0] + b[7]*c[1] + b[11]*c[2] + b[15]*c[3]

	matrix[4] = b[0]*c[4] + b[4]*c[5] + b[8]*c[6] + b[12]*c[7]
	matrix[5] = b[1]*c[4] + b[5]*c[5] + b[9]*c[6] + b[13]*c[7]
	matrix[6] = b[2]*c[4] + b[6]*c[5] + b[10]*c[6] + b[14]*c[7]
	matrix[7] = b[3]*c[4] + b[7]*c[5] + b[11]*c[6] + b[15]*c[7]

	matrix[8] = b[0]*c[8] + b[4]*c[9] + b[8]*c[10] + b[12]*c[11]
	matrix[9] = b[1]*c[8] + b[5]*c[9] + b[9]*c[10] + b[13]*c[11]
	matrix[10] = b[2]*c[8] + b[6]*c[9] + b[10]*c[10] + b[14]*c[11]
	matrix[11] = b[3]*c[8] + b[7]*c[9] + b[11]*c[10] + b[15]*c[11]

	matrix[12] = b[0]*c[12] + b[4]*c[13] + b[8]*c[14] + b[12]*c[15]
	matrix[13] = b[1]*c[12] + b[5]*c[13] + b[9]*c[14] + b[13]*c[15]
	matrix[14] = b[2]*c[12] + b[6]*c[13] + b[10]*c[14] + b[14]*c[15]
	matrix[15] = b[3]*c[12] + b[7]*c[13] + b[11]*c[14] + b[15]*c[15]
}

func Inverse(matrix, b []float32) {
	for i := 0; i < 4; i++ {
		src[i+0] = b[i*4+0]
		src[i+4] = b[i*4+1]
		src[i+8] = b[i*4+2]
		src[i+12] = b[i*4+3]
	}

	tmp[0] = src[10] * src[15]
	tmp[1] = src[11] * src[14]
	tmp[2] = src[9] * src[15]
	tmp[3] = src[11] * src[13]
	tmp[4] = src[9] * src[14]
	tmp[5] = src[10] * src[13]
	tmp[6] = src[8] * src[15]
	tmp[7] = src[11] * src[12]
	tmp[8] = src[8] * src[14]
	tmp[9] = src[10] * src[12]
	tmp[10] = src[8] * src[13]
	tmp[11] = src[9] * src[12]

	dst[0] = tmp[0]*src[5] + tmp[3]*src[6] + tmp[4]*src[7]
	dst[0] -= tmp[1]*src[5] + tmp[2]*src[6] + tmp[5]*src[7]
	dst[1] = tmp[1]*src[4] + tmp[6]*src[6] + tmp[9]*src[7]
	dst[1] -= tmp[0]*src[4] + tmp[7]*src[6] + tmp[8]*src[7]
	dst[2] = tmp[2]*src[4] + tmp[7]*src[5] + tmp[10]*src[7]
	dst[2] -= tmp[3]*src[4] + tmp[6]*src[5] + tmp[11]*src[7]
	dst[3] = tmp[5]*src[4] + tmp[8]*src[5] + tmp[11]*src[6]
	dst[3] -= tmp[4]*src[4] + tmp[9]*src[5] + tmp[10]*src[6]
	dst[4] = tmp[1]*src[1] + tmp[2]*src[2] + tmp[5]*src[3]
	dst[4] -= tmp[0]*src[1] + tmp[3]*src[2] + tmp[4]*src[3]
	dst[5] = tmp[0]*src[0] + tmp[7]*src[2] + tmp[8]*src[3]
	dst[5] -= tmp[1]*src[0] + tmp[6]*src[2] + tmp[9]*src[3]
	dst[6] = tmp[3]*src[0] + tmp[6]*src[1] + tmp[11]*src[3]
	dst[6] -= tmp[2]*src[0] + tmp[7]*src[1] + tmp[10]*src[3]
	dst[7] = tmp[4]*src[0] + tmp[9]*src[1] + tmp[10]*src[2]
	dst[7] -= tmp[5]*src[0] + tmp[8]*src[1] + tmp[11]*src[2]

	tmp[0] = src[2] * src[7]
	tmp[1] = src[3] * src[6]
	tmp[2] = src[1] * src[7]
	tmp[3] = src[3] * src[5]
	tmp[4] = src[1] * src[6]
	tmp[5] = src[2] * src[5]
	tmp[6] = src[0] * src[7]
	tmp[7] = src[3] * src[4]
	tmp[8] = src[0] * src[6]
	tmp[9] = src[2] * src[4]
	tmp[10] = src[0] * src[5]
	tmp[11] = src[1] * src[4]

	dst[8] = tmp[0]*src[13] + tmp[3]*src[14] + tmp[4]*src[15]
	dst[8] -= tmp[1]*src[13] + tmp[2]*src[14] + tmp[5]*src[15]
	dst[9] = tmp[1]*src[12] + tmp[6]*src[14] + tmp[9]*src[15]
	dst[9] -= tmp[0]*src[12] + tmp[7]*src[14] + tmp[8]*src[15]
	dst[10] = tmp[2]*src[12] + tmp[7]*src[13] + tmp[10]*src[15]
	dst[10] -= tmp[3]*src[12] + tmp[6]*src[13] + tmp[11]*src[15]
	dst[11] = tmp[5]*src[12] + tmp[8]*src[13] + tmp[11]*src[14]
	dst[11] -= tmp[4]*src[12] + tmp[9]*src[13] + tmp[10]*src[14]
	dst[12] = tmp[2]*src[10] + tmp[5]*src[11] + tmp[1]*src[9]
	dst[12] -= tmp[4]*src[11] + tmp[0]*src[9] + tmp[3]*src[10]
	dst[13] = tmp[8]*src[11] + tmp[0]*src[8] + tmp[7]*src[10]
	dst[13] -= tmp[6]*src[10] + tmp[9]*src[11] + tmp[1]*src[8]
	dst[14] = tmp[6]*src[9] + tmp[11]*src[11] + tmp[3]*src[8]
	dst[14] -= tmp[10]*src[11] + tmp[2]*src[8] + tmp[7]*src[9]
	dst[15] = tmp[10]*src[10] + tmp[4]*src[8] + tmp[9]*src[9]
	dst[15] -= tmp[8]*src[9] + tmp[11]*src[10] + tmp[5]*src[8]

	det := 1.0 / (src[0]*dst[0] + src[1]*dst[1] + src[2]*dst[2] + src[3]*dst[3])

	for i := 0; i < 16; i++ {
		matrix[i] = dst[i] * det
	}
}

func MultiplyVector(matrix []float32, x, y, z float32) {
	vectorX = matrix[0]*x + matrix[4]*y + matrix[8]*z + matrix[12]
	vectorY = matrix[1]*x + matrix[5]*y + matrix[9]*z + matrix[13]
	vectorZ = matrix[2]*x + matrix[6]*y + matrix[10]*z + matrix[14]
}

func Pitch(matrix []float32, pitch float32) {
	cos := float32(math.Cos(float64(pitch)))
	sin := float32(math.Sin(float64(pitch)))

	matrix[0] = 1.0
	matrix[1] = 0.0
	matrix[2] = 0.0
	matrix[3] = 0.0

	matrix[4] = 0.0
	matrix[5] = cos
	matrix[6] = sin
	matrix[7] = 0.0

	matrix[8] = 0.0
	matrix[9] = -sin
	matrix[10] = cos
	matrix[11] = 0.0

	matrix[12] = 0.0
	matrix[13] = 0.0
	matrix[14] = 0.0
	matrix[15] = 1.0
}

func Yaw(matrix []float32, yaw float32) {
	cos := float32(math.Cos(float64(yaw)))
	sin := float32(math.Sin(float64(yaw)))

	matrix[0] = cos
	matrix[1] = 0.0
	matrix[2] = -sin
	matrix[3] = 0.0

	matrix[4] = 0.0
	matrix[5] = 1.0
	matrix[6] = 0.0
	matrix[7] = 0.0

	matrix[8] = sin
	matrix[9] = 0.0
	matrix[10] = cos
	matrix[11] = 0.0

	matrix[12] = 0.0
	matrix[13] = 0.0
	matrix[14] = 0.0
	matrix[15] = 1.0
}

func Roll(matrix []float32, roll float32) {
	cos := float32(math.Cos(float64(roll)))
	sin := float32(math.Sin(float64(roll)))

	matrix[0] = cos
	matrix[1] = sin
	matrix[2] = 0.0
	matrix[3] = 0.0

	matrix[4] = -sin
	matrix[5] = cos
	matrix[6] = 0.0
	matrix[7] = 0.0

	matrix[8] = 0.0
	matrix[9] = 0.0
	matrix[10] = 1.0
	matrix[11] = 0.0

	matrix[12] = 0.0
	matrix[13] = 0.0
	matrix[14] = 0.0
	matrix[15] = 1.0
}

package fixed

import (
	"math"
)

const (
	FracBits  uint32 = 16
	FracUnits        = 1 << FracBits

	sinPrecise = false
	sinP       = 14745
)

var (
	One = Whole(1)
	Two = Whole(2)

	Pi     = Number(3, 14159)
	Tau    = Mul(Pi, Two)
	PiHalf = Div(Pi, Two)

	RadianToDegree = Div(Whole(180), Pi)

	sinB = Div(Whole(4), Pi)
	sinC = Div(Whole(-4), Mul(Pi, Pi))
)

func Mul(a, b int32) int32 {
	return int32((int64(a) * int64(b)) >> FracBits)
}

func Div(a, b int32) int32 {
	return int32((int64(a) << FracBits) / int64(b))
}

func Whole(val int32) int32 {
	return val << FracBits
}

func Digits(x int32) int32 {
	if x < 0 {
		x = -x
	}

	if x < 10 {
		return 10
	} else if x < 100 {
		return 100
	} else if x < 1000 {
		return 1000
	} else if x < 10000 {
		return 10000
	} else {
		return 100000
	}
}

func Fraction(val int32) int32 {
	return (val << FracBits) / Digits(val)
}

func Number(num, frac int32) int32 {
	if num < 0 {
		return -(Whole(-num) + Fraction(frac))
	} else {
		return Whole(num) + Fraction(frac)
	}
}

func Integer(val int32) int32 {
	return val >> FracBits
}

func FloatToFixed(val float32) int32 {
	return int32(val * FracUnits)
}

func Floating(val int32) float32 {
	return float32(val) / FracUnits
}

func Sin(val int32) int32 {
	if val > Pi {
		val -= Tau
	} else if val < -Pi {
		val += Tau
	}

	absVal := val

	if val < 0 {
		absVal = -val
	}

	y := Mul(sinB, val) + Mul(Mul(sinC, val), absVal)

	if sinPrecise {
		absY := y

		if y < 0 {
			absY = -y
		}

		y = Mul(sinP, Mul(y, absY)-y) + y
	}

	return y
}

func Cos(val int32) int32 {
	return Sin(val + PiHalf)
}

func Sqrt(val int32) int32 {
	return FloatToFixed(float32(math.Sqrt(float64(Floating(val)))))
}

func Atan2(a, b int32) int32 {
	x := Floating(a)
	y := Floating(b)

	return FloatToFixed(float32(math.Atan2(float64(x), float64(y))))
}

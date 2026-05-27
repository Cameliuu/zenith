package overlay

import (
	"math"
)

const (
	ScreenW = 1536
	ScreenH = 940
)

/*
================================================================ Utilities ================================================================
*/
type ViewAngles struct {
	Yaw, Pitch float32
}
type ViewMatrix struct {
	Matrix [16]float32
}
type Vector3 struct {
	X float32
	Y float32
	Z float32
}

func AnglesToVectors(angles ViewAngles) (forward, right, up Vector3) {
	p := float64(angles.Pitch) * math.Pi / 180.0
	y := float64(angles.Yaw) * math.Pi / 180.0

	sp := math.Sin(p)
	cp := math.Cos(p)
	sy := math.Sin(y)
	cy := math.Cos(y)

	forward = Vector3{
		X: float32(cp * cy),
		Y: float32(cp * sy),
		Z: float32(-sp),
	}
	right = Vector3{
		X: float32(sy),
		Y: float32(-cy),
		Z: 0,
	}
	up = Vector3{
		X: float32(sp * cy),
		Y: float32(sp * sy),
		Z: float32(cp),
	}
	return
}
func WorldToScreen(viewMatrix ViewMatrix, point Vector3, screenW, screenH int) (float32, float32, bool) {
	x := point.X
	y := point.Y
	z := point.Z

	screenX := viewMatrix.Matrix[0]*x + viewMatrix.Matrix[4]*y + viewMatrix.Matrix[8]*z + viewMatrix.Matrix[12]
	screenY := viewMatrix.Matrix[1]*x + viewMatrix.Matrix[5]*y + viewMatrix.Matrix[9]*z + viewMatrix.Matrix[13]
	screenW2 := viewMatrix.Matrix[3]*x + viewMatrix.Matrix[7]*y + viewMatrix.Matrix[11]*z + viewMatrix.Matrix[15]

	if screenW2 <= 0.001 {
		return 0, 0, false
	}

	ndcX := screenX / screenW2
	ndcY := screenY / screenW2

	if ndcX >= 1 || ndcY >= 1 || ndcX <= -1 || ndcY <= -1 || (1.0/screenW2) <= 0 {
		return 0, 0, false
	}

	halfW := float32(screenW) / 2.0
	halfH := float32(screenH) / 2.0

	px := ndcX*halfW + halfW
	py := -ndcY*halfH + halfH

	return px, py, true
}

func Dot(a, b Vector3) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

package main

import (
	"math"
)

const (
	MessageTypeError = 13
)

const (
	frictionForce float32 = 0.03

	ScreenWidth  = 640
	ScreenHeight = 480
)

type (
	Rect struct {
		ID     string  `json:"id"`
		Width  float32 `json:"w"`
		Height float32 `json:"h"`
		Vector Vector  `json:"vc"`

		PrevX      float32 `json:"px"`
		PrevY      float32 `json:"py"`
		Speed      Vector  `json:"s"`
		SpeedLimit float32 `json:"sl"`
	}

	Vector struct {
		X, Y float32
	}
)

func (v Vector) Distance(in Vector) float32 {
	dx := v.X - in.X
	dy := v.Y - in.Y
	return float32(math.Sqrt(float64(dx*dx + dy*dy)))
}

func NewVector(x, y float32) Vector {
	return Vector{X: x, Y: y}
}

func (r *Rect) HasCollisionWith(rect Rect) bool {
	return float32(r.Vector.Distance(rect.Vector)) < r.Width*2
}

func (r *Rect) UpdateXY(x, y, screenHeight, screenWidth float32) {
	r.PrevX = r.Vector.X
	r.PrevY = r.Vector.Y

	r.Vector.X += x
	r.Vector.Y += y

	if r.ReflectFromScreen(screenHeight, screenWidth) {
		r.Vector.X -= x
		r.Vector.Y -= y
	}
}

func (r *Rect) CalculateSpeed() {
	r.Speed.X = r.Vector.X - r.PrevX
	r.Speed.Y = r.Vector.Y - r.PrevY

	if r.SpeedLimit > 0 {
		if r.Speed.X >= r.SpeedLimit {
			r.Speed.X = r.SpeedLimit
		}
		if r.Speed.Y >= r.SpeedLimit {
			r.Speed.Y = r.SpeedLimit
		}
	}
}

func (r *Rect) AddSpeed(speedX, speedY float32) {
	r.Speed.X += speedX
	r.Speed.Y += speedY
}

func (r *Rect) RestrictSpeedLimit() {
	if r.Speed.X >= r.SpeedLimit {
		r.Speed.X = r.SpeedLimit
	}
	if r.Speed.Y >= r.SpeedLimit {
		r.Speed.Y = r.SpeedLimit
	}
}

func (r *Rect) SlowDown() {
	if r.Speed.X > 0 {
		r.Speed.X -= frictionForce
	}

	if r.Speed.Y > 0 {
		r.Speed.Y -= frictionForce
	}
}

func (v *Rect) Heading() float64 {
	val := math.Atan2(float64(v.Speed.Y), float64(v.Speed.X))
	return val
}

func (v *Rect) SpeedMag() float64 {
	return math.Sqrt(float64(v.Speed.X*v.Speed.X) + float64(v.Speed.Y*v.Speed.Y))
}

func (v *Rect) ReflectFrom(other Rect) {
	other.Speed.X = -(v.Speed.X * 2)
	other.Speed.Y = -(v.Speed.Y * 2)

	var phi = math.Atan(float64((v.Vector.Y-other.Vector.Y)/v.Vector.X - other.Vector.X))

	var theta1 = v.Heading()
	var theta2 = other.Heading()

	newSpeedX := (v.SpeedMag()*math.Cos(theta1-phi)+2*2*other.SpeedMag()*math.Cos(theta2-phi))/(2)*math.Cos(phi) + v.SpeedMag()*math.Sin(theta1-phi)*math.Sin(phi)
	newSpeedY := (v.SpeedMag()*math.Cos(theta1-phi)+2*2*other.SpeedMag()*math.Cos(theta2-phi))/(2)*math.Sin(phi) + v.SpeedMag()*math.Sin(theta1-phi)*math.Cos(phi)

	v.Speed.X = float32(newSpeedX)
	v.Speed.Y = float32(newSpeedY)

	for v.HasCollisionWith(other) {
		v.Vector.X += float32(newSpeedX)
		v.Vector.Y += float32(newSpeedY)
	}
}

func (r *Rect) ReflectFromScreen(screenHeight, screenWidth float32) bool {
	if r.Vector.Y+r.Height >= screenHeight {
		r.Speed.Y = -r.Speed.Y
		r.Vector.Y += r.Speed.Y

		return true
	}

	if r.Vector.X+r.Width >= screenWidth {
		r.Speed.X = -r.Speed.X
		r.Vector.X += r.Speed.X

		return true
	}

	if r.Vector.Y <= 0 {
		r.Speed.Y = -r.Speed.Y
		r.Vector.Y += r.Speed.Y

		return true
	}

	if r.Vector.X <= 0 {
		r.Speed.X = -r.Speed.X
		r.Vector.X += r.Speed.X

		return true
	}

	return false
}

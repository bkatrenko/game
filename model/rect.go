package model

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
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
		ID     string        `json:"id"`
		Image  *ebiten.Image `json:"-"`
		Width  float32       `json:"w"`
		Height float32       `json:"h"`
		Vector Vector        `json:"vc"`

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

func DrawCircle(screen *ebiten.Image, x, y float64, radius int, clr color.Color, fill bool) {
	radius64 := float64(radius)
	minAngle := math.Acos(1 - 1/radius64)

	for angle := float64(0); angle <= 360; angle += minAngle {
		xDelta := radius64 * math.Cos(angle)
		yDelta := radius64 * math.Sin(angle)

		x1 := int(math.Round(x + xDelta))
		y1 := int(math.Round(y + yDelta))

		screen.Set(x1, y1, clr)
	}

	if fill && radius > 1 {
		for r := radius - 1; r >= 1; r-- {
			DrawCircle(screen, x, y, r, clr, false)
		}
	}
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

func (r *Rect) UpdateXYClient(x, y, screenHeight, screenWidth float32) {
	r.Vector.X += x
	r.Vector.Y += y
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

func (r *Rect) ReflectFromScreen(screenHeight, screenWidth float32) bool {
	if r.Vector.Y+BallRadius >= screenHeight {
		r.Speed.Y = -r.Speed.Y
		r.Vector.Y += r.Speed.Y

		return true
	}

	if r.Vector.X+BallRadius >= screenWidth {
		r.Speed.X = -r.Speed.X
		r.Vector.X += r.Speed.X

		return true
	}

	if r.Vector.Y-BallRadius <= 0 {
		r.Speed.Y = -r.Speed.Y
		r.Vector.Y += r.Speed.Y

		return true
	}

	if r.Vector.X-BallRadius <= 0 {
		r.Speed.X = -r.Speed.X
		r.Vector.X += r.Speed.X

		return true
	}

	return false
}

package source

import (
	"math"
	"math/rand"
)

// Scene represents the view
type Scene struct {
	source, corner, horizon, vertical Vector
	objects                           []Hittable
	aperture                          float64
}

// NewScene returns a new View
func NewScene(source, corner, horizon, vertical Vector, aperture float64) Scene {
	return Scene{source, corner, horizon, vertical, make([]Hittable, 0), aperture}
}

// Register registers an object to hit
func (scn *Scene) Register(obj Hittable) { scn.objects = append(scn.objects, obj) }

// Source returns the source of the scene
func (scn Scene) Source() Vector { return scn.source }

// Corner returns the corner of the scene
func (scn Scene) Corner() Vector { return scn.corner }

// Horizon returns the horizon of the scene
func (scn Scene) Horizon() Vector { return scn.horizon }

// Vertical returns the vertical of the scene
func (scn Scene) Vertical() Vector { return scn.vertical }

// Aperture returns the aperture of the scene
func (scn Scene) Aperture() float64 { return scn.aperture }

// SourceTo sets the source of the scene
func (scn *Scene) SourceTo(source Vector) { scn.source = source }

// CornerTo sets the corner of the scene
func (scn *Scene) CornerTo(corner Vector) { scn.corner = corner }

// HorizonTo sets the horizon of the scene
func (scn *Scene) HorizonTo(horizon Vector) { scn.horizon = horizon }

// VerticalTo sets the vertical of the scene
func (scn *Scene) VerticalTo(vertical Vector) { scn.vertical = vertical }

// ApertureTo returns the aperture of the scene
func (scn *Scene) ApertureTo(aperture float64) { scn.aperture = aperture }

// ColorTrace tracks the color of a path
func (scn Scene) ColorTrace(source, towards Vector, depth int, gen *rand.Rand) Vector {
	color := NewVector(1, 1, 1)
	for d := 0; d < depth; d++ {
		if data := scn.Ray(source, towards); data.Err == nil {
			matter := data.Matter
			reflected := matter.Scatter(towards, data.Normal.Unit(), gen)
			color.IMul(matter.Albedo())
			source, towards = data.Point, reflected
		} else {
			t := (towards.Unit().Y() + 1) * .5
			background := NewVector(1, 1, 1).MulS(1 - t).Add(NewVector(.5, .7, 1).MulS(t))
			return color.Mul(background)
		}
	}
	return Vector{0, 0, 0}
}

// Color determines the color at a given position
func (scn Scene) Color(x, y, ns, depth int, dx, dy float64, gen *rand.Rand) (r, g, b uint8) {
	i, j := RandomDisk(scn.aperture, gen)
	h, v := scn.horizon.Unit().MulS(i), scn.vertical.Unit().MulS(j)
	start := scn.source.Add(h).Add(v)

	var color Vector
	for s := 0; s < ns; s++ {
		i, j := (float64(x)+gen.Float64())/dx, (float64(y)+gen.Float64())/dy
		end := scn.corner.Add(scn.horizon.MulS(i)).Add(scn.vertical.MulS(j))
		towards := end.Sub(start)

		color.IAdd(scn.ColorTrace(start, towards, depth, gen))
	}

	pixel := color.DivS(float64(ns)).MulS(255.999)
	r = uint8(pixel.X())
	g = uint8(pixel.Y())
	b = uint8(pixel.Z())
	return
}

// Ray represents a ray
func (scn Scene) Ray(source, towards Vector) HitData {
	towards = towards.Unit()
	closest := HitData{T: math.Inf(1), Err: ErrNotHit}
	for _, obj := range scn.objects {
		data := obj.Hit(source, towards)
		if data.Err == nil && data.T < closest.T {
			closest = data
		}
	}
	return closest
}

// RandomDisk creates a random pair of tuple that lies in a unit disk
func RandomDisk(radius float64, gen *rand.Rand) (x, y float64) {
	for {
		x, y = gen.Float64(), gen.Float64()
		if x*x+y*y <= 1 {
			return x * radius, y * radius
		}
	}
}
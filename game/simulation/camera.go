package simulation

import "github.com/z46-dev/florr-legacy/shared/util"

// newCamera initializes a new Camera instance with default values for position, velocity, field of view, and friction coefficient.
func newCamera() (camera *Camera) {
	camera = &Camera{
		X:               0,
		Y:               0,
		PrevX:           0,
		PrevY:           0,
		RealX:           0,
		RealY:           0,
		FieldOfView:     1024.0,
		PrevFieldOfView: 1024.0,
		RealFieldOfView: 1024.0,
	}

	return
}

// Update interpolates the camera from the previous world snapshot to the current one.
func (c *Camera) Update(lerpCoefficient float32) {
	c.X = util.Lerp32(c.PrevX, c.RealX, lerpCoefficient)
	c.Y = util.Lerp32(c.PrevY, c.RealY, lerpCoefficient)
	c.FieldOfView = util.Lerp32(c.PrevFieldOfView, c.RealFieldOfView, lerpCoefficient)
}

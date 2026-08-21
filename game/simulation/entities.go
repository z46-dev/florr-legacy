package simulation

import (
	"math"

	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/game/simulation/renders"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/florr-legacy/shared/util"
	"github.com/z46-dev/gctx2d"
)

func (c *GenericWorldObject) Update(lerpAlpha float32, deltaSeconds float64) {
	c.X = util.Lerp32(c.PrevX, c.RealX, lerpAlpha)
	c.Y = util.Lerp32(c.PrevY, c.RealY, lerpAlpha)
	c.Size = util.Lerp32(c.PrevSize, c.RealSize, lerpAlpha)
	c.Rotation = util.LerpAngle32(c.PrevRotation, c.RealRotation, lerpAlpha)

	if c.Damaged && !c.DamagedFlashWasActive {
		c.DamagedFlashCurrent = 1
	}

	c.DamagedFlashWasActive = c.Damaged
	c.DamagedFlashCurrent = max(0, c.DamagedFlashCurrent-float32(deltaSeconds/(3.0/20.0)))

	// the death delta is used to track 1 second of death animation time. During this time,
	// the entity will grow by 15% total over the second and fade out to 0% alpha.
	// After 1 second, the entity will be removed from the simulation.
	if c.IsDead {
		c.DeathDelta += float32(deltaSeconds)
		c.RealSize *= 1.025
		c.Alpha *= 0.9
	}
}

func (hasHealth *HasHealth) Update(lerpAlpha float32, deltaSeconds float64) {
	hasHealth.GenericWorldObject.Update(lerpAlpha, deltaSeconds)

	hasHealth.HealthRatio = util.Lerp32(hasHealth.HealthRatio, hasHealth.RealHealthRatio, lerpAlpha)

	if hasHealth.LastDamagedAt > 0 && !hasHealth.IsDead {
		hasHealth.LastDamagedAt -= float32(deltaSeconds)
		if hasHealth.LastDamagedAt < 0 {
			hasHealth.LastDamagedAt = 0
		}
	} else {
		hasHealth.SubHealthRatio = util.Lerp32(hasHealth.SubHealthRatio, hasHealth.HealthRatio, lerpAlpha*0.25)
	}
}

func (hasHealth *HasHealth) DrawHealth(ctx *gctx2d.Context, scale float32) {
	if hasHealth.HealthRatio >= 1 {
		return
	}

	var (
		padding float32 = 5 * scale
		y       float32 = hasHealth.Y + hasHealth.Size + 25*scale
	)

	ctx.SetLineCap(gctx2d.LineCapRound)
	ctx.SetLineWidth(6.5 / scale)
	ctx.SetStrokeStyle(gctx2d.Color{R: 0.1, G: 0.1, B: 0.1, A: 1})

	ctx.BeginPath()
	ctx.MoveTo(hasHealth.X-hasHealth.Size-padding, y)
	ctx.LineTo(hasHealth.X+hasHealth.Size+padding, y)
	ctx.Stroke()

	ctx.SetLineWidth(4 / scale)
	ctx.SetStrokeStyle(configs.RarityLegendary.Color)

	ctx.BeginPath()
	ctx.MoveTo(hasHealth.X-hasHealth.Size-padding, y)
	ctx.LineTo(hasHealth.X-hasHealth.Size-padding+(2*(hasHealth.Size+padding))*hasHealth.SubHealthRatio, y)
	ctx.Stroke()

	ctx.SetLineWidth(4 / scale)
	ctx.SetStrokeStyle(configs.RarityCommon.Color)

	ctx.BeginPath()
	ctx.MoveTo(hasHealth.X-hasHealth.Size-padding, y)
	ctx.LineTo(hasHealth.X-hasHealth.Size-padding+(2*(hasHealth.Size+padding))*hasHealth.HealthRatio, y)
	ctx.Stroke()
}

// Update advances flower interpolation and the archived expression easing.
func (f *Flower) Update(lerpAlpha float32, deltaSeconds float64) {
	f.HasHealth.Update(lerpAlpha, deltaSeconds)

	var (
		target        float32 = 5
		eyeClipTarget float32
	)

	switch f.AttackMode {
	case configs.AttackMode_Attack:
		target = -5
		eyeClipTarget = 1
	case configs.AttackMode_Defend:
		target = -3
	}

	f.MouthCurve = target + (f.MouthCurve-target)*float32(math.Pow(0.8, deltaSeconds*60))
	f.EyeClip = eyeClipTarget + (f.EyeClip-eyeClipTarget)*float32(math.Pow(0.8, deltaSeconds*60))
}

var (
	//rgba(251, 231, 80)
	flowerYellow gctx2d.Color = gctx2d.RGBA(251.0/255.0, 231.0/255.0, 80.0/255.0, 1)
	flowerFace   gctx2d.Color = gctx2d.RGBA(34.0/255.0, 34.0/255.0, 34.0/255.0, 1)
	flowerEye    gctx2d.Color = gctx2d.RGBA(238.0/255.0, 238.0/255.0, 238.0/255.0, 1)
)

// Draw renders the archived flower body and expression geometry.
func (f *Flower) Draw(ctx *gctx2d.Context) {
	var (
		scale     float32 = f.Size / 25
		fillColor         = flowerYellow
	)

	ctx.SetGlobalAlpha(f.Alpha)

	fillColor = renders.DamageFlashColor(fillColor, f.DamagedFlashCurrent)

	ctx.SetFillStyle(fillColor)
	ctx.SetStrokeStyle(render.MixColors(fillColor, gctx2d.ColorBlack, 0.2))
	ctx.SetLineWidth(3 * scale)
	ctx.FillCircle(f.X, f.Y, f.Size)
	ctx.StrokeCircle(f.X, f.Y, f.Size)

	var (
		lookSin float32 = float32(math.Sin(float64(f.Rotation)))
		lookCos float32 = float32(math.Cos(float64(f.Rotation)))
		lookX   float32 = lookCos * 2.35 * scale
		lookY   float32 = lookSin * (1.1 + 0.65*(1-float32(math.Abs(float64(lookCos))))) * scale
	)

	for _, eyeX := range []float32{-7, 7} {
		var x, y float32 = f.X + eyeX*scale, f.Y - 5*scale
		ctx.Save()
		if f.EyeClip > 0.001 {
			ctx.BeginPath()
			ctx.MoveTo(f.X, f.Y+util.Lerp32(-13, -5, f.EyeClip)*scale)
			ctx.LineTo(f.X+eyeX*2*scale, f.Y+util.Lerp32(-13, -14, f.EyeClip)*scale)
			ctx.LineTo(f.X+eyeX*2*scale, f.Y+6*scale)
			ctx.LineTo(f.X, f.Y+6*scale)
			ctx.ClosePath()
			ctx.Clip()
		}

		ctx.BeginPath()
		ctx.Ellipse(x, y, 3*scale, 6*scale, 0, 0, renders.TAU, false)
		ctx.SetFillStyle(flowerFace)
		ctx.Fill()

		ctx.Save()
		ctx.Clip()
		ctx.BeginPath()
		ctx.Arc(x+lookX, y+lookY, 3*scale, 0, renders.TAU, false)
		ctx.SetFillStyle(flowerEye)
		ctx.Fill()
		ctx.Restore()

		ctx.BeginPath()
		ctx.Ellipse(x, y, 3*scale, 6*scale, 0, 0, renders.TAU, false)
		ctx.SetStrokeStyle(flowerFace)
		ctx.SetLineWidth(0.5 * scale)
		ctx.Stroke()
		ctx.Restore()
	}

	ctx.BeginPath()
	ctx.MoveTo(f.X-6*scale, f.Y+10*scale)
	ctx.QuadraticCurveTo(f.X, f.Y+(10+f.MouthCurve)*scale, f.X+6*scale, f.Y+10*scale)
	ctx.SetStrokeStyle(flowerFace)
	ctx.SetLineWidth(1.5 * scale)
	ctx.SetLineCap(gctx2d.LineCapRound)
	ctx.Stroke()
}

func (p *Petal) Draw(ctx *gctx2d.Context, time, scale float32) {
	ctx.SetGlobalAlpha(p.Alpha)
	renders.Petals[p.PetalID].Active(ctx, time, scale, p.X, p.Y, p.Size, p.Rotation, p.ID, p.Alpha, p.DamagedFlashCurrent, false)
}

func (m *Mob) Draw(ctx *gctx2d.Context, time, scale float32) {
	ctx.SetGlobalAlpha(m.Alpha)
	renders.Mobs[m.MobID].Active(ctx, time, scale, m.X, m.Y, m.Size, m.Rotation, m.ID, m.Alpha, m.DamagedFlashCurrent, m.Attacking)
}

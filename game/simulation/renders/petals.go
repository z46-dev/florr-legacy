package renders

import (
	"math"

	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/gctx2d"
)

var baseIconFunc IconDrawFunc = func(self *PetalRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var (
		width, height float32        = float32(canvas.Width), float32(canvas.Height)
		cfg           *configs.Petal = configs.Petals[self.ID]
	)

	if cfg == nil || cfg.Special.Count <= 1 {
		self.StaticGenerator(&self.baseRender, canvas, ctx)
		return
	}

	for i := cfg.Special.Count - 1; i >= 0; i-- {
		var (
			angle float32 = TAU / float32(cfg.Special.Count) * float32(i)
			dist  float32 = 0.7 * width / 2
		)

		ctx.Save()
		ctx.SetObjectTransform(width/2+dist*float32(math.Cos(float64(angle))), height/2+dist*float32(math.Sin(float64(angle))), width/2, angle)
		// self.StaticGenerator(&self.baseRender, canvas, ctx)
		self.Active(ctx, 0, 1, 0, 0, 1, 0, 0, 1, 0, false)
		ctx.Restore()
	}
}

var PetalBasic *PetalRender = newPetal(configs.PetalID_Basic, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	ctx.SetFillStyle(DamageFlashColor(gctx2d.ColorWhite, damageFlashCurrent))
	ctx.FillCircle(x, y, size)

	ctx.SetLineWidth(.25 * size)
	ctx.SetStrokeStyle(render.MixColors(DamageFlashColor(gctx2d.ColorWhite, damageFlashCurrent), gctx2d.ColorBlack, .2))
	ctx.StrokeCircle(x, y, size)
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
}, baseIconFunc, 1)

var PetalFast *PetalRender = PetalBasic.Clone(configs.PetalID_Fast)

var heavyColor = gctx2d.RGBA(187.0/255.0, 187.0/255.0, 187.0/255.0, 1)
var PetalHeavy *PetalRender = newPetal(configs.PetalID_Heavy, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	ctx.SetFillStyle(DamageFlashColor(heavyColor, damageFlashCurrent))
	ctx.FillCircle(x, y, size)

	ctx.SetLineWidth(.25 * size)
	ctx.SetStrokeStyle(render.MixColors(DamageFlashColor(heavyColor, damageFlashCurrent), gctx2d.ColorBlack, .2))
	ctx.StrokeCircle(x, y, size)
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
}, baseIconFunc, 1)

var rockColor = gctx2d.RGBA(119.0/255.0, 119.0/255.0, 119.0/255.0, 1)
var PetalRock *PetalRender = newPetal(configs.PetalID_Rock, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	polygon(ctx, 5, x, y, size, rotation)
	fillStroke(ctx, DamageFlashColor(rockColor, damageFlashCurrent), .25*size)
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
}, baseIconFunc, 1)

var leafColor = render.HexToColor("#48b63e")
var PetalLeaf *PetalRender = newPetal(configs.PetalID_Leaf, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	var (
		fillColor   gctx2d.Color = DamageFlashColor(leafColor, damageFlashCurrent)
		strokeColor gctx2d.Color = render.MixColors(fillColor, gctx2d.ColorBlack, 0.2)
	)

	ctx.Save()
	ctx.SetObjectTransform(x, y, size*1.4285, rotation)
	ctx.SetLineWidth(0.157895)
	ctx.SetLineCap(gctx2d.LineCapRound)

	ctx.BeginPath()
	ctx.MoveTo(-0.657895, 0)
	ctx.LineTo(-0.921053, 0)
	ctx.SetStrokeStyle(strokeColor)
	ctx.Stroke()

	ctx.BeginPath()
	ctx.MoveTo(-0.657895, 0)
	ctx.BezierCurveTo(-0.394737, -0.631579, 0.394737, -0.631579, 0.921053, 0)
	ctx.BezierCurveTo(0.394737, 0.631579, -0.394737, 0.631579, -0.657895, 0)
	ctx.ClosePath()
	ctx.SetFillStyle(fillColor)
	ctx.SetStrokeStyle(strokeColor)
	ctx.Fill()
	ctx.Stroke()

	ctx.BeginPath()
	ctx.MoveTo(-0.342105, 0)
	ctx.QuadraticCurveTo(0.131579, -0.078947, 0.526316, 0)
	ctx.SetStrokeStyle(strokeColor)
	ctx.Stroke()
	ctx.Restore()
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width*.35, 0, 0, 1, 0, false)
}, baseIconFunc, 2)

var roseColor = render.HexToColor("#f993cb")
var PetalRose *PetalRender = newPetal(configs.PetalID_Rose, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	ctx.SetFillStyle(DamageFlashColor(roseColor, damageFlashCurrent))
	ctx.FillCircle(x, y, size)

	ctx.SetLineWidth(.25 * size)
	ctx.SetStrokeStyle(render.MixColors(DamageFlashColor(roseColor, damageFlashCurrent), gctx2d.ColorBlack, .2))
	ctx.StrokeCircle(x, y, size)
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
}, baseIconFunc, 1)

var fasterColor = render.HexToColor("#fdffc5")
var PetalFaster *PetalRender = newPetal(configs.PetalID_Faster, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	ctx.SetFillStyle(DamageFlashColor(fasterColor, damageFlashCurrent))
	ctx.FillCircle(x, y, size)

	ctx.SetLineWidth(.25 * size)
	ctx.SetStrokeStyle(render.MixColors(DamageFlashColor(fasterColor, damageFlashCurrent), gctx2d.ColorBlack, .2))
	ctx.StrokeCircle(x, y, size)
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
}, baseIconFunc, 1)

var PetalTwin *PetalRender = PetalBasic.Clone(configs.PetalID_Twin)
var PetalTriplet *PetalRender = PetalBasic.Clone(configs.PetalID_Triplet)

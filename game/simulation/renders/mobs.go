package renders

import (
	"math"

	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/gctx2d"
)

var MobRock *MobRender = newMob(configs.MobID_Rock, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	var sides int = 5 + int(id%4)
	ctx.BeginPath()

	var sub float32 = TAU / float32(sides)
	for i := range sides {
		var (
			angle, distance float32 = sub*float32(i) + rotation, size * (0.9 + float32((id+uint64(i))%3)*0.1)
			xx, yy          float32 = float32(math.Cos(float64(angle))) * distance, float32(math.Sin(float64(angle))) * distance
		)

		if i == 0 {
			ctx.MoveTo(x+xx, y+yy)
		} else {
			ctx.LineTo(x+xx, y+yy)
		}
	}

	ctx.ClosePath()
	fillStroke(ctx, DamageFlashColor(rockColor, damageFlashCurrent), 4/scale)
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
})

var MobBoulder *MobRender = newMob(configs.MobID_Boulder, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	var sides int = 8 + int(id%6)
	ctx.BeginPath()

	var sub float32 = TAU / float32(sides)
	for i := range sides {
		var (
			angle, distance float32 = sub*float32(i) + rotation, size * (0.95 + float32((id+uint64(i))%5)*0.025)
			xx, yy          float32 = float32(math.Cos(float64(angle))) * distance, float32(math.Sin(float64(angle))) * distance
		)

		if i == 0 {
			ctx.MoveTo(x+xx, y+yy)
		} else {
			ctx.LineTo(x+xx, y+yy)
		}
	}

	ctx.ClosePath()
	fillStroke(ctx, DamageFlashColor(rockColor, damageFlashCurrent), 4/scale)
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
})

var ladybugRed, ladybugBlack = gctx2d.RGBA(1, 4.0/255.0, 3.0/255.0, 1), gctx2d.RGBA(0.05, 0.05, 0.05, 1)
var MobLadybug *MobRender = newMob(configs.MobID_Ladybug, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	var (
		rarity      float32 = 0
		strokeWidth float32 = 0.2
		spotSeed    float64 = float64(id % 100000)
		red, black          = DamageFlashColor(ladybugRed, damageFlashCurrent), DamageFlashColor(ladybugBlack, damageFlashCurrent)
	)

	ctx.Save()
	ctx.SetObjectTransform(x, y, size, rotation)

	ctx.BeginPath()
	ctx.Arc(0.45, 0, 0.55, 0, TAU, false)
	ctx.SetFillStyle(black)
	ctx.Fill()
	ctx.SetLineWidth(strokeWidth)
	ctx.SetStrokeStyle(render.MixColors(black, gctx2d.ColorBlack, 0.2))
	ctx.Stroke()

	ctx.SetFillStyle(red)
	ctx.BeginPath()
	ctx.Arc(0, 0, 1, 0.575, TAU-0.575, false)
	ctx.QuadraticCurveTo(1.0/3.0, 0, float32(math.Cos(float64(0.575))), float32(math.Sin(float64(0.575))))
	ctx.ClosePath()
	ctx.Fill()

	ctx.Save()
	ctx.Clip()

	ctx.SetFillStyle(black)

	var amount int = 2 + int(math.Abs(math.Sin(spotSeed)*0.667)*12)
	amount += int(math.Floor(float64(rarity / 1.5)))

	for i := range amount {
		var (
			i64    float64 = float64(i)
			angle  float32 = float32(math.Sin(i64*100+spotSeed)) * TAU
			radius float32 = float32(math.Sin(i64*1000+spotSeed))*0.2 + 0.2
			d      float32 = float32(math.Sin(i64*10000+spotSeed))*0.3 + 0.7
		)

		radius += min(10, rarity) / 100

		ctx.BeginPath()
		ctx.Arc(float32(math.Cos(float64(angle)))*d, float32(math.Sin(float64(angle)))*d, radius, 0, TAU, false)
		ctx.Fill()
	}

	ctx.Restore()

	ctx.SetLineWidth(strokeWidth)
	ctx.SetStrokeStyle(render.MixColors(red, gctx2d.ColorBlack, 0.2))
	ctx.BeginPath()
	ctx.Arc(0, 0, 1, 0.575, TAU-0.575, false)
	ctx.QuadraticCurveTo(1.0/3.0, 0, float32(math.Cos(float64(0.575))), float32(math.Sin(float64(0.575))))
	ctx.ClosePath()
	ctx.Stroke()

	ctx.Restore()
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
})

var MobBee *MobRender = newMob(configs.MobID_Bee, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	var (
		black            gctx2d.Color = DamageFlashColor(gctx2d.RGBA(51.0/255.0, 51.0/255.0, 51.0/255.0, 1), damageFlashCurrent)
		bodyColor        gctx2d.Color = DamageFlashColor(gctx2d.RGBA(231.0/255.0, 231.0/255.0, 99.0/255.0, 1), damageFlashCurrent)
		bodyRadiusX      float32      = 1
		bodyRadiusY      float32      = 2.0 / 3.0
		bodyStrokeWidth  float32      = 1.0 / 6.0
		stingerTipX      float32      = -37.0 / 30.0
		stingerBaseX     float32      = -5.0 / 6.0
		stingerBaseY     float32      = 3.0 / 10.0
		stripeWidth      float32      = 1.0 / 3.0
		antennaStartX    float32      = 5.0 / 6.0
		antennaControlX  float32      = 7.0 / 6.0
		antennaEndX      float32      = 4.0 / 3.0
		antennaStartY    float32      = 1.0 / 6.0
		antennaEndY      float32      = 1.0 / 2.0
		antennaRadius    float32      = 1.0 / 6.0
		antennaLineWidth float32      = 1.0 / 10.0
		drawBody         func()
	)

	drawBody = func() {
		ctx.BeginPath()
		ctx.Ellipse(0, 0, bodyRadiusX, bodyRadiusY, 0, 0, TAU, false)
	}

	ctx.Save()
	ctx.SetObjectTransform(x, y, size, rotation)
	ctx.SetLineCap(gctx2d.LineCapRound)
	ctx.SetLineJoin(gctx2d.LineJoinRound)

	ctx.SetFillStyle(black)
	ctx.SetStrokeStyle(black)
	ctx.SetLineWidth(bodyStrokeWidth)
	ctx.BeginPath()
	ctx.MoveTo(stingerTipX, 0)
	ctx.LineTo(stingerBaseX, -stingerBaseY)
	ctx.LineTo(stingerBaseX, stingerBaseY)
	ctx.ClosePath()
	ctx.Fill()
	ctx.Stroke()

	ctx.SetFillStyle(bodyColor)
	drawBody()
	ctx.Fill()

	ctx.Save()
	ctx.Clip()
	ctx.SetFillStyle(black)
	for _, stripeX := range [...]float32{1.0 / 3.0, -1.0 / 3.0, -1} {
		ctx.BeginPath()
		ctx.Rect(stripeX, -bodyRadiusY, stripeWidth, bodyRadiusY*2)
		ctx.Fill()
	}

	ctx.Restore()

	ctx.SetStrokeStyle(render.MixColors(bodyColor, gctx2d.ColorBlack, 0.2))
	ctx.SetLineWidth(bodyStrokeWidth)
	drawBody()
	ctx.Stroke()

	ctx.SetFillStyle(black)
	ctx.SetStrokeStyle(black)
	ctx.SetLineWidth(antennaLineWidth)
	for direction := float32(-1); direction <= 1; direction += 2 {
		ctx.BeginPath()
		ctx.MoveTo(antennaStartX, antennaStartY*direction)
		ctx.QuadraticCurveTo(antennaControlX, antennaStartY*direction, antennaEndX, antennaEndY*direction)
		ctx.Stroke()

		ctx.BeginPath()
		ctx.Arc(antennaEndX, antennaEndY*direction, antennaRadius, 0, TAU, false)
		ctx.Fill()
	}

	ctx.Restore()
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/3.2, 0, 0, 1, 0, false)
})

var antBody, antMandible = render.HexToColor("#555555"), render.HexToColor("#292929")
var MobBabyAnt *MobRender = newMob(configs.MobID_BabyAnt, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	ctx.Save()
	ctx.SetObjectTransform(x, y, size, rotation)

	ctx.SetStrokeStyle(DamageFlashColor(antMandible, damageFlashCurrent))
	ctx.SetLineWidth(0.5)

	// const rotation = Math.sin(date * rotSpd + id / 4) * 0.075;
	var timeRatio float32 = 1.5
	if attacking {
		timeRatio = 4
	}

	var rot float32 = float32(math.Sin(float64(time*timeRatio+float32(id)/4)) * 0.075)

	ctx.BeginPath()
	ctx.MoveTo(0, -0.5)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.QuadraticCurveTo(0.785714, -0.714286, 1.57143, -0.357143)
	ctx.SetObjectTransform(0, 0, 1, rot)

	ctx.MoveTo(-0, 0.5)
	ctx.SetObjectTransform(0, 0, 1, rot)
	ctx.QuadraticCurveTo(0.785714, 0.714286, 1.57143, 0.357143)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.Stroke()

	ctx.BeginPath()
	ctx.Arc(0, 0, 1, 0, TAU, false)
	fillStroke(ctx, DamageFlashColor(antBody, damageFlashCurrent), 0.5)

	ctx.Restore()
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
})

var antWing = render.HexToColor("#EEEEEE")
var MobWorkerAnt *MobRender = newMob(configs.MobID_WorkerAnt, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	/*// mob: Worker Ant (registry 7)
	// source: live gameplay mob table[207]; coordinates are renderer-local
	// original bounds: (-21.5, -17.5) to (33.5, 17.5)
	// normalized with center (6, 0) and scale 0.036364
	function draw(ctx) {
	  const baseTransform = ctx.getTransform();
	  const normalization = [0.036364, 0, 0, 0.036364, -0.218182, 0];
	  const setLocalTransform = ([a, b, c, d, e, f]) => ctx.setTransform(
	    baseTransform.a * a + baseTransform.c * b,
	    baseTransform.b * a + baseTransform.d * b,
	    baseTransform.a * c + baseTransform.c * d,
	    baseTransform.b * c + baseTransform.d * d,
	    baseTransform.a * e + baseTransform.c * f + baseTransform.e,
	    baseTransform.b * e + baseTransform.d * f + baseTransform.f,
	  );
	  const setObjectTransform = (x, y, width, height, rotation) => {
	    setLocalTransform(normalization);
	    ctx.translate(x, y);
	    ctx.rotate(rotation);
	    ctx.scale(width, height);
	  };
	  setLocalTransform(normalization);
	  ctx.lineWidth = 7;
	  ctx.lineJoin = "round";
	  ctx.lineCap = "round";
	  ctx.fillStyle = "#555555";
	  ctx.strokeStyle = "#454545";
	  ctx.beginPath();
	  ctx.arc(-8, 0, 10, 0, TAU, false);
	  ctx.fill();
	  ctx.stroke();
	  ctx.strokeStyle = "#292929";
	  ctx.beginPath();
	  setObjectTransform(8, 0, 1, 1, 0);
	  ctx.moveTo(0, -7);
	  ctx.quadraticCurveTo(11, -9.563633, 22, -4.563632);
	  ctx.stroke();
	  ctx.beginPath();
	  ctx.moveTo(0, 7);
	  ctx.quadraticCurveTo(11, 9.563633, 22, 4.563632);
	  ctx.stroke();
	  ctx.fillStyle = "#555555";
	  ctx.strokeStyle = "#454545";
	  ctx.beginPath();
	  ctx.arc(0, 0, 14, 0, TAU, false);
	  ctx.fill();
	  ctx.stroke();
	}*/

	ctx.Save()
	ctx.SetObjectTransform(x, y, size, rotation)

	ctx.BeginPath()
	ctx.Arc(-0.571429, 0, 0.714286, 0, TAU, false)
	fillStroke(ctx, DamageFlashColor(antBody, damageFlashCurrent), 0.5)

	ctx.SetStrokeStyle(DamageFlashColor(antMandible, damageFlashCurrent))
	var timeRatio float32 = 1.5
	if attacking {
		timeRatio = 4
	}

	ctx.SetObjectTransform(0.571429, 0, 1, 0)

	var rot float32 = float32(math.Sin(float64(time*timeRatio+float32(id)/4)) * 0.075)

	ctx.BeginPath()
	ctx.MoveTo(0, -0.5)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.QuadraticCurveTo(0.785714, -0.714286, 1.57143, -0.357143)
	ctx.SetObjectTransform(0, 0, 1, rot)

	ctx.MoveTo(-0, 0.5)
	ctx.SetObjectTransform(0, 0, 1, rot)
	ctx.QuadraticCurveTo(0.785714, 0.714286, 1.57143, 0.357143)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.Stroke()

	ctx.BeginPath()
	ctx.Arc(0, 0, 1, 0, TAU, false)
	fillStroke(ctx, DamageFlashColor(antBody, damageFlashCurrent), 0.5)

	ctx.Restore()
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
})

var MobSoldierAnt *MobRender = newMob(configs.MobID_SoldierAnt, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	ctx.Save()
	ctx.SetObjectTransform(x, y, size, rotation)

	ctx.BeginPath()
	ctx.Arc(-0.571429, 0, 0.714286, 0, TAU, false)
	fillStroke(ctx, DamageFlashColor(antBody, damageFlashCurrent), 0.5)

	var timeRatio float32 = 3
	if attacking {
		timeRatio = 6
	}

	var rot float32 = float32(math.Sin(float64(time*timeRatio+float32(id)/4)) * 0.15)

	ctx.SetGlobalAlpha(0.3 * currentAlpha)
	ctx.SetFillStyle(DamageFlashColor(antWing, damageFlashCurrent))
	ctx.BeginPath()
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.Ellipse(-0.5, -0.571429, 1.07143, 0.5, 0.0224399, 0, TAU, false)
	ctx.SetObjectTransform(0, 0, 1, rot)
	ctx.Fill()
	ctx.BeginPath()
	ctx.SetObjectTransform(0, 0, 1, rot)
	ctx.Ellipse(-0.5, 0.571429, 1.07143, 0.5, -0.0224399, 0, TAU, false)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.Fill()
	ctx.SetGlobalAlpha(currentAlpha)

	ctx.SetObjectTransform(0.571429, 0, 1, 0)

	timeRatio = 1.5
	if attacking {
		timeRatio = 4
	}

	rot = float32(math.Sin(float64(time*timeRatio+float32(id)/4)) * 0.075)

	ctx.SetStrokeStyle(DamageFlashColor(antMandible, damageFlashCurrent))
	ctx.SetLineWidth(0.5)
	ctx.BeginPath()
	ctx.MoveTo(0, -0.5)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.QuadraticCurveTo(0.785714, -0.714286, 1.57143, -0.357143)
	ctx.SetObjectTransform(0, 0, 1, rot)

	ctx.MoveTo(-0, 0.5)
	ctx.SetObjectTransform(0, 0, 1, rot)
	ctx.QuadraticCurveTo(0.785714, 0.714286, 1.57143, 0.357143)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.Stroke()

	ctx.BeginPath()
	ctx.Arc(0, 0, 1, 0, TAU, false)
	fillStroke(ctx, DamageFlashColor(antBody, damageFlashCurrent), 0.5)

	ctx.Restore()
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
})

var MobQueenAnt *MobRender = newMob(configs.MobID_QueenAnt, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	/*// mob: Queen Ant (registry 9)
	// source: synthetic mob table[211]; coordinates are renderer-local
	// renderer transforms are baked into the affected geometry and stroke widths
	// original bounds: (-29.25, -17.862121) to (30.75, 17.862121)
	// normalized with center (0.75, 0) and scale 0.033333
	function draw(ctx) {
	  const baseTransform = ctx.getTransform();
	  const normalization = [0.033333, 0, 0, 0.033333, -0.025, 0];
	  const setLocalTransform = ([a, b, c, d, e, f]) => ctx.setTransform(
	    baseTransform.a * a + baseTransform.c * b,
	    baseTransform.b * a + baseTransform.d * b,
	    baseTransform.a * c + baseTransform.c * d,
	    baseTransform.b * c + baseTransform.d * d,
	    baseTransform.a * e + baseTransform.c * f + baseTransform.e,
	    baseTransform.b * e + baseTransform.d * f + baseTransform.f,
	  );
	  setLocalTransform(normalization);
	  ctx.lineCap = "round";
	  ctx.fillStyle = "#555555";
	  ctx.strokeStyle = "#454545";
	  ctx.beginPath();
	  ctx.arc(-12.5, 0, 15, 0, TAU, false);
	  ctx.fill();
	  ctx.lineWidth = 3.5;
	  ctx.stroke();
	  ctx.fillStyle = "#555555";
	  ctx.strokeStyle = "#454545";
	  ctx.beginPath();
	  ctx.arc(0, 0, 12.5, 0, TAU, false);
	  ctx.fill();
	  ctx.stroke();
	  ctx.save();
	    // ctx.globalAlpha is runtime-state-dependent (raw probe value 0)
	    ctx.fillStyle = "#EEEEEE";
	    ctx.strokeStyle = "#EEEEEE";
	    ctx.beginPath();
	    ctx.ellipse(-7, -8, 15, 7, 0.314159, 0, TAU, false);
	    ctx.fill();
	    ctx.beginPath();
	    ctx.ellipse(-7, 8, 15, 7, -0.314159, 0, TAU, false);
	    ctx.fill();
	    // ctx.globalAlpha is runtime-state-dependent (raw probe value 0)
	  ctx.restore();
	  ctx.strokeStyle = "#292929";
	  ctx.beginPath();
	  ctx.moveTo(12.5, -5.25);
	  ctx.quadraticCurveTo(20.75, -7.5, 29, -3.75);
	  ctx.lineWidth = 3.5;
	  ctx.stroke();
	  ctx.beginPath();
	  ctx.moveTo(12.5, 5.25);
	  ctx.quadraticCurveTo(20.75, 7.5, 29, 3.75);
	  ctx.stroke();
	  ctx.fillStyle = "#555555";
	  ctx.strokeStyle = "#454545";
	  ctx.beginPath();
	  ctx.arc(12.5, 0, 10.5, 0, TAU, false);
	  ctx.fill();
	  ctx.stroke();
	}*/

	ctx.Save()
	ctx.SetObjectTransform(x, y, size, rotation)

	ctx.BeginPath()
	ctx.Arc(-12.5/10.5, 0, 15/10.5, 0, TAU, false)
	fillStroke(ctx, DamageFlashColor(antBody, damageFlashCurrent), 3.5/10.5)

	ctx.BeginPath()
	ctx.Arc(0, 0, 12.5/10.5, 0, TAU, false)
	fillStroke(ctx, DamageFlashColor(antBody, damageFlashCurrent), 3.5/10.5)

	var timeRatio float32 = 5
	if attacking {
		timeRatio = 15
	}

	var rot float32 = float32(math.Sin(float64(time*timeRatio+float32(id)/4)) * 0.075)

	ctx.SetGlobalAlpha(0.3 * currentAlpha)
	ctx.SetFillStyle(DamageFlashColor(antWing, damageFlashCurrent))
	ctx.BeginPath()
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.Ellipse(-7/10.5, -8/10.5, 15/10.5, 7/10.5, 0.314159/10.5, 0, TAU, false)
	ctx.SetObjectTransform(0, 0, 1, rot)
	ctx.Fill()

	ctx.BeginPath()
	ctx.SetObjectTransform(0, 0, 1, rot)
	ctx.Ellipse(-7/10.5, 8/10.5, 15/10.5, 7/10.5, -0.314159/10.5, 0, TAU, false)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.Fill()
	ctx.SetGlobalAlpha(currentAlpha)
	ctx.SetObjectTransform(12.5/10.5, 0, 1, 0)

	ctx.SetStrokeStyle(DamageFlashColor(antMandible, damageFlashCurrent))
	ctx.SetLineWidth(0.5)
	ctx.BeginPath()
	ctx.MoveTo(0, -0.5)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.QuadraticCurveTo(0.785714, -0.714286, 1.57143, -0.357143)
	ctx.SetObjectTransform(0, 0, 1, rot)

	ctx.MoveTo(-0, 0.5)
	ctx.SetObjectTransform(0, 0, 1, rot)
	ctx.QuadraticCurveTo(0.785714, 0.714286, 1.57143, 0.357143)
	ctx.SetObjectTransform(0, 0, 1, -rot)
	ctx.Stroke()

	ctx.BeginPath()
	ctx.Arc(0, 0, 1, 0, TAU, false)
	fillStroke(ctx, DamageFlashColor(antBody, damageFlashCurrent), 3.5/10.5)

	ctx.Restore()
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2.2, 0, 0, 1, 0, false)
})

var ah1, ah2, ah3 = render.HexToColor("#B58500"), render.HexToColor("#946D00"), render.HexToColor("#6B4F00")
var MobAntHole *MobRender = newMob(configs.MobID_AntHole, func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool) {
	ctx.BeginPath()
	ctx.Arc(x, y, size, 0, TAU, false)
	ctx.SetFillStyle(DamageFlashColor(ah1, damageFlashCurrent))
	ctx.Fill()

	ctx.BeginPath()
	ctx.Arc(x, y, size*0.6666667, 0, TAU, false)
	ctx.SetFillStyle(DamageFlashColor(ah2, damageFlashCurrent))
	ctx.Fill()

	ctx.BeginPath()
	ctx.Arc(x, y, size*0.3333333, 0, TAU, false)
	ctx.SetFillStyle(DamageFlashColor(ah3, damageFlashCurrent))
	ctx.Fill()
}, func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context) {
	var width, height float32 = float32(canvas.Width), float32(canvas.Height)
	self.Active(ctx, 0, 0, width/2, height/2, width/2, 0, 0, 1, 0, false)
})

package renders

import (
	"math"

	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/gctx2d"
)

const TAU float32 = 2 * math.Pi

func fillStroke(ctx *gctx2d.Context, fill gctx2d.Color, lineWidth float32) {
	ctx.SetFillStyle(fill)
	ctx.Fill()

	if lineWidth > 0 {
		ctx.SetLineWidth(lineWidth)
		ctx.SetStrokeStyle(render.MixColors(fill, gctx2d.ColorBlack, 0.2))
		ctx.Stroke()
	}
}

func polygon(ctx *gctx2d.Context, sides int, x, y, size, rotation float32) {
	ctx.BeginPath()

	for i := range sides {
		var (
			angle  float32 = float32(i)*TAU/float32(sides) + rotation
			xx, yy float32 = x + float32(math.Cos(float64(angle)))*size, y + float32(math.Sin(float64(angle)))*size
		)

		if i == 0 {
			ctx.MoveTo(xx, yy)
		} else {
			ctx.LineTo(xx, yy)
		}
	}

	ctx.ClosePath()
}

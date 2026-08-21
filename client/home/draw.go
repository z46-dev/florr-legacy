package home

import (
	"github.com/z46-dev/florr-legacy/game/simulation/renders"
	"github.com/z46-dev/gctx2d"
)

var (
	outsideColorDarkening float32      = 0.1
	worldColor            gctx2d.Color = gctx2d.RGBA(54.0/255.0, 168.0/255.0, 92.0/255.0, 1.0)
	gridColor             gctx2d.Color = gctx2d.RGBA(0.0, 0.0, 0.0, 0.1)
)

// drawBackground draws the archived green field and screenshot-measured particles.
func (h *Home) drawBackground(ctx *gctx2d.Context, width, height float32) {
	ctx.SetFillStyle(worldColor)
	ctx.BeginPath()
	ctx.Rect(0, 0, width, height)
	ctx.Fill()

	if h.petalAnimation.Ready {
		h.petalAnimation.Width, h.petalAnimation.Height = width, height

		for _, petal := range h.petalAnimation.Petals {
			var render *renders.PetalRender = renders.Petals[petal.ID]
			if render == nil || render.CachedImage == nil {
				continue
			}

			ctx.SetObjectTransform(petal.X, petal.Y, petal.S, petal.R)
			ctx.DrawImageBounds(render.CachedImage, -0.5, -0.5, 1, 1)
			ctx.ResetTransform()
		}

		ctx.SetGlobalAlpha(1)
		ctx.ResetTransform()
	}
}

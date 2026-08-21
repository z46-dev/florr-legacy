package simulation

import (
	"fmt"
	"math"
	"slices"

	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/gctx2d"
)

var (
	outsideColorDarkening float32      = 0.1
	worldColor            gctx2d.Color = gctx2d.RGBA(54.0/255.0, 168.0/255.0, 92.0/255.0, 1.0)
	gridColor             gctx2d.Color = gctx2d.RGBA(0.0, 0.0, 0.0, 0.1)
)

func gameScale(size, fov float32) float32 {
	return size / fov
}

// screenShakeOffset returns a small screen-space camera shake while the local flower is damaged.
func (s *Simulation) screenShakeOffset() (x, y float32) {
	if s.Player == nil || s.Player.IsDead || (!s.Player.Damaged && s.Player.LastDamagedAt <= 0) {
		return
	}

	var strength float32 = 1
	if s.Player.LastDamagedAt > 0 {
		strength = max(0.2, s.Player.LastDamagedAt/HEALTH_DELTA_DURATION)
	}

	x = float32(math.Sin(float64(s.Time*91.0))) * 5 * strength
	y = float32(math.Sin(float64(s.Time*137.0+1.7))) * 4 * strength
	return
}

func (s *Simulation) Draw(ctx *gctx2d.Context, width, height float32) {
	if ctx == nil {
		return
	}

	s.Flowers.Lock.RLock()
	s.Petals.Lock.RLock()
	s.Mobs.Lock.RLock()
	defer s.Mobs.Lock.RUnlock()
	defer s.Petals.Lock.RUnlock()
	defer s.Flowers.Lock.RUnlock()

	s.Width, s.Height = width, height

	var (
		scale          float32 = gameScale(max(width, height), s.Camera.FieldOfView)
		shakeX, shakeY float32 = s.screenShakeOffset()
	)

	ctx.SetFillStyle(render.MixColors(worldColor, gctx2d.ColorBlack, outsideColorDarkening))
	ctx.BeginPath()
	ctx.Rect(0, 0, width, height)
	ctx.Fill()

	ctx.SetTransform(
		scale, 0,
		0, scale,
		width*0.5-float32(s.Camera.X)*scale+shakeX,
		height*0.5-float32(s.Camera.Y)*scale+shakeY,
	)

	// Draw the world boundary
	ctx.SetFillStyle(worldColor)
	ctx.BeginPath()
	ctx.Rect(-s.WorldWidth*0.5, -s.WorldHeight*0.5, s.WorldWidth, s.WorldHeight)
	ctx.Fill()

	{
		var (
			gridSpacing    float32 = 50
			halfViewWidth  float32 = width * 0.5 / scale
			halfViewHeight float32 = height * 0.5 / scale
			minX           float32 = s.Camera.X - halfViewWidth
			maxX           float32 = s.Camera.X + halfViewWidth
			minY           float32 = s.Camera.Y - halfViewHeight
			maxY           float32 = s.Camera.Y + halfViewHeight
			startX         float32 = float32(math.Floor(float64(minX/gridSpacing))) * gridSpacing
			startY         float32 = float32(math.Floor(float64(minY/gridSpacing))) * gridSpacing
		)

		ctx.SetStrokeStyle(gridColor)
		ctx.SetLineWidth(1 / scale)
		ctx.BeginPath()

		for x := startX; x <= maxX; x += gridSpacing {
			ctx.MoveTo(x, minY)
			ctx.LineTo(x, maxY)
		}

		for y := startY; y <= maxY; y += gridSpacing {
			ctx.MoveTo(minX, y)
			ctx.LineTo(maxX, y)
		}

		ctx.Stroke()
	}

	ctx.SetLineCap(gctx2d.LineCapRound)
	ctx.SetLineJoin(gctx2d.LineJoinRound)

	var (
		flowersOrder []uint64 = make([]uint64, 0, len(s.Flowers.Items))
		petalsOrder  []uint64 = make([]uint64, 0, len(s.Petals.Items))
		mobsOrder    []uint64 = make([]uint64, 0, len(s.Mobs.Items))
	)

	for _, flower := range s.Flowers.Items {
		flowersOrder = append(flowersOrder, flower.ID)
	}

	for _, petal := range s.Petals.Items {
		petalsOrder = append(petalsOrder, petal.ID)
	}

	for _, mob := range s.Mobs.Items {
		mobsOrder = append(mobsOrder, mob.ID)
	}

	slices.SortStableFunc(flowersOrder, func(a, b uint64) int {
		var (
			flowerA *Flower = s.Flowers.Items[a]
			flowerB *Flower = s.Flowers.Items[b]
		)

		if flowerA == nil || flowerB == nil {
			return 0
		}

		if flowerA.Size < flowerB.Size {
			return 1
		} else if flowerA.Size > flowerB.Size {
			return -1
		} else {
			if flowerA.ID < flowerB.ID {
				return -1
			} else if flowerA.ID > flowerB.ID {
				return 1
			} else {
				return 0
			}
		}
	})

	slices.SortStableFunc(petalsOrder, func(a, b uint64) int {
		var (
			petalA *Petal = s.Petals.Items[a]
			petalB *Petal = s.Petals.Items[b]
		)

		if petalA == nil || petalB == nil {
			return 0
		}

		if petalA.Size < petalB.Size {
			return 1
		} else if petalA.Size > petalB.Size {
			return -1
		} else {
			if petalA.ID < petalB.ID {
				return -1
			} else if petalA.ID > petalB.ID {
				return 1
			} else {
				return 0
			}
		}
	})

	slices.SortStableFunc(mobsOrder, func(a, b uint64) int {
		var (
			mobA *Mob = s.Mobs.Items[a]
			mobB *Mob = s.Mobs.Items[b]
		)

		if mobA == nil || mobB == nil {
			return 0
		}

		if mobA.Size < mobB.Size {
			return 1
		} else if mobA.Size > mobB.Size {
			return -1
		} else {
			if mobA.ID < mobB.ID {
				return -1
			} else if mobA.ID > mobB.ID {
				return 1
			} else {
				return 0
			}
		}
	})

	for _, id := range mobsOrder {
		var mob *Mob = s.Mobs.Items[id]
		if mob != nil {
			mob.DrawHealth(ctx, scale)
		}
	}

	for _, id := range flowersOrder {
		var flower *Flower = s.Flowers.Items[id]
		if flower != nil {
			flower.DrawHealth(ctx, scale)
		}
	}

	for _, id := range mobsOrder {
		var mob *Mob = s.Mobs.Items[id]
		if mob != nil {
			mob.Draw(ctx, s.Time, scale)
		}
	}

	for _, id := range petalsOrder {
		var petal *Petal = s.Petals.Items[id]
		if petal != nil {
			petal.Draw(ctx, s.Time, scale)
		}
	}

	for _, id := range flowersOrder {
		var flower *Flower = s.Flowers.Items[id]
		if flower != nil {
			flower.Draw(ctx)
		}
	}

	ctx.ResetTransform()
	ctx.SetGlobalAlpha(1)

	ctx.SetTextAlign(gctx2d.TextAlignLeft)
	ctx.SetTextBaseline(gctx2d.TextBaselineTop)
	render.Text(ctx, fmt.Sprintf("Client: FPS: %d, Flowers: %d, Petals: %d, Mobs: %d", s.RenderState.Renderer.GetFPS(), len(s.Flowers.Items), len(s.Petals.Items), len(s.Mobs.Items)), 5, 5, 16, nil, gctx2d.ColorWhite)
	render.Text(ctx, fmt.Sprintf("Server: MSPT: %.2f, FPS: %d, Flowers: %d, Petals: %d, Mobs: %d", s.serverStats.mspt, s.serverStats.fps, s.serverStats.numFlowers, s.serverStats.numPetals, s.serverStats.numMobs), 5, 25, 16, nil, gctx2d.ColorWhite)
}

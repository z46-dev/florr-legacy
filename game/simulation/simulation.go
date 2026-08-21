package simulation

import (
	"fmt"
	"time"

	"github.com/gogpu/gogpu"
	"github.com/gogpu/gpucontext"
	"github.com/gogpu/wgpu"
	"github.com/z46-dev/arctic"
	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/florr-legacy/shared/util"
	"github.com/z46-dev/golog"
)

func NewSimulation(renderer *render.Renderer, client *arctic.Client) (simulation *Simulation, err error) {
	if renderer == nil {
		err = fmt.Errorf("renderer is nil")
		return
	}

	simulation = &Simulation{
		WorldLerpAlpha: 1.0,
		Camera:         newCamera(),
		Flowers:        shared.NewCollection[*Flower](),
		Petals:         shared.NewCollection[*Petal](),
		Mobs:           shared.NewCollection[*Mob](),
		Client:         client,
		WorldWidth:     50000,
		WorldHeight:    3000,
		L:              golog.New().Prefix("[CLIENT]", golog.BoldPurple).Timestamp(),
	}

	if simulation.RenderState, err = renderer.NewState(); err != nil {
		simulation = nil
		return
	}

	simulation.RenderState.OnUpdate = func(dt float64) {
		simulation.Update(dt)
	}

	simulation.RenderState.Draw = func(dc *gogpu.Context) {
		var (
			width, height             int               = dc.Size()
			targetWidth, targetHeight uint32            = dc.SurfaceSize()
			surface                   *wgpu.TextureView = dc.SurfaceView()
		)

		if surface == nil {
			return
		}

		simulation.RenderState.Ctx.BeginTarget(width, height, int(targetWidth), int(targetHeight))
		simulation.Draw(simulation.RenderState.Ctx, float32(width), float32(height))
		simulation.RenderState.Ctx.Flush(surface, nil)
	}

	var (
		attack, defend bool    = false, false
		angle, mag     float64 = 0, 0
	)

	var sendInput = func() {
		var writer *shared.Writer = &shared.Writer{}
		writer.SetU8(shared.GameMsgTypeInput)
		writer.SetF32(float32(angle))
		writer.SetF32(float32(mag))

		switch {
		case attack:
			writer.SetU8(1)
		case defend:
			writer.SetU8(2)
		default:
			writer.SetU8(0)
		}

		simulation.Client.Send(writer.Bytes)
	}

	simulation.RenderState.OnMousePress = func(button gpucontext.MouseButton, x, y float64) {
		switch button {
		case gpucontext.MouseButtonLeft:
			attack = true
		case gpucontext.MouseButtonRight:
			defend = true
		}

		sendInput()
	}

	simulation.RenderState.OnMouseRelease = func(button gpucontext.MouseButton, x, y float64) {
		switch button {
		case gpucontext.MouseButtonLeft:
			attack = false
		case gpucontext.MouseButtonRight:
			defend = false
		}

		sendInput()
	}

	simulation.RenderState.OnMouseMove = func(x, y float64) {
		var (
			mouse, center *util.Vec2 = util.NewVec2(x, y), util.NewVec2(float64(simulation.Width)*0.5, float64(simulation.Height)*0.5)
			dist          float64
		)

		angle, dist = center.AngleTo(mouse), mouse.Dist(center)

		// 0 to 100px radius = 0 mag, 100px to 25% max(width, height) radius = 0..1 mag, > 25% max(width, height) radius = 1 mag
		var (
			q    float64 = max(center.X, center.Y)
			safe float64 = q * .005
		)

		if dist > safe {
			mag = min(1, (dist-safe)/(0.25*q-safe))
		} else {
			mag = 0
		}

		sendInput()
	}

	return
}

func (s *Simulation) Update(deltaSeconds float64) {
	s.Time += float32(deltaSeconds)
	if !s.LastWorldStateAt.IsZero() {
		s.WorldLerpAlpha = min(1, float32(time.Since(s.LastWorldStateAt))/float32(time.Second/20))
	}

	// ==================================================
	// ================= Update  Camera =================
	// ==================================================

	s.Camera.Update(s.WorldLerpAlpha)

	// ==================================================
	// ================= Update Flowers =================
	// ==================================================

	s.Flowers.Lock.Lock()

	for _, toAppend := range s.Flowers.ToAppend {
		s.Flowers.Items[toAppend.ID] = toAppend
	}

	for _, toRemove := range s.Flowers.ToRemove {
		delete(s.Flowers.Items, toRemove)
	}

	s.Flowers.ToAppend = s.Flowers.ToAppend[:0]
	s.Flowers.ToRemove = s.Flowers.ToRemove[:0]

	for _, flower := range s.Flowers.Items {
		flower.Update(s.WorldLerpAlpha, deltaSeconds)
		if flower.DeathDelta >= DEATH_DELTA_DURATION {
			s.Flowers.Remove(flower.ID)
		}
	}

	s.Player = s.Flowers.Items[s.PlayerID]

	s.Flowers.Lock.Unlock()

	// ==================================================
	// ================= Update  Petals =================
	// ==================================================

	s.Petals.Lock.Lock()

	for _, toAppend := range s.Petals.ToAppend {
		s.Petals.Items[toAppend.ID] = toAppend
	}

	for _, toRemove := range s.Petals.ToRemove {
		delete(s.Petals.Items, toRemove)
	}

	s.Petals.ToAppend = s.Petals.ToAppend[:0]
	s.Petals.ToRemove = s.Petals.ToRemove[:0]

	for _, petal := range s.Petals.Items {
		petal.Update(s.WorldLerpAlpha, deltaSeconds)
		if petal.DeathDelta >= DEATH_DELTA_DURATION {
			s.Petals.Remove(petal.ID)
		}
	}

	s.Petals.Lock.Unlock()

	// ==================================================
	// ================== Update  Mobs ==================
	// ==================================================

	s.Mobs.Lock.Lock()

	for _, toAppend := range s.Mobs.ToAppend {
		s.Mobs.Items[toAppend.ID] = toAppend
	}

	for _, toRemove := range s.Mobs.ToRemove {
		delete(s.Mobs.Items, toRemove)
	}

	s.Mobs.ToAppend = s.Mobs.ToAppend[:0]
	s.Mobs.ToRemove = s.Mobs.ToRemove[:0]

	for _, mob := range s.Mobs.Items {
		mob.Update(s.WorldLerpAlpha, deltaSeconds)
		if mob.DeathDelta >= DEATH_DELTA_DURATION {
			s.Mobs.Remove(mob.ID)
		}
	}

	s.Mobs.Lock.Unlock()
}

func (s *Simulation) Release() {
	if s.RenderState != nil {
		s.RenderState.Release()
		s.RenderState = nil
	}
}

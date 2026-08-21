package renders

import (
	"fmt"

	"github.com/gogpu/gputypes"
	"github.com/gogpu/wgpu"
	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/gctx2d"
	"github.com/z46-dev/gctx2d/gimage"
)

const (
	STATIC_CACHE_SIZE = 128
)

type (
	ActiveDrawFunc func(ctx *gctx2d.Context, time, scale, x, y, size, rotation float32, id uint64, currentAlpha, damageFlashCurrent float32, attacking bool)
	StaticDrawFunc func(self *baseRender, canvas *gctx2d.OffscreenCanvas, ctx *gctx2d.Context)
	// Petal          struct {
	// 	ID              configs.PetalID
	// 	Active          ActiveDrawFunc
	// 	StaticGenerator StaticDrawFunc
	// 	StaticCached    *gctx2d.OffscreenCanvas
	// 	CachedImage     *gimage.Image
	// }
	// Mob struct {
	// 	ID              configs.MobID
	// 	Active          ActiveDrawFunc
	// 	StaticGenerator StaticDrawFunc
	// 	StaticCached    *gctx2d.OffscreenCanvas
	// 	CachedImage     *gimage.Image
	// }

	baseRender struct {
		Active          ActiveDrawFunc
		StaticGenerator StaticDrawFunc
		StaticCached    *gctx2d.OffscreenCanvas
		CachedImage     *gimage.Image
	}

	PetalRender struct {
		ID configs.PetalID
		baseRender
	}

	MobRender struct {
		ID configs.MobID
		baseRender
	}
)

// ==================================================
// ====================== Base ======================
// ==================================================

func (br *baseRender) Generate(device *wgpu.Device, round bool) (err error) {
	if br.StaticCached, err = gctx2d.NewOffscreenCanvas(device, STATIC_CACHE_SIZE, STATIC_CACHE_SIZE); err != nil {
		return
	}

	var ctx *gctx2d.Context
	if ctx, err = br.StaticCached.GetContext(); err != nil {
		return
	}

	ctx.Begin(br.StaticCached.Width, br.StaticCached.Height)

	if round {
		ctx.SetLineCap(gctx2d.LineCapRound)
		ctx.SetLineJoin(gctx2d.LineJoinRound)
	} else {
		ctx.SetLineCap(gctx2d.LineCapButt)
		ctx.SetLineJoin(gctx2d.LineJoinBevel)
	}

	br.StaticGenerator(br, br.StaticCached, ctx)

	if err = br.StaticCached.Flush(&gputypes.Color{R: 0, G: 0, B: 0, A: 0}); err != nil {
		return
	}

	br.CachedImage = br.StaticCached.Image("cache")
	return
}

// ==================================================
// ===================== Petals =====================
// ==================================================

func newPetal(id configs.PetalID, active ActiveDrawFunc, static StaticDrawFunc) (p *PetalRender) {
	p = &PetalRender{ID: id}
	p.Active = active
	p.StaticGenerator = static
	Petals[id] = p
	return
}

func (p *PetalRender) Clone(id configs.PetalID) (clone *PetalRender) {
	clone = &PetalRender{ID: id}
	clone.Active = p.Active
	clone.StaticGenerator = p.StaticGenerator
	Petals[id] = clone
	return
}

// ==================================================
// ====================== Mobs ======================
// ==================================================

func newMob(id configs.MobID, active ActiveDrawFunc, static StaticDrawFunc) (m *MobRender) {
	m = &MobRender{ID: id}
	m.Active = active
	m.StaticGenerator = static
	Mobs[id] = m
	return
}

func (m *MobRender) Clone(id configs.MobID) (clone *MobRender) {
	clone = &MobRender{ID: id}
	clone.Active = m.Active
	clone.StaticGenerator = m.StaticGenerator
	Mobs[id] = clone
	return
}

var (
	Petals map[configs.PetalID]*PetalRender = make(map[configs.PetalID]*PetalRender)
	Mobs   map[configs.MobID]*MobRender     = make(map[configs.MobID]*MobRender)
)

func Init(device *wgpu.Device, round bool) (err error) {
	for _, petal := range Petals {
		if err = petal.Generate(device, round); err != nil {
			err = fmt.Errorf("Failed to generate petal %s[%d]: %w", configs.Petals[petal.ID].Meta.Name, petal.ID, err)
			return
		}
	}

	for _, mob := range Mobs {
		if err = mob.Generate(device, round); err != nil {
			err = fmt.Errorf("Failed to generate mob %s[%d]: %w", configs.Mobs[mob.ID].Meta.Name, mob.ID, err)
			return
		}
	}

	return
}

var damageFlashRed = gctx2d.RGBA(1, 0.2, 0.2, 1)

func DamageFlashColor(base gctx2d.Color, progress float32) gctx2d.Color {
	if progress <= 0 {
		return base
	}

	if progress > 2.0/3.0 {
		return render.MixColors(base, gctx2d.ColorWhite, 0.5)
	}

	if progress > 1.0/3.0 {
		return render.MixColors(base, damageFlashRed, 0.5)
	}

	return render.MixColors(base, damageFlashRed, 0.25)
}

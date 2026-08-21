package home

import (
	"math"
	"math/rand/v2"

	"github.com/gogpu/gogpu"
	"github.com/gogpu/gpucontext"
	"github.com/gogpu/wgpu"
	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/shared"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/gctx2d"
	"github.com/z46-dev/spiffy"
)

func (h *Home) newLabel(text string, size float64) (label *spiffy.Label) {
	label = h.UI.NewLabel(spiffy.Rect{}, text)
	label.TextStyle.Color = gctx2d.ColorWhite
	label.TextStyle.FontSize = size
	label.TextStyle.Align = gctx2d.TextAlignCenter
	label.TextStyle.Baseline = gctx2d.TextBaselineMiddle
	label.TextStyle.BorderColor = gctx2d.ColorBlack
	label.TextStyle.Font = h.RenderState.Renderer.LegacyFont
	label.TextStyle.BorderWidth = float32(size) / 15
	return
}

func (h *Home) newButton(text string) (button *spiffy.Button) {
	button = h.UI.NewButton(spiffy.Rect{}, text, nil)
	button.Style.FillColor = gctx2d.RGBA(90.0/255.0, 159.0/255.0, 219.0/255.0, 1)
	button.Style.BorderColor = render.MixColors(button.Style.FillColor, gctx2d.ColorBlack, 0.2)
	button.Style.BorderWidth = 3
	button.Style.BorderRadius = 3
	button.Style.Click.Type = spiffy.ClickAnimationNone
	button.Style.HoverMixAmount = .1
	button.Style.PressMixAmount = .18
	button.TextStyle.Color = gctx2d.ColorWhite
	button.TextStyle.Font = h.RenderState.Renderer.LegacyFont
	button.TextStyle.FontSize = 11
	button.TextStyle.BorderColor = gctx2d.ColorBlack
	button.TextStyle.BorderWidth = 1
	return
}

// buildUI creates the archived screen controls through Spiffy.
func (h *Home) buildUI() {
	h.UI = spiffy.New(h.RenderState.Ctx)
	h.UI.SetClipboard(h.RenderState.Renderer.App)
	h.title = h.newLabel("florr.io", 56)

	h.input.prompt = h.newLabel("This pretty little flower is called...", 14)
	h.input.submitHint = h.newLabel("(press enter to spawn)", 10)
	h.input.errorLabel = h.newLabel("", 10)

	h.input.nickname = h.UI.NewTextInput(spiffy.Rect{}, "Nickname", nil)
	h.input.nickname.SetMaxLength(shared.NameLengthLimit)
	h.input.nickname.Style.FillColor = gctx2d.RGBA(238.0/255.0, 238.0/255.0, 238.0/255.0, 1)
	h.input.nickname.Style.BorderColor = gctx2d.ColorBlack
	h.input.nickname.Style.FocusBorderColor = gctx2d.ColorBlack
	h.input.nickname.Style.BorderWidth = 3.5
	h.input.nickname.Style.BorderRadius = 0
	h.input.nickname.TextStyle.FontSize = 18
	h.input.nickname.TextStyle.Font = h.RenderState.Renderer.LegacyFont
	h.input.nickname.TextStyle.Align = gctx2d.TextAlignCenter
	h.input.nickname.Padding = spiffy.Insets{Left: 8, Right: 8}
	h.input.nickname.OnSubmit = func(_ *spiffy.TextInput, value string) {
		h.OnSubmit(value)
	}

	h.tips.howTo = h.newLabel("How to play", 21)
	h.tips.move = h.newLabel("Use mouse to move", 11.2)
	h.tips.attack = h.newLabel("Left click to attack", 11.2)
	h.tips.defend = h.newLabel("Right click to defend", 11.2)

	h.menus.settings = h.newButton("Settings")
	h.menus.petals = h.newButton("Petals")
	h.menus.mobs = h.newButton("Mobs")
	h.menus.changelog = h.newButton("Changelog")
}

func (h *Home) layoutUI(width, height float32) {
	h.title.Bounds = spiffy.Rect{X: 0, Y: 116, Width: width, Height: 108}
	h.tips.howTo.Bounds = spiffy.Rect{X: 0, Y: 408, Width: width, Height: 40}
	h.tips.move.Bounds = spiffy.Rect{X: 0, Y: 442, Width: width, Height: 20}
	h.tips.attack.Bounds = spiffy.Rect{X: 0, Y: 458, Width: width, Height: 20}
	h.tips.defend.Bounds = spiffy.Rect{X: 0, Y: 474, Width: width, Height: 20}
	h.input.prompt.Bounds = spiffy.Rect{X: 0, Y: 276, Width: width, Height: 28}
	h.input.nickname.Bounds = spiffy.Rect{X: width*.5 - 140, Y: 310, Width: 280, Height: 35}
	h.input.submitHint.Bounds = spiffy.Rect{X: 0, Y: 352, Width: width, Height: 22}
	h.input.errorLabel.Bounds = spiffy.Rect{X: 0, Y: 376, Width: width, Height: 20}

	var bottomY float32 = height - 31.2
	h.menus.settings.Bounds = spiffy.Rect{X: 7, Y: bottomY, Width: 70, Height: 24.5}
	h.menus.petals.Bounds = spiffy.Rect{X: 84, Y: bottomY, Width: 70, Height: 24.5}
	h.menus.mobs.Bounds = spiffy.Rect{X: 161, Y: bottomY, Width: 70, Height: 24.5}
	h.menus.changelog.Bounds = spiffy.Rect{X: 238, Y: bottomY, Width: 84, Height: 24.5}
}

// bindRendererCallbacks routes application events into the active Spiffy UI.
func (h *Home) bindRendererCallbacks() {
	h.RenderState.Cleanup = func() { h.UI.Clear() }
	h.RenderState.OnUpdate = func(dt float64) {
		if !h.petalAnimation.Ready {
			return
		}

		h.petalAnimation.time += dt

		var numPetals int = int(h.petalAnimation.Height * h.petalAnimation.Width * 0.00005)
		if len(h.petalAnimation.Petals) < numPetals && rand.Float64() > .99 {
			var id configs.PetalID = h.petalAnimation.PossiblePetals[rand.IntN(len(h.petalAnimation.PossiblePetals))]
			h.petalAnimation.Petals = append(h.petalAnimation.Petals, &PetalAnimationPetal{
				ID:   id,
				X:    rand.Float32() * -100,
				Y:    rand.Float32() * max(h.petalAnimation.Height, 1),
				S:    (16 + rand.Float32()*16) * float32(configs.Petals[id].Stats.Radius/configs.PetalBasic.Stats.Radius),
				R:    rand.Float32() * 2 * math.Pi,
				Time: 0,
				Seed: rand.Float32(),
			})
		}

		for i, petal := range h.petalAnimation.Petals {
			petal.X += float32(dt) * (25 + 75*petal.Seed)
			petal.Y += float32(dt) * float32(math.Sin(h.petalAnimation.time*0.5+float64(petal.Seed)*10)) * 1.5

			if petal.X-petal.S > h.petalAnimation.Width {
				if len(h.petalAnimation.Petals) > numPetals {
					h.petalAnimation.Petals = append(h.petalAnimation.Petals[:i], h.petalAnimation.Petals[i+1:]...)
					continue
				}

				petal.X = -petal.S * 1.5
				petal.Y = rand.Float32() * max(h.petalAnimation.Height, 1)
				petal.ID = h.petalAnimation.PossiblePetals[rand.IntN(len(h.petalAnimation.PossiblePetals))]
				petal.S = (16 + rand.Float32()*16) * float32(configs.Petals[petal.ID].Stats.Radius/configs.PetalBasic.Stats.Radius)
			}
		}
	}

	h.RenderState.Draw = func(dc *gogpu.Context) {
		var width, height int = dc.Size()
		var surface *wgpu.TextureView = dc.SurfaceView()
		if surface == nil {
			return
		}

		var targetWidth, targetHeight uint32 = dc.SurfaceSize()
		h.RenderState.Ctx.BeginTarget(width, height, int(targetWidth), int(targetHeight))
		h.RenderState.Ctx.SetEdgeAAMode(gctx2d.EdgeAAModeOn)
		h.RenderState.Ctx.SetImageSmoothingMode(gctx2d.ImageSmoothingModeLinear)
		h.drawBackground(h.RenderState.Ctx, float32(width), float32(height))
		h.UI.DrawTo(h.RenderState.Ctx)
		var _ error = h.RenderState.Ctx.Flush(surface, nil)
	}

	h.RenderState.OnMouseMove = func(x, y float64) { h.UI.OnMouseMove(x, y) }
	h.RenderState.OnMousePress = func(button gpucontext.MouseButton, x, y float64) { h.UI.OnMousePress(button, x, y) }
	h.RenderState.OnMouseRelease = func(button gpucontext.MouseButton, x, y float64) { h.UI.OnMouseRelease(button, x, y) }
	h.RenderState.OnScroll = func(deltaX, deltaY float64) { h.UI.OnScroll(deltaX, deltaY) }
	h.RenderState.OnKeyPress = func(key gpucontext.Key, mods gpucontext.Modifiers) { h.UI.OnKeyPress(key, mods) }
	h.RenderState.OnKeyRelease = func(key gpucontext.Key, mods gpucontext.Modifiers) { h.UI.OnKeyRelease(key, mods) }
	h.RenderState.OnTextInput = func(text string) { h.UI.OnTextInput(text) }
	h.RenderState.OnResize = func(width, height int) {
		h.layoutUI(float32(width), float32(height))
	}
}

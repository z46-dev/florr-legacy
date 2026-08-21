package home

import (
	"math"
	"slices"
	"strconv"
	"strings"

	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/game/simulation/renders"
	"github.com/z46-dev/florr-legacy/shared/configs"
	"github.com/z46-dev/gctx2d"
	"github.com/z46-dev/spiffy"
)

const (
	menuNone MenuKind = iota
	menuSettings
	menuPetals
	menuMobs
	menuChangelog

	mobGalleryRotation float32 = -math.Pi * 3 / 4
)

type scrollbarMetrics struct {
	Track, Thumb spiffy.Rect
	Maximum      float32
}

var (
	menuBlue       gctx2d.Color = gctx2d.RGBA(90.0/255.0, 159.0/255.0, 219.0/255.0, 1)
	menuBlueDark   gctx2d.Color = gctx2d.RGBA(73.0/255.0, 129.0/255.0, 177.0/255.0, 1)
	menuBlueDarker gctx2d.Color = gctx2d.RGBA(62.0/255.0, 113.0/255.0, 157.0/255.0, 1)
	menuRed        gctx2d.Color = gctx2d.RGBA(201.0/255.0, 83.0/255.0, 76.0/255.0, 1)
)

// sortedPetals returns only definitions which also have a static renderer.
func sortedPetals() (ids []configs.PetalID) {
	for id, definition := range configs.Petals {
		if definition != nil && definition.Meta != nil && renders.Petals[id] != nil {
			ids = append(ids, id)
		}
	}

	slices.Sort(ids)
	return
}

// sortedMobs returns only definitions which also have a static renderer.
func sortedMobs() (ids []configs.MobID) {
	for id, definition := range configs.Mobs {
		if definition != nil && definition.Meta != nil && renders.Mobs[id] != nil {
			ids = append(ids, id)
		}
	}

	slices.Sort(ids)
	return
}

// parseChangelog converts the deliberately small Markdown format into dated sections.
func parseChangelog(markdown string) (entries []ChangelogEntry) {
	var current *ChangelogEntry
	for line := range strings.SplitSeq(markdown, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "## ") {
			entries = append(entries, ChangelogEntry{Date: strings.TrimSpace(strings.TrimPrefix(line, "## "))})
			current = &entries[len(entries)-1]
		} else if current != nil && strings.HasPrefix(line, "- ") {
			current.Changes = append(current.Changes, strings.TrimSpace(strings.TrimPrefix(line, "- ")))
		}
	}

	return
}

// toggleMenu opens one anchored archive panel at a time.
func (h *Home) toggleMenu(kind MenuKind) {
	if h.menu.Open == kind {
		kind = menuNone
	}

	h.menu.Open = kind
	h.menu.Scroll = 0
	h.menu.DraggingScrollbar = false
	h.setSettingsVisible(kind == menuSettings)
}

func (h *Home) setSettingsVisible(visible bool) {
	h.menus.keyboardMovement.Visible = visible
	h.menus.movementHelper.Visible = visible
}

func (h *Home) panelBounds() (bounds spiffy.Rect) {
	var width, height, anchor float32
	switch h.menu.Open {
	case menuSettings:
		width, height, anchor = 192, 132, 49
	case menuPetals:
		width, height, anchor = 272, 315, 138
	case menuMobs:
		width, height, anchor = 280, 300, 227
	case menuChangelog:
		width, height, anchor = 400, 320, 324
	default:
		return
	}

	bounds = spiffy.Rect{
		X:      max(float32(5), min(anchor-width*.5, h.menu.Width-width-5)),
		Y:      max(float32(5), h.menu.Height-45-height),
		Width:  width,
		Height: height,
	}
	return
}

func pointIn(bounds spiffy.Rect, x, y float32) bool {
	return x >= bounds.X && x <= bounds.X+bounds.Width && y >= bounds.Y && y <= bounds.Y+bounds.Height
}

// pressMenu begins dragging the current menu's scrollbar, including track jumps.
func (h *Home) pressMenu(x, y float32) {
	var (
		metrics scrollbarMetrics
		visible bool
	)
	if metrics, visible = h.currentScrollbarMetrics(); !visible {
		return
	}

	var hitBounds spiffy.Rect = spiffy.Rect{
		X: metrics.Track.X - 5, Y: metrics.Track.Y,
		Width: 15, Height: metrics.Track.Height,
	}
	if !pointIn(hitBounds, x, y) {
		return
	}

	h.menu.DraggingScrollbar = true
	if pointIn(metrics.Thumb, x, y) {
		h.menu.ScrollbarGrabOffset = y - metrics.Thumb.Y
	} else {
		h.menu.ScrollbarGrabOffset = metrics.Thumb.Height * .5
		h.moveMenu(y)
	}
}

// moveMenu updates scrolling while the scrollbar thumb is being dragged.
func (h *Home) moveMenu(y float32) {
	if !h.menu.DraggingScrollbar {
		return
	}

	var (
		metrics scrollbarMetrics
		visible bool
	)
	if metrics, visible = h.currentScrollbarMetrics(); !visible {
		h.menu.DraggingScrollbar = false
		return
	}

	var travel float32 = metrics.Track.Height - metrics.Thumb.Height
	var thumbY float32 = min(metrics.Track.Y+travel, max(metrics.Track.Y, y-h.menu.ScrollbarGrabOffset))
	h.menu.Scroll = metrics.Maximum * (thumbY - metrics.Track.Y) / travel
}

// releaseMenu ends scrollbar dragging and handles the custom close affordance.
func (h *Home) releaseMenu(x, y float32) {
	h.menu.DraggingScrollbar = false
	var panel spiffy.Rect = h.panelBounds()
	if h.menu.Open != menuNone && pointIn(spiffy.Rect{X: panel.X + panel.Width - 28, Y: panel.Y + 10, Width: 20, Height: 20}, x, y) {
		h.toggleMenu(h.menu.Open)
	}
}

// scrollMenu scrolls gallery and changelog content while the pointer is over it.
func (h *Home) scrollMenu(delta float32) {
	var panel spiffy.Rect = h.panelBounds()
	if h.menu.Open == menuNone || !pointIn(panel, h.menu.MouseX, h.menu.MouseY) {
		return
	}

	var (
		metrics scrollbarMetrics
		visible bool
	)
	if metrics, visible = h.currentScrollbarMetrics(); !visible {
		return
	}

	h.menu.Scroll = min(metrics.Maximum, max(float32(0), h.menu.Scroll+delta*36))
}

// drawMenuPanel draws the surface below retained settings inputs.
func (h *Home) drawMenuPanel(ctx *gctx2d.Context) {
	if h.menu.Open == menuNone {
		return
	}

	var panel spiffy.Rect = h.panelBounds()
	ctx.SetFillStyle(menuBlue)
	ctx.FillRoundedRect(panel.X, panel.Y, panel.Width, panel.Height, 2)
	ctx.SetLineWidth(5)
	ctx.SetStrokeStyle(menuBlueDark)
	ctx.StrokeRoundedRect(panel.X, panel.Y, panel.Width, panel.Height, 2)
}

// drawMenuContent draws static gallery art, text, tooltips, and the Petals progress pie.
func (h *Home) drawMenuContent(ctx *gctx2d.Context) {
	h.drawPetalProgress(ctx)
	if h.menu.Open == menuNone {
		return
	}

	var panel spiffy.Rect = h.panelBounds()
	switch h.menu.Open {
	case menuSettings:
		h.outlinedText(ctx, "Settings", panel.X+panel.Width*.5, panel.Y+20, 18, gctx2d.TextAlignCenter)
		h.drawSettings(ctx, panel)
	case menuPetals:
		h.outlinedText(ctx, "Petal Gallery", panel.X+panel.Width*.5, panel.Y+20, 18, gctx2d.TextAlignCenter)
		h.drawPetalGallery(ctx, panel)
	case menuMobs:
		h.outlinedText(ctx, "Mob Gallery", panel.X+panel.Width*.5, panel.Y+20, 18, gctx2d.TextAlignCenter)
		h.drawMobGallery(ctx, panel)
	case menuChangelog:
		h.outlinedText(ctx, "Changelog", panel.X+panel.Width*.5, panel.Y+20, 18, gctx2d.TextAlignCenter)
		h.drawChangelog(ctx, panel)
	}

	h.drawClose(ctx, panel)
}

func (h *Home) drawPetalProgress(ctx *gctx2d.Context) {
	var bounds spiffy.Rect = h.menus.petals.Bounds
	var (
		x, y     float32 = bounds.X + bounds.Width - 11, bounds.Y + bounds.Height*.5
		progress float32 = float32(math.Mod(h.petalAnimation.time, 8) / 8)
	)

	ctx.SetFillStyle(render.MixColors(menuBlue, gctx2d.ColorBlack, .15))
	ctx.FillCircle(x, y, 7)
	if progress > 0 {
		ctx.BeginPath()
		ctx.MoveTo(x, y)
		ctx.Arc(x, y, 7, -math.Pi/2, -math.Pi/2+float32(math.Pi*2)*progress, false)
		ctx.ClosePath()
		ctx.SetFillStyle(render.MixColors(menuBlueDark, gctx2d.ColorBlack, .35))
		ctx.Fill()
	}
}

func (h *Home) drawClose(ctx *gctx2d.Context, panel spiffy.Rect) {
	var bounds spiffy.Rect = spiffy.Rect{X: panel.X + panel.Width - 28, Y: panel.Y + 10, Width: 20, Height: 20}
	ctx.SetFillStyle(menuRed)
	ctx.FillRoundedRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, 3)
	ctx.SetStrokeStyle(gctx2d.ColorWhite)
	ctx.SetLineWidth(3)
	ctx.BeginPath()
	ctx.MoveTo(bounds.X+5, bounds.Y+5)
	ctx.LineTo(bounds.X+bounds.Width-5, bounds.Y+bounds.Height-5)
	ctx.MoveTo(bounds.X+bounds.Width-5, bounds.Y+5)
	ctx.LineTo(bounds.X+5, bounds.Y+bounds.Height-5)
	ctx.Stroke()
}

func (h *Home) drawSettings(ctx *gctx2d.Context, panel spiffy.Rect) {
	h.outlinedText(ctx, "Keyboard movement", panel.X+42, panel.Y+65, 11, gctx2d.TextAlignLeft)
	h.outlinedText(ctx, "Movement helper", panel.X+42, panel.Y+100, 11, gctx2d.TextAlignLeft)
}

func (h *Home) galleryClip(ctx *gctx2d.Context, panel spiffy.Rect) {
	var bounds spiffy.Rect = galleryBounds(panel)
	ctx.Save()
	ctx.BeginPath()
	ctx.Rect(bounds.X, bounds.Y, bounds.Width, bounds.Height)
	ctx.Clip()
}

func galleryBounds(panel spiffy.Rect) (bounds spiffy.Rect) {
	bounds = spiffy.Rect{X: panel.X + 7, Y: panel.Y + 43, Width: panel.Width - 14, Height: panel.Height - 50}
	return
}

func (h *Home) drawPetalGallery(ctx *gctx2d.Context, panel spiffy.Rect) {
	var (
		hovered       *configs.Petal
		hoveredBounds spiffy.Rect
		contentArea   spiffy.Rect = galleryBounds(panel)
	)

	h.galleryClip(ctx, panel)
	for index, id := range h.menu.Petals {
		var bounds spiffy.Rect = spiffy.Rect{
			X:      panel.X + 24 + float32(index%4)*56,
			Y:      panel.Y + 50 + float32(index/4)*56 - h.menu.Scroll,
			Width:  47,
			Height: 47,
		}

		var isHovered bool = pointIn(contentArea, h.menu.MouseX, h.menu.MouseY) && pointIn(bounds, h.menu.MouseX, h.menu.MouseY)
		var petal *configs.Petal = configs.Petals[id]
		var rarity *configs.Rarity = configs.Rarities[petal.Meta.Rarity]
		ctx.SetFillStyle(render.MixColors(rarity.Color, gctx2d.ColorWhite, hoverMix(isHovered)))
		ctx.FillRoundedRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, 2)
		ctx.SetLineWidth(2)
		ctx.SetStrokeStyle(render.MixColors(rarity.Color, gctx2d.ColorBlack, .18))
		ctx.StrokeRoundedRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, 2)

		var petalRender *renders.PetalRender = renders.Petals[id]
		if petalRender != nil && petalRender.IconImage != nil {
			ctx.SetObjectTransform(bounds.X+bounds.Width*.5, bounds.Y+18, 17, 0)
			ctx.DrawImageBounds(petalRender.IconImage, -.5, -.5, 1, 1)
			ctx.ResetTransform()
		}
		h.outlinedText(ctx, petal.Meta.Name, bounds.X+bounds.Width*.5, bounds.Y+39, 8, gctx2d.TextAlignCenter)

		if isHovered {
			hovered = petal
			hoveredBounds = bounds
		}
	}
	ctx.Restore()
	h.drawScrollbar(ctx)
	if hovered != nil {
		h.drawPetalTooltip(ctx, hovered, hoveredBounds)
	}
}

func (h *Home) drawMobGallery(ctx *gctx2d.Context, panel spiffy.Rect) {
	var (
		hovered       *configs.Petal
		hoveredBounds spiffy.Rect
		contentArea   spiffy.Rect = galleryBounds(panel)
	)
	h.galleryClip(ctx, panel)
	for index, id := range h.menu.Mobs {
		var bounds spiffy.Rect = spiffy.Rect{
			X: panel.X + 9, Y: panel.Y + 55 + float32(index)*154 - h.menu.Scroll,
			Width: panel.Width - 28, Height: 145,
		}
		var mob *configs.Mob = configs.Mobs[id]
		var rarity *configs.Rarity = configs.Rarities[mob.Meta.Rarity]
		ctx.SetFillStyle(menuBlueDark)
		ctx.FillRoundedRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, 2)

		h.outlinedText(ctx, mob.Meta.Name, bounds.X+11, bounds.Y+17, 14, gctx2d.TextAlignLeft)
		h.coloredText(ctx, rarity.Name, bounds.X+11, bounds.Y+34, 10, gctx2d.TextAlignLeft, rarity.Color)
		h.drawWrappedText(ctx, mob.Meta.Description, bounds.X+11, bounds.Y+52, bounds.Width-72, 10, 13)

		var renderBounds spiffy.Rect = spiffy.Rect{X: bounds.X + bounds.Width - 57, Y: bounds.Y + 10, Width: 47, Height: 47}
		ctx.SetFillStyle(render.MixColors(menuBlue, gctx2d.ColorWhite, .08))
		ctx.FillRoundedRect(renderBounds.X, renderBounds.Y, renderBounds.Width, renderBounds.Height, 3)
		var mobRender *renders.MobRender = renders.Mobs[id]
		if mobRender != nil && mobRender.CachedImage != nil {
			ctx.SetObjectTransform(renderBounds.X+renderBounds.Width*.5, renderBounds.Y+renderBounds.Height*.5, 34, mobGalleryRotation)
			ctx.DrawImageBounds(mobRender.CachedImage, -.5, -.5, 1, 1)
			ctx.ResetTransform()
		}

		for dropIndex, drop := range mob.Drops {
			var dropBounds spiffy.Rect = spiffy.Rect{
				X: bounds.X + 11 + float32(dropIndex)*45, Y: bounds.Y + 91,
				Width: 40, Height: 38,
			}
			var petal *configs.Petal = configs.Petals[drop.PetalID]
			if petal == nil {
				continue
			}
			var dropRarity *configs.Rarity = configs.Rarities[petal.Meta.Rarity]
			var isHovered bool = pointIn(contentArea, h.menu.MouseX, h.menu.MouseY) && pointIn(dropBounds, h.menu.MouseX, h.menu.MouseY)
			ctx.SetFillStyle(render.MixColors(dropRarity.Color, gctx2d.ColorWhite, hoverMix(isHovered)))
			ctx.FillRoundedRect(dropBounds.X, dropBounds.Y, dropBounds.Width, dropBounds.Height, 2)
			ctx.SetLineWidth(2)
			ctx.SetStrokeStyle(render.MixColors(dropRarity.Color, gctx2d.ColorBlack, .18))
			ctx.StrokeRoundedRect(dropBounds.X, dropBounds.Y, dropBounds.Width, dropBounds.Height, 2)
			if petalRender := renders.Petals[drop.PetalID]; petalRender != nil && petalRender.CachedImage != nil {
				ctx.SetObjectTransform(dropBounds.X+dropBounds.Width*.5, dropBounds.Y+13, 14, 0)
				ctx.DrawImageBounds(petalRender.CachedImage, -.5, -.5, 1, 1)
				ctx.ResetTransform()
			}
			h.outlinedText(ctx, petal.Meta.Name, dropBounds.X+dropBounds.Width*.5, dropBounds.Y+29, 6.5, gctx2d.TextAlignCenter)
			h.outlinedText(ctx, formatPercentage(drop.Chance*100), dropBounds.X+dropBounds.Width*.5, dropBounds.Y+46, 7, gctx2d.TextAlignCenter)
			if isHovered {
				hovered = petal
				hoveredBounds = dropBounds
			}
		}
	}
	ctx.Restore()
	h.drawScrollbar(ctx)
	if hovered != nil {
		h.drawPetalTooltip(ctx, hovered, hoveredBounds)
	}
}

func (h *Home) drawPetalTooltip(ctx *gctx2d.Context, petal *configs.Petal, anchor spiffy.Rect) {
	var (
		width  float32         = min(float32(252), h.menu.Width-10)
		height float32         = 70
		rarity *configs.Rarity = configs.Rarities[petal.Meta.Rarity]
	)
	if h.menu.Open == menuMobs {
		width = min(float32(290), h.menu.Width-10)
	}

	var x float32 = anchor.X + anchor.Width*.5 - width*.5
	x = min(h.menu.Width-width-5, max(float32(5), x))
	var y float32 = anchor.Y - height - 6
	if y < 5 {
		y = anchor.Y + anchor.Height + 6
	}
	y = min(h.menu.Height-height-5, max(float32(5), y))

	var bounds spiffy.Rect = spiffy.Rect{X: x, Y: y, Width: width, Height: height}
	ctx.SetFillStyle(gctx2d.RGBA(18.0/255.0, 31.0/255.0, 39.0/255.0, .82))
	ctx.FillRoundedRect(bounds.X, bounds.Y, bounds.Width, bounds.Height, 3)
	h.outlinedText(ctx, petal.Meta.Name, x+8, y+15, 15, gctx2d.TextAlignLeft)
	if rarity != nil {
		h.coloredText(ctx, rarity.Name, x+8, y+31, 10, gctx2d.TextAlignLeft, rarity.Color)
	}
	h.drawWrappedText(ctx, petal.Meta.Description, x+8, y+51, bounds.Width-16, 10, 12)
}

func (h *Home) drawChangelog(ctx *gctx2d.Context, panel spiffy.Rect) {
	var y float32 = panel.Y + 55 - h.menu.Scroll
	h.galleryClip(ctx, panel)
	for _, entry := range h.menu.Changelog {
		h.outlinedText(ctx, entry.Date, panel.X+14, y, 15, gctx2d.TextAlignLeft)
		y += 25
		for _, change := range entry.Changes {
			var lines []string = h.wrapText(ctx, "- "+change, panel.Width-42, 10.5)
			for _, line := range lines {
				h.outlinedText(ctx, line, panel.X+21, y, 10.5, gctx2d.TextAlignLeft)
				y += 15
			}
		}
		y += 9
		ctx.SetStrokeStyle(menuBlueDark)
		ctx.SetLineWidth(4)
		ctx.BeginPath()
		ctx.MoveTo(panel.X+16, y)
		ctx.LineTo(panel.X+panel.Width-24, y)
		ctx.Stroke()
		y += 19
	}
	ctx.Restore()
	h.drawScrollbar(ctx)
}

// currentScrollbarMetrics calculates one shared layout for drawing and pointer interaction.
func (h *Home) currentScrollbarMetrics() (metrics scrollbarMetrics, visible bool) {
	var (
		contentHeight  float32
		viewportHeight float32
	)
	switch h.menu.Open {
	case menuPetals:
		contentHeight, viewportHeight = float32((len(h.menu.Petals)+3)/4)*56, 257
	case menuMobs:
		contentHeight, viewportHeight = float32(len(h.menu.Mobs))*154, 242
	case menuChangelog:
		contentHeight, viewportHeight = h.changelogHeight(), 262
	default:
		return
	}
	if contentHeight <= viewportHeight {
		return
	}

	var (
		panel                       spiffy.Rect = h.panelBounds()
		trackX, trackY, trackHeight float32     = panel.X + panel.Width - 13, panel.Y + 47, panel.Height - 57
		thumbHeight                 float32     = max(float32(24), trackHeight*viewportHeight/contentHeight)
		thumbY                      float32     = trackY + (trackHeight-thumbHeight)*h.menu.Scroll/(contentHeight-viewportHeight)
	)
	metrics = scrollbarMetrics{
		Track:   spiffy.Rect{X: trackX, Y: trackY, Width: 5, Height: trackHeight},
		Thumb:   spiffy.Rect{X: trackX, Y: thumbY, Width: 5, Height: thumbHeight},
		Maximum: contentHeight - viewportHeight,
	}
	visible = true
	return
}

func (h *Home) drawScrollbar(ctx *gctx2d.Context) {
	var (
		metrics scrollbarMetrics
		visible bool
	)
	if metrics, visible = h.currentScrollbarMetrics(); !visible {
		return
	}

	ctx.SetFillStyle(render.MixColors(menuBlueDark, menuBlue, .35))
	ctx.FillRoundedRect(metrics.Track.X, metrics.Track.Y, metrics.Track.Width, metrics.Track.Height, 2.5)
	ctx.SetFillStyle(menuBlueDarker)
	ctx.FillRoundedRect(metrics.Thumb.X, metrics.Thumb.Y, metrics.Thumb.Width, metrics.Thumb.Height, 2.5)
}

func (h *Home) changelogHeight() (height float32) {
	for _, entry := range h.menu.Changelog {
		height += 53 + float32(len(entry.Changes))*23
	}
	return
}

func formatPercentage(value float64) (text string) {
	text = strings.TrimRight(strings.TrimRight(strconv.FormatFloat(value, 'f', 2, 64), "0"), ".") + "%"
	if value > 0 && value < 1 {
		text = strings.TrimPrefix(text, "0")
	}
	return
}

func hoverMix(hovered bool) (amount float32) {
	if hovered {
		amount = .12
	}
	return
}

func (h *Home) outlinedText(ctx *gctx2d.Context, text string, x, y float32, size float64, align gctx2d.TextAlign) {
	h.coloredText(ctx, text, x, y, size, align, gctx2d.ColorWhite)
}

func (h *Home) coloredText(ctx *gctx2d.Context, text string, x, y float32, size float64, align gctx2d.TextAlign, color gctx2d.Color) {
	ctx.SetFont(size, h.RenderState.Renderer.LegacyFont)
	ctx.SetTextAlign(align)
	ctx.SetTextBaseline(gctx2d.TextBaselineMiddle)
	ctx.SetLineWidth(float32(size) / 15)
	ctx.SetStrokeStyle(gctx2d.ColorBlack)
	ctx.StrokeText(text, x, y)
	ctx.SetFillStyle(color)
	ctx.FillText(text, x, y)
}

func (h *Home) drawWrappedText(ctx *gctx2d.Context, text string, x, y, width float32, size float64, lineHeight float32) {
	for _, line := range h.wrapText(ctx, text, width, size) {
		h.outlinedText(ctx, line, x, y, size, gctx2d.TextAlignLeft)
		y += lineHeight
	}
}

func (h *Home) wrapText(ctx *gctx2d.Context, text string, width float32, size float64) (lines []string) {
	ctx.SetFont(size, h.RenderState.Renderer.LegacyFont)
	var line string
	for _, word := range strings.Fields(text) {
		var candidate string = word
		if line != "" {
			candidate = line + " " + word
		}
		var metrics gctx2d.TextLineMetrics
		metrics, _ = ctx.MeasureTextLine(candidate, nil)
		if line != "" && metrics.Width > width {
			lines = append(lines, line)
			line = word
		} else {
			line = candidate
		}
	}
	if line != "" {
		lines = append(lines, line)
	}
	return
}

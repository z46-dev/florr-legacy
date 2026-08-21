package home

import (
	"fmt"
	"maps"
	"slices"

	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/game/simulation/renders"
)

func New(renderer *render.Renderer) (home *Home, err error) {
	if renderer == nil {
		err = fmt.Errorf("renderer is nil")
		return
	}

	home = &Home{
		petalAnimation: &PetalAnimation{
			PossiblePetals: slices.Collect(maps.Keys(renders.Petals)),
		},
	}

	if home.RenderState, err = renderer.NewState(); err != nil {
		home = nil
		return
	}

	home.buildUI()
	home.bindRendererCallbacks()
	home.layoutUI(1280, 720)

	go func() {
		var err error
		if err = renders.Init(renderer.Provider.Device(), true); err != nil {
			panic(fmt.Errorf("Failed to initialize petal renders: %w", err))
		}

		home.petalAnimation.Ready = true
	}()

	return
}

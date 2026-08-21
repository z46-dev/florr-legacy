package main

import (
	"fmt"

	"github.com/z46-dev/florr-legacy/client/home"
	"github.com/z46-dev/florr-legacy/client/render"
	"github.com/z46-dev/florr-legacy/client/singleplayer"
	"github.com/z46-dev/florr-legacy/game/simulation"
	"github.com/z46-dev/golog"
)

func main() {
	var (
		renderer *render.Renderer
		err      error
		log      *golog.Logger = golog.New()
	)

	if renderer, err = render.NewRenderer(log); err != nil {
		log.Panicf("create renderer: %v", err)
	}

	var localServer *singleplayer.LocalServer
	if localServer, err = singleplayer.ListenForSingleplayerLobbies(); err != nil {
		log.Panicf("listen for singleplayer lobbies: %v", err)
	}

	// var homeScreen *home.Home
	// if homeScreen, err = home.New(renderer, func(nickname string) (playErr error) {
	// 	var view *gameview.View
	// 	if view, playErr = singleplayer.NewSingleplayer(renderer, localServer.Address()); playErr != nil {
	// 		playErr = fmt.Errorf("connect to local game: %w", playErr)
	// 		return
	// 	}

	// 	if playErr = view.Spawn(nickname); playErr != nil {
	// 		view.Release()
	// 		playErr = fmt.Errorf("spawn flower: %w", playErr)
	// 		return
	// 	}
	// 	view.SetOnExit(func() {
	// 		view.Release()
	// 		homeScreen.Resume()
	// 		renderer.SetState(homeScreen.RenderState)
	// 	})

	// 	renderer.SetState(view.RenderState)
	// 	return
	// }); err != nil {
	// 	var _ error = localServer.Close()
	// 	log.Panicf("create home screen: %v", err)
	// }

	var homeScreen *home.Home
	if homeScreen, err = home.New(renderer); err != nil {
		log.Panicf("Create home screen: %v", err)
	}

	homeScreen.OnSubmit = func(username string) {
		var (
			sim *simulation.Simulation
			err error
		)

		if sim, err = singleplayer.NewSingleplayer(renderer, localServer.Address()); err != nil {
			err = fmt.Errorf("connect to local game: %w", err)
			return
		}

		if err = sim.Spawn(username); err != nil {
			sim.Release()
			err = fmt.Errorf("spawn flower: %w", err)
			return
		}

		// sim.SetOnExit(func() {
		// 	sim.Release()
		// 	homeScreen.Resume()
		// 	renderer.SetState(homeScreen.RenderState)
		// })

		renderer.SetState(sim.RenderState)
	}

	renderer.SetState(homeScreen.RenderState)

	if err = renderer.WaitAndCleanup(); err != nil {
		log.Panicf("wait and cleanup renderer: %v", err)
	}

	if err = localServer.Close(); err != nil {
		log.Errorf("close singleplayer server: %v", err)
	}
}

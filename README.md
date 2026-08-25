![Tests](https://img.shields.io/github/actions/workflow/status/z46-dev/florr-legacy/ci.yml?branch=main&event=push&label=CI)
![Made with Golang](https://img.shields.io/badge/-Made_with_Golang-007d9c?logo=go&logoColor=white)

# florr-legacy

A rewrite of https://static.florr.io/old/ to support multiplayer and potentially more gamemodes.

The evidence-gated roadmap to full behavioral and visual preservation is in
[`plan.md`](plan.md). Archived inputs under `static.florr.io/src` are immutable;
all reverse-engineering outputs and provenance live under
`static.florr.io/analysis`.

## Tech Stack

- Everything is in Go.
- The client application is built using [GoGPU](https://github.com/z46-dev/gogpu)
  - Helper libraries for rendering are:
    - [gctx2d](https://github.com/z46-dev/gctx2d): A HTML Canvas style 2D rendering context for GoGPU.
    - [spiffy](https://github.com/z46-dev/spiffy): A UI library that works with gctx2d and GoGPU.
- Communication between the game instance and the client is done using [arctic](https://github.com/z46-dev/arctic), which helps us reuse the same codebase for both singleplayer and multiplayer lobbies.

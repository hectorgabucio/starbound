# Starbound Story — Ebiten WASM Proof of Concept

This project is a proof of concept for using [Ebiten](https://ebiten.org/) to create a WebAssembly (WASM) module and deploy a web-based videogame written in Go. It demonstrates how Go and Ebiten can be used to build interactive, pixel-art games that run natively in the browser.



## Project Structure
- `main.go` — Entry point; handles game loop, stage transitions, and audio
- `internal/object/object.go` — Core object and cursor logic
- `assets/` — Contains all images and audio used in the game
- `site/` — (Optional) For web or static assets if needed

## Requirements
- Go 1.18 or newer
- Ebiten v2

To install dependencies:
```sh
go mod tidy
```

## Building for WebAssembly (WASM)
To build the game as a WASM module and serve it in the browser:

```sh
go run github.com/hajimehoshi/ebiten/v2/cmd/wasmserve@latest
```

Then open your browser at the provided local address (usually http://localhost:8080).

## Running the Game Natively
```sh
go run main.go
```

## Credits
- Developed by Héctor Gabucio
- Powered by [Ebiten](https://ebiten.org/)
- Inspired by Starbound

---
Enjoy exploring the Starbound Storyline WASM experience!

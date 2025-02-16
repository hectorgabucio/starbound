

serve: 
	go build ./... && go run github.com/hajimehoshi/wasmserve@latest . && sleep 2 && open http://localhost:8080
local:
	air
site-compile:
	env GOOS=js GOARCH=wasm go build -o site/starbound-story.wasm starbound-story


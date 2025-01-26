

serve: 
	go build ./cmd && go run github.com/hajimehoshi/wasmserve@latest . && sleep 2 && open http://localhost:8080
local:
	air

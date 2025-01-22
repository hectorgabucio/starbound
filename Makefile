

serve: 
	go build && go run github.com/hajimehoshi/wasmserve@latest . && sleep 2 && open http://localhost:8080
local:
	air

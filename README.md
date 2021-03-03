# wassex

This compiles Golang to wasm (webassembly) that opens a websocket back to the server, and then says so on the browser page. 

```
GOARCH=wasm GOOS=js go build -o client.wasm client.go
go run wwws.go 
```

Then open http://localhost:9090 

The websocket is requested via a call to http://localhost:9090/ws

This has been tested in Google Chrome on macOS using go1.16 

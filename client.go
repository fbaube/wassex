// +build js,wasm

package main

import (
	"context"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	// "os"
	"syscall/js"
	"time"
)

func main() {
	fmt.Println("Hello, WebAssembly! WTF")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	fmt.Println("Websocket TRYING...")
	// Dial(ctx, u string, opts *DialOpts) (*Conn, *http.Response, error)
	// Dial performs a WebSocket handshake on url.
	// The response is the WebSocket handshake response from the server.
	// You never need to close resp.Body yourself.
	// If an error occurs, the returned response may be non nil,
	// but you can only read the first 1024 bytes of the body.
	// http:// & https:// URLs work and are interpreted as ws/wss.
	c, _, err := websocket.Dial(ctx, "ws://localhost:9090/ws", nil)
	if err == nil {
		fmt.Println("Websocket OKAY")
	} else {
		fmt.Println("Websocket FAILED")
		TryJSWS()
	}
	defer c.Close(websocket.StatusInternalError,
		"Client sez: Ouch, defer'd Close!")

	err = wsjson.Write(ctx, c, "Hello from client!")
	if err != nil {
		fmt.Println("wsjson.Write FAILED")
	}
	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		fmt.Println("BARF")
	}
	fmt.Printf("received: %v \n", v)
	doDomStuff()
	ch := make(chan bool)
	<-ch
	// c.Close(websocket.StatusNormalClosure, "Client says: Close'ing OK")
}

func TryJSWS() {
	ws := js.Global().Get("WebSocket").New("ws://localhost:9090/ws")
	fmt.Printf("JS WS :: %+v \n", ws)
}

func doDomStuff() {
	//1. Adding an <h1> element in the HTML document
	document := js.Global().Get("document")
	p := document.Call("createElement", "h1")
	p.Set("innerHTML", "Hello from Golang!")
	document.Get("body").Call("appendChild", p)
	//2. Exposing go functions/values in javascript variables.
	js.Global().Set("goVar", "I am a variable set from Go")
	js.Global().Set("sayHello", js.FuncOf(sayHello))
}

func sayHello(this js.Value, inputs []js.Value) interface{} {
	firstArg := inputs[0].String()
	return "Hi " + firstArg + " from Go!"
}

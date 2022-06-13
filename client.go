// +build js,wasm

package main

import (
	"context"
	// "encoding/xml"
	"fmt"
	SU "github.com/fbaube/stringutils"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	S "strings"
	"syscall/js"
	"time"
)

// DOC, HED, BOD are immutable global singletons,
// so let's name them that way.
var DOC, HED, BOD js.Value

// SrvStat and SrvMsgs are for meta from the server.
var SrvStat, SrvMsgs js.Value

// init sets up some JS basics.
func init() {
	DOC = js.Global().Get("document")
	HED = DOC.Get("head")
	BOD = DOC.Get("body")
	if !(DOC.Truthy() && HED.Truthy()) {
		// return errors.New("Unable to get document object")
		panic("CAN'T GET webpage <html> and/or <head>")
	}
	if !BOD.Truthy() {
		panic("CAN'T GET webpage <body>")
	}
	SrvStat = DOC.Call("getElementById", "serverstatus")
	SrvMsgs = DOC.Call("getElementById", "servermessages")
	if !(SrvStat.Truthy() && SrvMsgs.Truthy()) {
		panic("CAN'T GET <serverstatus> and/or <servermessages>")
	}
}

func main() {
	sTime := time.Now().Local().Format(time.RFC3339)
	fmt.Println("Execution at", sTime)
	// Replace(str, oldstr, newstr string, m int) string
	sTime = S.Replace(sTime, "T", " ", 1)
	sTime = S.Replace(sTime, "+", " GMT+", 1)
	sTime = S.TrimSuffix(sTime, ":00")
	SrvStat.Set("innerHTML",
		SU.EmojiGreenlite+"localhost OK at local time "+sTime)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	HED = DOC.Get("head")
	BOD = DOC.Get("body")
	if !HED.Truthy() {
		panic("CAN'T GET JS HED")
	}
	if !BOD.Truthy() {
		panic("CAN'T GET JS DOC")
	}
	appendServerMessage("Trying websocket...")
	// Dial(ctx, u string, opts *DialOpts) (*Conn, *http.Response, error)
	// Dial performs a WebSocket handshake on url.
	// The *http.Response is the WebSocket handshake response from
	// the server, but you never need to close the rsp.Body yourself.
	// If an error occurs, the returned response may be non nil,
	// but you can only read the first 1024 bytes of the body.
	// http:// & https:// URLs work and are interpreted as ws/wss.
	c, _, err := websocket.Dial(ctx, "ws://localhost:9090/ws", nil)
	if err == nil {
		appendServerMessage("Websocket initialized OK")
	} else {
		appendServerMessage("Websocket FAILED, tryng to fetch external JS websocket...")
		TryJSWS()
	}
	defer c.Close(websocket.StatusInternalError,
		"Client sez: Ouch, defer'd Close!")

	appendServerMessage("<ul><li>Hello1</li><li>Hello2<li></ul>")
	appendServerMessage("<details open><summary>Item One</summary>jkljkljkl jkljkljkl <br/> jkljkljkl jkljkljkl</summary></details>")
	appendServerMessage("<details><summary>Item Two</summary>jkljkljkl jkljkljkl <br/> jkljkljkl jkljkljkl</summary></details>")
	appendServerMessage("<details><summary>Item Thri</summary>jkljkljkl jkljkljkl <br/> jkljkljkl jkljkljkl</summary></details>")
	appendServerMessage("<details open><summary>Item Four</summary>jkljkljkl jkljkljkl <br/> jkljkljkl jkljkljkl</summary></details>")

	err = wsjson.Write(ctx, c, "Hello from client!")
	if err != nil {
		fmt.Println("wsjson.Write FAILED")
	}
	var v interface{}
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		fmt.Println("BARF")
	}
	appendServerMessage(fmt.Sprintf("SERVER SAID: %v", v))
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
	var p js.Value
	// 1. Add an <h1>
	p = DOC.Call("createElement", "h1")
	p.Set("innerHTML", "This is an H1")
	BOD.Call("appendChild", p)
	p = DOC.Call("createElement", "h2")
	p.Set("innerHTML", "This is an H2")
	BOD.Call("appendChild", p)
	p = DOC.Call("createElement", "h3")
	p.Set("innerHTML", "This is an H3")
	BOD.Call("appendChild", p)
	p = DOC.Call("createElement", "h4")
	p.Set("innerHTML", "This is an H4")
	BOD.Call("appendChild", p)
	p = DOC.Call("createElement", "h5")
	p.Set("innerHTML", "This is an H5")
	BOD.Call("appendChild", p)
	p = DOC.Call("createElement", "h6")
	p.Set("innerHTML", "This is an H6")
	BOD.Call("appendChild", p)
	// 2. Expose Go functions/values in JS variables.
	// js.Global().Set("goVar", "I am a variable set from Go")
	// js.Global().Set("sayHello", js.FuncOf(sayHello))
}

func sayHello(this js.Value, inputs []js.Value) interface{} {
	firstArg := inputs[0].String()
	return "Hi " + firstArg + " from Go!"
}

func setServerStatus(s string) {
	SrvStat.Set("innerHTML", s)
}

func appendServerMessage(s string) {
	b4 := SrvMsgs.Get("innerHTML")
	if b4.String() != "" {
		s = b4.String() + "<br/>" + s
	}
	SrvMsgs.Set("innerHTML", s)
}

// +build !js

package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"context"
	"fmt"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	// "syscall/js"
	"time"
)

var (
	listen = flag.String("p", ":9090", "port (note the colon!)")
	dir    = flag.String("dir", ".", "directory to serve")
)

func wasmCheck(w http.ResponseWriter, req *http.Request) {
	if strings.HasSuffix(req.URL.Path, ".wasm") {
		w.Header().Set("content-type", "application/wasm")
	}
	http.FileServer(http.Dir(*dir)).ServeHTTP(w, req)
}

func socketer(w http.ResponseWriter, req *http.Request) {

	// Accept(w, req, opts *AcceptOptions) (*Conn, error)
	// Accept accepts a WebSocket handshake from a client
	//    and upgrades the the connection to a WebSocket.
	// Accept will not allow cross origin requests by default.
	// Accept will write a response to w on all errors.
	c, err := websocket.Accept(w, req, nil)
	if err != nil {
		fmt.Println("BARF-1")
	}
	// (c *Conn) Close(code StatusCode, reason string) error
	// Close performs the WebSocket close handshake with the
	// given status code and reason.
	// Details: It writes a WebSocket close frame with a 5s
	// timeout and then wait 5s for the peer to send a close
	// frame. All data messages received from the peer during
	// the close handshake are discarded. The connection can
	// only be closed once; later calls to Close are no-ops.
	// The max length of reason is 125 bytes. Avoid sending
	// a dynamic reason. Close unblocks all goroutines inter-
	// acting with the connection once complete.
	defer c.Close(websocket.StatusInternalError,
		"Server sez: Ouch, defer'd Close!")

	ctx, cancel := context.WithTimeout(req.Context(), time.Second*10)
	defer cancel()

	var v interface{}
	// Read(ctx, c *websocket.Conn, v interface{}) error
	// Read reads a JSON message from c into v.
	err = wsjson.Read(ctx, c, &v)
	if err != nil {
		fmt.Println("BARF-2")
	}
	log.Printf("received: %v", v)
	err = wsjson.Write(ctx, c, "Hello BACK from server")
	c.Close(websocket.StatusNormalClosure, "Server says: Close'ing OK")
}

func main() {
	flag.Parse()
	log.Printf("Will listen on port %q...", *listen)
	http.HandleFunc("/", wasmCheck)
	http.HandleFunc("/ws", socketer)
	http.ListenAndServe(*listen, nil)
}

package main

import (
	wbs "code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

const (
	WEB_ROOT = "./www/"
	WebSocketPath = "/ws"
)

func webListen(addr string) {
	http.Handle(WebSocketPath, wbs.Handler(wsHandler))
	http.Handle("/", http.FileServer(http.Dir(WEB_ROOT)))
	log.Println("Listening for HTTP on", addr)
	log.Printf("Listening for websockets on %s", addr + WebSocketPath)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func wsHandler(ws *wbs.Conn) {
	defer ws.Close()
	
	log.Println("Websocket connection recieved.")
	log.Println("Handeling user.")
	handleUser(ws)
}

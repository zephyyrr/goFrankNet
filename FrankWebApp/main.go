package main

import (
	"github.com/Zephyyrr/goFrankNet/frank"

	hash "crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"io"
)

const VERSION = 1.2

var (
	frankAddr = flag.String("a", "localhost:1342", "Address to FrankNet server")
	httpAddr  = flag.String("http", "0.0.0.0:2337", "Address to listen on for web server")
	
	clientDisconnectPackage, _ = json.Marshal(&frank.ServerPacket{CommandType: frank.DISCONNECTED})
)

func main() {
	flag.Parse()
	log.Printf("FrankNet WebApp Server v%g", VERSION)
	log.Println("Coupled with FrankNet at address", *frankAddr)
	webListen(*httpAddr)
}

func handleUser(user io.ReadWriteCloser) {
	closing := false
	defer func() {
		user.Write(clientDisconnectPackage)
		closing = true
		user.Close()
	}()
	
	dec := json.NewDecoder(user)
	username, password, newAcc := "", "", false
	
	for username == "" {
		cp := new(frank.ClientPacket)
		dec.Decode(cp)
		//fmt.Println("cp:", cp)
		username, password = cp.Username, cp.Password
	}
	server, err := frank.NewFrankConn(*frankAddr, username, HashPass(password), newAcc)
	defer func() {
		server.Close()
	}()

	if err != nil {
		user.Write(clientDisconnectPackage)
		if server != nil {
			server.Close()
		}
		user.Close()
		return
	}
	go func() {
		for sp := range server.Incoming {
			if sp.CommandType == frank.PING {
				cp := server.MakeClientPackage(frank.PONG, "")
				server.SendClientPackage(cp)
				continue
			}
			b, _ := json.Marshal(sp)
			_, err := user.Write(b)
			if err != nil || closing {
				closing = true
				return
			}
		}
	}()

	for {
		if closing {
			return
		}
		cp := new(frank.ClientPacket)
		err := dec.Decode(cp)
		if err != nil {
			closing = true
			return
		}
		if cp.CommandType != 0 {
			cp = server.MakeClientPackage(cp.CommandType, cp.Data)
			server.SendClientPackage(cp)
		}
	}
}

func HashPass(password string) string {
	hasher := hash.New()
	res := make([]byte, 0, hasher.Size())
	hasher.Write([]byte(password))
	res = hasher.Sum(res)
	return HashAsString(res)
}

func HashAsString(in []byte) string {
	out := ""
	for _, b := range in {
		if b < 0x10 {
			out = fmt.Sprintf("%s-0%X", out, b)
		} else {
			out = fmt.Sprintf("%s-%X", out, b)
		}
	}
	out = out[1:]
	return out
}

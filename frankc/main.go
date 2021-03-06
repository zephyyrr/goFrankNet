package main

import (
	"github.com/Zephyyrr/goFrankNet/frank"

	"crypto/sha1"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"time"
)

//might not be necessary...
type State struct {
	Current  string
	Playlist []string
	Users    map[string]bool
	Voting   string
}

func NewState() *State {
	s := new(State)
	s.Playlist = make([]string, 0, 25)
	s.Users = make(map[string]bool)
	s.Voting = "No voting in progress"
	return s
}

func (s State) IsPlaying(song string) bool {
	return song == s.Current
}

var (
	//Modifiers
	newAcc = flag.Bool("n", false, "Create new account.")

	//Settings
	address     = flag.String("a", "localhost:1342", "Address to FrankController.")
	password    = flag.String("p", "", "Password for FrankController.")
	youtube     = flag.String("y", "", "Youtube link to que.")
	grooveshark = flag.String("g", "", "Grooveshark link to que.")
	state       = NewState()
)

var conn *frank.FrankConn

func main() {
	sys_user, err := user.Current()
	sys_username := ""
	if err == nil {
		sys_username = sys_user.Username
	}
	username := flag.String("u", sys_username, "Username on FrankController. Defaults to logged in user's username.")
	flag.Parse()

	ps := HashPass(*password)

	log.Printf("Connecting to: %s", *address)

	conn, err = frank.NewFrankConn(*address, *username, ps, *newAcc)
	if err != nil {
		log.Fatalf("Error on connect: %s", err)
	}
	log.Println("Connection established!")

	go func() {
		for sp := range conn.Incoming {
			switch sp.CommandType {
			case frank.KICKED:
				log.Printf("%s kicked!", sp.Data)
			case frank.BANNED:
				log.Printf("%s banned!", sp.Data)
			case frank.SERVER_SHUTDOWN:
				log.Println("Server is shutting down!")
				os.Exit(0)
			case frank.MESSAGE:
				messageParser(sp.Data)
			case frank.SONGUPDATE:
				state.Current_Update(sp.NowPlaying)
			case frank.CLEAR_PLAYLIST:
				state.Playlist = make([]string, 0, 25)
			case frank.USER_DISCONNECT:
				state.Users_Remove(sp.Data)
			case frank.USER_CONNECT:
				state.Users_Add(sp.Data)
			case frank.FULL_UPDATE:
				state.Full_Update(sp)
				log.Println("Full update recived!")
			case frank.ADDNEWSONG:
				state.Playlist = append(state.Playlist, sp.Data)
			case frank.YOURADMIN:
				log.Println("You are an Admin!")
			case frank.ADMINLOG:
				log.Printf("Admin: %s", sp.Data)
			case frank.PING:
				conn.SendClientPackage(conn.MakeClientPackage(frank.PONG, ""))
			case frank.STARTVOTE:
				state.SetVoting(sp)
			default:
				log.Printf("%d: %s", sp.CommandType, sp.Data)
			}
		}

	}()

	switch {
	case *youtube != "":
		AddYoutube(*youtube)
	case *grooveshark != "": //sendGrooveshark(conn, *grooveshark)
	}

	time.Sleep(2 * time.Second)
	conn.Close()
}

func HashPass(password string) string {
	hasher := sha1.New()
	res := make([]byte, 0, hasher.Size())
	hasher.Write([]byte(password))
	res = hasher.Sum(res)
	return asString(res)
}

func asString(in []byte) string {
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

func messageParser(msg string) {
	switch msg {
	case "Password or username incorrect!":
		log.Fatalln(msg)
	default:
		log.Printf("Message: %s", msg)
	}
}

func AddYoutube(link string) {
	conn.SendClientPackage(conn.MakeClientPackage(frank.YOUTUBE, link))
}

func vote(vote byte) {
	conn.SendClientPackage(conn.MakeClientPackage(vote, ""))
}

/*
The functions below are subject of further research, if they really are necessary.
They might be removed if I find no use of the State struct.
*/

func SendTestPackets(conn *frank.FrankConn) {
	p := conn.MakeClientPackage(frank.YOUTUBE, "http://www.youtube.com/watch?v=FCARADb9asE&")
	log.Printf("Sending type %d packet with payload \"%s\"", p.CommandType, p.Data)
	conn.SendClientPackage(p)
	time.Sleep(3 * time.Second)
}

func TestHash() {
	ps := HashPass(*password)
	log.Printf("HashedPassword: %X", ps)
}

func (state *State) Current_Update(current string) {
	if current == "" {
		state.Current = "Not Playing!"
	} else {
		state.Current = current
	}
}

func (state *State) Users_Add(user string) {
	state.Users[user] = true
}

func (state *State) Users_Remove(user string) {
	delete(state.Users, user)
}

func (state *State) Full_Update(sp *frank.ServerPacket) {
	state.Current_Update(sp.NowPlaying)
	state.Playlist = sp.PlayList
	state.Users = make(map[string]bool)
	for _, s := range sp.Users {
		state.Users[s] = true
	}
}

func (state *State) SetVoting(sp *frank.ServerPacket) {
	switch sp.NowPlaying {
	case "prew":
		state.Voting = sp.Data + " votes for previous song"
	case "next":
		state.Voting = sp.Data + " votes for next song"
	case "clear":
		state.Voting = sp.Data + " votes for clearing of playlist"
	}
	go func() {
		time.Sleep(15 * time.Second)
		state.Voting = "No voting in progress"
	}()
}

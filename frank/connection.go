package frank

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"net"
	"time"
	//"log"
)

const VERSION = "5A"

var random = rand.New(rand.NewSource(time.Now().Unix()))

//Connection to a FrankNet server
//Incomming packets from the server can be collected from the Incomming channel.
type FrankConn struct {
	net.Conn
	Incoming      chan *ServerPacket
	username, mac string
	password      string
}

//Creates a new Frank Connection to the address.
//If newAcc is true, it will attempt to create a new user on the server.
//Password should be hashed with SHA512.
func NewFrankConn(addr, uname, pass string, newAcc bool) (*FrankConn, error) {
	fc := new(FrankConn)
	fc.username = uname
	fc.password = pass
	fc.mac = getMACAdress()
	fc.Incoming = make(chan *ServerPacket, 5)

	var err error
	fc.Conn, err = net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			sizeSlice := make([]byte, 4, 4)
			n, err := fc.Conn.Read(sizeSlice)
			if err != nil {
				break
			}

			if n != 4 {
				//log.Printf("only read %d bytes\n", n)
				continue
			}
			//log.Printf("Size bytes: %v", sizeSlice)
			size := int32(sizeSlice[0]) | int32(sizeSlice[1])<<8 | int32(sizeSlice[2])<<16 | int32(sizeSlice[3]<<24)
			//log.Printf("next packetsize: %d\n", size)

			dataSlice := make([]byte, size)
			_, err = fc.Conn.Read(dataSlice)

			if err != nil {
				break
			}
			//log.Println(string(dataSlice))
			sp := new(ServerPacket)
			xml.Unmarshal(dataSlice, sp)
			fc.Incoming <- sp
		}
	}()

	var p *ClientPacket
	if newAcc {
		p = fc.MakeClientPackage(NEWACCOUNT, VERSION)
	} else {
		p = fc.MakeClientPackage(CONNECT, VERSION)
	}
	fc.SendClientPackage(p)
	return fc, nil
}

//Connects to the FrankNet server. 
// Same as NewFrankConn(addr, uname, pass, false)
func Connect(addr, uname, pass string) (*FrankConn, error) {
	return NewFrankConn(addr, uname, pass, false)
}

// Creates a new package for sending to the server with some fields pre-filled in.
// These includes: Username, Password, MacAdress, CommandType and Data.
func (conn FrankConn) MakeClientPackage(ct byte, data string) *ClientPacket {
	cp := new(ClientPacket)
	cp.Username = conn.username
	cp.Password = conn.password
	cp.CommandType = ct
	cp.Data = data
	cp.MacAdress = conn.mac
	return cp
}

//Sends a packet to the server.
func (conn FrankConn) SendClientPackage(packet *ClientPacket) error {
	//log.Printf("Sending packet: %v\n", packet)
	b, err := xml.Marshal(packet) //XML-style! Henshin!

	if err != nil {
		return err
	}
	//log.Println(b)
	//log.Println(string(b))
	l := len(b)

	sizeSlice := []byte{byte(l & 0xFF), byte((l >> 8) & 0xFF), byte((l >> 16) & 0xFF), byte((l >> 24) & 0xFF)}
	_, err = conn.Write(sizeSlice) //Send length of packet first as four bytes in little-endian
	//log.Println("size as slice:", sizeSlice)
	//log.Println("size of packet:", l)
	if err != nil {
		return err
	}

	_, err = conn.Write(b) //Then send packet

	return err
}

//Closes the connection to the server
func (conn FrankConn) Close() error {
	p := conn.MakeClientPackage(DISCONNECT, "")
	conn.SendClientPackage(p)

	return conn.Conn.Close()
}

func getMACAdress() string {
	return fmt.Sprint(random.Int())
}

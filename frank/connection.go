package frank

import (
	"net"
	"encoding/xml"
	//"log"
)

const VERSION = "5A"

type FrankConn struct {
	net.Conn
	Incoming chan *ServerPacket 
	username, mac string
	password string
}

func NewFrankConn(addr, uname, pass string, newAcc bool) (*FrankConn, error) {
	fc := new(FrankConn)
	fc.username = uname
	fc.password = pass
	fc.mac = GetMACAdress()
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

func Connect(addr, uname, pass string) (*FrankConn, error) {
	return NewFrankConn(addr, uname, pass, false)
}

func (conn FrankConn) MakeClientPackage(ct byte, data string) *ClientPacket {
	cp := new(ClientPacket)
	cp.Username = conn.username
	cp.Password = conn.password
	cp.CommandType = ct
	cp.Data = data
	cp.MacAdress = conn.mac
	return cp
}

func (conn FrankConn) SendClientPackage(packet *ClientPacket) error {
	//log.Printf("Sending packet: %v\n", packet)
	b, err := xml.Marshal(packet) //XML-style! Henshin!
	
	if err != nil {
		return err
	}
	//log.Println(b)
	//log.Println(string(b))
	l := len(b)
	
	sizeSlice := []byte{byte(l&0xFF), byte((l>>8)&0xFF), byte((l>>16)&0xFF), byte((l>>24)&0xFF)}
	_, err = conn.Write(sizeSlice) //Send length of packet first as four bytes in little-endian
	//log.Println("size as slice:", sizeSlice)
	//log.Println("size of packet:", l)
	if err != nil {
		return err
	}
	
	_, err = conn.Write(b) //Then send packet
	
	return err
}

func (conn FrankConn) Close() {
	p := conn.MakeClientPackage(DISCONNECT, "")
	conn.SendClientPackage(p)
	
	conn.Conn.Close()
}

func GetMACAdress() string {
	return "Xenoblade"
}

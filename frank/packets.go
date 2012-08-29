package frank

import "encoding/xml"

type ClientPacket struct {
	Username string `xml:"Username"`
	Password string `xml:"Password"`
	CommandType byte `xml:"CommandType"`
	Data string `xml:"Data"`
	MacAdress string `xml:"MacAdress"`
	XMLName   xml.Name `xml:"ClientToServerPackage"`
}

type ServerPacket struct {
	NowPlaying string `xml:"nowPlaying"`
	Duration int32 `xml:"duration"`
	PlayList []string `xml:"playList>string"`
	Users []string `xml:"users>string"`
	CommandType byte `xml:"CommandType"`
	Data string `xml:"Data"`
	XMLName   xml.Name `xml:"ServerToClientPackage"`
}

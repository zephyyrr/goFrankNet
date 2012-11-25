package frank

const (
	CONNECT     = iota
	DISCONNECT  //Done
	NEWACCOUNT  //Done
	YOUTUBE     //Done
	GROOVESHARK //Done
	KICK        //Done
	BAN         //Done
	NEXT        //Done
	PREW        //Done
	STOP        //Done
	RESUME      //Done
	FPLAY       //Done
	CLEAR       //Done
	VOTENEXT    //Done
	VOTEPREW    //Done
	VOTECLEAR   //Done
	PONG        //Done
	VOLUME      //Done
	VOICE
	//Add move playlist items thing... later.
)

const (
	CONNECTED = iota
	KICKED
	BANNED
	DISCONNECTED
	SERVER_SHUTDOWN
	MESSAGE
	SONGUPDATE
	CLEAR_PLAYLIST
	USER_DISCONNECT
	USER_CONNECT
	FULL_UPDATE
	ADDNEWSONG
	YOURADMIN
	ADMINLOG
	VIDEORESPONSE
	PING
	STARTVOTE
)

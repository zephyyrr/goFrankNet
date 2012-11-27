function Show(id) {
	document.getElementById(id).style.display = 'block'
}
function Hide(id) {
	document.getElementById(id).style.display = 'none'
}

var ClientPacket = function(ct) {
	this.CommandType=ct
	this.Data=""
}

var SizePacket = function(s) {
	this.Size=s
}

function goOfflineFuncGenerator(message) {
	return function(sp) {
		Show("loginArea")
		Hide("loggedin")
		alert(message)
	}
}

var handlers = []

//Connected
handlers[0] = function(sp) {
	Hide("loginArea")
	Show("loggedin")
}

//Kicked
handlers[1] = goOfflineFuncGenerator("You have been kicked out!")

//Banned
handlers[2] = goOfflineFuncGenerator("You have been banned!")

//Disconnect
handlers[3] = goOfflineFuncGenerator("You have been disconnected.")

//Shutdown
handlers[4] = goOfflineFuncGenerator("Server is shutting down.")

//Message
handlers[5] = function(sp) {alert(sp.Data)}

//Song Update
handlers[6] = function(sp) {
	if (document.getElementById("playing")) {
		document.getElementById("playing").setAttribute("id", "")
	}
	var rows = document.getElementById('playlistTableBody').rows;
	for (var i=0; i<rows.length; i++) {
		if (rows[i].childNodes[1].innerHTML == sp.NowPlaying) {
			rows[i].setAttribute("id", "playing")
			return
		}
	}
}

//Clear Playlist
handlers[7] = function(sp) {
	var rows = document.getElementById('playlistTableBody').rows;
	while( rows[0] ) {
		rows[0].parentNode.removeChild( rows[0] );
	}
}

//User Disconnect
handlers[8] = function(sp) {
	var list = document.getElementById('users')
	for(var i=0; i<list.childNodes.length;i++) {
		if (list.childNodes[i].innerHTML==sp.Data) {
			list.removeChild(list.childNodes[i])
		}
	}
}

//User Connect
handlers[9] = function(sp) {
	handlers[8](sp)
	var list = document.getElementById('users')
	var newli=document.createElement("li")
	newli.appendChild(document.createTextNode(sp.Data))
	list.appendChild(newli)
}

//Full Update
handlers[10] = function(sp) {
	handlers[7](sp)
	if (sp.PlayList) {
		for (var i=0; i < sp.PlayList.length; i++){
			var row = document.getElementById('playlistTableBody').insertRow(-1)
			var cell1=row.insertCell(0);
			var cell2=row.insertCell(1);
			cell1.innerHTML=i+1;
			cell2.innerHTML=sp.PlayList[i];
		}
	}
	if (sp.Users) {
		for (var i=0; i < sp.Users.length; i++){
			var list = document.getElementById('users')
			var newli=document.createElement("li")
			newli.appendChild(document.createTextNode(sp.Users[i]))
			list.appendChild(newli)
		}
	}
	handlers[6](sp)
}

//Add New Song
handlers[11] = function(sp) {
	var row = document.getElementById('playlistTableBody').insertRow(-1)
	var cell1=row.insertCell(0);
	var cell2=row.insertCell(1);
	cell1.innerHTML=document.getElementById('playlistTableBody').childNodes.length-1;
	cell2.innerHTML=sp.Data;
}

//You Are Admin
//handlers[12]

//Admin Log
handlers[13] = function(sp) {console.log(sp.Data)}

//Video Response
//handlers[14]

//Ping (implemented serverside)
//handlers[15]

//Voting started
handlers[16] = function(sp) {
	document.getElementById('voteData').innerHTML=sp.Data + " votes for "
	switch (sp.NowPlaying) {
		case "next": 
			document.getElementById('voteData').innerHTML+="next song"
			break
		case "prew": 
			document.getElementById('voteData').innerHTML+="previous song"
			break
		case "clear": 
			document.getElementById('voteData').innerHTML+="clearing of playlist"
			break
	}
	setTimeout(function() {document.getElementById('voteData').innerHTML="No voting currently..."}, 15*1000)
}

function login(register) {
	document.getElementById('loginButton').disabled = true
	if (localStorage) {
		if (document.getElementById("remember").checked) {
			localStorage.setItem("username", document.getElementById("username").value)
			localStorage.setItem("password", document.getElementById("password").value)
			localStorage.setItem("remember", document.getElementById("remember").checked)
			
			if (document.getElementById("auto").checked) {
				localStorage.setItem("auto", document.getElementById("auto").checked)
			}
		} else {
			localStorage.clear()
		}
	}
	
	
	if (WebSocket) {
		ws = new WebSocket("ws://" + window.location.host + "/ws")
		ws.onopen = function() {
			var cp = new ClientPacket(register ? 2 : 0)
			cp.Username = document.getElementById("username").value
			cp.Password = document.getElementById("password").value
			SendClientPacket(cp)
		}
		ws.onmessage = function(evt) {
			var sp = JSON.parse(evt.data);
			if (sp.CommandType === undefined) {
				sp.CommandType = -1
			}
			if (handlers[sp.CommandType] !== undefined) {
				//console.log(sp.CommandType)
				handlers[sp.CommandType](sp)
			}
		}
		ws.onclose = function() {
			ws.close()
		}
	} else {
		Hide("loginArea")
		Show("WSnotSupported")
		
	}
	document.getElementById('loginButton').disabled = false
	return false
}

function AddSong(type, uri) {
	switch (type) {
		case "Youtube":
			var cp = new ClientPacket(3)
			cp.Data = uri
			SendClientPacket(cp)
			break
	}
	document.getElementById('addSongForm').reset()
}

function Vote(type) {
	switch (type) {
		case "Next":
			SendClientPacket(new ClientPacket(13))
			break
		case "Prew":
			SendClientPacket(new ClientPacket(14))
			break
		case "Clear":
			SendClientPacket(new ClientPacket(15))
			break
	}
}

function SendClientPacket(cp) {
	ws.send(JSON.stringify(cp))
} 

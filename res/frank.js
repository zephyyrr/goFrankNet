types=new Object
types.Youtube="http://www.youtube.com?v="

Next="Next"
Prew="Prew"
Clear="Clear"


function loadXMLDoc(url, cfunc) {
	var xmlhttp
	if (window.XMLHttpRequest)
	{// code for IE7+, Firefox, Chrome, Opera, Safari
	xmlhttp=new XMLHttpRequest();
	}
	else
	{// code for IE6, IE5
	xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
	}
	xmlhttp.onreadystatechange=cfunc(xmlhttp);
	xmlhttp.open("GET",url,true);
	xmlhttp.send();
}

function postAJAX(url, data, cfunc) {
	var xmlhttp
	if (window.XMLHttpRequest) {// code for IE7+, Firefox, Chrome, Opera, Safari
		xmlhttp=new XMLHttpRequest();
	} else {// code for IE6, IE5
		xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
	}
	if (cfunc != undefined) {
		xmlhttp.onreadystatechange=cfunc(xmlhttp);
	}
	xmlhttp.open("POST",url,true);
	xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
	//alert("New Request for " + data + " to " + url +"!")
	xmlhttp.send(data);
}

function setupForm() {
	x=document.getElementById("addSongForm")
	x.reset= function() {
		document.getElementById("addForm_link").value=types[document.getElementById("addForm_type").value]
	}
	x.reset()
}

function update_playlist() {
	loadXMLDoc("playlist/table",function(xmlhttp) {
		return function() {
			if (xmlhttp.readyState==4 && xmlhttp.status==200) {
				document.getElementById("playlist").innerHTML=xmlhttp.responseText;
			}
		}
	})
}

function update_current() {
	loadXMLDoc("current",function(xmlhttp) {
		return function() {
			if (xmlhttp.readyState==4 && xmlhttp.status==200) {
				document.getElementById("current").innerHTML=xmlhttp.responseText;
				if (xmlhttp.responseText == "") {
					document.title="Radio Frank - Not Playing!"
				} else {
					document.title="Radio Frank - " + xmlhttp.responseText
				}
			}
		}
	})
}

function update_voting() {
	loadXMLDoc("vote", function(xmlhttp) {
		return function() {
			if (xmlhttp.readyState==4 && xmlhttp.status==200) {
				document.getElementById("voteData").innerHTML=xmlhttp.responseText;
			}
		}
	})
}

function addSong() {
	var form = ""
	var x=document.getElementById("addSongForm")
	for (var i= 0; i< x.length-1;i++) {
		var elem=x.elements[i]
		if (i != 0) {
			form+="&"
		}
		form+=elem.name+"="+elem.value
	}
	x.reset()
	postAJAX("addsong", form)
}

function vote(type) {
	loadXMLDoc("vote?vote=" + type, function(){})
}
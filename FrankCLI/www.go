package main

import (
	"net/http"
	"html/template"
	"log"
	"flag"
	//"io"
	//"os"
)
const MAIN_HTML_NI="frank_noninteractive.html"
const MAIN_HTML="frank.html"
const CURR_HTML = "{{.Current}}"
const VOTE_HTML = "{{.Voting}}"
const PL_TABLE_HTML = `
{{$state:=.}}
<table id="playlist_table">
	<thead><tr><th colspan=2>Playlist</th></tr></thead>
	<tbody>
		{{range $k, $v := .Playlist}}<tr {{if $state.IsPlaying $v }}id="playing"{{end}}>
		<td>{{indexOne $k}}.</td><td>{{$v}}</td></tr> 
		{{end}}
	</tbody>
</table>
`

var funcMap = template.FuncMap{
	"indexOne": func(x int) int {return x+1},
}

var templ_f func() interface{}
var ninteractive = flag.Bool("nih", false, "Non-Interactive HTTP")

var root_templ = template.Must(template.New(MAIN_HTML).Funcs(funcMap).ParseFiles(MAIN_HTML))
var plTable_templ = template.Must(template.New("plTable").Funcs(funcMap).Parse(PL_TABLE_HTML))
var curr_templ = template.Must(template.New("current").Funcs(funcMap).Parse(CURR_HTML))
var vote_templ = template.Must(template.New("vote").Funcs(funcMap).Parse(VOTE_HTML))

func ListenAndServe(addr string, f func() interface{}) error {
	templ_f = f
	flag.Parse()
	http.HandleFunc("/", rootHandler)
	http.Handle("/res/", http.FileServer(http.Dir(".")))
	http.HandleFunc("/playlist/table", templ_Handl_Gen(plTable_templ))
	http.HandleFunc("/current", templ_Handl_Gen(curr_templ))
	http.HandleFunc("/vote", voteHandler)
	if *ninteractive {
		log.Println("Running in Non-Interactive mode")
		temp_templ:=loadTemplate(MAIN_HTML_NI)
		if temp_templ != nil {
			root_templ=temp_templ
			log.Println("Using Non-Interactive page instead")
		} else {
			log.Println("Could not use Non-Interactive page. ")
		}
		
	} else {
		http.HandleFunc("/addsong", addSongHandler)
	}
	
	return http.ListenAndServe(addr, nil)
}

func loadTemplate(path string) (t *template.Template) {
	defer func() {
		if x:=recover(); x!=nil {
			t=nil
		}
	}()
	t = template.Must(template.New(path).Funcs(funcMap).ParseFiles(path))
	return
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	//root_templ = template.Must(template.ParseFiles(MAIN_HTML)) //DEBUGGING
	err := root_templ.Execute(w, templ_f())
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func templ_Handl_Gen(template *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		err := template.Execute(w, templ_f())
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func addSongHandler(w http.ResponseWriter, req *http.Request) {
	songtype := req.FormValue("type")
	link := req.FormValue("link")
	
	log.Printf("%s requested %s as %s", req.RemoteAddr, link, songtype)
	
	if link != "" {
		switch songtype {
			case "Youtube": AddYoutube(link)
		}
	}
	w.Write([]byte("Thank you.\n"))
}

func voteHandler(w http.ResponseWriter, req *http.Request) {
	if !*ninteractive {
		vote := req.FormValue("vote")
		switch vote {
			case "Next": VoteNext()
			case "Prew": VotePrew()
			case "Clear": VoteClear()
		}
	}
	templ_Handl_Gen(vote_templ)(w, req)
}


package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

var dirFlag = flag.String("dir", "", "Directory with static content.")
var dir string
var templates *template.Template

func handler(rw http.ResponseWriter, req *http.Request) {
	err := templates.ExecuteTemplate(rw, "index.html", nil)
	if err != nil {
		glog.Errorln(err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Move struct {
	X      int
	Y      int
	Player Player
}

type GameLink struct {
	Source string
	Target string
}

type Server struct {
	Entrants  chan User
	Exeunt    chan User
	Players   []User
	Broadcast chan bool
	Request   chan GameLink
	Accept    chan GameLink
}

type Player string

const (
	White Player = "white"
	Black Player = "black"
	None  Player = "empty"
)

func (p Player) Switch() Player {
	switch p {
	case White:
		return Black
	case Black:
		return White
	}
	return None
}

type User struct {
	Name       string
	Connection *websocket.Conn
}

func boardhandler(w http.ResponseWriter, r *http.Request) {

}
func gamehandler(w http.ResponseWriter, r *http.Request) {
	player := Black
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	for {
		move := Move{}
		err := conn.ReadJSON(&move)
		if err != nil {
			glog.Errorln(err)
			continue
		}
		glog.V(1).Infoln(move)
		move.Player = player
		player = player.Switch()

		err = conn.WriteJSON(move)
		if err != nil {
			glog.Errorln(err)
			continue
		}
	}
}

func (s *Server) Lobby() {
	for {
		glog.V(1).Infoln("Checking lobby")
		var u User
		var gl GameLink
		select {
		case u = <-s.Entrants:
			glog.V(1).Infoln("New player")
			s.Players = append(s.Players, u)
			broadcast(s.Players)
		case u = <-s.Exeunt:
			glog.V(1).Infoln("Player left")
			s.Exit(u)
		case gl = <-s.Request:
			glog.V(1).Infoln("Requesting game:", gl)
			target, err := findUser(s.Players, gl.Target)
			if err != nil {
				glog.Errorln(err)
				continue
			}
			target.Connection.WriteJSON(gl)
		case gl = <-s.Accept:
			glog.V(1).Infoln("Accepting game:", gl)
			target, err := findUser(s.Players, gl.Target)
			if err != nil {
				continue
			}
			src, err := findUser(s.Players, gl.Source)
			if err != nil {
				continue
			}
			s.Exit(target)
			s.Exit(src)
		}
	}
}

func (s *Server) Exit(u User) {
	if i, err := indexOf(s.Players, u); err != nil {
		glog.Errorln(nil)
	} else {
		s.Players[i], s.Players = s.Players[len(s.Players)-1], s.Players[:len(s.Players)-1]
		broadcast(s.Players)
	}
}

func findUser(slice []User, name string) (User, error) {
	for _, v := range slice {
		if v.Name == name {
			return v, nil
		}
	}
	return User{}, fmt.Errorf("Element not found")
}

func indexOf(slice []User, elem User) (int, error) {
	for i, v := range slice {
		if v == elem {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Element not found")
}

func broadcast(players []User) {
	glog.V(1).Infoln("Broadcasting")
	glog.V(1).Infoln(players)
	names := []string{}
	for _, p := range players {
		names = append(names, p.Name)
	}
	for _, p := range players {
		go p.Connection.WriteJSON(names)
	}
}

func (s *Server) lobbyhandler(rw http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		glog.Errorln(err)
		return
	}
	_, message, _ := conn.ReadMessage()
	glog.Infoln("Player entered lobby:", string(message))
	s.Entrants <- User{Connection: conn, Name: string(message)}
	for {
		glog.V(1).Infoln("Awaiting request")
		gamelink := GameLink{}
		err = conn.ReadJSON(&gamelink)
		glog.V(1).Infoln("Request received:", gamelink)
		if err != nil {
			glog.Errorln(err)
			s.Exeunt <- User{Connection: conn, Name: string(message)}
			return
		}
		s.Request <- gamelink
	}
}

func main() {

	flag.Parse()
	dir, err := filepath.Abs(*dirFlag)
	if err != nil {
		glog.Errorln("Failed to aquire static directory:", err)
		return
	}
	templates, err = template.New("Templates").Delims("{[", "]}").ParseGlob(filepath.Join(dir, "templates", "*"))
	if err != nil {
		glog.Errorln("Failed to parse templates:", err)
		return
	}

	glog.V(1).Infoln("Starting Server at", dir)
	js := filepath.Join(dir, "js")
	css := filepath.Join(dir, "css")
	img := filepath.Join(dir, "img")

	glog.V(2).Infoln(js)
	glog.V(2).Infoln(css)
	glog.V(2).Infoln(img)

	server := Server{
		Players:  []User{},
		Exeunt:   make(chan User),
		Entrants: make(chan User),
		Request:  make(chan GameLink),
		Accept:   make(chan GameLink),
	}

	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir(js))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir(css))))
	http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir(img))))

	http.HandleFunc("/ws/lobby", server.lobbyhandler)
	http.HandleFunc("/ws/game", gamehandler)
	http.HandleFunc("/board", boardhandler)
	http.HandleFunc("/", handler)

	go server.Lobby()

	http.ListenAndServe(":8080", nil)
}

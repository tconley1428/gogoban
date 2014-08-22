package lobby

import (
	"fmt"
	"net/http"

	"gogoban/game"
	"gogoban/users"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Status string

const (
	Request Status = "request"
	Accept  Status = "accept"
	Cancel  Status = "cancel"
	Decline Status = "decline"
)

type GameLink struct {
	Source string
	Target string
	Status Status
}

type Lobby struct {
	Entrants  chan users.User
	Exeunt    chan users.User
	Players   []users.User
	Broadcast chan bool
	GameLinks chan GameLink
	Games     map[string]*game.Game
}

func CreateLobby() *Lobby {
	return &Lobby{
		Players:   []users.User{},
		Exeunt:    make(chan users.User),
		Entrants:  make(chan users.User),
		GameLinks: make(chan GameLink),
		Games:     make(map[string]*game.Game),
	}
}

func (s *Lobby) Run() {
	for {
		glog.V(1).Infoln("Checking lobby")
		var u users.User
		var gl GameLink
		select {
		case u = <-s.Entrants:
			glog.V(1).Infoln("New player")
			err := s.Enter(u)
			if err != nil {
				u.Connection.WriteMessage(0, []byte("Could not join lobby"))
			}
		case u = <-s.Exeunt:
			glog.V(1).Infoln("Player left")
			s.Exit(u)
		case gl = <-s.GameLinks:
			target, err := findusers(s.Players, gl.Target)
			if err != nil {
				glog.Errorln(err)
				continue
			}
			source, err := findusers(s.Players, gl.Source)
			if err != nil {
				glog.Errorln(err)
				continue
			}

			switch gl.Status {
			case Request:
				glog.V(1).Infoln("Requesting game:", gl)
				target.Connection.WriteJSON(gl)
			case Accept:
				glog.V(1).Infoln("Accepting game:", gl)
				glog.V(1).Infoln(source)
				glog.V(1).Infoln(target)

				g := game.CreateGame(source.Name, target.Name)
				glog.V(1).Infoln(g)
				s.Games[g.SessionID] = g

				source.Connection.WriteJSON(g)
				target.Connection.WriteJSON(g)

				glog.V(1).Infoln(g)

				s.Exit(target)
				s.Exit(source)
			case Cancel:
				target.Connection.WriteJSON(gl)
			case Decline:
				source.Connection.WriteJSON(gl)
			}
		}
	}
}

func (s *Lobby) Exit(u users.User) {
	if i, err := indexOf(s.Players, u); err != nil {
		glog.Errorln(err)
	} else {
		s.Players[i], s.Players = s.Players[len(s.Players)-1], s.Players[:len(s.Players)-1]
		broadcast(s.Players)
	}
}
func (s *Lobby) Enter(u users.User) error {
	for _, v := range s.Players {
		if v.Name == u.Name {
			return fmt.Errorf("Player already in lobby")
		}
	}
	s.Players = append(s.Players, u)
	broadcast(s.Players)
	return nil
}

func findusers(slice []users.User, name string) (users.User, error) {
	for _, v := range slice {
		if v.Name == name {
			return v, nil
		}
	}
	return users.User{}, fmt.Errorf("Element not found")
}

func indexOf(slice []users.User, elem users.User) (int, error) {
	for i, v := range slice {
		if v == elem {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Element not found")
}

func broadcast(players []users.User) {
	glog.V(1).Infoln("Broadcasting")
	glog.V(1).Infoln(players)
	names := []string{}
	for _, p := range players {
		names = append(names, p.Name)
	}
	for _, p := range players {
		go func(p users.User) { p.Connection.WriteJSON(names) }(p)
	}
}

func (s *Lobby) GameHandler(rw http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		glog.Errorln(err)
		return
	}
	userCookie, err := req.Cookie("username")
	username := userCookie.Value
	if err != nil {
		glog.Errorln("Player not logged in")
		return
	}
	session := mux.Vars(req)["session"]

	glog.V(2).Infoln(s.Games)
	glog.V(2).Infoln("Player ", username, " entered game:", session)
	game, exists := s.Games[session]
	glog.V(1).Infoln(game, " ", exists)
	if !exists {
		glog.Errorf("Game does not exist: %v in %v", session, s.Games)
		return
	}
	if game.Black.Name == username {
		game.Black.Connection = conn
		game.BlackJoined = true
	} else if game.White.Name == username {
		game.White.Connection = conn
		game.WhiteJoined = true
	} else {
		glog.Errorln("Not player's game: ", username)
	}
	glog.Infoln(game)
	if game.WhiteJoined && game.BlackJoined {
		game.Run()
	}
}

func (s *Lobby) Handler(rw http.ResponseWriter, req *http.Request) {
	conn, err := upgrader.Upgrade(rw, req, nil)
	if err != nil {
		glog.Errorln(err)
		return
	}

	_, message, _ := conn.ReadMessage()
	glog.V(2).Infoln("Player entered lobby:", string(message))
	s.Entrants <- users.User{Connection: conn, Name: string(message)}
	for {
		glog.V(1).Infoln("Awaiting request")
		gamelink := GameLink{}
		err = conn.ReadJSON(&gamelink)
		glog.V(1).Infoln("Request received:", gamelink)
		if err != nil {
			glog.Errorln(err)
			s.Exeunt <- users.User{Connection: conn, Name: string(message)}
			return
		}
		s.GameLinks <- gamelink
	}
}

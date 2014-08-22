package game

import (
	"fmt"
	"gogoban/users"
	"math/rand"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

type Game struct {
	SessionID     string
	White         users.User
	Black         users.User
	WhiteJoined   bool
	BlackJoined   bool
	CurrentPlayer users.User
}

type Move struct {
	X      int
	Y      int
	Player Player
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

func (g Game) NextPlayer() users.User {
	if g.CurrentPlayer == g.Black {
		return g.White
	}
	return g.Black
}

func CreateGame(white string, black string) *Game {
	return &Game{
		SessionID: fmt.Sprint(rand.Int()),
		Black:     users.User{Name: black},
		White:     users.User{Name: white},
	}
}

func (g Game) Run() {
	glog.Infoln("Game running")
	g.CurrentPlayer = g.Black
	color := Black
	for {
		move := Move{}
		err := g.CurrentPlayer.Connection.WriteMessage(websocket.TextMessage, []byte("Your Turn"))
		if err != nil {
			glog.Errorln(err)
			break
		}
		err = g.CurrentPlayer.Connection.ReadJSON(&move)
		if err != nil {
			glog.Errorln(err)
			break
		}
		move.Player = color
		color = color.Switch()
		glog.V(1).Infoln(move)
		g.CurrentPlayer = g.NextPlayer()

		err = g.Black.Connection.WriteJSON(move)
		if err != nil {
			glog.Errorln(err)
			break
		}
		err = g.White.Connection.WriteJSON(move)
		if err != nil {
			glog.Errorln(err)
			break
		}
	}
}

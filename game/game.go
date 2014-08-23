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
	Board         Board
}

type Board [][]Player

type Point struct {
	X int
	Y int
}

type Move struct {
	Loc    Point
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
	game := &Game{
		SessionID: fmt.Sprint(rand.Int()),
		Black:     users.User{Name: black},
		White:     users.User{Name: white},
		Board:     make([][]Player, 19),
	}
	for i := range game.Board {
		game.Board[i] = make([]Player, 19)
		for j := range game.Board[i] {
			game.Board[i][j] = None
		}
	}
	return game
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
		glog.V(1).Infoln(g.Board)
		move.Player = color
		if !g.Board.IsValid(move) {
			g.CurrentPlayer.Connection.WriteMessage(websocket.TextMessage, []byte("InvalidMove"))
			continue
		}
		g.Board.Apply(move)

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

func (b Board) Apply(m Move) {
	b[m.Loc.X][m.Loc.Y] = m.Player
	for _, loc := range b.AdjacentSquares(m.Loc) {
		if b[loc.X][loc.Y] == m.Player.Switch() {
			adj := Move{Loc: loc, Player: m.Player.Switch()}
			if b.Liberties(adj) == 0 {
				b.TakeGroup(adj)
			}
		}
	}
}

func (b Board) TakeGroup(m Move) {
	for _, loc := range b.AdjacentSquares(m.Loc) {
		if b[loc.X][loc.Y] == m.Player {
			b[loc.X][loc.Y] = None
			b.TakeGroup(Move{Loc: loc, Player: m.Player})
		}
	}
}

func (b Board) OnBoard(loc Point) bool {
	if loc.X < 0 {
		return false
	} else if loc.X > len(b) {
		return false
	} else if loc.Y < 0 {
		return false
	} else if loc.Y > len(b[0]) {
		return false
	}
	return true
}

func (b Board) IsValid(m Move) bool {
	glog.V(1).Infoln(b[m.Loc.X][m.Loc.Y])
	if b[m.Loc.X][m.Loc.Y] != None {
		glog.V(1).Infoln("Taken")
		return false
	}
	if b.Liberties(m) == 0 {
		glog.V(1).Infoln("0 Liberties")
		return false
	}
	if b.Ko(m) {
		return false
	}
	return true
}

func (b Board) Ko(m Move) bool {
	return false
}

func (b Board) LibertiesExcept(m Move, counted *map[Point]bool) int {
	liberties := 0
	(*counted)[m.Loc] = true
	for _, loc := range b.AdjacentSquares(m.Loc) {
		if (*counted)[loc] {
			continue
		}
		switch b[loc.X][loc.Y] {
		case None:
			liberties++
		case m.Player:
			if !(*counted)[m.Loc] {
				liberties += b.LibertiesExcept(m, counted)
			}
		}
	}
	return liberties
}
func (b Board) Liberties(m Move) int {
	visited := make(map[Point]bool)
	return b.LibertiesExcept(m, &visited)
}

func (b Board) AdjacentSquares(loc Point) []Point {
	locs := []Point{}
	if b.OnBoard(Point{loc.X - 1, loc.Y}) {
		locs = append(locs, Point{loc.X - 1, loc.Y})
	}
	if b.OnBoard(Point{loc.X + 1, loc.Y}) {
		locs = append(locs, Point{loc.X + 1, loc.Y})
	}
	if b.OnBoard(Point{loc.X, loc.Y - 1}) {
		locs = append(locs, Point{loc.X, loc.Y - 1})
	}
	if b.OnBoard(Point{loc.X, loc.Y + 1}) {
		locs = append(locs, Point{loc.X, loc.Y + 1})
	}
	return locs
}

func (b Board) String() string {
	s := ""
	for _, r := range b {
		for _, v := range r {
			s += fmt.Sprint(v, " ")
		}
		s += fmt.Sprintln("")
	}
	return s
}

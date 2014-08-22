package users

import "github.com/gorilla/websocket"

type User struct {
	Name       string
	Connection *websocket.Conn
}

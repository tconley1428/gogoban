package users

import "github.com/gorilla/websocket"

var Users = []string{"Tim", "Firefox"}

type User struct {
	Name       string
	Connection *websocket.Conn
}

func Exists(s string) bool {
	for _, v := range Users {
		if v == s {
			return true
		}
	}
	return false
}

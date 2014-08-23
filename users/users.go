package users

import (
	"fmt"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/websocket"
)

var Users = []string{"Tim", "Firefox"}
var CookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type User struct {
	Name       string
	Connection *websocket.Conn
}

func GetName(req *http.Request) (string, error) {
	if cookie, err := req.Cookie("session"); err == nil {
		value := make(map[string]string)
		if err = CookieHandler.Decode("session", cookie.Value, &value); err != nil {
			return "", fmt.Errorf("Cookie invalid")
		}
		glog.V(1).Infoln(value)
		return value["name"], nil
	} else {
		return "", fmt.Errorf("Cookie missing")
	}
}

func Exists(s string) bool {
	for _, v := range Users {
		if v == s {
			return true
		}
	}
	return false
}

package gogoban

import (
	"gogoban/lobby"
	"gogoban/users"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"

	"github.com/golang/glog"
)

var templates *template.Template

func handler(rw http.ResponseWriter, req *http.Request) {
	glog.V(1).Infoln(req.URL.String())
	err := templates.ExecuteTemplate(rw, "index.html", nil)
	if err != nil {
		glog.Errorln(err)
	}
}

func board(rw http.ResponseWriter, req *http.Request) {
	glog.Infoln("Board")
	err := templates.ExecuteTemplate(rw, "board.html", nil)
	if err != nil {
		glog.Errorln(err)
	}
}

func loginHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	redirectTarget := "/#/"
	if name != "" && pass != "" {
		glog.V(1).Infoln(name)
		glog.V(1).Infoln(pass)

		if users.Exists(name) {
			setSession(name, response)
			redirectTarget = "/#/lobby"
		} else {
			redirectTarget = "/#/?failed=" + name
		}
		glog.V(1).Infoln(redirectTarget)
	}
	http.Redirect(response, request, redirectTarget, 302)
}
func registerHandler(response http.ResponseWriter, request *http.Request) {
	name := request.FormValue("name")
	pass := request.FormValue("password")
	glog.V(1).Infoln(name)
	glog.V(1).Infoln(pass)
}
func logoutHandler(response http.ResponseWriter, request *http.Request) {
	clearSession(response)
	http.Redirect(response, request, "/", 302)
}

func setSession(userName string, response http.ResponseWriter) {
	value := map[string]string{
		"name": userName,
	}
	if encoded, err := cookieHandler.Encode("session", value); err == nil {
		http.SetCookie(response, &http.Cookie{
			Name:  "session",
			Value: encoded,
			Path:  "/",
		})
		http.SetCookie(response, &http.Cookie{
			Name:  "username",
			Value: userName,
			Path:  "/",
		})
	}
}

func clearSession(response http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(response, cookie)
}

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))
var router = mux.NewRouter()

func Start(dir string) {
	var err error
	templates, err = template.New("Templates").Delims("{[", "]}").ParseGlob(filepath.Join(dir, "templates", "*"))
	if err != nil {
		glog.Errorln("Failed to parse templates:", err)
		return
	}

	glog.V(1).Infoln("Starting Server at", dir)

	lob := lobby.CreateLobby()
	go lob.Run()

	router.HandleFunc("/register", registerHandler).Methods("POST")
	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/logout", logoutHandler).Methods("POST")
	router.HandleFunc("/wss/lobby", lob.Handler)
	router.HandleFunc("/wss/game/{session}", lob.GameHandler)
	router.PathPrefix("/www/").Handler(http.FileServer(http.Dir(dir)))
	router.HandleFunc("/", handler)

	http.Handle("/", router)

	err = http.ListenAndServeTLS(":1443", filepath.Join(dir, "public_key"), filepath.Join(dir, "private_key"), nil)
	if err != nil {
		glog.Errorln(err)
	}
}

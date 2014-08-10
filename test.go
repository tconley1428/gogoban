package main

import (
	"flag"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/golang/glog"
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

	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir(js))))
	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir(css))))
	http.Handle("/img/", http.StripPrefix("/img", http.FileServer(http.Dir(img))))

	http.HandleFunc("/", handler)

	http.ListenAndServe(":8080", nil)
}

package main

import (
	"flag"
	"path/filepath"

	"gogoban/gogoban"

	"github.com/golang/glog"
)

var dirFlag = flag.String("dir", "", "Directory with static content.")

func main() {

	flag.Parse()
	dir, err := filepath.Abs(*dirFlag)
	if err != nil {
		glog.Errorln("Failed to aquire static directory:", err)
		return
	}
	gogoban.Start(dir)
}

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Wes974/inventory/gui"
	"github.com/Wes974/inventory/lib/guiapi"
	"github.com/Wes974/inventory/parts"
)

var root string

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	gopath := os.Getenv("GOPATH")
	paths := filepath.SplitList(gopath)
	for _, p := range paths {
		project := filepath.Join(p, "src", "github.com", "Wes974", "inventory")
		info, err := os.Stat(project)
		if err == nil && info.IsDir() {
			root = project
			break
		}
	}
	if root == "" {
		log.Fatal("couldn't find the project in GOPATH")
	}
}

func fileServer(url string, path ...string) http.Handler {
	dir := http.Dir(filepath.Join(path...))
	return http.StripPrefix(url,
		http.FileServer(dir))
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		log.Fatal("first argument needs to be a database file path")
	}
	err := parts.SetupDB(args[0])
	if err != nil {
		log.Fatal(err)
	}
	log.Println("inventory server at :5080")
	http.Handle("/js/", fileServer("/js/", root, "gui", "js"))
	http.Handle("/guiapi/", guiapi.DefaultHandler)
	http.Handle("/", gui.Router())
	log.Println(http.ListenAndServe(":5080", nil))
}

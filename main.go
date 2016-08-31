package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/howeyc/fsnotify"
)

var (
	pdfpath   *string
	thumbpath *string
)

func main() {

	pdfpath = flag.String("pdfpath", "", "the absolute path to our pdfs")
	thumbpath = flag.String("thumbpath", "", "the absolute path to where we place our generated thumbs")
	port := flag.String("port", "9001", "the port on which to present our server")
	flag.Parse()

	if *pdfpath == "" || *thumbpath == "" {
		fmt.Println("Invalid -pdfpath or -thumbpath yo")
		os.Exit(1)
	}

	go watch()

	http.HandleFunc("/", servePage)
	http.Handle("/pdf/", http.StripPrefix("/pdf/", http.FileServer(http.Dir(*pdfpath))))
	http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir(*thumbpath))))
	http.ListenAndServe(":"+*port, nil)
}

func servePage(w http.ResponseWriter, r *http.Request) {

	filemap := make(map[string]string) //filemap maps png path => pdf path

	thumbs, _ := ioutil.ReadDir(*thumbpath)
	for _, thumb := range thumbs {
		if filepath.Ext(thumb.Name()) == ".png" {
			filemap["/img/"+thumb.Name()] = "/pdf/" + thumb.Name()[0:len(thumb.Name())-len(filepath.Ext(thumb.Name()))]
		}
	}

	html := `<html><head>
	<style type="text/css">
	body {
		text-align: center;
	}
	img {
		max-width: 150px;
		height:auto;
	}
	</style>
	</head><body>`
	for png, pdf := range filemap {
		html += "<a style=\"display:inline-block;margin:10px;\" href=\"" + pdf + "\"><img src=\"" + png + "\"></a>"
	}
	html += "</body></html>"

	w.Write([]byte(html))

}

//Watch grabs all pdfs in our *pdfpath folder, and for each pdf that
//does not have a corresponding .png in *thumbpath, creates one and places it there
func watch() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if filepath.Ext(ev.Name) == ".pdf" {
					go handleFileChange(ev)
				}
			case err := <-watcher.Error:
				fmt.Println("watcher error:", err)
			}
		}
	}()

	err = watcher.Watch(*pdfpath)
	if err != nil {
		panic(err)
	}

	done := make(chan bool) //Just keep going
	<-done
	watcher.Close()
}

func handleFileChange(ev *fsnotify.FileEvent) {
	if ev.IsCreate() {
		exec.Command("gs", "-q", "-o", *thumbpath+filepath.Base(ev.Name)+".png", "-dUseCropBox", "-sDEVICE=pngalpha", "-dLastPage=1", ev.Name).Output()
	}
	if ev.IsDelete() {
		//TODO delete unnecessary thumbs
	}
}

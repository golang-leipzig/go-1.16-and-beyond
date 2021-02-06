package main

import (
	"embed"
	"log"
	"net/http"
	"os"
	"text/template"
)

//go:embed index.gohtml
var indexGoHTML string

type templateData struct {
	Title,
	Description string
}

//go:embed assets
// Note that embed.FS is read-only.
// This also means that it is safe to share between goroutines.
var assets embed.FS

func main() {
	tmpl := template.Must(template.New("").Parse(indexGoHTML))

	http.Handle("/assets/", http.FileServer(http.FS(assets)))
	http.Handle("/", http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		tmpl.Execute(rw, templateData{
			Title:       "16th Leipzig Golang Meetup",
			Description: "An example on how to use the new embed package.",
		})
	}))

	log.Printf("Listening on %q", os.Args[1])
	err := http.ListenAndServe(os.Args[1], nil)
	if err != nil {
		log.Fatal(err)
	}
}

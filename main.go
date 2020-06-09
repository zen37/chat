package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {

	addr := flag.String("addr", ":8080", "The address of the application.")
	flag.Parse()

	/*

		html := "<html><head><title>Chat</title></head><body>Let's chat!</body></html>"

			http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte(html))
			})
	*/

	r := newRoom()

	//log.Println("room: ", r)

	http.Handle("/", &templateHandler{filename: "chat.html"})

	http.Handle("/room", r)

	go r.run()

	//start the web server
	log.Println("Starting web server on", *addr)
	//if err := http.ListenAndServeTLS(":8080", "certFile", "keyFile", nil); err != nil {
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServeTLS:", err)
	}
}

package main

import (
	_ "embed"
	"flag"
	"html/template"
	"log"
	"net/http"
)

var indexTmpl = template.Must(template.ParseFiles("index.html"))

//go:embed style.css
var styleCss []byte

//go:embed script.js
var scriptJs []byte

//go:embed htmx.min.js
var htmxMinJs []byte

var AppFlags struct {
	Port     string
	DataRoot string
}

var eventStore Store[Event]
var taskStore Store[Task]

func main() {
	// Parse App Flags
	flag.StringVar(&AppFlags.Port, "port", "8080", "server port")
	flag.StringVar(&AppFlags.DataRoot, "dataRoot", "data", "app's data root directory")
	flag.Parse()

	// Init Store
	eventStore = NewStore[Event]()
	taskStore = NewStore[Task]()

	router := NewRouter()
	// View
	router.Get("/", indexView)

	// Public Resources
	router.Get("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Write(styleCss)
	})
	router.Get("/script.js", func(w http.ResponseWriter, r *http.Request) { w.Write(scriptJs) })
	router.Get("/htmx.min.js", func(w http.ResponseWriter, r *http.Request) { w.Write(htmxMinJs) })

	// API
	router.Get("/events", findEventsHandler)
	router.Post("/events", createEventHandler)
	router.Post("/events/{id}/done", markEventAsDone)

	router.Get("/tasks", findTasksHandler)
	router.Post("/tasks", createTaskHandler)

	log.Fatal(http.ListenAndServe(":"+AppFlags.Port, nil))
}

func indexView(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{
		"PageTitle": "Index Page",
		"Events":    eventStore.FindAll(),
		"Tasks":     taskStore.FindAll(),
	}
	if err := indexTmpl.Execute(w, content); err != nil {
		panic(err)
	}
}

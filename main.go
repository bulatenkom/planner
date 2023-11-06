package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//go:embed script.js
var scriptJs []byte

//go:embed public/htmx.min.js
var htmxMinJs []byte

//go:embed public/missing/missing.min.css
var missingCss []byte

//go:embed public/missing/19.js
//go:embed public/missing/menu.js
//go:embed public/missing/overflow-nav.js
//go:embed public/missing/tabs.js
var missingCss_JSFiles embed.FS

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
	// Views
	router.Get("/", indexView)
	router.Get("/views/components/backlog.html", backlogView)
	router.Get("/views/components/events.html", eventsView)
	router.Get("/backlog", backlogView)
	router.Get("/calendar", calendarView)
	router.Get("/events", eventsView)
	router.Get("/readme", readmeView)
	router.Get("/changelog", changelogView)

	// Public Resources
	router.Get("/public/missing/missing.min.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Write(missingCss)
	})
	router.Get("/script.js", func(w http.ResponseWriter, r *http.Request) { w.Write(scriptJs) })
	router.Get("/public/htmx.min.js", func(w http.ResponseWriter, r *http.Request) { w.Write(htmxMinJs) })
	fs.WalkDir(missingCss_JSFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Ext(path) == ".js" {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			router.Get(path, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Add("Content-Type", "text/javascript")
				w.Write(data)
			})
		}
		return nil
	})

	// API
	router.Get("/events", findEventsHandler)
	router.Post("/events", createEventHandler)
	router.Post("/events/{id}/done", markEventAsDone)

	router.Get("/tasks", findTasksHandler)
	router.Post("/tasks", createTaskHandler)

	log.Fatal(http.ListenAndServe(":"+AppFlags.Port, nil))
}

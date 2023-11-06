package main

import (
	_ "embed"
	"html/template"
	"net/http"
)

// TODO Ideally these htmls should be embedded into exe
var templates = template.Must(template.ParseGlob("views/components/*.html"))

func indexView(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{
		"PageTitle": "Index Page", // Obsolete
	}
	if err := templates.ExecuteTemplate(w, "app.html", content); err != nil {
		panic(err)
	}
}

func backlogView(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{
		"Tasks": taskStore.FindAll(),
	}
	if err := templates.ExecuteTemplate(w, "backlog.html", content); err != nil {
		panic(err)
	}
}

func eventsView(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{
		"Events": eventStore.FindAll(),
	}
	if err := templates.ExecuteTemplate(w, "events.html", content); err != nil {
		panic(err)
	}
}

func calendarView(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{}
	if err := templates.ExecuteTemplate(w, "calendar.html", content); err != nil {
		panic(err)
	}
}

//go:embed readme.md
var readmeMd string

func readmeView(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{
		"Markdown": readmeMd,
	}
	if err := templates.ExecuteTemplate(w, "readme.html", content); err != nil {
		panic(err)
	}
}

//go:embed release-notes.txt
var releaseNotesTxt string

func changelogView(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{
		"Markdown": releaseNotesTxt,
	}
	if err := templates.ExecuteTemplate(w, "changelog.html", content); err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func markEventAsDone(w http.ResponseWriter, r *http.Request) {
	vars, ok := r.Context().Value("pathvars").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "could not access path variable")
		return
	}

	// Business logic. Just figured out a face of Spring developer. LOL :)
	event, err := eventStore.FindById(vars["{id}"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	event.Done = true
	// End of business logic

	json, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Fprintf(w, "%v", string(json))
}

func findEventsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("HX-Request") {
	case "true":
		findEventsHtmlHandler(w, r)
	default:
		findEventsJsonHandler(w, r)
	}
}

func findEventsHtmlHandler(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{
		"Events": eventStore.FindAll(),
	}
	if err := templates.ExecuteTemplate(w, "events.html", content); err != nil {
		panic(err)
	}
}

func findEventsJsonHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(eventStore.FindAll())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(json))
}

func createEventHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	rcvEvent := EventDto{}
	if err := json.Unmarshal(body, &rcvEvent); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(resolveUnmarshalErr(body, err)))
		return
	}

	plannedOn, err := time.Parse("2006-01-02T15:04", rcvEvent.PlannedOn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	duration, err := time.ParseDuration(rcvEvent.Duration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	eventType, err := ParseEventType(rcvEvent.Type)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	event := NewEvent(
		rcvEvent.Title,
		rcvEvent.Description,
		plannedOn,
		duration,
		eventType,
	)

	event, err = eventStore.Create(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%v", string(json))
}

func findTaskByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars, ok := r.Context().Value("pathvars").(map[string]string)
	if !ok {
		http.Error(w, "could not access path variable", http.StatusInternalServerError)
		return
	}
	switch r.Header.Get("HX-Request") {
	case "true":
		findTaskByIdHtmlHandler(w, r, vars["{id}"])
	default:
		findTaskByIdJsonHandler(w, r, vars["{id}"])
	}
}

func findTaskByIdHtmlHandler(w http.ResponseWriter, r *http.Request, id string) {
	task, err := taskStore.FindById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := templates.ExecuteTemplate(w, "backlog-task.html", task); err != nil {
		panic(err)
	}
}

func findTaskByIdJsonHandler(w http.ResponseWriter, r *http.Request, id string) {
	task, err := taskStore.FindById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(json))
}

func findTasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("HX-Request") {
	case "true":
		findTasksHtmlHandler(w, r)
	default:
		findTasksJsonHandler(w, r)
	}
}

func findTasksHtmlHandler(w http.ResponseWriter, r *http.Request) {
	status := strings.TrimSpace(r.URL.Query().Get("status"))

	content := map[string]any{
		"Tasks":      FindAllWithQuery(taskStore, TaskQuery{taskStatus(status)}),
		"StatusDict": getStatusDict(),
	}
	if err := templates.ExecuteTemplate(w, "backlog-find-all.html", content); err != nil {
		panic(err)
	}
}

type TaskQuery struct {
	Status taskStatus
}

func FindAllWithQuery(store Store[Task], query TaskQuery) map[string]*Task {
	if query.Status == "" {
		return store.store
	}

	resItems := make(map[string]*Task, len(store.store))

	if query.Status != "" {
		for k, v := range store.store {
			if v.Status == query.Status {
				resItems[k] = v
			}
		}
	}
	return resItems
}

func findTasksJsonHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(taskStore.FindAll())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(json))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("HX-Request") {
	case "true":
		createTaskHtmlHandler(w, r)
	default:
		createTaskJsonHandler(w, r)
	}
}

func createTaskHtmlHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(16 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rcvTask := TaskDto{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
	}

	// validation should go here...

	task := NewTask(
		rcvTask.Title,
		rcvTask.Description,
	)

	task, err = taskStore.Create(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render view
	w.Header().Add("HX-Trigger", "task-created")
	if err := templates.ExecuteTemplate(w, "backlog-task.html", task); err != nil {
		panic(err)
	}
}

func createTaskJsonHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	rcvTask := TaskDto{}
	if err := json.Unmarshal(body, &rcvTask); err != nil {
		fmt.Println(resolveUnmarshalErr(body, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	task := NewTask(
		rcvTask.Title,
		rcvTask.Description,
	)

	task, err = taskStore.Create(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%v", string(json))
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars, ok := r.Context().Value("pathvars").(map[string]string)
	if !ok {
		http.Error(w, "could not access path variable", http.StatusInternalServerError)
		return
	}

	switch r.Header.Get("HX-Request") {
	case "true":
		updateTaskHtmlHandler(w, r, vars["{id}"])
	default:
		http.Error(w, "unsuppported content-type", http.StatusNotAcceptable)
		// updateTaskJsonHandler(w, r)
	}
}

func updateTaskHtmlHandler(w http.ResponseWriter, r *http.Request, id string) {
	err := r.ParseMultipartForm(16 * 1024)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	rcvTask := TaskDto{
		Title:       r.FormValue("title"),
		Description: r.FormValue("description"),
		Status:      r.FormValue("status"),
	}

	taskStatus, err := ParseTaskStatus(rcvTask.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validation should go here...

	task, err := taskStore.FindById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// update
	task.Title = rcvTask.Title
	task.Description = rcvTask.Description
	task.Status = taskStatus

	task, err = taskStore.Update(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// render view
	if err := templates.ExecuteTemplate(w, "backlog-task.html", task); err != nil {
		panic(err)
	}
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	vars, ok := r.Context().Value("pathvars").(map[string]string)
	if !ok {
		http.Error(w, "could not access path variable", http.StatusInternalServerError)
		return
	}

	switch r.Header.Get("HX-Request") {
	case "true":
		deleteTaskHtmlHandler(w, r, vars["{id}"])
	default:
		http.Error(w, "unsuppported content-type", http.StatusNotAcceptable)
		// updateTaskJsonHandler(w, r)
	}
}

func deleteTaskHtmlHandler(w http.ResponseWriter, r *http.Request, id string) {
	if err := taskStore.deleteById(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

// Grabbed from this thread: https://groups.google.com/g/golang-nuts/c/RKbb4P_psP4
func resolveUnmarshalErr(data []byte, err error) string {
	if e, ok := err.(*json.UnmarshalTypeError); ok {
		// grab stuff ahead of the error
		var i int
		for i = int(e.Offset) - 1; i != -1 && data[i] != '\n' && data[i] != ','; i-- {
		}
		info := strings.TrimSpace(string(data[i+1 : int(e.Offset)]))
		s := fmt.Sprintf("%s - at: %s", e.Error(), info)
		return s
	}
	if e, ok := err.(*json.UnmarshalFieldError); ok {
		return e.Error()
	}
	if e, ok := err.(*json.InvalidUnmarshalError); ok {
		return e.Error()
	}
	return err.Error()
}

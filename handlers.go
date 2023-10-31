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
	if err != err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Fprintf(w, "%v", string(json))
}

func findEventsHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(eventStore.FindAll())
	if err != err {
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
		fmt.Println(resolveUnmarshalErr(body, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
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
	if err != err {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(event)
	if err != err {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%v", string(json))
}

func findTasksHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(taskStore.FindAll())
	if err != err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(json))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
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
	if err != err {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(task)
	if err != err {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%v", string(json))
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

package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

var slotStore = NewSlotStore()

var indexTmpl = template.Must(template.ParseFiles("index.html"))

//go:embed style.css
var styleCss []byte

//go:embed script.js
var scriptJs []byte

func main() {
	router := NewRouter()
	// View
	router.Get("/", indexView)

	// Public Resources
	router.Get("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		w.Write(styleCss)
	})
	router.Get("/script.js", func(w http.ResponseWriter, r *http.Request) { w.Write(scriptJs) })

	// API
	router.Get("/slots", findSlotsHandler)
	router.Post("/slots", createSlotHandler)
	router.Post("/slots/{id}/done", markSlotAsDone)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func indexView(w http.ResponseWriter, r *http.Request) {
	content := map[string]any{
		"PageTitle": "Index Page",
		"Slots":     slotStore.FindAll(),
	}
	if err := indexTmpl.Execute(w, content); err != nil {
		panic(err)
	}
}

func markSlotAsDone(w http.ResponseWriter, r *http.Request) {
	vars, ok := r.Context().Value("pathvars").(map[string]string)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "could not access path variable")
		return
	}

	// Business logic. Just figured out a face of Spring developer. LOL :)
	slot, err := slotStore.FindById(vars["{id}"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, err)
		return
	}

	slot.Done = true
	// End of business logic

	json, err := json.Marshal(slot)
	if err != err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Fprintf(w, "%v", string(json))
}

func findSlotsHandler(w http.ResponseWriter, r *http.Request) {
	json, err := json.Marshal(slotStore.FindAll())
	if err != err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "%v", string(json))
}

func createSlotHandler(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	rcvSlot := SlotDto{}
	if err := json.Unmarshal(body, &rcvSlot); err != nil {
		fmt.Println(resolveUnmarshalErr(body, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	plannedOn, err := time.Parse("2006-01-02T15:04", rcvSlot.PlannedOn)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	duration, err := time.ParseDuration(rcvSlot.Duration)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	slotType, err := ParseSlotType(rcvSlot.Type)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	slot := NewSlot(
		rcvSlot.Title,
		rcvSlot.Description,
		plannedOn,
		duration,
		slotType,
	)

	slot, err = slotStore.Create(slot)
	if err != err {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	json, err := json.Marshal(slot)
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

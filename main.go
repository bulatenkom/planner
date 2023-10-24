package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

var slotStore = NewSlotStore()

func main() {
	router := NewRouter()
	router.Get("/slots/", findSlotsHandler)
	router.Post("/slots/", createSlotHandler)
	router.Post("/slots/{id}/done", markSlotAsDone)

	log.Fatal(http.ListenAndServe(":8080", nil))
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
	w.WriteHeader(http.StatusCreated)

	slotType, err := ParseSlotType(r.URL.Query().Get("slotType"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	plannedOn, err := time.Parse(time.RFC3339, r.URL.Query().Get("plannedOn"))
	if err != err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	duration, err := time.ParseDuration(r.URL.Query().Get("duration"))
	if err != err {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	slot := NewSlot(
		r.URL.Query().Get("title"),
		r.URL.Query().Get("description"),
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	fmt.Fprintf(w, "%v", string(json))
}

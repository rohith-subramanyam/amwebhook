package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/alertmanager/template"
)

type responseJSON struct {
	Status  int
	Message string
}

func asJSON(w http.ResponseWriter, status int, message string) {
	data := responseJSON{
		Status:  status,
		Message: message,
	}
	bytes, _ := json.Marshal(data)
	json := string(bytes[:])

	w.WriteHeader(status)
	fmt.Fprint(w, json)
}

func webhook(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// https://godoc.org/github.com/prometheus/alertmanager/template#Data.
	data := template.Data{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		asJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("Alerts: GroupLabels=%v, CommonLabels=%v", data.GroupLabels,
		data.CommonLabels)
	for _, alert := range data.Alerts {
		log.Printf("Alert: status=%s,Labels=%v,Annotations=%v",
			alert.Status, alert.Labels, alert.Annotations)
	}

	asJSON(w, http.StatusOK, "success")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok!")
}

func main() {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/webhook", webhook)

	listenAddress := ":8080"
	if os.Getenv("PORT") != "" {
		listenAddress = ":" + os.Getenv("PORT")
	}

	log.Printf("Listening on: %v", listenAddress)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
}

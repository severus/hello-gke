package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// A Message represents the text and the timestamp.
type Message struct {
	Text      string    `json:"text"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

// SetTimestamp adds timestamp to the message.
func (m *Message) SetTimestamp() {
	m.Timestamp = time.Now().UTC()
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	m := new(Message)

	if err := json.NewDecoder(r.Body).Decode(m); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	m.Timestamp = time.Now().UTC()

	if err := json.NewEncoder(w).Encode(m); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

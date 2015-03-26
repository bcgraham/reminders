package main

import (
	"encoding/xml"
	"net/http"
)

func main() {
	http.HandleFunc("/twiml", twiml)
	http.ListenAndServe(":3001", nil)
}

func twiml(w http.ResponseWriter, r *http.Request) {
	twiml := TwiML{Say: "Hello, Travis!"}
	x, err := xml.Marshal(twiml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write(x)
}

type TwiML struct {
	XMLName xml.Name `xml:"Response"`

	Say string `xml:",omitempty"`
}

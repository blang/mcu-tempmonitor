package main

import (
	"container/ring"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	listenUDP  = flag.String("listen-udp", ":10001", "Listen udp address")
	listenHTTP = flag.String("listen-http", ":8080", "Listen http address")
	staticDir  = flag.String("static-dir", filepath.Dir(os.Args[0]), "Path to static files")
	high       = flag.Float64("high", 40, "Maximum temperature")
	low        = flag.Float64("low", -10, "Minimum temperature")
)

var tempBuf = TempBuffer{
	buf: ring.New(100000),
}

func main() {
	flag.Parse()

	handler := func(f float32) {
		if float64(f) > *high || float64(f) < *low {
			return
		}
		tempBuf.Append(
			Temp{
				TS:    time.Now(),
				Value: f,
			})
		log.Printf("Temp: %.02f", f)
	}
	http.HandleFunc("/list", List)
	http.HandleFunc("/", Index)
	go func() {
		if err := http.ListenAndServe(*listenHTTP, nil); err != nil {
			log.Fatalf("Error HTTP: %s", err)
		}
	}()
	if err := ListenAndServeUDP(*listenUDP, handler); err != nil {
		log.Fatalf("Error UDP: %s", err)
	}
}

type JSONTemp struct {
	Value float32 `json:"temp"`
	TS    int64   `json:"ts"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(*staticDir, "index.html"))
}
func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var temps []JSONTemp
	for _, t := range tempBuf.Slice() {
		temps = append(temps, JSONTemp{
			Value: t.Value,
			TS:    t.TS.Unix(),
		})
	}
	enc := json.NewEncoder(w)
	err := enc.Encode(temps)
	if err != nil {
		http.Error(w, "Error", http.StatusInternalServerError)
	}
}

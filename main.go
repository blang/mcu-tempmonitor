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
	interval   = flag.Duration("interval", 30*time.Second, "Probe interval")
	buffersize = flag.Int("size", 3000, "Size of records")
)

var tempBuf TempBuffer

func main() {
	flag.Parse()

	tempBuf = TempBuffer{
		buf: ring.New(*buffersize),
	}

	var intervalBuf []float32
	nextI := nextInterval()

	handler := func(f float32) {
		if float64(f) > *high || float64(f) < *low {
			return
		}
		intervalBuf = append(intervalBuf, f)
		if time.Now().After(nextI) {
			tempBuf.Append(
				Temp{
					TS:    time.Now(),
					Value: mean(intervalBuf),
				})
			intervalBuf = intervalBuf[:0]
			nextI = nextInterval()
		}
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

func mean(b []float32) float32 {
	var sum float32 = 0
	for _, e := range b {
		sum += e
	}
	return sum / float32(len(b))
}

func nextInterval() time.Time {
	return time.Now().Add(*interval)
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

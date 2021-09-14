package main

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/DataDog/datadog-agent/pkg/trace/pb"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s", r.Method, r.URL)
	w.WriteHeader(http.StatusOK)

	buf := &bytes.Buffer{}
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		log.Println("failed to get body")
		return
	}

	var traces pb.Traces
	_, err = traces.UnmarshalMsg(buf.Bytes())
	if err != nil {
		log.Println("failed parse traces")
		return
	}

	for _, t := range traces {
		for _, s := range t {
			log.Printf(
				"[%s] %s (%+v)",
				s.Service,
				s.Name,
				s.Meta,
			)
		}
	}
}

func main() {
	const listen = "0.0.0.0:8126"
	log.Printf("start server: %s", listen)

	http.HandleFunc("/", handler)
	err := http.ListenAndServe(listen, nil)
	log.Print(err)
}

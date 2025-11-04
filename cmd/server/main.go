package main

import (
    "log"
    "net/http"

    "github.com/peakied/ocr_server/pkg/api"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/ocr", api.OCRHandler)

    addr := ":8080"
    log.Printf("starting server on %s", addr)
    if err := http.ListenAndServe(addr, mux); err != nil {
        log.Fatalf("server failed: %v", err)
    }
}

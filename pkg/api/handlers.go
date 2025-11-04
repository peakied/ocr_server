package api

import (
    "encoding/json"
    "io"
    "log"
    "net/http"

    "github.com/example/ocr_server/internal/ocr"
)

type ocrRequest struct {
    Images []string `json:"images"`
    Type   string   `json:"type"` // "base64" or "url" (only base64 implemented)
}

type ocrResponse struct {
    Results []string `json:"results"`
    Error   string   `json:"error,omitempty"`
}

// OCRHandler handles POST /ocr
func OCRHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "failed to read body", http.StatusBadRequest)
        return
    }

    var req ocrRequest
    if err := json.Unmarshal(body, &req); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }

    if req.Type == "" {
        req.Type = "base64"
    }

    if req.Type != "base64" {
        http.Error(w, "only type=\"base64\" is supported in this version", http.StatusBadRequest)
        return
    }

    results := make([]string, 0, len(req.Images))
    for i, b64 := range req.Images {
        text, err := ocr.OCRFromBase64(b64)
        if err != nil {
            log.Printf("ocr error for image %d: %v", i, err)
            results = append(results, "")
            continue
        }
        results = append(results, text)
    }

    resp := ocrResponse{Results: results}
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(resp)
}

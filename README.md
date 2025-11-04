# OCR Server (Go) — simple Tesseract wrapper

This project provides a minimal Go HTTP server with a POST /ocr endpoint that accepts an array of images (base64) and returns OCR results using the tesseract CLI.

Prerequisites
- Go 1.20+
- Tesseract OCR installed and available on PATH

Windows (PowerShell) install notes:
- Install via Chocolatey: `choco install tesseract` (run PowerShell as Administrator)
- Or download official installer from tesseract project and ensure `tesseract.exe` is on PATH

Build & run

Open PowerShell and run:

```powershell
# from repository root (where go.mod is)
go build -o bin/ocr_server ./cmd/server
.
bin\ocr_server.exe
```

API

POST /ocr
- Content-Type: application/json
- Body shape:

```json
{
  "images": ["<base64-string>", "..."],
  "type": "base64"
}
```

- Response:

```json
{
  "results": ["text for first image", "text for second image"]
}
```

PowerShell example (one image)

```powershell
$b = [Convert]::ToBase64String((Get-Content -Path 'sample.png' -Encoding byte))
$body = @{ images = @($b); type = 'base64' } | ConvertTo-Json
Invoke-RestMethod -Method Post -ContentType 'application/json' -Body $body -Uri 'http://localhost:8080/ocr'
```

Notes
- This implementation invokes the `tesseract` CLI. Make sure Tesseract is installed.
- For production, consider `github.com/otiai10/gosseract/v2` which wraps the Tesseract API directly.

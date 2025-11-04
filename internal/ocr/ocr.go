package ocr

import (
    "encoding/base64"
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "strings"
)

// OCRFromBase64 decodes a base64 image and runs tesseract CLI on it.
// Returns the recognized text or an error. Requires `tesseract` installed and on PATH.
func OCRFromBase64(b64 string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(b64)
    if err != nil {
        return "", fmt.Errorf("failed to decode base64: %w", err)
    }

    tmpFile, err := ioutil.TempFile("", "ocr-*.png")
    if err != nil {
        return "", fmt.Errorf("failed to create temp file: %w", err)
    }
    defer os.Remove(tmpFile.Name())
    if _, err := tmpFile.Write(data); err != nil {
        tmpFile.Close()
        return "", fmt.Errorf("failed to write temp file: %w", err)
    }
    tmpFile.Close()

    // Run tesseract <tmpfile> stdout
    cmd := exec.Command("tesseract", tmpFile.Name(), "stdout")
    out, err := cmd.Output()
    if err != nil {
        if ee, ok := err.(*exec.ExitError); ok {
            return "", fmt.Errorf("tesseract failed: %s", string(ee.Stderr))
        }
        return "", fmt.Errorf("failed to run tesseract: %w", err)
    }

    text := strings.TrimSpace(string(out))
    return text, nil
}

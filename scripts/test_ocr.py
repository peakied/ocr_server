#!/usr/bin/env python3
"""
Simple Python test client for POST /ocr

Usage:
  python scripts/test_ocr.py [path-to-image]

If no image is provided, a tiny 1x1 PNG is used (no readable text, but tests the request flow).
This script uses the requests library if available, otherwise falls back to urllib from the stdlib.
"""
import sys
import json
import base64
import os

DEFAULT_PNG_B64 = (
    # 1x1 white PNG
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAA" \
    "AAC0lEQVR42mP8z8AARgAB/6D+gQAAAABJRU5ErkJggg=="
)

def read_image_b64(path):
    with open(path, "rb") as f:
        return base64.b64encode(f.read()).decode('ascii')

def post_with_requests(b64_list):
    import requests
    payload = {"images": b64_list, "type": "base64"}
    r = requests.post("http://localhost:8080/ocr", json=payload, timeout=30)
    r.raise_for_status()
    return r.text

def post_with_urllib(b64_list):
    from urllib import request
    payload = json.dumps({"images": b64_list, "type": "base64"}).encode('utf-8')
    req = request.Request("http://localhost:8080/ocr", data=payload, headers={"Content-Type": "application/json"})
    with request.urlopen(req, timeout=30) as resp:
        return resp.read().decode('utf-8')

def main():
    # Priority: explicit arg > ./sample.png (repo) > embedded default
    b64_list = []
    if len(sys.argv) > 1:
        # read all provided image paths
        for img_path in sys.argv[1:]:
            try:
                b64_list.append(read_image_b64(img_path))
            except Exception as e:
                print(f"failed to read image '{img_path}': {e}")
                sys.exit(2)
    elif os.path.exists(os.path.join(os.getcwd(), "sample.png")):
        try:
            b64_list.append(read_image_b64(os.path.join(os.getcwd(), "sample.png")))
            print("using sample.png from repository root")
        except Exception as e:
            print(f"failed to read sample.png: {e}")
            sys.exit(2)
    else:
        b64_list = [DEFAULT_PNG_B64]

    try:
        # prefer requests if installed
        import requests  # type: ignore
        out = post_with_requests(b64_list)
    except Exception:
        try:
            out = post_with_urllib(b64_list)
        except Exception as e:
            print(f"request failed: {e}")
            sys.exit(3)

    print(out)

if __name__ == '__main__':
    main()

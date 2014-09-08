
package main

import (
  "bytes"
  "fmt"
  "io"
  "log"
  "mime/multipart"
  "net/http"
  "os"
  "path/filepath"
)

const SERVER = "http://api.us.faceplusplus.com/detection/detect"
const API_KEY = "77bc802c35114a591719044f3e41d6ed"
const API_SECRET = "7Eww90BemlP6kRY66RZZaz-dHx0LVi9L"

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
  file, err := os.Open(path)
  if err != nil {
      return nil, err
  }
  defer file.Close()

  body := bytes.Buffer{}
  writer := multipart.NewWriter(&body)

  part, err := writer.CreateFormFile(paramName, filepath.Base(path))
  if err != nil {
      return nil, err
  }
  _, err = io.Copy(part, file)

  for key, val := range params {
      err = writer.WriteField(key, val)
      if err != nil {
        log.Fatal(err)
      }
  }

  log.Println("%+v", &body)
  err = writer.Close()
  if err != nil {
      return nil, err
  }

  req, _ := http.NewRequest("POST", uri, &body)
  req.Header.Set("Content-Type", writer.FormDataContentType())
  return req, nil
}

func main() {
  extraParams := map[string]string{
      "api_key": API_KEY,
      "api_secret": API_SECRET,
      "attribute": "glass,pose,gender,age,race,smiling",
  }

  filePath, err := filepath.Abs(os.Args[1])

  if err != nil {
    log.Fatal(err)
  }

  request, err := newfileUploadRequest(SERVER, extraParams, "img", filePath)
  if err != nil {
      log.Fatal(err)
  }

  client := &http.Client{}
  resp, err := client.Do(request)
  if err != nil {
      log.Fatal(err)
  } else {
      body := &bytes.Buffer{}
      _, err := body.ReadFrom(resp.Body)
    if err != nil {
          log.Fatal(err)
      }
    resp.Body.Close()
      fmt.Println(resp.StatusCode)
      fmt.Println(resp.Header)
      fmt.Println(body)
  }
}

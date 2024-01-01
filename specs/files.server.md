

```golang
package main

import (
    "fmt"
    "io"
    "mime/multipart"
    "net/http"
    "os"
)

const (
    uploadDir = "./uploads"
)

func main() {
    http.HandleFunc("/upload", uploadHandler)
    http.HandleFunc("/files", listFilesHandler)

    port := 8080
    fmt.Printf("Starting server on port: %d\n", port)
    http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        err := r.ParseMultipartForm(32 << 20)
        if err != nil {
            fmt.Println(err)
            return
        }

        // Get the file from the request
        file, _, err := r.FormFile("file")
        if err != nil {
            fmt.Println(err)
            return
        }
        defer file.Close()

        // Save the file to the upload directory
        filename := file.Filename()
        dst, err := os.Create(uploadDir + "/" + filename)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer dst.Close()

        // Copy the file to the upload directory
        io.Copy(dst, file)
        fmt.Printf("File uploaded: %s\n", filename)

        w.Write([]byte("File uploaded successfully"))
    }
}

func listFilesHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        // Get the list of files from the upload directory
        files, err := os.ListDir(uploadDir)
        if err != nil {
            fmt.Println(err)
            return
        }

        // Convert the file list to an HTML string
        html := "<h1>List of Uploaded Files</h1>"
        for _, filename := range files {
            html += fmt.Sprintf("<p><a href=\"/files/%s\">%s</a></p>", filename, filename)
        }

        // Send the HTML string as the response
        w.Write([]byte(html))
    }
}

```

```bash
curl -F "file=@path/to/file" http://localhost:8080/upload
curl -F "file=@myfile.txt" http://localhost:8080/upload
```

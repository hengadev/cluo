# Here is the MVP handler for the golang API

```go
func HandleAudio(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // 32MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "missing audio", http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Received audio: %s (%d bytes)", header.Filename, header.Size)

	// MVP: save locally
	out, _ := os.Create("./uploads/" + header.Filename)
	defer out.Close()
	io.Copy(out, file)

	w.WriteHeader(http.StatusOK)
}
```



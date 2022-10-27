package recon

import (
	"io"
	"log"
	"net/http"
	"os"
)

func CheckIfExists(path string) bool {
	// Check if target info file exists, and if return
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func RequestUrl(url string) []byte {
	// Request url
	r, er := http.Get(url)
	if er != nil {
		println(er.Error())
	}
	defer r.Body.Close()

	// Return data
	data, er := io.ReadAll(r.Body)
	if er != nil {
		log.Fatal("Unable to download resource from: " + url)
	}

	return data
}

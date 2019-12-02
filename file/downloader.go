package file

import (
	"io"
	"net/http"
	"os"
)

// Download http url conntent convert to file
func Download(toFile, url string) error {
	f, err := os.Create(toFile)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

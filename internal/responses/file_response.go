package responses

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func LoadFileFromURL(url string) (int, string, []byte) {
	prefix := "/var/www/"
	home, _ := os.UserHomeDir()
	page := filepath.Join(home, prefix, url)

	info, err := os.Stat(page)
	if os.IsNotExist(err) {
		return notFound(url)
	}
	if err != nil {
		return serverError(url)
	}
	if info.IsDir() {
		page = filepath.Join(page, "index.html")
	}

	data, err := os.ReadFile(page)
	if err != nil {
		return notFound(url)
	}
	return http.StatusOK, "Status OK", data
}

func notFound(url string) (int, string, []byte) {
	return http.StatusNotFound, fmt.Sprintf("Page %s not found", url), nil
}

func serverError(url string) (int, string, []byte) {
	return http.StatusInternalServerError, fmt.Sprintf("Could not retrieve page %s", url), nil
}

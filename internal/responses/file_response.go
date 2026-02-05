package responses

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func LoadFileFromURL(url string, headers map[string][]string) (int, string, []byte) {
	prefix := "/var/www/"
	home, _ := os.UserHomeDir()
	page := filepath.Join(home, prefix, url)

	info, err := os.Stat(page)
	if os.IsNotExist(err) {
		return notFound(url, headers)
	}
	if err != nil {
		return serverError(url, headers)
	}
	if info.IsDir() {
		page = filepath.Join(page, "index.html")
	}

	data, err := os.ReadFile(page)
	if err != nil {
		return notFound(url, headers)
	}

	ext := filepath.Ext(page)

	switch ext {
	case ".html":
		headers["Content-Type"] = []string{"text/html"}
	case ".css":
		headers["Content-Type"] = []string{"text/css"}
	case ".js":
		headers["Content-Type"] = []string{"application/javascript"}
	case ".json":
		headers["Content-Type"] = []string{"application/json"}
	case ".txt":
		headers["Content-Type"] = []string{"text/plain"}
	default:
		headers["Content-Type"] = []string{"application/octet-stream"}
	}

	return http.StatusOK, "Status OK", data
}

func notFound(url string, headers map[string][]string) (int, string, []byte) {
	headers["Content-Type"] = []string{"text/plain"}
	return http.StatusNotFound, fmt.Sprintf("Page %s not found", url), nil
}

func serverError(url string, headers map[string][]string) (int, string, []byte) {
	headers["Content-Type"] = []string{"text/plain"}
	return http.StatusInternalServerError, fmt.Sprintf("Could not retrieve page %s", url), nil
}

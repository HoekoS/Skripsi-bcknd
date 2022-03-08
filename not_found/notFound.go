package not_found

import (
	"fmt"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, status int, path string) {
	if r.URL.Path != path {
		w.WriteHeader(status)
		if status == http.StatusNotFound {
			fmt.Fprint(w, "custom 404")
		}
	}
}

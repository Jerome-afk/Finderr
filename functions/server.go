package functions

import (
	"net/http"
)

func Router(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		IndexHandler(w, r)
	} else if r.URL.Path == "/400" {
		//BadRequest(w, r)
		return
	} else {
		//PageNotFound(w, "404.html", nil, http.StatusNotFound)
		return
	}
}

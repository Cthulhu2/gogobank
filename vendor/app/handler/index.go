package handler

import (
	"net/http"
)

// IndexGET displays the home page
func IndexGET(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World!"))
}

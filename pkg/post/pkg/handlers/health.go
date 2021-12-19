package handlers

import "net/http"

func Health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("healthy"))
	w.WriteHeader(http.StatusOK)
}

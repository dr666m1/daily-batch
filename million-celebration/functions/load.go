package celebrate

import (
	"net/http"
)

func Load(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("hello"))
}

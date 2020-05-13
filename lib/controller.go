package lib

import "net/http"

// Controller コントローラ
type Controller interface {
	RegistControllers(mux *http.ServeMux)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

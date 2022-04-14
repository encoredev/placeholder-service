package http

import (
	"net/http"
)

func handleHealthRoute(w http.ResponseWriter, _ *http.Request) {
	setHeaders(w, http.StatusOK)

	_, _ = w.Write([]byte(`{ "ok": true }`))
}

func handleOtherRoutes(w http.ResponseWriter, _ *http.Request) {
	setHeaders(w, http.StatusNotFound)

	_, _ = w.Write([]byte(`{
  "code": "app_not_deployed",
  "message": "the encore application has not been deployed to this environment yet",
  "details": null
}`))
}

func setHeaders(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Encore-Service", "placeholder")
	w.Header().Set("X-Powered-By", "https://encore.dev")

	w.WriteHeader(statusCode)
}

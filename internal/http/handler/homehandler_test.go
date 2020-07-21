package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/", HomeHandler)

	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Bad status code, got: %v, want: %v", rr.Code, http.StatusOK)
	}
	data := rr.Body.String()
	if data != "TTTGameServer\n" {
		t.Errorf("Unexpected response, got: %v, want %v", data, "TTTGameServer\n")
	}
}

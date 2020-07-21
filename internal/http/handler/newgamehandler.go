package handler

import (
	"fmt"
	"net/http"
)

//NewGameHandler handles creating a new game
func NewGameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "New Game")
	}
}
package handler

import (
	"fmt"
	"net/http"
)

//JoinGameHandler handles game joining
func JoinGameHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Join Game")
	}
}

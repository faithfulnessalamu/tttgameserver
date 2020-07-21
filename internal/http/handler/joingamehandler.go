package handler

import (
	"fmt"
	"net/http"

	"github.com/thealamu/internal/engine"
)

//JoinGameHandler handles game joining
func JoinGameHandler() http.HandlerFunc {
	gE := engine.New()
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Join Game")
	}
}

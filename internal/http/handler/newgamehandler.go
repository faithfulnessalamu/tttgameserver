package handler

import (
	"fmt"
	"net/http"
	"github.com/thealamu/internal/engine"
)

//NewGameHandler handles creating a new game
func NewGameHandler() http.HandlerFunc {
	gE := engine.New()
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "New Game")
	}
}
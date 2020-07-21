package handler

import (
	"fmt"
	"net/http"

	"github.com/patrickmn/go-cache"
)

//JoinGameHandler handles game joining
func JoinGameHandler(db *cache.Cache) http.HandlerFunc {
	//gE := engine.New()
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Join Game")
	}
}

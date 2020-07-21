package handler

import (
	"fmt"
	"net/http"
)

//HomeHandler handles the default path
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "TTTGameServer")
}

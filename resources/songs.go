package resources

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// GetSongs gets the list of songs in the library (HTTP GET)
func GetSongs(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
	fmt.Fprint(w, "Hhihihi")
}

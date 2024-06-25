package serv

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// WriteJSON sends JSON back with the right headers to the [http.ResponseWriter].
//
// In case there's an error, it writes the error back and returns a 500.
func WriteJSON(w http.ResponseWriter, a any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(a); err != nil {
		slog.Error("error encoding response", "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

package helpers

import (
	"encoding/json"
	"io"
	"net/http"
)

func Ummarshal(r *http.Request, x interface{}) error {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		PanicHelper(err)
	}
	if err := json.Unmarshal(body, x); err != nil {
		return err
	}
	return nil
}

func WriteJSON(w http.ResponseWriter, x interface{}) {
	res, err := json.MarshalIndent(x, "", "  ")
	PanicHelper(err)
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}

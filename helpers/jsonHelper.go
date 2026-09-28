package helpers

import (
	"encoding/json"
	"io"
	"net/http"

	"go-food-store/json/response"
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

// WriteJSON menulis envelope response sekaligus status HTTP-nya.
//
// Sengaja menerima response.WebResponse, bukan interface{}: sebelumnya semua
// handler membalas 200 karena tidak ada yang memanggil WriteHeader, dan kode
// error hanya hidup di body. Client tidak bisa membedakan sukses dari gagal,
// dan CI pun tidak — `curl -f` lolos untuk route yang rusak. Dengan tipe
// konkret, status selalu berasal dari satu sumber: Code.
func WriteJSON(w http.ResponseWriter, res response.WebResponse) {
	body, err := json.MarshalIndent(res, "", "  ")
	PanicHelper(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	w.Write(body)
}

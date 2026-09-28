package request

import (
	"net/http"
	"strconv"
)

const (
	defaultLimit = 20
	maxLimit     = 100
)

// Page adalah potongan daftar yang diminta client.
//
// ponytail: tanpa jumlah total. Menghitungnya berarti satu query COUNT
// tambahan di setiap permintaan daftar, dan yang dibutuhkan halaman menu
// hanyalah "masih ada lagi atau tidak" — yang sudah terjawab oleh jumlah baris
// yang kembali sama dengan Limit. Tambahkan total kalau UI benar-benar perlu
// menampilkan nomor halaman terakhir.
type Page struct {
	Number int
	Limit  int
	Offset int
}

// PageFrom membaca ?page= dan ?limit= dari query string.
//
// Nilai yang tidak masuk akal dibetulkan diam-diam alih-alih ditolak: daftar
// adalah endpoint baca, dan membalas 400 untuk "?page=abc" hanya menyulitkan
// tanpa melindungi apa pun. Yang penting Limit punya batas atas, supaya satu
// permintaan tidak bisa menarik seluruh tabel.
func PageFrom(r *http.Request) Page {
	number, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || number < 1 {
		number = 1
	}
	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = defaultLimit
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return Page{
		Number: number,
		Limit:  limit,
		Offset: (number - 1) * limit,
	}
}

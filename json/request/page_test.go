package request

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPageFrom(t *testing.T) {
	cases := []struct {
		query  string
		number int
		limit  int
		offset int
	}{
		{"", 1, 20, 0},
		{"?page=3", 3, 20, 40},
		{"?limit=5", 1, 5, 0},
		{"?page=2&limit=5", 2, 5, 5},
		// Nilai tidak masuk akal dibetulkan, bukan ditolak.
		{"?page=0", 1, 20, 0},
		{"?page=-4", 1, 20, 0},
		{"?limit=0", 1, 20, 0},
		{"?page=abc&limit=xyz", 1, 20, 0},
		// Batas atas menahan satu permintaan menarik seluruh tabel.
		{"?limit=100000", 1, 100, 0},
	}

	for _, c := range cases {
		t.Run(c.query, func(t *testing.T) {
			got := PageFrom(httptest.NewRequest(http.MethodGet, "/food/"+c.query, nil))
			if got.Number != c.number || got.Limit != c.limit || got.Offset != c.offset {
				t.Errorf("= %+v, mau {Number:%d Limit:%d Offset:%d}", got, c.number, c.limit, c.offset)
			}
		})
	}
}

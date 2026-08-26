package config

import (
	"os"
	"strings"
	"testing"
)

// clearDBEnv menghapus semua variabel DB_* dan memulihkannya setelah test.
// t.Setenv tidak bisa dipakai di sini: ia hanya bisa men-SET, bukan menghapus,
// sehingga tidak bisa membentuk environment kosong yang jadi premis test ini.
//
// Tanpa ini, developer yang baru menyalin .env.example ke shell-nya akan
// melihat suite merah di checkout yang bersih — justru orang yang paling
// mungkin sedang mencoba repo ini.
func clearDBEnv(t *testing.T) {
	t.Helper()
	for _, key := range []string{"DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME"} {
		old, ok := os.LookupEnv(key)
		if !ok {
			continue
		}
		// key dan old di-capture eksplisit supaya tidak bergantung pada
		// aturan per-iterasi Go 1.22+ untuk bisa dibaca dengan benar.
		key, old := key, old
		t.Cleanup(func() { os.Setenv(key, old) })
		os.Unsetenv(key)
	}
}

func TestDSNMemakaiDefaultSaatEnvKosong(t *testing.T) {
	clearDBEnv(t)

	got := DSN()

	for _, want := range []string{"root", "tcp(localhost:3306)", "go-food-store", "parseTime=True"} {
		if !strings.Contains(got, want) {
			t.Errorf("DSN() = %q, tidak memuat %q", got, want)
		}
	}
}

func TestDSNMembacaEnv(t *testing.T) {
	t.Setenv("DB_USER", "appuser")
	t.Setenv("DB_PASSWORD", "s3cret")
	t.Setenv("DB_HOST", "mysql")
	t.Setenv("DB_PORT", "3307")
	t.Setenv("DB_NAME", "foodstore")

	got := DSN()

	want := "appuser:s3cret@tcp(mysql:3307)/foodstore?charset=utf8mb4&parseTime=True&loc=Local"
	if got != want {
		t.Errorf("DSN() = %q, mau %q", got, want)
	}
}

func TestDSNMengizinkanPasswordKosong(t *testing.T) {
	t.Setenv("DB_USER", "root")
	t.Setenv("DB_PASSWORD", "")

	got := DSN()

	if !strings.HasPrefix(got, "root:@tcp(") {
		t.Errorf("DSN() = %q, mau berawalan %q", got, "root:@tcp(")
	}
}

// Test di atas mengunci bentuk DSN untuk password kosong, tapi TIDAK bisa
// membedakan os.LookupEnv dari os.Getenv: fallback DB_PASSWORD memang ""
// juga, jadi kedua implementasi menghasilkan string yang sama.
//
// Yang membedakan keduanya adalah variabel yang fallback-nya TIDAK kosong.
// DB_HOST di-set kosong secara eksplisit harus tetap kosong, bukan diganti
// "localhost" — itulah arti "string kosong yang di-set dianggap disengaja".
// Dengan os.Getenv test ini gagal.
func TestDSNMenghormatiNilaiKosongYangEksplisit(t *testing.T) {
	// Assertion "@tcp(:3306)/" membaca DB_PORT dari environment, jadi test ini
	// ikut terpapar pollution meskipun ia men-set DB_HOST sendiri.
	clearDBEnv(t)
	t.Setenv("DB_HOST", "")

	got := DSN()

	if !strings.Contains(got, "@tcp(:3306)/") {
		t.Errorf("DSN() = %q, mau memuat %q", got, "@tcp(:3306)/")
	}
	if strings.Contains(got, "localhost") {
		t.Errorf("DSN() = %q, DB_HOST kosong yang eksplisit malah diganti fallback", got)
	}
}

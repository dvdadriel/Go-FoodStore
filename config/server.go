package config

import "time"

// ServerAddr mengembalikan alamat bind HTTP. Default ":8080" (semua interface),
// bukan "localhost:8080" — di dalam container, mengikat loopback membuat port
// mapping Docker tidak tembus dari host.
func ServerAddr() string { return env("SERVER_ADDR", ":8080") }

// CORSOrigin menentukan origin mana yang boleh memanggil API ini dari browser.
// Default "*" cocok untuk pengembangan; di produksi set ke origin frontend.
func CORSOrigin() string { return env("CORS_ORIGIN", "*") }

// Batas waktu server. Tanpa ini sebuah koneksi yang menggantung memegang
// goroutine dan file descriptor selamanya — satu klien lambat cukup untuk
// menghabiskan keduanya (slowloris).
const (
	ReadHeaderTimeout = 5 * time.Second
	ReadTimeout       = 15 * time.Second
	WriteTimeout      = 15 * time.Second
	IdleTimeout       = 60 * time.Second
	ShutdownTimeout   = 10 * time.Second
)

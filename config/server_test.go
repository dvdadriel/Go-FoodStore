package config

import (
	"os"
	"testing"
)

// Isolasi yang sama seperti clearDBEnv, tapi untuk SERVER_ADDR: test default
// harus benar-benar membentuk environment tanpa SERVER_ADDR, bukan hanya
// mengasumsikannya. Pemakainya tidak boleh memanggil t.Parallel().
func clearServerEnv(t *testing.T) {
	t.Helper()
	const key = "SERVER_ADDR"
	old, ok := os.LookupEnv(key)
	if !ok {
		return
	}
	t.Cleanup(func() {
		if err := os.Setenv(key, old); err != nil {
			t.Errorf("memulihkan %s: %v", key, err)
		}
	})
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("menghapus %s: %v", key, err)
	}
}

func TestServerAddrDefault(t *testing.T) {
	clearServerEnv(t)

	if got := ServerAddr(); got != ":8080" {
		t.Errorf("ServerAddr() = %q, mau %q", got, ":8080")
	}
}

func TestServerAddrDariEnv(t *testing.T) {
	t.Setenv("SERVER_ADDR", "127.0.0.1:9000")

	if got := ServerAddr(); got != "127.0.0.1:9000" {
		t.Errorf("ServerAddr() = %q, mau %q", got, "127.0.0.1:9000")
	}
}

// Default-nya ":8080" (semua interface), bukan "localhost:8080": di dalam
// container, mengikat loopback membuat port mapping Docker tidak tembus dari
// host. Ini yang justru jadi alasan task ini ada, jadi dikunci eksplisit.
func TestServerAddrDefaultTidakTerikatLoopback(t *testing.T) {
	clearServerEnv(t)

	if got := ServerAddr(); got == "localhost:8080" || got == "127.0.0.1:8080" {
		t.Errorf("ServerAddr() = %q, terikat loopback sehingga port mapping Docker tidak tembus", got)
	}
}

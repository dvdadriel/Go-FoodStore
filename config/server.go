package config

// ServerAddr mengembalikan alamat bind HTTP. Default ":8080" (semua interface),
// bukan "localhost:8080" — di dalam container, mengikat loopback membuat port
// mapping Docker tidak tembus dari host.
func ServerAddr() string { return env("SERVER_ADDR", ":8080") }

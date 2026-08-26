package config

import "os"

// env mengembalikan nilai environment variable, atau fallback kalau tidak
// di-set.
//
// Sengaja os.LookupEnv, bukan os.Getenv: string kosong yang di-set secara
// eksplisit dianggap nilai yang disengaja, bukan "tidak di-set". Password
// kosong itu sah untuk MySQL root di lokal, dan lebih umum, DB_HOST="" harus
// tetap kosong alih-alih diam-diam diganti "localhost".
//
// Helper ini dipakai bersama oleh DSN() dan ServerAddr(), karena itu ia tinggal
// di file sendiri: tidak jadi milik database.go maupun server.go.
func env(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

#!/usr/bin/env bash
# Menembak API yang sedang berjalan dan mencetak request/response apa adanya.
# Tidak ada assertion di sini selain HTTP request-nya berhasil dikirim:
# tugas file ini menampilkan yang terjadi, bukan menilainya. Yang menilai
# adalah `go test`.
set -euo pipefail

BASE=${BASE:-http://localhost:8080}
ADMIN_USERNAME=${ADMIN_USERNAME:-admin}
ADMIN_PASSWORD=${ADMIN_PASSWORD:-admin12345}

TOKEN=""

show() { # show METHOD PATH [JSON_BODY]
  local method=$1 path=$2 body=${3:-}

  printf '$ curl -i -X %s %s%s' "$method" "$BASE" "$path"
  [ -n "$TOKEN" ] && printf " \\\\\n    -H 'Authorization: Bearer <token>'"
  [ -n "$body" ] && printf " \\\\\n    -H 'Content-Type: application/json' \\\\\n    -d '%s'" "$body"
  printf '\n\n'

  local args=(-sS -i -X "$method" "$BASE$path")
  [ -n "$TOKEN" ] && args+=(-H "Authorization: Bearer $TOKEN")
  [ -n "$body" ] && args+=(-H 'Content-Type: application/json' -d "$body")

  curl "${args[@]}"

  printf '\n\n---\n\n'
}

echo "Dijalankan: $(date -u '+%Y-%m-%d %H:%M UTC')"
echo "Commit:     ${GITHUB_SHA:-$(git rev-parse --short HEAD)}"
echo
echo '---'
echo

# Tanpa token dulu: membuktikan endpoint tertutup memang tertutup, dan
# etalase menu memang terbuka.
show GET /health
show GET /food/
show POST /food/ '{"FoodName":"Nasi Goreng","FoodPrice":25000}'

# Login sebagai admin. Token tidak dicetak; yang penting statusnya.
TOKEN=$(curl -sS -X POST "$BASE/auth/login" \
  -H 'Content-Type: application/json' \
  -d "{\"Username\":\"$ADMIN_USERNAME\",\"Password\":\"$ADMIN_PASSWORD\"}" \
  | sed -n 's/.*"Token": "\([^"]*\)".*/\1/p')

if [ -z "$TOKEN" ]; then
  echo 'Login admin gagal, tidak bisa melanjutkan.' >&2
  exit 1
fi
echo 'Login admin berhasil, permintaan berikutnya membawa bearer token.'
echo
echo '---'
echo

# Alur utuh, bukan endpoint terpisah: buat customer, buat makanan,
# buat transaksi yang merujuk keduanya, lalu terima pembayarannya.
# Nama field mengikuti json/request/ apa adanya — PascalCase, bukan snake_case.
show POST /cust/ '{"CustomerName":"Bukti CI"}'
show GET /cust/

show POST /food/ '{"FoodName":"Nasi Goreng","FoodPrice":25000}'
show GET /food/

show POST /transaction/ '{"CustomerId":1,"Food":[{"FoodId":1,"Quantity":2}]}'
show GET /transaction/unpaid/
show PUT /transaction/acc/1
show GET /transaction/paid/

# Transaksi yang sudah lunas terkunci.
show DELETE /transaction/1

# Permintaan yang ditolak. Ini bagian terpenting halaman ini: lihat status
# barisnya, lalu lihat isi body-nya.
show GET /food/999999
show POST /transaction/ '{"CustomerId":1,"Food":[{"FoodId":1,"Quantity":0}]}'
show POST /transaction/ '{"CustomerId":1,"Food":[{"FoodId":1,"Quantity":1},{"FoodId":1,"Quantity":2}]}'

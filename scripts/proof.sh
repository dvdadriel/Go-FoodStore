#!/usr/bin/env bash
# Menembak API yang sedang berjalan dan mencetak request/response apa adanya.
# Tidak ada assertion di sini selain HTTP request-nya berhasil dikirim:
# tugas file ini menampilkan yang terjadi, bukan menilainya. Yang menilai
# adalah `go test`.
set -euo pipefail

BASE=${BASE:-http://localhost:8080}

show() { # show METHOD PATH [JSON_BODY]
  local method=$1 path=$2 body=${3:-}

  printf '$ curl -i -X %s %s%s' "$method" "$BASE" "$path"
  [ -n "$body" ] && printf " \\\\\n    -H 'Content-Type: application/json' \\\\\n    -d '%s'" "$body"
  printf '\n\n'

  if [ -n "$body" ]; then
    curl -sS -i -X "$method" "$BASE$path" \
      -H 'Content-Type: application/json' -d "$body"
  else
    curl -sS -i -X "$method" "$BASE$path"
  fi

  printf '\n\n---\n\n'
}

echo "Dijalankan: $(date -u '+%Y-%m-%d %H:%M UTC')"
echo "Commit:     ${GITHUB_SHA:-$(git rev-parse --short HEAD)}"
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

# Sengaja diminta yang tidak ada. Ini bagian terpenting halaman ini:
# lihat status barisnya, lalu lihat isi body-nya.
show GET /food/999999

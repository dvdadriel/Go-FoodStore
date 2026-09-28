# Go-FoodStore

REST API pemesanan makanan dengan Go, disusun berlapis: **controller → service → repository**. Setiap layer punya interface terpisah dari implementasinya, sehingga service bisa diuji tanpa database sama sekali.

[![CI](https://github.com/dvdadriel/Go-FoodStore/actions/workflows/ci.yml/badge.svg)](https://github.com/dvdadriel/Go-FoodStore/actions/workflows/ci.yml)

**Stack:** Go 1.24 · gorm · gorilla/mux · MySQL 8.4 · go-playground/validator

## Cara menjalankan

```bash
git clone https://github.com/dvdadriel/Go-FoodStore.git
cd Go-FoodStore
docker compose up --build
```

API tersedia di `http://localhost:8080`. Skema database dibuat otomatis lewat `AutoMigrate` saat start, jadi tidak ada langkah migrasi manual.

Tanpa Docker, dengan MySQL lokal:

```bash
cp .env.example .env
set -a; source .env; set +a    # aplikasi membaca environment proses, bukan file .env
go run .
```

## Arsitektur

```
main.go
  └── config.SetupModel(db, validate)        wiring dependensi (DI manual)
        ├── repositories/                    akses data lewat gorm
        ├── services/                        validasi + aturan bisnis
        └── controllers/                     parsing HTTP + penulisan response
              └── routes.NewRouter(...)      pemetaan method + path
```

Setiap layer punya pasangan interface dan implementasi (`FoodRepo` / `FoodRepoImpl`, `FoodService` / `FoodServiceImpl`). Yang dipakai konsumen adalah interface-nya, jadi service bisa diuji dengan stub repository — tanpa database, tanpa library mock.

Konfigurasi dibaca dari environment (`config/env.go`), bukan hardcoded:

| Variabel | Default | Keterangan |
|---|---|---|
| `DB_USER` | `root` | |
| `DB_PASSWORD` | *(kosong)* | Nilai kosong yang di-set eksplisit dihormati, tidak diganti default |
| `DB_HOST` | `localhost` | |
| `DB_PORT` | `3306` | |
| `DB_NAME` | `go-food-store` | |
| `SERVER_ADDR` | `:8080` | Semua interface, bukan `localhost` — kalau loopback, port mapping Docker tidak tembus |
| `CORS_ORIGIN` | `*` | Origin frontend yang diizinkan |
| `JWT_SECRET` | *(wajib)* | Minimal 32 karakter. Aplikasi menolak start tanpa ini |
| `TOKEN_TTL_SECONDS` | `3600` | Masa berlaku token akses |
| `ADMIN_USERNAME` | *(kosong)* | Admin pertama, dibuat saat start kalau belum ada |
| `ADMIN_PASSWORD` | *(kosong)* | Kalau salah satunya kosong, admin tidak dibuat |

## Endpoint

Semua response memakai envelope `{ Code, Status, Message, Data }`.

| Method | Path | Akses | Keterangan |
|---|---|---|---|
| `GET` | `/health` | publik | Status proses dan koneksi database |
| `POST` | `/auth/register` | publik | Daftar akun baru (selalu peran `customer`) |
| `POST` | `/auth/login` | publik | Tukar kredensial dengan token |
| `GET` | `/food/` | publik | Semua menu |
| `GET` | `/food/{foodId}` | publik | Menu per id |
| `GET` | `/cust/` | admin | Semua customer |
| `POST` | `/cust/` | admin | Buat customer |
| `GET` | `/cust/{custId}` | admin | Customer per id |
| `PUT` | `/cust/{custId}` | admin | Update customer |
| `DELETE` | `/cust/{custId}` | admin | Hapus customer |
| `POST` | `/food/` | admin | Buat menu |
| `PUT` | `/food/{foodId}` | admin | Update menu |
| `DELETE` | `/food/{foodId}` | admin | Hapus menu |
| `GET` | `/transaction/unpaid/` | admin | Transaksi belum lunas |
| `GET` | `/transaction/paid/` | admin | Transaksi sudah lunas |
| `POST` | `/transaction/` | login | Buat transaksi |
| `GET` | `/transaction/{transactionId}` | login | Transaksi per id |
| `PUT` | `/transaction/{transactionId}` | login | Update transaksi |
| `DELETE` | `/transaction/{transactionId}` | login | Hapus transaksi |
| `PUT` | `/transaction/acc/{transactionId}` | admin | Tandai transaksi lunas |

Perhatikan **trailing slash** pada route koleksi (`/food/`, bukan `/food`).

Semua endpoint daftar menerima `?page=` dan `?limit=` (default 20, maksimum 100). Nilai yang tidak masuk akal dibetulkan diam-diam, bukan ditolak — `?page=abc` sama dengan halaman pertama. Response tidak membawa jumlah total: menghitungnya berarti satu query `COUNT` tambahan di tiap permintaan, sementara "masih ada lagi atau tidak" sudah terjawab oleh jumlah baris yang kembali sama dengan `limit`.

```bash
$ curl -s 'localhost:8080/food/?page=2&limit=5'
```

Postman collection: [`docs/Go-FoodStore.postman_collection.json`](docs/Go-FoodStore.postman_collection.json)

## Autentikasi

Token bertanda tangan HMAC-SHA256, dikirim lewat header `Authorization: Bearer <token>`. Kata sandi disimpan sebagai PBKDF2-HMAC-SHA256 (210.000 iterasi, salt acak 16 byte per pengguna).

Keduanya memakai pustaka standar saja — `crypto/pbkdf2` masuk stdlib di Go 1.24 — jadi tidak ada dependensi baru untuk ini. Formatnya bukan JWT: satu algoritma, tanpa header algoritma yang bisa dipalsukan, sekitar lima puluh baris di `auth/token.go`. Ganti ke JWT betulan kalau nanti ada layanan lain yang harus ikut memverifikasi token yang sama.

```bash
$ curl -s -X POST localhost:8080/auth/login -H 'Content-Type: application/json' \
    -d '{"Username":"admin","Password":"admin12345"}'
{
  "Code": 200, "Status": "OK", "Message": "Successfully login",
  "Data": {
    "Token": "eyJzdWIiOjEsInVzciI6ImFkbWluIiwicm9sZSI6ImFkbWluIiwiZXhwIjoxNzU5MDAwMDAwfQ.T0lz...",
    "ExpiresIn": 3600,
    "User": { "Id": 1, "Username": "admin", "Role": "admin" }
  }
}
```

Dua peran: **admin** mengelola menu, customer, dan pembayaran; **customer** membuat dan melihat transaksi. `/auth/register` selalu menghasilkan peran customer — kalau peran bisa diminta sendiri saat mendaftar, siapa pun tinggal mendaftar sebagai admin. Admin pertama lahir dari `ADMIN_USERNAME`/`ADMIN_PASSWORD` saat aplikasi start.

Tanpa token, endpoint terjaga membalas `401`; dengan token yang sah tapi peran yang salah, `403`:

```bash
$ curl -s -o /dev/null -w '%{http_code}\n' -X POST localhost:8080/food/
401
$ curl -s -o /dev/null -w '%{http_code}\n' -X POST localhost:8080/food/ -H "Authorization: Bearer $TOKEN_CUSTOMER"
403
```

### Contoh alur lengkap

Output di bawah ini diambil dari stack yang benar-benar berjalan, bukan dikarang.

```bash
$ curl -s -X POST localhost:8080/food/ -H "Authorization: Bearer $TOKEN" \
    -H 'Content-Type: application/json' -d '{"FoodName":"Nasi Goreng","FoodPrice":25000}'
{ "Code": 200, "Status": "OK", "Message": "Successfully create new food", "Data": null }

$ curl -s -X POST localhost:8080/cust/ -H "Authorization: Bearer $TOKEN" \
    -H 'Content-Type: application/json' -d '{"CustomerName":"David Adriel"}'
{ "Code": 200, "Status": "OK", "Message": "Successfully create new customer", "Data": null }

$ curl -s -X POST localhost:8080/transaction/ -H "Authorization: Bearer $TOKEN" \
    -H 'Content-Type: application/json' -d '{"CustomerId":1,"Food":[{"FoodId":1,"Quantity":2}]}'
{ "Code": 200, "Status": "OK", "Message": "Successfully create new transaction", "Data": { "TransactionId": 1 } }

$ curl -s localhost:8080/transaction/unpaid/ -H "Authorization: Bearer $TOKEN"
{
  "Code": 200, "Status": "OK", "Message": "Successfully get all unpaid data",
  "Data": [{
    "TransactionId": 1,
    "Customer": { "Id": 1, "CustomerName": "David Adriel" },
    "Food": [{ "FoodId": 1, "FoodName": "Nasi Goreng", "FoodPrice": 25000, "Quantity": 2, "Subtotal": 50000 }],
    "TotalPrice": 50000,
    "AlreadyPay": false
  }]
}

$ curl -s -X PUT localhost:8080/transaction/acc/1 -H "Authorization: Bearer $TOKEN"
{ "Code": 200, "Status": "OK", "Message": "Successfully update transaction status", "Data": null }
```

Setelahnya transaksi berpindah dari `/transaction/unpaid/` ke `/transaction/paid/`.

Id yang tidak dikenal, route tak terdaftar, dan method salah masing-masing tertangani:

```bash
$ curl -s localhost:8080/food/999
{ "Code": 404, "Status": "Not Found", "Message": "Food not found", "Data": null }

$ curl -s localhost:8080/nope
{ "Code": 404, "Status": "Not Found", "Message": "Data not found", "Data": null }

$ curl -s -X PATCH localhost:8080/food/
{ "Code": 405, "Status": "Method Not Allowed", "Message": "Wrong method", "Data": null }
```

## Test

```bash
go test ./... -race -cover
```

Layer `config` dan `services/food_service` berada di **100% statement coverage**; controller dan repository warisan belum punya test, sehingga total repo sekitar 5,9%. CI mencetak kedua angka itu — mencetak hanya totalnya menutupi kerja yang sudah dilakukan, mencetak hanya yang 100% terbaca sebagai cherry-picking.

Unit test service memakai **stub repository yang ditulis tangan** (`services/food_service/stub_repo_test.go`), bukan library mock. Dua hal yang membuatnya lebih dari sekadar mengejar angka coverage:

- **Field fungsi yang dibiarkan nil akan memanggil `t.Fatalf`** kalau method-nya tersentuh. Dengan begitu test bisa membuktikan sebuah method justru **tidak** dipanggil — misalnya `DeleteFood` tidak boleh jalan saat data tidak ditemukan. `t.Helper()` membuat kegagalannya diatribusikan ke call site di kode produksi, bukan ke baris test.
- **Setiap stub mengembalikan `Message` yang berbeda** (`"dari CreateFood"`, `"dari GetFoodById"`, …) dan test meng-assert mana yang kembali. Tanpa ini, "service meneruskan response repository apa adanya" tidak teruji — service bisa mengarang response 200 sendiri dan semua assertion `Code` tetap lolos.

Suite ini diuji dengan mutation testing: dari 25 mutasi yang disuntikkan ke kode produksi, 20 tertangkap, dan setiap mutasi alur kontrol tertangkap.

### Transaksi yang sudah lunas tidak bisa diubah

`PUT /transaction/{id}` dan `DELETE /transaction/{id}` mengembalikan `409 Conflict` kalau transaksi sudah ditandai lunas:

```bash
$ curl -s -X DELETE localhost:8080/transaction/1 -H "Authorization: Bearer $TOKEN"
{ "Code": 409, "Status": "Conflict", "Message": "Can't delete paid transaction", "Data": null }
```

Ini logika bisnis yang disengaja, bukan error. Sebelumnya kasus ini membalas `304 Not Modified`; 304 dilarang membawa body oleh spesifikasi HTTP, jadi alasan penolakannya hilang di jalan.

### Harga disalin saat transaksi dibuat

`transaction_foods.unit_price` menyimpan harga yang berlaku pada saat pesanan dibuat. Mengubah harga menu tidak mengubah nilai transaksi yang sudah terjadi — struk kemarin tetap berbunyi sama hari ini. `TotalPrice` dihitung dari harga tercatat itu, bukan dari harga menu sekarang.

## Batasan yang diketahui

Dicantumkan terbuka karena ini keputusan sadar atau cacat yang sudah diketahui, bukan hal yang terlewat.

- **Repository mengembalikan `response.WebResponse`.** Artinya layer data mengetahui status code HTTP — pencampuran tanggung jawab. Yang lebih benar: repository mengembalikan `(data, error)` dan pemetaan ke HTTP dilakukan di controller. Belum diubah karena menyentuh ketiga domain sekaligus.
- **`FindById` mengabaikan flag `found`** dari repository dan meneruskan response apa pun yang diterima. Hasilnya kebetulan benar karena repository sudah mengisi response 404, tapi kebenarannya bergantung pada kebetulan itu. Ada test yang mendokumentasikan perilaku ini — bukan membenarkannya.
- **Pesan validasi tidak konsisten**: `Create` memakai `"Please check the request"`, `Update` memakai `"Please check your request"`. Sekarang keduanya terpaku oleh test, jadi menyeragamkannya kelak adalah perubahan kontrak yang disengaja, bukan diam-diam.
- **`go-playground/validator` masih v9** dan sudah tidak dipelihara. Migrasi ke v10 menyentuh seluruh service dan controller, jadi sengaja di luar cakupan. Efek sampingnya: v9 adalah module `+incompatible` tanpa `go.mod` sendiri, sehingga `go mod tidy` menarik dependensi test-only miliknya (`gopkg.in/go-playground/assert.v1`) menjadi entri indirect di sini. Ia tidak ada di build graph (`go list -deps ./...` tidak memuatnya); menghapusnya manual akan dibatalkan oleh pemeriksaan `go mod tidy` di CI.
- **Kepemilikan transaksi belum dicek.** Setiap pengguna yang login bisa membaca dan mengubah transaksi milik siapa pun, karena tabel `users` belum terhubung ke tabel `customers`. Peran sudah memisahkan admin dari customer; kepemilikan per baris belum.
- **Token tidak bisa dicabut.** Token yang bocor tetap berlaku sampai `TOKEN_TTL_SECONDS` habis. Selama masa berlakunya satu jam ini masih sepadan; kalau dipanjangkan, perlu daftar cabut.
- **`helpers.PanicHelper` melakukan panic pada error koneksi**, sehingga aplikasi mati saat start kalau database belum siap. Karena itu `docker-compose.yml` memakai healthcheck pada MySQL dan `api` menunggu `service_healthy` — tanpa itu container `api` panic sebelum MySQL menerima koneksi.

## Lisensi

MIT — lihat [LICENSE](LICENSE).

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

## Endpoint

Semua response memakai envelope `{ Code, Status, Message, Data }`.

| Method | Path | Keterangan |
|---|---|---|
| `GET` | `/cust/` | Semua customer |
| `POST` | `/cust/` | Buat customer |
| `GET` | `/cust/{custId}` | Customer per id |
| `PUT` | `/cust/{custId}` | Update customer |
| `DELETE` | `/cust/{custId}` | Hapus customer |
| `GET` | `/food/` | Semua menu |
| `POST` | `/food/` | Buat menu |
| `GET` | `/food/{foodId}` | Menu per id |
| `PUT` | `/food/{foodId}` | Update menu |
| `DELETE` | `/food/{foodId}` | Hapus menu |
| `GET` | `/transaction/unpaid/` | Transaksi belum lunas |
| `GET` | `/transaction/paid/` | Transaksi sudah lunas |
| `POST` | `/transaction/` | Buat transaksi |
| `GET` | `/transaction/{transactionId}` | Transaksi per id |
| `PUT` | `/transaction/{transactionId}` | Update transaksi |
| `DELETE` | `/transaction/{transactionId}` | Hapus transaksi |
| `PUT` | `/transaction/acc/{transactionId}` | Tandai transaksi lunas |

Perhatikan **trailing slash** pada route koleksi (`/food/`, bukan `/food`).

Postman collection: [`docs/Go-FoodStore.postman_collection.json`](docs/Go-FoodStore.postman_collection.json)

### Contoh alur lengkap

Output di bawah ini diambil dari stack yang benar-benar berjalan, bukan dikarang.

```bash
$ curl -s -X POST localhost:8080/food/ -H 'Content-Type: application/json' \
    -d '{"FoodName":"Nasi Goreng","FoodPrice":25000}'
{ "Code": 200, "Status": "OK", "Message": "Successfully create new food", "Data": null }

$ curl -s -X POST localhost:8080/cust/ -H 'Content-Type: application/json' \
    -d '{"CustomerName":"David Adriel"}'
{ "Code": 200, "Status": "OK", "Message": "Successfully create new customer", "Data": null }

$ curl -s -X POST localhost:8080/transaction/ -H 'Content-Type: application/json' \
    -d '{"CustomerId":1,"Food":[{"FoodId":1,"Quantity":2}]}'
{ "Code": 200, "Status": "OK", "Message": "Successfully create new transaction", "Data": null }

$ curl -s localhost:8080/transaction/unpaid/
{
  "Code": 200, "Status": "OK", "Message": "Successfully get all unpaid data",
  "Data": [{
    "TransactionId": 1,
    "Customer": { "Id": 1, "CustomerName": "David Adriel" },
    "Food": [{ "FoodId": 1, "FoodName": "Nasi Goreng", "FoodPrice": 25000, "Quantity": 2 }],
    "AlreadyPay": false
  }]
}

$ curl -s -X PUT localhost:8080/transaction/acc/1
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

`PUT /transaction/{id}` dan `DELETE /transaction/{id}` mengembalikan `304` kalau transaksi sudah ditandai lunas:

```bash
$ curl -s -X DELETE localhost:8080/transaction/1
{ "Code": 304, "Status": "Not modified", "Message": "Can't delete paid transaction", "Data": null }
```

Ini logika bisnis yang disengaja, bukan error.

## Batasan yang diketahui

Dicantumkan terbuka karena ini keputusan sadar atau cacat yang sudah diketahui, bukan hal yang terlewat.

- **🔴 Status code HTTP selalu 200.** Tidak ada satu pun pemanggilan `w.WriteHeader` di seluruh repo, jadi status code hanya ditulis ke dalam body JSON, bukan ke response HTTP-nya:

  ```bash
  $ curl -s -o /dev/null -w '%{http_code}\n' localhost:8080/food/999
  200
  $ curl -s localhost:8080/food/999
  { "Code": 404, ... }
  ```

  Konsekuensinya nyata: klien yang memeriksa status HTTP — hampir semua HTTP client, load balancer, dan tooling monitoring — **tidak akan pernah bisa mendeteksi error di API ini**. Ini cacat paling serius yang tersisa. Perbaikannya terpusat (satu helper penulis response yang memanggil `WriteHeader(res.Code)` sebelum menulis body), tapi mengubahnya adalah perubahan kontrak yang memengaruhi setiap endpoint, jadi dipisahkan dari pekerjaan ini agar bisa diuji tersendiri.

- **Pesan `304` pada jalur update menyebut "delete".** `repositories/transaction_repository/transaction_repo_impl.go:394` mengembalikan `"Can't delete paid transaction"` padahal itu operasi update — salah tempel. `Status`-nya juga beda kapitalisasi antar dua tempat (`"Not modified"` di jalur delete, `"Not Modified"` di jalur update).

- **Repository mengembalikan `response.WebResponse`.** Artinya layer data mengetahui status code HTTP — pencampuran tanggung jawab. Yang lebih benar: repository mengembalikan `(data, error)` dan pemetaan ke HTTP dilakukan di controller. Belum diubah karena menyentuh ketiga domain sekaligus.
- **`FindById` mengabaikan flag `found`** dari repository dan meneruskan response apa pun yang diterima. Hasilnya kebetulan benar karena repository sudah mengisi response 404, tapi kebenarannya bergantung pada kebetulan itu. Ada test yang mendokumentasikan perilaku ini — bukan membenarkannya.
- **Pesan validasi tidak konsisten**: `Create` memakai `"Please check the request"`, `Update` memakai `"Please check your request"`. Sekarang keduanya terpaku oleh test, jadi menyeragamkannya kelak adalah perubahan kontrak yang disengaja, bukan diam-diam.
- **`go-playground/validator` masih v9** dan sudah tidak dipelihara. Migrasi ke v10 menyentuh seluruh service dan controller, jadi sengaja di luar cakupan. Efek sampingnya: v9 adalah module `+incompatible` tanpa `go.mod` sendiri, sehingga `go mod tidy` menarik dependensi test-only miliknya (`gopkg.in/go-playground/assert.v1`) menjadi entri indirect di sini. Ia tidak ada di build graph (`go list -deps ./...` tidak memuatnya); menghapusnya manual akan dibatalkan oleh pemeriksaan `go mod tidy` di CI.
- **Tidak ada autentikasi.** Semua endpoint terbuka.
- **`helpers.PanicHelper` melakukan panic pada error koneksi**, sehingga aplikasi mati saat start kalau database belum siap. Karena itu `docker-compose.yml` memakai healthcheck pada MySQL dan `api` menunggu `service_healthy` — tanpa itu container `api` panic sebelum MySQL menerima koneksi.

## Lisensi

MIT — lihat [LICENSE](LICENSE).

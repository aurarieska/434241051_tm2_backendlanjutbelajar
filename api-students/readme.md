# API Mahasiswa (Go + Fiber)

REST API data mahasiswa, disimpan di memori (tanpa database).
Server: `http://localhost:3000`
Jalankan dengan: `go run .`

## Kontrak API

1. **Daftar mahasiswa**
   - Metode: GET
   - Endpoint: `/api/v1/students`
   - Parameter (query, semua opsional):
     - `page`: nomor halaman, default 1
     - `limit`: jumlah data per halaman, default 10, maksimal 100
     - `search`: cari pada nama, tidak peka huruf besar/kecil
     - `sort`: `id`, `nim`, `name`, atau `grade` (default `id`)
     - `order`: `asc` atau `desc` (default `asc`)
     - `is_active`: `true` atau `false`
     - `min_grade` dan `max_grade`: rentang nilai
   - Contoh body permintaan: tidak ada
   - Status: 200
   - Contoh respons (200):
```json
     {"success":true,"message":"daftar mahasiswa berhasil diambil","data":[{"id":1,"nim":"2024001","name":"Budi Santoso","grade":85,"is_active":true}],"meta":{"page":1,"limit":10,"total":1,"total_pages":1}}
```

2. **Satu mahasiswa**
   - Metode: GET
   - Endpoint: `/api/v1/students/:id`
   - Parameter: `id` (path, angka positif)
   - Contoh body permintaan: tidak ada
   - Status: 200, 400, 404
   - Contoh respons:
     - 200:
```json
       {"success":true,"message":"mahasiswa ditemukan","data":{"id":1,"nim":"2024001","name":"Budi Santoso","grade":85,"is_active":true}}
```
     - 400:
```json
       {"success":false,"message":"id harus berupa angka positif"}
```
     - 404:
```json
       {"success":false,"message":"mahasiswa tidak ditemukan"}
```

3. **Tambah mahasiswa**
   - Metode: POST
   - Endpoint: `/api/v1/students`
   - Parameter: header `Content-Type: application/json`
   - Contoh body permintaan (`id` dan `is_active` diisi server, `is_active` otomatis `true`):
```json
     {"nim":"2024001","name":"Budi Santoso","grade":85}
```
   - Status: 201, 400, 409, 415, 422
   - Contoh respons:
     - 201 (disertai header `Location: /api/v1/students/1`):
```json
       {"success":true,"message":"mahasiswa berhasil dibuat","data":{"id":1,"nim":"2024001","name":"Budi Santoso","grade":85,"is_active":true}}
```
     - 400 (body bukan JSON sah):
```json
       {"success":false,"message":"body harus berupa JSON yang valid"}
```
     - 409 (NIM ganda):
```json
       {"success":false,"message":"NIM sudah dipakai mahasiswa lain"}
```
     - 415 (Content-Type salah):
```json
       {"success":false,"message":"Content-Type harus application/json"}
```
     - 422 (validasi gagal):
```json
       {"success":false,"message":"validasi gagal","errors":{"nim":"wajib diisi","grade":"harus antara 0 dan 100"}}
```

4. **Ganti seluruh data mahasiswa**
   - Metode: PUT
   - Endpoint: `/api/v1/students/:id`
   - Parameter: `id` (path, angka positif) dan header `Content-Type: application/json`
   - Contoh body permintaan (semua field wajib dikirim):
```json
     {"nim":"2024001","name":"Budi Baru","grade":90,"is_active":false}
```
   - Status: 200, 400, 404, 409, 415, 422
   - Contoh respons:
     - 200:
```json
       {"success":true,"message":"mahasiswa berhasil diganti seluruhnya","data":{"id":1,"nim":"2024001","name":"Budi Baru","grade":90,"is_active":false}}
```
     - 404:
```json
       {"success":false,"message":"mahasiswa tidak ditemukan"}
```
     - 409 (NIM dipakai mahasiswa lain):
```json
       {"success":false,"message":"NIM sudah dipakai mahasiswa lain"}
```
     - 422 (field tidak lengkap):
```json
       {"success":false,"message":"validasi gagal","errors":{"grade":"wajib diisi pada PUT","is_active":"wajib diisi pada PUT"}}
```

5. **Ubah sebagian data mahasiswa**
   - Metode: PATCH
   - Endpoint: `/api/v1/students/:id`
   - Parameter: `id` (path, angka positif) dan header `Content-Type: application/json`
   - Contoh body permintaan (minimal satu field, hanya yang dikirim yang diubah):
```json
     {"grade":95}
```
   - Status: 200, 400, 404, 409, 415, 422
   - Contoh respons:
     - 200:
```json
       {"success":true,"message":"mahasiswa berhasil diperbarui sebagian","data":{"id":1,"nim":"2024001","name":"Budi Baru","grade":95,"is_active":false}}
```
     - 400 (body `{}`):
```json
       {"success":false,"message":"tidak ada field yang diubah"}
```
     - 422 (nilai tidak valid):
```json
       {"success":false,"message":"validasi gagal","errors":{"grade":"harus antara 0 dan 100"}}
```

6. **Hapus mahasiswa**
   - Metode: DELETE
   - Endpoint: `/api/v1/students/:id`
   - Parameter: `id` (path, angka positif)
   - Contoh body permintaan: tidak ada
   - Status: 204, 400, 404
   - Contoh respons:
     - 204: tanpa body
     - 404:
```json
       {"success":false,"message":"mahasiswa tidak ditemukan"}
```

7. **Endpoint tidak dikenal**
   - Metode: semua
   - Endpoint: alamat apa pun selain di atas
   - Parameter: tidak ada
   - Contoh body permintaan: tidak ada
   - Status: 404
   - Contoh respons:
```json
     {"success":false,"message":"endpoint tidak ditemukan"}
```

## Arti Status

1. **200**: pengambilan atau perubahan berhasil
2. **201**: penambahan berhasil, disertai header `Location`
3. **204**: penghapusan berhasil, tanpa body
4. **400**: body bukan JSON sah, atau id bukan angka
5. **404**: data atau endpoint tidak ditemukan
6. **409**: bertentangan dengan data yang ada (NIM ganda)
7. **415**: Content-Type bukan `application/json`
8. **422**: validasi isi gagal, dengan rincian per field

## Aturan Validasi

1. `nim` dan `name` tidak boleh kosong
2. `grade` harus antara 0 dan 100
3. `nim` harus unik
4. Pada PUT, semua field wajib dikirim
5. Pada PATCH, minimal satu field harus dikirim

## Sumber Bantuan

Perancangan API, penulisan seluruh kode (`main.go`, `model.go`, `helper.go`, `handler.go`), pengujian endpoint, dan pengambilan tangkapan layar saya kerjakan sendiri. Kode yang saya berikan sudah final dan tidak diubah oleh alat bantu. Selama pengerjaan, saya memakai alat bantu berikut.

1. **Alat bantu AI: Claude (Anthropic)**
   - Peran: alat bantu penulisan dokumentasi. Saya yang menentukan isi, arah, dan batasan pekerjaan, lalu Claude membantu menuangkannya ke dalam tulisan.
   - Bagian yang dibantu, berurutan sesuai pengerjaan:
     1. **Laporan enam endpoint dasar** (GET daftar, GET satu data, POST, PUT, PATCH, DELETE): menyusun format laporan berisi potongan kode dan satu paragraf penjelasan untuk setiap endpoint, berdasarkan kode handler yang saya berikan.
     2. **Laporan pengembangan endpoint daftar** (`page` dan `limit` dengan batas atas, `search`, `sort` dan `order` dengan daftar putih, filter `is_active` dan rentang `grade`, serta `meta`): menjelaskan alur kode yang saya tulis dan alasan batas `limit` maksimal 100.
     3. **Laporan penerapan status code** (200, 201, 204, 400, 404, 409, 415, 422): memetakan setiap situasi ke function dan bagian kodenya, dengan pola situasi, function, lalu letak kode.
     4. **Kontrak API di README.md**: menyusun daftar metode, endpoint, parameter, contoh body, status, dan contoh respons dalam format bernomor.
     5. **Bagian Sumber Bantuan ini**: membantu menyusun redaksi tulisan.
   - Bagian yang tidak dibantu: seluruh kode program, pemilihan batas `limit`, aturan validasi, penentuan status code, serta pengujian di Postman dan tangkapan layarnya.
   - Verifikasi: semua penjelasan dan contoh respons dari Claude saya cocokkan kembali dengan kode dan hasil pengujian saya sendiri sebelum dipakai.

2. **Dokumentasi resmi** (hapus jika tidak dibuka)
   - Dokumentasi Go (`go.dev/doc`) untuk paket `sort`, `strings`, dan `strconv`.
   - Dokumentasi Fiber v2 (`docs.gofiber.io`) untuk routing, group, middleware, dan `BodyParser`.

3. **Alat pengujian**
   - Postman, untuk mengirim request dan mengambil tangkapan layar hasil pengujian.
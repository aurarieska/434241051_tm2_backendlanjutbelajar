# hasil komunikasi http bukan code

# http request
> POST /api/v1/users HTTP/1.1        <- baris permintaan: metode, jalur, versi 
> Host: localhost:3000 *server/host tujuan req*              <- mulai di sini: header permintaan
> User-Agent: curl/8.5.0
> Accept: *
> Content-Type: application/json
> Content-Length: 70
>                                     <- baris kosong: pemisah header dan body
> {"username":"sari","email":"sari@example.com","password":"rahasia123"} *body : data yang dikirim client ke server*

# http response
< HTTP/1.1 201 Created                <- baris status: versi, kode, keterangan
< Date: Mon, 17 Aug 2026 09:29:05 GMT <- mulai di sini: header respons
< Content-Type: application/json
< Content-Length: 173
< X-Request-Id: 1ae63cf-7575-47bb-a7b8-9d7af9278228
< Location: /api/v1/users/1
<                                     <- baris kosong
< {"success":true,"message":"user berhasil dibuat","data":{...}}
```
/*
jalankan : yang dikasih contoh pengujian di word gak cocok di powershell saya

1. buat user
$body = '{"username":"johndoe","email":"john@example.com","password":"rahasia123"}'
$body | curl.exe -i -X POST "http://localhost:3000/api/v1/users" -H "Content-Type: application/json" --data-binary "@-"

2. paginasi dan pengurutan
curl.exe -i "http://localhost:3000/api/v1/users?page=1&limit=2&sort=username&order=desc"

3. pencarian dan penyaringan
curl.exe -i "http://localhost:3000/api/v1/users?search=an&is_active=true"
atau bisa juga
curl.exe -i "localhost:3000/api/v1/users?search=an&is_active=true"
atau bisa juga
"http://localhost:3000/api/v1/users?is_active=true"  

intinya klo curl.exe -i "http://localhost:3000/api/v1/users" semua, trs klo mau tambahin 
kriteria lain pake ? cari based on keyword --> search, kolom=value untuk cari based on kriteria tertentu

4. PUT -- seluruh field wajib
$body = '{"username":"john_baru","email":"jb@example.com","is_active":false}'
$body | curl.exe -i -X PUT "http://localhost:3000/api/v1/users/1" -H "Content-Type: application/json" --data-binary "@-"

5. PUT tanpa email --> harus pake email
$body = '{"kolom1":"value1","kolomdst":"valuedst"}'

$body | curl.exe -i -X PUT "http://localhost:3000/api/v1/users/1" -H "Content-Type: application/json" --data-binary "@-"

$body = '{"username":"john_biru","email":"jb@example.com"}'
$body | curl.exe -i -X PUT "http://localhost:3000/api/v1/users/1" -H "Content-Type: application/json" --data-binary "@-"

6. PATCH sebagaian
$body = '{"kolom1":value1}'
$body | curl.exe -X PATCH "http://localhost:3000/api/v1/namaGroupApi/parameter" -H "Content-Type: application/json" --data-binary "@-"

$body = '{"email":"johnpusing@gmail.com"}'
$body | curl.exe -X PATCH "http://localhost:3000/api/v1/users/1" -H "Content-Type: application/json" --data-binary "@-"

7. Tanpa content-type
untuk metode post, put, patch harus pake content-type -> semua turunan users dikenai function requireJSON yg didalamnya periksa content type untuk methodberbody yaitu POST, PUT, PATCH

8. delete
curl.exe -i -X DELETE "http://localhost:3000/api/v1/users/ID"
*/


package main
 
import (
    "fmt"
    "log"
    "strings"
    "time"
 
    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/requestid"
)
 
var metodeBerbody = map[string]bool{
    fiber.MethodPost:  true,
    fiber.MethodPut:   true,
    fiber.MethodPatch: true,
}
 
// requireJSON menolak request berisi body yang Content-Type-nya bukan JSON.
// Status yang tepat untuk kasus ini adalah 415, bukan 400.
func requireJSON(c *fiber.Ctx) error {
    if metodeBerbody[c.Method()] {
        ct := c.Get("Content-Type")
        if !strings.HasPrefix(ct, fiber.MIMEApplicationJSON) {
            return fail(c, fiber.StatusUnsupportedMediaType,
                "Content-Type harus application/json")
        }
    }
    return c.Next()
}
 
func main() {
    app := fiber.New(fiber.Config{
        AppName: "Praktikum Backend Lanjut - Pertemuan 2",
        ErrorHandler: func(c *fiber.Ctx, err error) error {
            status := fiber.StatusInternalServerError
            pesan := "terjadi kesalahan pada server"
            if e, ok := err.(*fiber.Error); ok {
                status = e.Code
                pesan = e.Message
            }
            return fail(c, status, pesan)
        },
    })
 
    // Middleware global — urutan pemasangan menentukan urutan eksekusi
    app.Use(requestid.New())
    app.Use(logger.New(logger.Config{
        Format: "[${time}] ${locals:requestid} ${method} ${path} ${status} ${latency}\n",
    }))
    app.Use(cors.New())
 
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
 
    api := app.Group("/api/v1")
 
    api.Get("/health", func(c *fiber.Ctx) error {
        return ok(c, "server berjalan", fiber.Map{"timestamp": time.Now()})
    })
 
    // requireJSON dipasang khusus pada grup ini, bukan global
    u := api.Group("/users", requireJSON) // require json dipasang untuk semua turunan users
    u.Get("/", listUsers)
    u.Get("/:id", getUser)
    u.Post("/", createUser)
    u.Put("/:id", replaceUser)
    u.Patch("/:id", patchUser)
    u.Delete("/:id", deleteUser)
 
    // Endpoint yang tidak dikenal
    app.Use(func(c *fiber.Ctx) error {
        return fail(c, fiber.StatusNotFound, "endpoint tidak ditemukan")
    })
 
    fmt.Println("Server berjalan di http://localhost:3000")
    log.Fatal(app.Listen(":3000"))
}




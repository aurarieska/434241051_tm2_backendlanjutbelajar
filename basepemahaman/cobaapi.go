/*
untuk jalankan api : curl.exe -v http://localhost:3000/api/v1/users -H "Content-Type: application/json" --data-raw '{"username":"sari","email":"sari@example.com","password":"rahasia123"}'

format : //contoh yg dikasih di word gabisa dipake di power shell windows. code jadi succes: false terus alias GAGAL
curl.exe -v linkendpoint(http://server/endpoint) -H "c" "Content-Type: application/json" --data-raw '{"kolom": "value", dst}

$body = '{"username":"johndoe","email":"john@example.com","password":"rahasia123"}'
$body | curl.exe -i -X POST "http://localhost:3000/api/v1/users" -H "Content-Type: application/json" --data-binary "@-"
*/


package main

import "github.com/gofiber/fiber/v2"

func main() {
	app := fiber.New()

	app.Post("/api/v1/users", func(c *fiber.Ctx) error {
		// c = context reqeuest, untuk ambil atau atur hal" yg berhubungan dengan reqeuest 
		// ex : c.Status, c.Body = ambil body request
		return c.Status(201).JSON(fiber.Map{
			"success": true,
			"message": "user berhasil dibuat",
			"data": fiber.Map{ //data disimpan dalam bentuk map
				"id":       1,
				"username": "sari",
				"email":    "sari@example.com",
			},
		})
	})

	app.Listen(":3000")
}
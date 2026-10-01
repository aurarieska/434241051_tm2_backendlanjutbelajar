package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Penyimpanan sementara di memori
var students = []Student{}
var nextID = 1

func ok(c *fiber.Ctx, message string, data any) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true, Message: message, Data: data,
	})
}

func okList(c *fiber.Ctx, message string, data any, meta *Meta) error {
	return c.Status(fiber.StatusOK).JSON(WebResponse{
		Success: true, Message: message, Data: data, Meta: meta,
	})
}

func created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(WebResponse{
		Success: true, Message: message, Data: data,
	})
}

func noContent(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNoContent)
}

func fail(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(WebResponse{Success: false, Message: message})
}

func failValidation(c *fiber.Ctx, errs map[string]string) error {
	return c.Status(fiber.StatusUnprocessableEntity).JSON(WebResponse{
		Success: false, Message: "validasi gagal", Errors: errs,
	})
}

func findStudentIndex(id int) int {
	for i := range students {
		if students[i].ID == id {
			return i
		}
	}
	return -1
}

// nimTerpakai true jika NIM sudah dipakai mahasiswa lain (kecualiID diabaikan).
func nimTerpakai(nim string, kecualiID int) bool {
	for _, s := range students {
		if s.ID != kecualiID && s.NIM == nim {
			return true
		}
	}
	return false
}

func paramID(c *fiber.Ctx) (int, bool) {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func gradeValid(g int) bool {
	return g >= 0 && g <= 100
}

const (
	defaultLimit = 10
	maxLimit     = 100 // batas atas agar satu request tidak menarik seluruh data
)

// Daftar putih field yang boleh dipakai untuk mengurutkan.
var allowedSort = map[string]bool{
	"id": true, "nim": true, "name": true, "grade": true,
}

// parseListQuery membaca query string dan memberi nilai bawaan yang aman.
func parseListQuery(c *fiber.Ctx) ListQuery {
	q := ListQuery{
		Page:   c.QueryInt("page", 1),
		Limit:  c.QueryInt("limit", defaultLimit),
		Search: strings.TrimSpace(c.Query("search")),
		Sort:   c.Query("sort", "id"),
		Order:  strings.ToLower(c.Query("order", "asc")),
	}

	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = defaultLimit
	}
	if q.Limit > maxLimit {
		q.Limit = maxLimit
	}
	if !allowedSort[q.Sort] {
		q.Sort = "id"
	}
	if q.Order != "desc" {
		q.Order = "asc"
	}

	if raw := c.Query("is_active"); raw != "" {
		if v, err := strconv.ParseBool(raw); err == nil {
			q.IsActive = &v
		}
	}
	if raw := c.Query("min_grade"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			q.MinGrade = &v
		}
	}
	if raw := c.Query("max_grade"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil {
			q.MaxGrade = &v
		}
	}

	return q
}

// cocokNama true jika kata kunci muncul di nama (huruf besar/kecil diabaikan).
func cocokNama(s Student, kata string) bool {
	return strings.Contains(strings.ToLower(s.Name), strings.ToLower(kata))
}
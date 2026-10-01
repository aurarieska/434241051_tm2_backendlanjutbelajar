package main

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GET /api/v1/students?page=1&limit=10
func listStudents(c *fiber.Ctx) error {
	page, limit := parsePagination(c)

	total := len(students)
	totalPages := (total + limit - 1) / limit

	mulai := (page - 1) * limit
	if mulai > total {
		mulai = total
	}
	akhir := mulai + limit
	if akhir > total {
		akhir = total
	}

	return okList(c, "daftar mahasiswa berhasil diambil", students[mulai:akhir], &Meta{
		Page: page, Limit: limit, Total: total, TotalPages: totalPages,
	})
}

// GET /api/v1/students/:id
func getStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	return ok(c, "mahasiswa ditemukan", students[i])
}

// POST /api/v1/students -> 201 + header Location
func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi"
	} else if nimTerpakai(req.NIM, 0) {
		errs["nim"] = "sudah dipakai"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi"
	}
	if !gradeValid(req.Grade) {
		errs["grade"] = "harus antara 0 dan 100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	baru := Student{
		ID:       nextID,
		NIM:      req.NIM,
		Name:     req.Name,
		Grade:    req.Grade,
		IsActive: true,
	}
	students = append(students, baru)
	nextID++

	return created(c, "mahasiswa berhasil dibuat", baru,
		"/api/v1/students/"+strconv.Itoa(baru.ID))
}

// PUT /api/v1/students/:id -> ganti seluruh isi, semua field wajib
func replaceStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	var req ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)

	errs := map[string]string{}
	if req.NIM == "" {
		errs["nim"] = "wajib diisi pada PUT"
	} else if nimTerpakai(req.NIM, id) {
		errs["nim"] = "sudah dipakai"
	}
	if req.Name == "" {
		errs["name"] = "wajib diisi pada PUT"
	}
	if req.Grade == nil {
		errs["grade"] = "wajib diisi pada PUT"
	} else if !gradeValid(*req.Grade) {
		errs["grade"] = "harus antara 0 dan 100"
	}
	if req.IsActive == nil {
		errs["is_active"] = "wajib diisi pada PUT"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	students[i].NIM = req.NIM
	students[i].Name = req.Name
	students[i].Grade = *req.Grade
	students[i].IsActive = *req.IsActive

	return ok(c, "mahasiswa berhasil diganti seluruhnya", students[i])
}

// PATCH /api/v1/students/:id -> hanya field yang dikirim
func patchStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "body harus berupa JSON yang valid")
	}

	if req.NIM == nil && req.Name == nil && req.Grade == nil && req.IsActive == nil {
		return fail(c, fiber.StatusBadRequest, "tidak ada field yang diubah")
	}

	// Validasi semua dulu, baru terapkan, supaya tidak ada perubahan setengah jalan
	errs := map[string]string{}
	if req.NIM != nil {
		nim := strings.TrimSpace(*req.NIM)
		req.NIM = &nim
		if nim == "" {
			errs["nim"] = "tidak boleh kosong"
		} else if nimTerpakai(nim, id) {
			errs["nim"] = "sudah dipakai"
		}
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		req.Name = &name
		if name == "" {
			errs["name"] = "tidak boleh kosong"
		}
	}
	if req.Grade != nil && !gradeValid(*req.Grade) {
		errs["grade"] = "harus antara 0 dan 100"
	}
	if len(errs) > 0 {
		return failValidation(c, errs)
	}

	if req.NIM != nil {
		students[i].NIM = *req.NIM
	}
	if req.Name != nil {
		students[i].Name = *req.Name
	}
	if req.Grade != nil {
		students[i].Grade = *req.Grade
	}
	if req.IsActive != nil {
		students[i].IsActive = *req.IsActive
	}

	return ok(c, "mahasiswa berhasil diperbarui sebagian", students[i])
}

// DELETE /api/v1/students/:id -> 204 tanpa body
func deleteStudent(c *fiber.Ctx) error {
	id, valid := paramID(c)
	if !valid {
		return fail(c, fiber.StatusBadRequest, "id harus berupa angka positif")
	}

	i := findStudentIndex(id)
	if i == -1 {
		return fail(c, fiber.StatusNotFound, "mahasiswa tidak ditemukan")
	}

	students = append(students[:i], students[i+1:]...)

	return noContent(c)
}
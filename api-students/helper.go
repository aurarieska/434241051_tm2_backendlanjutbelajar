package main

import (
	"strconv"

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

func created(c *fiber.Ctx, message string, data any, location string) error {
	c.Set("Location", location)
	return c.Status(fiber.StatusCreated).JSON(WebResponse{
		Success: true, Message: message, Data: data,
	})
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
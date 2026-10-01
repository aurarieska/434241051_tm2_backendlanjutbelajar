package main

type Student struct {
	ID       int    `json:"id"`
	NIM      string `json:"nim"` // penanda unik
	Name     string `json:"name"`
	Grade    int    `json:"grade"`
	IsActive bool   `json:"is_active"`
}

// POST: ID dan IsActive diisi server
type CreateStudentRequest struct {
	NIM   string `json:"nim"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

// PUT: ganti seluruh isi, pointer agar "tidak dikirim" (nil) beda dari 0/false
type ReplaceStudentRequest struct {
	NIM      string `json:"nim"`
	Name     string `json:"name"`
	Grade    *int   `json:"grade"`
	IsActive *bool  `json:"is_active"`
}

// PATCH: ubah sebagian, hanya field non-nil yang diubah
type PatchStudentRequest struct {
	NIM      *string `json:"nim,omitempty"`
	Name     *string `json:"name,omitempty"`
	Grade    *int    `json:"grade,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// Amplop respons yang sama untuk semua endpoint, sukses maupun gagal
type WebResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}
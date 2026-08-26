package foodservice

import (
	"testing"

	"go-food-store/json/response"
	"go-food-store/models"
)

// stubFoodRepo adalah FoodRepo yang dikendalikan test.
//
// Konvensi "method ini TIDAK boleh dipanggil" — satu cara saja, dipakai
// konsisten di semua test:
//
//	Assertion-nya lewat calls/called(). Itu yang menyatakan maksud test.
//	Field fungsi yang dibiarkan nil adalah backstop: kalau method-nya
//	ternyata dipanggil, test gagal lewat t.Fatalf dengan pesan jelas,
//	bukan lewat panic yang mematikan seluruh test binary.
type stubFoodRepo struct {
	t *testing.T

	createFoodFn  func(models.Food) response.WebResponse
	updateFoodFn  func(models.Food) response.WebResponse
	deleteFoodFn  func(uint) response.WebResponse
	getAllFoodFn  func() response.WebResponse
	getFoodByIdFn func(uint) (response.WebResponse, bool)

	// Perekam panggilan, untuk assertion.
	createdFood models.Food
	updatedFood models.Food
	deletedID   uint
	calls       []string
}

func (s *stubFoodRepo) CreateFood(food models.Food) response.WebResponse {
	if s.createFoodFn == nil {
		s.t.Helper()
		s.t.Fatalf("CreateFood dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "CreateFood")
	s.createdFood = food
	return s.createFoodFn(food)
}

func (s *stubFoodRepo) UpdateFood(food models.Food) response.WebResponse {
	if s.updateFoodFn == nil {
		s.t.Helper()
		s.t.Fatalf("UpdateFood dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "UpdateFood")
	s.updatedFood = food
	return s.updateFoodFn(food)
}

func (s *stubFoodRepo) DeleteFood(id uint) response.WebResponse {
	if s.deleteFoodFn == nil {
		s.t.Helper()
		s.t.Fatalf("DeleteFood dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "DeleteFood")
	s.deletedID = id
	return s.deleteFoodFn(id)
}

func (s *stubFoodRepo) GetAllFood() response.WebResponse {
	if s.getAllFoodFn == nil {
		s.t.Helper()
		s.t.Fatalf("GetAllFood dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "GetAllFood")
	return s.getAllFoodFn()
}

func (s *stubFoodRepo) GetFoodById(id uint) (response.WebResponse, bool) {
	if s.getFoodByIdFn == nil {
		s.t.Helper()
		s.t.Fatalf("GetFoodById dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "GetFoodById")
	return s.getFoodByIdFn(id)
}

// called melaporkan apakah sebuah method pernah dipanggil.
func (s *stubFoodRepo) called(name string) bool {
	for _, c := range s.calls {
		if c == name {
			return true
		}
	}
	return false
}

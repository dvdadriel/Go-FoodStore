package foodservice

import (
	"go-food-store/json/response"
	"go-food-store/models"
)

// stubFoodRepo adalah FoodRepo yang dikendalikan test. Setiap field fungsi
// boleh nil; kalau dipanggil saat nil, test akan gagal lewat panic dengan
// pesan yang jelas — itu justru cara kita membuktikan sebuah method TIDAK
// seharusnya dipanggil.
type stubFoodRepo struct {
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
	s.calls = append(s.calls, "CreateFood")
	s.createdFood = food
	if s.createFoodFn == nil {
		panic("CreateFood dipanggil padahal test tidak mengharapkannya")
	}
	return s.createFoodFn(food)
}

func (s *stubFoodRepo) UpdateFood(food models.Food) response.WebResponse {
	s.calls = append(s.calls, "UpdateFood")
	s.updatedFood = food
	if s.updateFoodFn == nil {
		panic("UpdateFood dipanggil padahal test tidak mengharapkannya")
	}
	return s.updateFoodFn(food)
}

func (s *stubFoodRepo) DeleteFood(id uint) response.WebResponse {
	s.calls = append(s.calls, "DeleteFood")
	s.deletedID = id
	if s.deleteFoodFn == nil {
		panic("DeleteFood dipanggil padahal test tidak mengharapkannya")
	}
	return s.deleteFoodFn(id)
}

func (s *stubFoodRepo) GetAllFood() response.WebResponse {
	s.calls = append(s.calls, "GetAllFood")
	if s.getAllFoodFn == nil {
		panic("GetAllFood dipanggil padahal test tidak mengharapkannya")
	}
	return s.getAllFoodFn()
}

func (s *stubFoodRepo) GetFoodById(id uint) (response.WebResponse, bool) {
	s.calls = append(s.calls, "GetFoodById")
	if s.getFoodByIdFn == nil {
		panic("GetFoodById dipanggil padahal test tidak mengharapkannya")
	}
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

package foodservice

import (
	"slices"
	"testing"

	"go-food-store/json/request"
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
//
// Selalu dibuat lewat newStub agar field t tidak pernah nil.
type stubFoodRepo struct {
	t *testing.T

	createFoodFn  func(models.Food) response.WebResponse
	updateFoodFn  func(models.Food) response.WebResponse
	deleteFoodFn  func(uint) response.WebResponse
	getAllFoodFn  func(request.Page) response.WebResponse
	getFoodByIdFn func(uint) (response.WebResponse, bool)

	// Perekam panggilan, untuk assertion.
	createdFood models.Food
	updatedFood models.Food
	deletedID   uint
	// lookedUpIDs merekam id yang diteruskan ke GetFoodById. Tanpa ini,
	// service bisa memeriksa keberadaan food X lalu memutasi food Y dan
	// tidak ada yang menyadarinya.
	lookedUpIDs []uint
	calls       []string
}

// newStub membuat stub yang siap pakai. Lewat konstruktor ini supaya field t
// tidak mungkin lupa di-set: kalau t nil, backstop nil-fn merosot dari
// t.Fatalf yang rapi menjadi nil pointer dereference plus goroutine dump.
func newStub(t *testing.T) *stubFoodRepo {
	return &stubFoodRepo{t: t}
}

func (s *stubFoodRepo) CreateFood(food models.Food) response.WebResponse {
	s.t.Helper()
	if s.createFoodFn == nil {
		s.t.Fatalf("CreateFood dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "CreateFood")
	s.createdFood = food
	return s.createFoodFn(food)
}

func (s *stubFoodRepo) UpdateFood(food models.Food) response.WebResponse {
	s.t.Helper()
	if s.updateFoodFn == nil {
		s.t.Fatalf("UpdateFood dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "UpdateFood")
	s.updatedFood = food
	return s.updateFoodFn(food)
}

func (s *stubFoodRepo) DeleteFood(id uint) response.WebResponse {
	s.t.Helper()
	if s.deleteFoodFn == nil {
		s.t.Fatalf("DeleteFood dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "DeleteFood")
	s.deletedID = id
	return s.deleteFoodFn(id)
}

func (s *stubFoodRepo) GetAllFood(page request.Page) response.WebResponse {
	s.t.Helper()
	if s.getAllFoodFn == nil {
		s.t.Fatalf("GetAllFood dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "GetAllFood")
	return s.getAllFoodFn(page)
}

func (s *stubFoodRepo) GetFoodById(id uint) (response.WebResponse, bool) {
	s.t.Helper()
	if s.getFoodByIdFn == nil {
		s.t.Fatalf("GetFoodById dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "GetFoodById")
	s.lookedUpIDs = append(s.lookedUpIDs, id)
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

// assertLookedUp memastikan GetFoodById dipanggil dengan urutan id tertentu.
func (s *stubFoodRepo) assertLookedUp(want ...uint) {
	s.t.Helper()
	if !slices.Equal(s.lookedUpIDs, want) {
		s.t.Errorf("id yang di-lookup = %v, mau %v", s.lookedUpIDs, want)
	}
}

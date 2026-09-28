package foodservice

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-playground/validator"

	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
)

// Pesan pada response stub dibuat berbeda-beda dengan sengaja. Kode HTTP
// saja tidak cukup: kalau semua stub membalas 200, test tidak bisa
// membedakan "service meneruskan response repository" dari "service
// mengarang response 200 sendiri". Message-lah yang jadi penanda asal.
func okResponse(msg string) response.WebResponse {
	return response.WebResponse{Code: http.StatusOK, Status: "OK", Message: msg}
}

func notFoundResponse() response.WebResponse {
	return response.WebResponse{
		Code:    http.StatusNotFound,
		Status:  "Not Found",
		Message: "Food not found",
	}
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name         string
		req          request.CreateFoodReq
		wantCode     int
		wantStatus   string
		wantMessage  string
		wantRepoCall bool
	}{
		{
			name:     "request valid diteruskan ke repository",
			req:      request.CreateFoodReq{FoodName: "Nasi Goreng", FoodPrice: 25000},
			wantCode: http.StatusOK,
			// Datang dari stub CreateFood, bukan dikarang oleh service.
			wantStatus:   "OK",
			wantMessage:  "dari CreateFood",
			wantRepoCall: true,
		},
		{
			name:         "nama kosong ditolak sebelum menyentuh repository",
			req:          request.CreateFoodReq{FoodName: "", FoodPrice: 25000},
			wantCode:     http.StatusBadRequest,
			wantStatus:   "Bad Request",
			wantMessage:  "Please check the request",
			wantRepoCall: false,
		},
		{
			name:         "harga nol ditolak karena bertanda required",
			req:          request.CreateFoodReq{FoodName: "Nasi Goreng", FoodPrice: 0},
			wantCode:     http.StatusBadRequest,
			wantStatus:   "Bad Request",
			wantMessage:  "Please check the request",
			wantRepoCall: false,
		},
		{
			name:         "nama 201 karakter ditolak karena melewati max=200",
			req:          request.CreateFoodReq{FoodName: strings.Repeat("a", 201), FoodPrice: 25000},
			wantCode:     http.StatusBadRequest,
			wantStatus:   "Bad Request",
			wantMessage:  "Please check the request",
			wantRepoCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := newStub(t)
			repo.createFoodFn = func(f models.Food) response.WebResponse {
				return okResponse("dari CreateFood")
			}
			svc := NewFoodService(repo, validator.New())

			got := svc.Create(tt.req)

			if got.Code != tt.wantCode {
				t.Errorf("Code = %d, mau %d", got.Code, tt.wantCode)
			}
			if got.Status != tt.wantStatus {
				t.Errorf("Status = %q, mau %q", got.Status, tt.wantStatus)
			}
			if got.Message != tt.wantMessage {
				t.Errorf("Message = %q, mau %q", got.Message, tt.wantMessage)
			}
			if repo.called("CreateFood") != tt.wantRepoCall {
				t.Errorf("CreateFood dipanggil = %v, mau %v", repo.called("CreateFood"), tt.wantRepoCall)
			}
		})
	}
}

func TestCreateMemetakanFieldKeModel(t *testing.T) {
	repo := newStub(t)
	repo.createFoodFn = func(f models.Food) response.WebResponse {
		return okResponse("dari CreateFood")
	}
	svc := NewFoodService(repo, validator.New())

	svc.Create(request.CreateFoodReq{FoodName: "Sate Ayam", FoodPrice: 30000})

	if !repo.called("CreateFood") {
		t.Fatal("CreateFood tidak dipanggil")
	}
	if repo.createdFood.FoodName != "Sate Ayam" {
		t.Errorf("FoodName = %q, mau %q", repo.createdFood.FoodName, "Sate Ayam")
	}
	if repo.createdFood.FoodPrice != 30000 {
		t.Errorf("FoodPrice = %v, mau %v", repo.createdFood.FoodPrice, 30000.0)
	}
}

func TestDeleteTidakMenghapusSaatDataTidakAda(t *testing.T) {
	repo := newStub(t)
	// deleteFoodFn sengaja nil — lihat konvensi di stub_repo_test.go.
	repo.getFoodByIdFn = func(id uint) (response.WebResponse, bool) {
		return notFoundResponse(), false
	}
	svc := NewFoodService(repo, validator.New())

	got := svc.Delete(99)

	if got.Code != http.StatusNotFound {
		t.Errorf("Code = %d, mau %d", got.Code, http.StatusNotFound)
	}
	if got.Message != "Food not found" {
		t.Errorf("Message = %q, mau response dari GetFoodById", got.Message)
	}
	if repo.called("DeleteFood") {
		t.Error("DeleteFood dipanggil padahal data tidak ditemukan")
	}
	repo.assertLookedUp(99)
}

func TestDeleteMenghapusSaatDataAda(t *testing.T) {
	repo := newStub(t)
	repo.getFoodByIdFn = func(id uint) (response.WebResponse, bool) {
		return okResponse("dari GetFoodById"), true
	}
	repo.deleteFoodFn = func(id uint) response.WebResponse {
		return okResponse("dari DeleteFood")
	}
	svc := NewFoodService(repo, validator.New())

	got := svc.Delete(7)

	if got.Code != http.StatusOK {
		t.Errorf("Code = %d, mau %d", got.Code, http.StatusOK)
	}
	if got.Status != "OK" {
		t.Errorf("Status = %q, mau %q", got.Status, "OK")
	}
	// Kedua stub membalas 200; Message-nya yang membuktikan response mana
	// yang diteruskan.
	if got.Message != "dari DeleteFood" {
		t.Errorf("response berasal dari %q, mau dari DeleteFood", got.Message)
	}
	if !repo.called("DeleteFood") {
		t.Fatal("DeleteFood tidak dipanggil")
	}
	if repo.deletedID != 7 {
		t.Errorf("id yang dihapus = %d, mau 7", repo.deletedID)
	}
	// Lookup harus memeriksa food yang sama dengan yang dihapus.
	repo.assertLookedUp(7)
}

func TestUpdate(t *testing.T) {
	t.Run("request tidak valid ditolak sebelum lookup", func(t *testing.T) {
		repo := newStub(t) // semua fn nil — lihat konvensi di stub_repo_test.go
		svc := NewFoodService(repo, validator.New())

		got := svc.Update(request.UpdateFoodReq{Id: 1, FoodName: "", FoodPrice: 100})

		if got.Code != http.StatusBadRequest {
			t.Errorf("Code = %d, mau %d", got.Code, http.StatusBadRequest)
		}
		if got.Status != "Bad Request" {
			t.Errorf("Status = %q, mau %q", got.Status, "Bad Request")
		}
		// Beda satu kata dari pesan Create ("the" vs "your") — keduanya
		// bagian dari kontrak API, jadi dikunci apa adanya.
		if got.Message != "Please check your request" {
			t.Errorf("Message = %q, mau %q", got.Message, "Please check your request")
		}
		if len(repo.calls) != 0 {
			t.Errorf("repository tersentuh: %v", repo.calls)
		}
	})

	t.Run("id nol ditolak karena bertanda required", func(t *testing.T) {
		repo := newStub(t)
		svc := NewFoodService(repo, validator.New())

		got := svc.Update(request.UpdateFoodReq{Id: 0, FoodName: "Bakso", FoodPrice: 20000})

		if got.Code != http.StatusBadRequest {
			t.Errorf("Code = %d, mau %d", got.Code, http.StatusBadRequest)
		}
		if got.Message != "Please check your request" {
			t.Errorf("Message = %q, mau %q", got.Message, "Please check your request")
		}
		if len(repo.calls) != 0 {
			t.Errorf("repository tersentuh: %v", repo.calls)
		}
	})

	t.Run("tidak meng-update saat data tidak ada", func(t *testing.T) {
		repo := newStub(t)
		repo.getFoodByIdFn = func(id uint) (response.WebResponse, bool) {
			return notFoundResponse(), false
		}
		svc := NewFoodService(repo, validator.New())

		got := svc.Update(request.UpdateFoodReq{Id: 42, FoodName: "Bakso", FoodPrice: 20000})

		if got.Code != http.StatusNotFound {
			t.Errorf("Code = %d, mau %d", got.Code, http.StatusNotFound)
		}
		if got.Message != "Food not found" {
			t.Errorf("Message = %q, mau response dari GetFoodById", got.Message)
		}
		if repo.called("UpdateFood") {
			t.Error("UpdateFood dipanggil padahal data tidak ditemukan")
		}
		repo.assertLookedUp(42)
	})

	t.Run("meneruskan id ke gorm.Model saat data ada", func(t *testing.T) {
		repo := newStub(t)
		repo.getFoodByIdFn = func(id uint) (response.WebResponse, bool) {
			return okResponse("dari GetFoodById"), true
		}
		repo.updateFoodFn = func(f models.Food) response.WebResponse {
			return okResponse("dari UpdateFood")
		}
		svc := NewFoodService(repo, validator.New())

		got := svc.Update(request.UpdateFoodReq{Id: 5, FoodName: "Bakso", FoodPrice: 20000})

		if got.Code != http.StatusOK {
			t.Errorf("Code = %d, mau %d", got.Code, http.StatusOK)
		}
		// Kedua stub membalas 200; Message-nya yang membuktikan response
		// mana yang diteruskan.
		if got.Message != "dari UpdateFood" {
			t.Errorf("response berasal dari %q, mau dari UpdateFood", got.Message)
		}
		if !repo.called("UpdateFood") {
			t.Fatal("UpdateFood tidak dipanggil")
		}
		if repo.updatedFood.ID != 5 {
			t.Errorf("ID = %d, mau 5", repo.updatedFood.ID)
		}
		if repo.updatedFood.FoodName != "Bakso" {
			t.Errorf("FoodName = %q, mau %q", repo.updatedFood.FoodName, "Bakso")
		}
		if repo.updatedFood.FoodPrice != 20000 {
			t.Errorf("FoodPrice = %v, mau %v", repo.updatedFood.FoodPrice, 20000.0)
		}
		// Lookup harus memeriksa food yang sama dengan yang di-update.
		repo.assertLookedUp(5)
	})
}

func TestFindAllDiteruskanApaAdanya(t *testing.T) {
	want := response.WebResponse{
		Code:    http.StatusOK,
		Status:  "OK",
		Message: "dari GetAllFood",
		Data:    []response.FoodResponse{{Id: 1, FoodName: "Mie Ayam", FoodPrice: 15000}},
	}
	repo := newStub(t)
	repo.getAllFoodFn = func(request.Page) response.WebResponse { return want }
	svc := NewFoodService(repo, validator.New())

	got := svc.FindAll(request.Page{Number: 1, Limit: 20})

	if got.Code != want.Code {
		t.Errorf("Code = %d, mau %d", got.Code, want.Code)
	}
	if got.Status != want.Status {
		t.Errorf("Status = %q, mau %q", got.Status, want.Status)
	}
	if got.Message != want.Message {
		t.Errorf("response berasal dari %q, mau dari GetAllFood", got.Message)
	}
	foods, ok := got.Data.([]response.FoodResponse)
	if !ok {
		t.Fatalf("Data bertipe %T, mau []response.FoodResponse", got.Data)
	}
	if len(foods) != 1 || foods[0].FoodName != "Mie Ayam" {
		t.Errorf("Data = %+v", foods)
	}
}

// FindById mengabaikan flag "found" dari repository dan mengembalikan
// response apa pun yang diberikan. Test ini mendokumentasikan perilaku yang
// ada sekarang, bukan membenarkannya — lihat "Batasan yang diketahui" di README.
func TestFindByIdMeneruskanResponseNotFound(t *testing.T) {
	repo := newStub(t)
	repo.getFoodByIdFn = func(id uint) (response.WebResponse, bool) {
		return notFoundResponse(), false
	}
	svc := NewFoodService(repo, validator.New())

	got := svc.FindById(123)

	if got.Code != http.StatusNotFound {
		t.Errorf("Code = %d, mau %d", got.Code, http.StatusNotFound)
	}
	if got.Message != "Food not found" {
		t.Errorf("Message = %q, mau response dari GetFoodById", got.Message)
	}
	repo.assertLookedUp(123)
}

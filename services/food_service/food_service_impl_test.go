package foodservice

import (
	"net/http"
	"testing"

	"github.com/go-playground/validator"

	"go-food-store/json/request"
	"go-food-store/json/response"
	"go-food-store/models"
)

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
		wantRepoCall bool
	}{
		{
			name:         "request valid diteruskan ke repository",
			req:          request.CreateFoodReq{FoodName: "Nasi Goreng", FoodPrice: 25000},
			wantCode:     http.StatusOK,
			wantRepoCall: true,
		},
		{
			name:         "nama kosong ditolak sebelum menyentuh repository",
			req:          request.CreateFoodReq{FoodName: "", FoodPrice: 25000},
			wantCode:     http.StatusBadRequest,
			wantRepoCall: false,
		},
		{
			name:         "harga nol ditolak karena bertanda required",
			req:          request.CreateFoodReq{FoodName: "Nasi Goreng", FoodPrice: 0},
			wantCode:     http.StatusBadRequest,
			wantRepoCall: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &stubFoodRepo{
				createFoodFn: func(f models.Food) response.WebResponse {
					return okResponse("Successfully create new food")
				},
			}
			svc := NewFoodService(repo, validator.New())

			got := svc.Create(tt.req)

			if got.Code != tt.wantCode {
				t.Errorf("Code = %d, mau %d", got.Code, tt.wantCode)
			}
			if repo.called("CreateFood") != tt.wantRepoCall {
				t.Errorf("CreateFood dipanggil = %v, mau %v", repo.called("CreateFood"), tt.wantRepoCall)
			}
		})
	}
}

func TestCreateMemetakanFieldKeModel(t *testing.T) {
	repo := &stubFoodRepo{
		createFoodFn: func(f models.Food) response.WebResponse {
			return okResponse("ok")
		},
	}
	svc := NewFoodService(repo, validator.New())

	svc.Create(request.CreateFoodReq{FoodName: "Sate Ayam", FoodPrice: 30000})

	if repo.createdFood.FoodName != "Sate Ayam" {
		t.Errorf("FoodName = %q, mau %q", repo.createdFood.FoodName, "Sate Ayam")
	}
	if repo.createdFood.FoodPrice != 30000 {
		t.Errorf("FoodPrice = %v, mau %v", repo.createdFood.FoodPrice, 30000.0)
	}
}

func TestDeleteTidakMenghapusSaatDataTidakAda(t *testing.T) {
	repo := &stubFoodRepo{
		getFoodByIdFn: func(id uint) (response.WebResponse, bool) {
			return notFoundResponse(), false
		},
		// deleteFoodFn sengaja nil: kalau dipanggil, stub akan panic.
	}
	svc := NewFoodService(repo, validator.New())

	got := svc.Delete(99)

	if got.Code != http.StatusNotFound {
		t.Errorf("Code = %d, mau %d", got.Code, http.StatusNotFound)
	}
	if repo.called("DeleteFood") {
		t.Error("DeleteFood dipanggil padahal data tidak ditemukan")
	}
}

func TestDeleteMenghapusSaatDataAda(t *testing.T) {
	repo := &stubFoodRepo{
		getFoodByIdFn: func(id uint) (response.WebResponse, bool) {
			return okResponse("found"), true
		},
		deleteFoodFn: func(id uint) response.WebResponse {
			return okResponse("Successfully delete food")
		},
	}
	svc := NewFoodService(repo, validator.New())

	got := svc.Delete(7)

	if got.Code != http.StatusOK {
		t.Errorf("Code = %d, mau %d", got.Code, http.StatusOK)
	}
	if repo.deletedID != 7 {
		t.Errorf("id yang dihapus = %d, mau 7", repo.deletedID)
	}
}

func TestUpdate(t *testing.T) {
	t.Run("request tidak valid ditolak sebelum lookup", func(t *testing.T) {
		repo := &stubFoodRepo{} // semua fn nil: panic kalau tersentuh
		svc := NewFoodService(repo, validator.New())

		got := svc.Update(request.UpdateFoodReq{Id: 1, FoodName: "", FoodPrice: 100})

		if got.Code != http.StatusBadRequest {
			t.Errorf("Code = %d, mau %d", got.Code, http.StatusBadRequest)
		}
		if len(repo.calls) != 0 {
			t.Errorf("repository tersentuh: %v", repo.calls)
		}
	})

	t.Run("tidak meng-update saat data tidak ada", func(t *testing.T) {
		repo := &stubFoodRepo{
			getFoodByIdFn: func(id uint) (response.WebResponse, bool) {
				return notFoundResponse(), false
			},
		}
		svc := NewFoodService(repo, validator.New())

		got := svc.Update(request.UpdateFoodReq{Id: 42, FoodName: "Bakso", FoodPrice: 20000})

		if got.Code != http.StatusNotFound {
			t.Errorf("Code = %d, mau %d", got.Code, http.StatusNotFound)
		}
		if repo.called("UpdateFood") {
			t.Error("UpdateFood dipanggil padahal data tidak ditemukan")
		}
	})

	t.Run("meneruskan id ke gorm.Model saat data ada", func(t *testing.T) {
		repo := &stubFoodRepo{
			getFoodByIdFn: func(id uint) (response.WebResponse, bool) {
				return okResponse("found"), true
			},
			updateFoodFn: func(f models.Food) response.WebResponse {
				return okResponse("Successfully update food data")
			},
		}
		svc := NewFoodService(repo, validator.New())

		svc.Update(request.UpdateFoodReq{Id: 5, FoodName: "Bakso", FoodPrice: 20000})

		if repo.updatedFood.ID != 5 {
			t.Errorf("ID = %d, mau 5", repo.updatedFood.ID)
		}
		if repo.updatedFood.FoodName != "Bakso" {
			t.Errorf("FoodName = %q, mau %q", repo.updatedFood.FoodName, "Bakso")
		}
	})
}

func TestFindAllDiteruskanApaAdanya(t *testing.T) {
	want := response.WebResponse{
		Code: http.StatusOK,
		Data: []response.FoodResponse{{Id: 1, FoodName: "Mie Ayam", FoodPrice: 15000}},
	}
	repo := &stubFoodRepo{
		getAllFoodFn: func() response.WebResponse { return want },
	}
	svc := NewFoodService(repo, validator.New())

	got := svc.FindAll()

	if got.Code != want.Code {
		t.Errorf("Code = %d, mau %d", got.Code, want.Code)
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
	repo := &stubFoodRepo{
		getFoodByIdFn: func(id uint) (response.WebResponse, bool) {
			return notFoundResponse(), false
		},
	}
	svc := NewFoodService(repo, validator.New())

	got := svc.FindById(123)

	if got.Code != http.StatusNotFound {
		t.Errorf("Code = %d, mau %d", got.Code, http.StatusNotFound)
	}
}

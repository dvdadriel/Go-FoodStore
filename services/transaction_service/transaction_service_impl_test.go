package transactionservice

import (
	"net/http"
	"testing"

	"github.com/go-playground/validator"

	"go-food-store/json/request"
	"go-food-store/json/response"
)

func newService(repo *stubTransactionRepo) TransactionService {
	return NewTransactionService(repo, validator.New())
}

func validRequest() request.CreateTransactionReq {
	return request.CreateTransactionReq{
		CustomerId: 1,
		Food:       []request.FoodDetailReq{{FoodId: 1, Quantity: 2}},
	}
}

// Permintaan yang tidak masuk akal harus berhenti di service, tanpa pernah
// menyentuh database.
func TestCreateMenolakPermintaanTidakValid(t *testing.T) {
	cases := []struct {
		nama string
		req  request.CreateTransactionReq
	}{
		{
			// Dua baris untuk makanan yang sama menabrak primary key
			// (transaction_id, food_id) dan dulu keluar sebagai 500.
			nama: "FoodId duplikat",
			req: request.CreateTransactionReq{
				CustomerId: 1,
				Food: []request.FoodDetailReq{
					{FoodId: 7, Quantity: 1},
					{FoodId: 7, Quantity: 3},
				},
			},
		},
		{
			nama: "quantity nol",
			req: request.CreateTransactionReq{
				CustomerId: 1,
				Food:       []request.FoodDetailReq{{FoodId: 1, Quantity: 0}},
			},
		},
		{
			// Dulu lolos dan mengurangi total belanja.
			nama: "quantity negatif",
			req: request.CreateTransactionReq{
				CustomerId: 1,
				Food:       []request.FoodDetailReq{{FoodId: 1, Quantity: -5}},
			},
		},
		{
			nama: "tanpa item",
			req: request.CreateTransactionReq{
				CustomerId: 1,
				Food:       []request.FoodDetailReq{},
			},
		},
		{
			nama: "tanpa customer",
			req: request.CreateTransactionReq{
				Food: []request.FoodDetailReq{{FoodId: 1, Quantity: 1}},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.nama, func(t *testing.T) {
			repo := newStub(t)
			res := newService(repo).Create(c.req)

			if res.Code != http.StatusBadRequest {
				t.Errorf("Code = %d, mau %d", res.Code, http.StatusBadRequest)
			}
			if repo.called("CreateTransaction") {
				t.Error("transaksi tetap dibuat padahal permintaannya tidak valid")
			}
		})
	}
}

// Customer yang tidak ada dibalas apa adanya dari repository, bukan ditelan
// dan diubah jadi sukses.
func TestCreateMeneruskanCustomerTidakDitemukan(t *testing.T) {
	repo := newStub(t)
	repo.getCustomerFn = func(uint) (response.WebResponse, bool) {
		return response.WebResponse{Code: http.StatusNotFound, Message: "Customer with id 1 not found"}, false
	}

	res := newService(repo).Create(validRequest())

	if res.Code != http.StatusNotFound {
		t.Errorf("Code = %d, mau %d", res.Code, http.StatusNotFound)
	}
	if repo.called("CreateTransaction") {
		t.Error("transaksi dibuat padahal customer-nya tidak ada")
	}
}

func TestCreateMeneruskanPermintaanYangSah(t *testing.T) {
	repo := newStub(t)
	repo.getCustomerFn = func(uint) (response.WebResponse, bool) { return response.WebResponse{}, true }
	repo.createFn = func(req request.CreateTransactionReq) response.WebResponse {
		if req.CustomerId != 1 || len(req.Food) != 1 {
			t.Errorf("repository menerima %+v, bukan permintaan aslinya", req)
		}
		return response.WebResponse{Code: http.StatusOK}
	}

	if res := newService(repo).Create(validRequest()); res.Code != http.StatusOK {
		t.Errorf("Code = %d, mau %d", res.Code, http.StatusOK)
	}
}

// Update memakai aturan yang sama dengan Create — gampang terlewat saat
// aturannya hanya dipasang di satu jalur.
func TestUpdateMenolakFoodIdDuplikat(t *testing.T) {
	repo := newStub(t)
	res := newService(repo).Update(request.UpdateTransactionReq{
		Id:         1,
		CustomerId: 1,
		Food: []request.FoodDetailReq{
			{FoodId: 3, Quantity: 1},
			{FoodId: 3, Quantity: 2},
		},
	})

	if res.Code != http.StatusBadRequest {
		t.Errorf("Code = %d, mau %d", res.Code, http.StatusBadRequest)
	}
	if repo.called("UpdateTransaction") {
		t.Error("transaksi tetap diubah padahal permintaannya tidak valid")
	}
}

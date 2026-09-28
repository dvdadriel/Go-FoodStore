package transactionservice

import (
	"testing"

	"go-food-store/json/request"
	"go-food-store/json/response"
)

// stubTransactionRepo adalah TransactionRepositories yang dikendalikan test.
// Pola dan konvensinya sama dengan stub di services/food_service: field fungsi
// yang dibiarkan nil adalah backstop — kalau method-nya ternyata dipanggil,
// test gagal dengan pesan jelas alih-alih panic.
type stubTransactionRepo struct {
	t *testing.T

	createFn      func(request.CreateTransactionReq) response.WebResponse
	updateFn      func(request.UpdateTransactionReq) response.WebResponse
	getCustomerFn func(uint) (response.WebResponse, bool)

	calls []string
}

func newStub(t *testing.T) *stubTransactionRepo { return &stubTransactionRepo{t: t} }

func (s *stubTransactionRepo) CreateTransaction(req request.CreateTransactionReq) response.WebResponse {
	s.t.Helper()
	if s.createFn == nil {
		s.t.Fatalf("CreateTransaction dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "CreateTransaction")
	return s.createFn(req)
}

func (s *stubTransactionRepo) UpdateTransaction(req request.UpdateTransactionReq) response.WebResponse {
	s.t.Helper()
	if s.updateFn == nil {
		s.t.Fatalf("UpdateTransaction dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "UpdateTransaction")
	return s.updateFn(req)
}

func (s *stubTransactionRepo) GetCustomerById(customerId uint) (response.WebResponse, bool) {
	s.t.Helper()
	if s.getCustomerFn == nil {
		s.t.Fatalf("GetCustomerById dipanggil padahal test tidak mengharapkannya")
	}
	s.calls = append(s.calls, "GetCustomerById")
	return s.getCustomerFn(customerId)
}

func (s *stubTransactionRepo) DeleteTransaction(uint) response.WebResponse {
	s.t.Fatalf("DeleteTransaction dipanggil padahal test tidak mengharapkannya")
	return response.WebResponse{}
}

func (s *stubTransactionRepo) GetAllPaidTransaction(request.Page) response.WebResponse {
	s.t.Fatalf("GetAllPaidTransaction dipanggil padahal test tidak mengharapkannya")
	return response.WebResponse{}
}

func (s *stubTransactionRepo) GetUnpaidTransaction(request.Page) response.WebResponse {
	s.t.Fatalf("GetUnpaidTransaction dipanggil padahal test tidak mengharapkannya")
	return response.WebResponse{}
}

func (s *stubTransactionRepo) GetTransactionById(uint) (response.WebResponse, bool) {
	s.t.Fatalf("GetTransactionById dipanggil padahal test tidak mengharapkannya")
	return response.WebResponse{}, false
}

func (s *stubTransactionRepo) AcceptPayment(uint) response.WebResponse {
	s.t.Fatalf("AcceptPayment dipanggil padahal test tidak mengharapkannya")
	return response.WebResponse{}
}

func (s *stubTransactionRepo) GetFoodById(uint) (response.WebResponse, bool) {
	s.t.Fatalf("GetFoodById dipanggil padahal test tidak mengharapkannya")
	return response.WebResponse{}, false
}

func (s *stubTransactionRepo) called(name string) bool {
	for _, c := range s.calls {
		if c == name {
			return true
		}
	}
	return false
}

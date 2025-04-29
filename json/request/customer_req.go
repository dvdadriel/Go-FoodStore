package request

type CreateCustReq struct {
	CustomerName string `validate:"required,min=1,max=200" json:"CustomerName"`
}

type UpdateCustReq struct {
	Id           uint   `validate:"required"`
	CustomerName string `validate:"required,min=1,max=200" json:"CustomerName"`
}

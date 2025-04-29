package response

type WebResponse struct {
	Code    int
	Status  string
	Message string
	Data    interface{}
}

package request

type RegisterReq struct {
	Username string `validate:"required,min=3,max=100" json:"Username"`
	// Panjang minimum, bukan aturan "harus ada simbol dan angka": aturan
	// komposisi mendorong orang membuat kata sandi pendek yang mudah ditebak
	// mesin tapi susah diingat manusia. Panjang yang menentukan.
	Password string `validate:"required,min=8,max=200" json:"Password"`
}

type LoginReq struct {
	Username string `validate:"required" json:"Username"`
	Password string `validate:"required" json:"Password"`
}

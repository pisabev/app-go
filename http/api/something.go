package api

type Request struct {
	Field string `json:"field" valid:"required"`
}

type Response struct {
	Data string `json:"data"`
}

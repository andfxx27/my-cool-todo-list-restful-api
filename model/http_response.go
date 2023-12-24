package model

type HTTPResponse struct {
	Status  int
	Message string
	Result  interface{}
}

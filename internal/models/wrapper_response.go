package models

type ResponseWrapper[T any] struct {
	Data T `json:"data"`
}

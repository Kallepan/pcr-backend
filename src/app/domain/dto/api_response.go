package dto

type APIResponse[T any] struct {
	ResponseKey     int    `json:"response_key"`
	ResponseMessage string `json:"response_message"`
	Data            T      `json:"data"`
}

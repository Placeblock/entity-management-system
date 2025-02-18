package rest

type CreateMessageDto struct {
	Message string `json:"message"`
}

type APIResponse[T any] struct {
	Data T `json:"data"`
}

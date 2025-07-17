package api

type APIResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type APIError struct {
	Message string `json:"message"`
}

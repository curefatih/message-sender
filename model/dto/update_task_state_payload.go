package dto

type TaskStateUpdateRequest struct {
	Active bool
}

type TaskStateUpdateResponse struct {
	Message string `json:"message"`
	Active  bool   `json:"active"`
}

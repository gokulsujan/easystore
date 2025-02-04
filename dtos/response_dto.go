package dtos

type SuccessResponse struct {
	Status  string      `json:"status" example:"success"`
	Message string      `json:"message" example:"Operation successful"`
	Result  interface{} `json:"result"`
}

type ErrorResponse struct {
	Status  string `json:"status" example:"failed"`
	Message string `json:"message" example:"Error message"`
}

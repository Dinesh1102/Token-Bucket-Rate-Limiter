package api

type CheckRequest struct {
	ClientKey string `json:"clientKey" binding:"required"`
}

type CreateClientRequest struct {
	ClientKey string  `json:"clientKey" binding:"required"`
	Algorithm string  `json:"algorithm" binding:"required"`
	Rate      float64 `json:"rate" binding:"required"`
	Burst     int     `json:"burst" binding:"required"`
}

type UpdateClientRequest struct {
	Algorithm string  `json:"algorithm" binding:"required"`
	Rate      float64 `json:"rate" binding:"required"`
	Burst     int     `json:"burst" binding:"required"`
}
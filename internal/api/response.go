package api

type CheckResponse struct {
	Allowed bool   `json:"allowed"`
	Message string `json:"message,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type ClientResponse struct {
	ClientKey string  `json:"clientKey"`
	Algorithm string  `json:"algorithm"`
	Rate      float64 `json:"rate"`
	Burst     int     `json:"burst"`
}
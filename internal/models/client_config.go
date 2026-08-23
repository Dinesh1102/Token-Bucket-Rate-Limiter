package models

type ClientConfig struct {
	ClientKey string
	Algorithm AlgorithmName
	Rate      float64
	Burst     int
}

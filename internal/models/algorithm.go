package models

type AlgorithmName string

const (
	AlgorithmTokenBucket   AlgorithmName = "TOKEN_BUCKET"
	AlgorithmSlidingWindow AlgorithmName = "SLIDING_WINDOW"
)

func (a AlgorithmName) String() string {
	return string(a)
}

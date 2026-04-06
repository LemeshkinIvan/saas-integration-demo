package common

type RegularResponseDTO[T any] struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
	Data        T      `json:"data,omitempty"`
}

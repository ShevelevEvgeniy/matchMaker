package dto

type User struct {
	Name    string  `json:"name" validate:"required"`
	Skill   float32 `json:"skill" validate:"required"`
	Latency float32 `json:"latency" validate:"required"`
}

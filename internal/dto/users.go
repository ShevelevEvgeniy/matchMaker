package dto

type User struct {
	Name    string  `json:"name" validate:"required"`
	Skill   float64 `json:"skill" validate:"required"`
	Latency float64 `json:"latency" validate:"required"`
}

type UserDistance struct {
	Index    int
	Distance float64
}

package dto

type User struct {
	Name    string  `json:"name" validate:"required,regexp=^[a-zA-Z_]+$"`
	Skill   float64 `json:"skill" validate:"required,gt=0,lt=11"`
	Latency float64 `json:"latency" validate:"required,gt=0,lt=11"`
}

type UserDistance struct {
	Index    int
	Distance float64
}

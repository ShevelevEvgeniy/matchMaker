package dto

type Users struct {
	User []User `json:"users" validate:"required"`
}

type User struct {
	Name    string  `json:"name" validate:"required"`
	Skill   float64 `json:"skill" validate:"required,gt=0,lt=11"`
	Latency float64 `json:"latency" validate:"required,gt=0,lt=11"`
}

type UserDistance struct {
	Index    int64
	Distance float64
}

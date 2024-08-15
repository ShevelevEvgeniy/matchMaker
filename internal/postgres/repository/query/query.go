package query

import _ "embed"

var (
	//go:embed sql/add_user.sql
	AddUser string
)

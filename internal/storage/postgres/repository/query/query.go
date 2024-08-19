package query

import _ "embed"

var (
	//go:embed sql/save_users.sql
	SaveUsers string

	//go:embed sql/get_users_in_search.sql
	GetUsersInSearch string

	//go:embed sql/unmark_search.sql
	UnmarkSearch string
)

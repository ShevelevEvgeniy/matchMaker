package custom_errors

import "github.com/pkg/errors"

var (
	NotEnoughUsers = errors.New("the number of users is less than the group size")
)

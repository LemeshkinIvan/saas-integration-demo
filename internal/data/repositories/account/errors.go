package account

import "errors"

var ErrPostgresArgument = errors.New("AccountRepository: postgres is nil")

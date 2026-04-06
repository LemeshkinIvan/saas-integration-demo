package auth

import "errors"

var ErrPostgresArgument = errors.New("AuthRepository: postgres is nil")
var ErrCacheArgument = errors.New("AuthRepository: cache is nil")

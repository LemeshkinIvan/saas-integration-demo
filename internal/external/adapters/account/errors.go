package account

import "errors"

var ErrAdapterInit = errors.New("AccountAdapter: config is nil")
var ErrUnauthorized = errors.New("AccountAdapter: the user is not logged in.")
var ErrBase = errors.New("AccountAdapter")

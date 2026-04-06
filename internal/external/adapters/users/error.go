package user

import "errors"

var ErrConfigArgument = errors.New("UserAdapter: amo config is nil")
var ErrBase = errors.New("UserAdapter")
var ErrStatusForbidden = errors.New("UserAdapter: don't have the rights to call this method")
var ErrUnathorized = errors.New("UserAdapter: the user is not logged in")

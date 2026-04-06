package chat

import "errors"

var ErrPostgresArgument = errors.New("ChatRepository: postgres is nil")

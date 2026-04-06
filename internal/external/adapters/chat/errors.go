package chat

import "errors"

var ErrAdapterInit = errors.New("ChatAdapter: config is nil")
var ErrBase = errors.New("ChatAdapter:")
var ErrUncorrectData = errors.New("ChatAdapter: incorrect data was passed. Details are available in the response body")
var ErrSignatureInvalid = errors.New("ChatAdapter: the request signature is incorrect")

package source

import "errors"

var ErrAdapterInit = errors.New("init failed")
var ErrUnautharized = errors.New("integration is not authorized")
var ErrUncorrectData = errors.New("incorrect data was transmitted")
var ErrNotFound = errors.New("The integration has no controlled sources or the source has not been found")
var ErrBase = errors.New("SourceAdapter")

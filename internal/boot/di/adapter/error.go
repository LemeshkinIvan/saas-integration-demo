package adapter

import "errors"

var ErrAdaptersInit = errors.New("AdapterContainer: init failed")

var ErrTelegramConfig = errors.New("AdapterContainer: telegram cfg is nil")
var ErrAmoConfig = errors.New("AdapterContainer: amo cfg is nil")

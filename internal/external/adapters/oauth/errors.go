package oauth

import "errors"

var ErrAdapterInit = errors.New("OauthAdapter: init failed")
var ErrUnautharized = errors.New("OauthAdapter: integration is not authorized")
var ErrUncorrectData = errors.New("OauthAdapter: incorrect data was transmitted")
var ErrNotFound = errors.New("OauthAdapter: The integration has no controlled sources or the source has not been found")
var ErrBase = errors.New("OauthAdapter")
var ErrIncorrectData = errors.New("OauthAdapter: Incorrect data was passed. Details are available in the response body")

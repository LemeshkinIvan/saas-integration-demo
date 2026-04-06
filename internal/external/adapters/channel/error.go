package channel

import "errors"

var ErrAdapterInit = errors.New("ChannelAdapters: config is nil")
var ErrSignature = errors.New("ChannelAdapters: The request signature is incorrect")
var ErrChannelExist = errors.New("ChannelAdapters: The channel does not exist")
var ErrUncorrectData = errors.New("ChannelAdapters: Incorrect data was transmitted. Details are available in the response body")
var ErrBase = errors.New("ChannelAdapters")

package pipeline

import "errors"

var ErrConfigArgument = errors.New("PipelineAdapter: amo config is nil")
var ErrUnauthorized = errors.New("PipelineAdaper: user not authorized")
var ErrBase = errors.New("PipelineAdapter")

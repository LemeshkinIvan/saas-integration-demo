package config

import "errors"

var ErrParseYAML = errors.New("ConfigContainer: failed to parse YAML")
var ErrReadFile = errors.New("ConfigContainer: failed to read config file")

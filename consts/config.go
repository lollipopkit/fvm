package consts

import _ "embed"

var (
	//go:embed settings.json
	VSC_CONFIG string
)

const (
	ConfigFileName = "config.json"
)

package config

import _ "embed"

//go:embed default-sync-strategy.yaml
var DefaultSyncStrategy []byte

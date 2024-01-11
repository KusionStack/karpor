package ui

import (
	"embed"
)

// Embedded contains embedded UI resources
//
//go:embed build/*
//nolint:typecheck
var Embedded embed.FS

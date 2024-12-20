package ui

import (
	"embed"
)

// Embedded contains embedded UI resources
//
//go:embed build/*
var Embedded embed.FS

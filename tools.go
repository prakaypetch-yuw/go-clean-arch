//go:build tools
// +build tools

package tools

// following https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module

import (
	_ "github.com/google/wire/cmd/wire"
)

package utils

import (
	"fmt"
	"runtime"
)

var (
	version = "0.0.1"
)

// Version return current version
func Version() string {
	return fmt.Sprintf("%v (%v %v)", version, runtime.GOOS, runtime.GOARCH)
}

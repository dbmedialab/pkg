// Package logruskimpy - provides short, skimpy log lines for console output
// requires logrus >= v1.4.0
package logruskimpy

import (
	"fmt"
	"path/filepath"
	"runtime"

	log "github.com/sirupsen/logrus"
)

// New - returns a logrus.Formatter that focuses on the log message
// and data, by minimizing the function and file information
func New() (lf log.Formatter) {
	return &log.TextFormatter{
		ForceColors:      true, // we don't want logrus to fall back on the non-color formatter
		FullTimestamp:    false,
		CallerPrettyfier: skimpyCaller,
	}
}

func skimpyCaller(rf *runtime.Frame) (function string, file string) {
	return "", fmt.Sprintf("%s:%d", filepath.Base(rf.File), rf.Line)
}

package logruskimpy

import (
	"bytes"
	"regexp"
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestOutput(t *testing.T) {
	var logrusBuf bytes.Buffer
	log.SetReportCaller(true)
	log.SetOutput(&logrusBuf)
	log.SetLevel(log.DebugLevel)

	log.SetFormatter(New())
	log.WithFields(log.Fields{
		"Muse":  "Hyper Music",
		"Mason": "Exceeder",
	}).Debug("test")

	logString := logrusBuf.String()
	t.Logf("Output:\n`%s`", logString)

	// very careful regex - ansi color codes are hard to match
	pattern := `^\x1b[37mDEBU\x1b[0m\[\d+\]logruskimpy_test.go:\d+\s+test\s+[^M]*Mason[^=]*=Exceeder\s+[^M]*Muse[^=]*=[^H]+Hyper Music`

	match, err := regexp.MatchString(pattern, logString)
	if err != nil {
		t.Errorf("Regex error: %s", err.Error())
		return
	}

	if !match {
		t.Errorf("Log output did not match!")
	}
}

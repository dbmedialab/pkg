package logrustic

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestElasticFields(t *testing.T) {
	f := NewFormatter("test.")

	b, err := f.Format(logrus.WithFields(logrus.Fields{
		"error":  errors.New("wild walrus"),
		"body":   "ladida",
		"caller": "Arnold Schwarzenegger",
	}))
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	t.Logf("Format output:\n%s", string(b))

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	val, ok := entry["test.error"]

	if !ok {
		t.Error("Entry not set")
	}

	expected := "wild walrus"
	if val != expected {
		t.Errorf("Expected %s, got %s", expected, val)
	}

	_, ok = entry["error"]
	if ok {
		t.Error("Original value still set")
	}
}

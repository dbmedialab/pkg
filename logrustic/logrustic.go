// Package logrustic : can be used to set up a logrus-instance
// that definitely won't break the elasticSearch parser.
// (whereas vanilla logrus with json encoding usually will)
//
// Usage examples:
//
//  // Functional style:
//  // Get a logrus.Logger instance
//  app.Logger = logrustic.NewLogger("logrus.", logrus.InfoLevel)
//
//  // Global Variables style:
//  // Get a Formatter and apply it to logrus's global Logger instance
//  logrus.SetFormatter(logrustic.NewFormatter("logrus."))
//
// Works with most versions of github.com/sirupsen/logrus
// (developed vs logrus v1.4.2)
//
// For a complete list of built in elasticSearch fields
// that your logrus fields may collide with,
// run the following command in your Kibana console:
// `GET filebeat-7.0.0-beta1/_mapping`
package logrustic

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger - returns a logrus-logger
// with elasticSearch-compatible formatting.
//
// `prefix` is prepended to all non-default data keys, to avoid collision with
// any of the predefined keys in ElasticSearch.
// example values: "logrus.", "myappname."
//
// `level` sets the log level. Logrus will only output log entries
// whose severity are equal to or greater than this level.
func NewLogger(prefix string, level logrus.Level) (ll *logrus.Logger) {
	ll = logrus.New()
	ll.Level = level

	ll.Formatter = NewFormatter(prefix)

	ll.SetReportCaller(true) // include the calling method in the log entry
	ll.SetOutput(os.Stderr)  // use stderr as output

	return ll
}

// NewFormatter - returns an elasticSearch-compatible logrus-formatter
func NewFormatter(prefix string) (lf logrus.Formatter) {
	lf = &ElasticFormatter{
		JSON: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyFile: prefix + "file",
				logrus.FieldKeyMsg:  "@message",   // overrides the "message" field of the kibana entry
				logrus.FieldKeyTime: "@timestamp", // the elasticSearch "date" datatype is compatible with logrus timestamp output
			},
		},
		Prefix: prefix,
	}

	return lf
}

// ElasticFormatter - A logrus json formatter that adds a prefix to all
// non-default fields
type ElasticFormatter struct {
	JSON   *logrus.JSONFormatter
	Prefix string
}

// Format - returns an ES-sanitized and json-encoded version of the log-entry
func (ef *ElasticFormatter) Format(entry *logrus.Entry) ([]byte, error) {

	// Copy data fields to a new fields-instance,
	// and add a prefix to each key
	nd := logrus.Fields{}
	for fieldName, content := range entry.Data {
		nd[ef.Prefix+fieldName] = content
	}
	entry.Data = nd

	// use the logrus.JSONFormatter to format the modified entry as JSON
	return ef.JSON.Format(entry)
}

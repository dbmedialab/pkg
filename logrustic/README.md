# Logrustic

This package can be used to set up a logrus-instance that definitely won't break the elasticSearch parser (whereas vanilla logrus with json encoding usually will).

In practice, this enables Kibana to receive structured logs from your go-app, allowing you to browse (and search) by default logrus fields like `@message`, `func`, `level` and `{yourprefix.}file`, as well as any custom fields you've added via the `WithFields` method.

### Usage examples

```go
import (
	"github.com/dbmedialab/pkg/logrustic"
	"github.com/sirupsen/logrus"
)

// Functional style:
// Get a logrus.Logger instance
app.Logger = logrustic.NewLogger("logrus.", logrus.InfoLevel)

// Global Variables style:
// Get a Formatter and apply it to logrus's global Logger instance
logrus.SetFormatter(logrustic.NewFormatter("logrus."))
```

Works with most versions of logrus (developed vs logrus v1.4.2)

### Elastic debugging
If you suspect your app may be running into elastic parser-problems, you can confirm these suspicions by using the following filter in Kibana: `kubernetes.namespace:logging and kubernetes.pod.name:filebeat* and message:"WARN"`

For a complete list of built in elasticSearch fields that log entry fields may collide with (and a demonstration of why a prefix was chosen as the solution), run the following command in your Kibana console: `GET filebeat-7.0.0-beta1/_mapping`

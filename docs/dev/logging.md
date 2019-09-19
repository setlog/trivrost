# Logging

## Log library

In trivrost and related projects, logging is done using the Logrus logging library, which can be found under https://github.com/sirupsen/logrus. Typically, it is imported using a named import `log`:

```
import log "github.com/sirupsen/logrus"
```

## Log setup

The package `pkg/logging` can configure a log formatter for Logrus and log all output to a file. It can be used by calling `logging.Initialize()`. This can be seen in `cmd/launcher/main.go:init()`. The library uses two incrementing variables in the log file name instead of one. This allows us to differentiate between occurences of the user starting trivrost (4 digit number) and occurrences of trivrost restarting itself (single letter of the alphabet). E.g., a set of log files named as follows would tell you that the user started trivrost two times, whereby during the first time trivrost restarted itself twice and during the second time it restarted once:

```
0000a.MyProduct.2019-06-27_15-05-03.log
0000b.MyProduct.2019-06-27_15-05-04.log
0000c.MyProduct.2019-06-27_15-05-05.log
0001a.MyProduct.2019-06-27_15-15-13.log
0001b.MyProduct.2019-06-27_15-15-14.log
```

## How to log

When importing the Logrus package, one can log using the various public functions of the package. However, we want to avoid logging messages of types other than `Debug` and `Warn` inside packages contained under `pkg/`, because we want these to be reusable and free from unwanted side-effects (such as printing to standard out). If you still want to log information, use Go's default log library. Our `logging` package sets the default log library's output to be relayed to and processed in `pkg/logging/relay.go` during initialization. If you want to log an error, you should return an `error` instead and let the caller handle it.

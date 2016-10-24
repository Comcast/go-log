# Log

Logging is an important part of the application and having a consistent logging mechanism and structure is mandatory. With several teams writing different components that talk to each other, being able to read each others logs could be the difference between finding bugs quickly or wasting hours.

## License
This software is licensed under the Apache license. For full text see [LICENSE](./LICENSE)

## Examples
The following code provides some examples of logging.

### Logging Without Log Levels
```
func main() {
	var buf bytes.Buffer
	log.Init("LogExample", 0, log.DevWriter{Device: log.DevAll, Writer: &buf})

	Square(5)

	log.Shutdown()
	fmt.Printf(buf.String())
}

func Square(n int) int {
	context := "SquareTwice"
	function := "Square"

	log.Start(context, function)

	sq := n * n
	log.Warnf(context, function, "Square 1: %d\n", sq)

	sq = 0
	for i := 0; i < n; i++ {
		sq += n
	}
	log.Tracef(context, function, "Square 2: %d\n", sq)

	log.Complete(context, function)

	return sq
}
```
The output produced
```
2016/08/22 13:05:42.596: LogExample[29218]: main.go#24: SquareTwice: Square: Started:
2016/08/22 13:05:42.596: LogExample[29218]: main.go#27: SquareTwice: Square: Warning: Square 1: 25
2016/08/22 13:05:42.596: LogExample[29218]: main.go#33: SquareTwice: Square: Trace: Square 2: 25
2016/08/22 13:05:42.596: LogExample[29218]: main.go#35: SquareTwice: Square: Completed:
```

### Logging With Log Levels
```
var warnLogger = log.NewLogger("squareLogger", func() int { return log.LevelWarning })

func main() {
	var buf bytes.Buffer
	log.Init("LogExample", 0, log.DevWriter{Device: log.DevAll, Writer: &buf})

	Square(5)

	log.Shutdown()
	fmt.Printf(buf.String())
}

func Square(n int) int {
	context := "SquareTwice"
	function := "Square"

	log.Start(context, function)

	sq := n * n
	warnLogger.Warnf(context, function, "Square 1: %d\n", sq)

	sq = 0
	for i := 0; i < n; i++ {
		sq += n
	}
	// NOTE: The following TRACE line will not be emitted because the log level
	// of the `warnLogger` is at the `Warning` level.
	warnLogger.Tracef(context, function, "Square 2: %d\n", sq)

	log.Complete(context, function)

	return sq
}
```
The output produced
```
2016/08/22 13:20:41.594: LogExample[29406]: main.go#26: SquareTwice: Square: Started:
2016/08/22 13:20:41.594: LogExample[29406]: main.go#29: SquareTwice: Square: Warning: Square 1: 25
2016/08/22 13:20:41.594: LogExample[29406]: main.go#37: SquareTwice: Square: Completed:
```

## Project Info
Package docs: https://godoc.org/github.com/Comcast/go-log/log

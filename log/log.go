/**
* Copyright 2016 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package log

import (
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// Date and time layout for each trace line.
const (
	layout       = "2006/01/02 15:04:05.000"
	emptyMessage = "**** LOG ERROR: MESSAGE IS EMPTY - PLEASE REPORT ****\n"
	loggingIsOff = "**** LOG WARNING: LOGGING WAS OFF - PLEASE REPORT ****\n"
)

// Formatter provide support for special formatting.
type Formatter interface {
	Format() string
}

// line is passed to the safe write goroutine
// as the string to write to the device.
type line struct {
	w io.Writer
	b []byte
}

// logger maintains internal state for our logger.
type logger struct {
	dest   map[int8]io.Writer
	destMu sync.RWMutex

	mu    sync.Mutex
	wg    sync.WaitGroup
	write chan line
	timer *time.Timer

	shutdown      bool
	loggingOff    bool
	pendingWrites int32
	prefix        string
	test          int32
}

// logger maintains a pointer to the single logger.
var l = logger{
	timer:  time.NewTimer(time.Hour),
	prefix: "PREFIX",
}

// Init initializes the logging system for use. It can be called
// multiple times to reset the destination.
func Init(prefix string, bufferSize int, dws ...DevWriter) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.write != nil {
		// We need to Unlock the mutex before
		// calling Shutdown and get back in.
		l.mu.Unlock()
		{
			// Shutdown the log with the current configuration.
			Shutdown()
		}
		l.mu.Lock()
	}

	// Set user defined values.
	l.prefix = prefix
	l.write = make(chan line, bufferSize)

	l.destMu.Lock()
	{
		// Create and init the map of devices.
		l.dest = map[int8]io.Writer{
			DevError:   os.Stderr,
			DevPanic:   os.Stderr,
			DevWarning: os.Stderr,

			DevStart:  os.Stdout,
			DevTrace:  os.Stdout,
			DevQuery:  os.Stdout,
			DevData:   os.Stdout,
			DevSplunk: os.Stdout,
		}
	}
	l.destMu.Unlock()

	// If a device is provided, update the writer.
	if dws != nil {
		for _, dw := range dws {
			// Were we asked to update all the devices.
			if dw.Device == DevAll {
				Dev.All(dw.Writer)
				continue
			}

			l.destMu.Lock()
			{
				// Just update the single device.
				l.dest[dw.Device] = dw.Writer
			}
			l.destMu.Unlock()
		}
	}

	// Set the flags.
	l.loggingOff = false
	l.shutdown = false

	// Create the safe writer goroutine to prevent the log
	// from causing the host application to block on log calls.
	l.wg.Add(1)
	go safeWrite()
}

// InitTest configures the logger for testing purposes.
func InitTest(prefix string, bufferSize int, dws ...DevWriter) {
	Init(prefix, bufferSize, dws...)
	atomic.StoreInt32(&l.test, 1)
}

// Shutdown will wait until all the pending writes are complete.
func Shutdown() {
	l.mu.Lock()
	{
		l.shutdown = true
		close(l.write)
		l.wg.Wait()
		l.write = nil

		atomic.StoreInt32(&l.test, 0)
	}
	l.mu.Unlock()
}

// dtFile returns the current time and file for logging.
func dtFile(calldepth int, function string) (dateTime string, file string, funcName string, pid int) {
	// Capture the name of the function logging if
	// a function was not provided.
	if function == "" {
		pc := make([]uintptr, calldepth+1)
		runtime.Callers(calldepth, pc)
		f := runtime.FuncForPC(pc[calldepth-1])
		_, funcName = path.Split(f.Name())
	} else {
		funcName = function
	}

	if atomic.LoadInt32(&l.test) == 1 {
		return time.Date(2009, time.November, 10, 15, 0, 0, 0, time.UTC).UTC().Format(layout), "file.go#512", funcName, 69910
	}

	dateTime = time.Now().UTC().Format(layout)

	_, filePath, line, ok := runtime.Caller(calldepth)
	if !ok {
		return dateTime, "unknown.go#0:", "missing", os.Getpid()
	}
	_, file = path.Split(filePath)

	return dateTime, fmt.Sprintf("%s#%d", file, line), funcName, os.Getpid()
}

// output performs the actual write to the destination device.
func output(w io.Writer, format string, a ...interface{}) {
	if w == nil {
		return
	}
	if format == "" {
		format = emptyMessage
	} else if a != nil {
		format = fmt.Sprintf(format, a...)
	}

	if format[len(format)-1] != '\n' {
		format = format + "\n"
	}

	// Create a slice from the string.
	b := []byte(format)

	l.mu.Lock()
	{
		// We are shutting down. Get out of town.
		if l.shutdown {
			l.mu.Unlock()
			return
		}

		// We have turned logging off. Wait here until the existing
		// buffer has been flushed and then we can start again.
		if l.loggingOff {
			if atomic.LoadInt32(&l.pendingWrites) > 0 {
				l.mu.Unlock()
				return
			}

			l.loggingOff = false
			fmt.Fprintf(w, loggingIsOff)
		}

		const waitTime = 25 * time.Millisecond

		// Found out the timer can be reset down but not up. Creating
		// the timer with one hour as the initial value.
		l.timer.Reset(waitTime)

		// If we can't perform the write within the wait time, then
		// let's not wait and turn off logging.
		select {
		case l.write <- line{w, b}:
			atomic.AddInt32(&l.pendingWrites, 1)
			l.timer.Stop()
		case <-l.timer.C:
			l.loggingOff = true
		}
	}
	l.mu.Unlock()
}

// safeWrite is run as a goroutine. It pulls a message from the
// channel and perform the write.
func safeWrite() {
	for ln := range l.write {
		ln.w.Write(ln.b)
		atomic.AddInt32(&l.pendingWrites, -1)
	}

	l.wg.Done()
}

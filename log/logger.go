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

// Set of levels that are compared for filtering tracing to
// the specific log levels.
const (
	LevelOff     = 0
	LevelError   = 1
	LevelWarning = 2
	LevelOutput  = 3
	LevelTrace   = 4
)

// Logger represents an individual logger with logging
// level permissions.
type Logger struct {
	Up1   UplevelLogger
	name  string
	level func() int
}

// NewLogger creates a logger for use of writting logs
// within the scope of a configured logging level.
func NewLogger(name string, level func() int) *Logger {
	l := &Logger{
		name:  name,
		level: level,
	}

	// Init the Up1 logger support.
	l.Up1.l = l
	l.Up1.up = 2

	return l
}

// Start is used for the entry into a function.
// Min logLevel required for logging: LevelTrace(4)
func (l *Logger) Start(context interface{}, function string) {
	if l.level() >= LevelTrace {
		Up1.Start(context, function)
	}
}

// Startf is used for the entry into a function with a formatted message.
// Min logLevel required for logging: LevelTrace(4)
func (l *Logger) Startf(context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelTrace {
		Up1.Startf(context, function, format, a...)
	}
}

// Complete is used for the exit of a function.
// Min logLevel required for logging: LevelTrace(4)
func (l *Logger) Complete(context interface{}, function string) {
	if l.level() >= LevelTrace {
		Up1.Complete(context, function)
	}
}

// Completef is used for the exit of a function with a formatted message.
// Min logLevel required for logging: LevelTrace(4)
func (l *Logger) Completef(context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelTrace {
		Up1.Completef(context, function, format, a...)
	}
}

// CompleteErr is used to write an error with complete into the trace.
// Min logLevel required for logging: LevelError(1)
func (l *Logger) CompleteErr(err error, context interface{}, function string) {
	if l.level() >= LevelError {
		Up1.CompleteErr(err, context, function)
	}
}

// CompleteErrf is used to write an error with complete into the trace with a formatted message.
// Min logLevel required for logging: LevelError(1)
func (l *Logger) CompleteErrf(err error, context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelError {
		Up1.CompleteErrf(err, context, function, format, a...)
	}
}

// Err is used to write an error into the trace.
// Min logLevel required for logging: LevelError(1)
func (l *Logger) Err(err error, context interface{}, function string) {
	if l.level() >= LevelError {
		Up1.Err(err, context, function)
	}
}

// Errf is used to write an error into the trace with a formatted message.
// Min logLevel required for logging: LevelError(1)
func (l *Logger) Errf(err error, context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelError {
		Up1.Errf(err, context, function, format, a...)
	}
}

// ErrFatal is used to write an error into the trace then terminate the program.
// Min logLevel required for logging: LevelError(1)
func (l *Logger) ErrFatal(err error, context interface{}, function string) {
	if l.level() >= LevelError {
		Up1.ErrFatal(err, context, function)
	}
}

// ErrFatalf is used to write an error into the trace with a formatted message then terminate the program.
// Min logLevel required for logging: LevelError(1)
func (l *Logger) ErrFatalf(err error, context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelError {
		Up1.ErrFatalf(err, context, function, format, a...)
	}
}

// ErrPanic is used to write an error into the trace then panic the program.
// Min logLevel required for logging: LevelError(1)
func (l *Logger) ErrPanic(err error, context interface{}, function string) {
	if l.level() >= LevelError {
		Up1.ErrPanic(err, context, function)
	}
}

// ErrPanicf is used to write an error into the trace with a formatted message then panic the program.
// Min logLevel required for logging: LevelError(1)
func (l *Logger) ErrPanicf(err error, context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelError {
		Up1.ErrPanicf(err, context, function, format, a...)
	}
}

// Tracef is used to write information into the trace with a formatted message.
// Min logLevel required for logging: LevelTrace(4)
func (l *Logger) Tracef(context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelTrace {
		Up1.Tracef(context, function, format, a...)
	}
}

// Warnf is used to write a warning into the trace with a formatted message.
// Min logLevel required for logging: LevelWarning(2)
func (l *Logger) Warnf(context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelWarning {
		Up1.Warnf(context, function, format, a...)
	}
}

// Queryf is used to write a query into the trace with a formatted message.
// Min logLevel required for logging: LevelTrace(4)
func (l *Logger) Queryf(context interface{}, function string, format string, a ...interface{}) {
	if l.level() >= LevelTrace {
		Up1.Queryf(context, function, format, a...)
	}
}

// DataKV is used to write a key/value pair into the trace.
// Min logLevel required for logging: LevelOutput(3)
func (l *Logger) DataKV(context interface{}, function string, key string, value interface{}) {
	if l.level() >= LevelOutput {
		Up1.DataKV(context, function, key, value)
	}
}

// DataBlock is used to write a block of data into the trace.
// Min logLevel required for logging: LevelOutput(3)
func (l *Logger) DataBlock(context interface{}, function string, block interface{}) {
	if l.level() >= LevelOutput {
		Up1.DataBlock(context, function, block)
	}
}

// DataString is used to write a string with CRLF each on their own line.
// Min logLevel required for logging: LevelOutput(3)
func (l *Logger) DataString(context interface{}, function string, message string) {
	if l.level() >= LevelOutput {
		Up1.DataString(context, function, message)
	}
}

// DataTrace is used to write a block of data from an io.Stringer respecting each line.
// Min logLevel required for logging: LevelOutput(3)
func (l *Logger) DataTrace(context interface{}, function string, formatters ...Formatter) {
	if l.level() >= LevelOutput {
		Up1.DataTrace(context, function, formatters...)
	}
}

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

// UplevelLogger controls the stack frame level for file name, line number
// and function name.  It can be used to embed logging calls in helper
// functions that report the file name, line number and function name of
// the routine that calls the helper.
type UplevelLogger struct {
	l  *Logger
	up Uplevel
}

// Start is used for the entry into a function.
// Min logLevel required for logging: LevelTrace(4)
func (lvl UplevelLogger) Start(context interface{}, function string) {
	if lvl.l.level() >= LevelTrace {
		lvl.up.Start(context, function)
	}
}

// Startf is used for the entry into a function with a formatted message.
// Min logLevel required for logging: LevelTrace(4)
func (lvl UplevelLogger) Startf(context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelTrace {
		lvl.up.Startf(context, function, format, a...)
	}
}

// Complete is used for the exit of a function.
// Min logLevel required for logging: LevelTrace(4)
func (lvl UplevelLogger) Complete(context interface{}, function string) {
	if lvl.l.level() >= LevelTrace {
		lvl.up.Complete(context, function)
	}
}

// Completef is used for the exit of a function with a formatted message.
// Min logLevel required for logging: LevelTrace(4)
func (lvl UplevelLogger) Completef(context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelTrace {
		lvl.up.Completef(context, function, format, a...)
	}
}

// CompleteErr is used to write an error with complete into the trace.
// Min logLevel required for logging: LevelError(1)
func (lvl UplevelLogger) CompleteErr(err error, context interface{}, function string) {
	if lvl.l.level() >= LevelError {
		lvl.up.CompleteErr(err, context, function)
	}
}

// CompleteErrf is used to write an error with complete into the trace with a formatted message.
// Min logLevel required for logging: LevelError(1)
func (lvl UplevelLogger) CompleteErrf(err error, context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelError {
		lvl.up.CompleteErrf(err, context, function, format, a...)
	}
}

// Err is used to write an error into the trace.
// Min logLevel required for logging: LevelError(1)
func (lvl UplevelLogger) Err(err error, context interface{}, function string) {
	if lvl.l.level() >= LevelError {
		lvl.up.Err(err, context, function)
	}
}

// Errf is used to write an error into the trace with a formatted message.
// Min logLevel required for logging: LevelError(1)
func (lvl UplevelLogger) Errf(err error, context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelError {
		lvl.up.Errf(err, context, function, format, a...)
	}
}

// ErrFatal is used to write an error into the trace then terminate the program.
// Min logLevel required for logging: LevelError(1)
func (lvl UplevelLogger) ErrFatal(err error, context interface{}, function string) {
	if lvl.l.level() >= LevelError {
		lvl.up.ErrFatal(err, context, function)
	}
}

// ErrFatalf is used to write an error into the trace with a formatted message then terminate the program.
// Min logLevel required for logging: LevelError(1)
func (lvl UplevelLogger) ErrFatalf(err error, context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelError {
		lvl.up.ErrFatalf(err, context, function, format, a...)
	}
}

// ErrPanic is used to write an error into the trace then panic the program.
// Min logLevel required for logging: LevelError(1)
func (lvl UplevelLogger) ErrPanic(err error, context interface{}, function string) {
	if lvl.l.level() >= LevelError {
		lvl.up.ErrPanic(err, context, function)
	}
}

// ErrPanicf is used to write an error into the trace with a formatted message then panic the program.
// Min logLevel required for logging: LevelError(1)
func (lvl UplevelLogger) ErrPanicf(err error, context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelError {
		lvl.up.ErrPanicf(err, context, function, format, a...)
	}
}

// Tracef is used to write information into the trace with a formatted message.
// Min logLevel required for logging: LevelTrace(4)
func (lvl UplevelLogger) Tracef(context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelTrace {
		lvl.up.Tracef(context, function, format, a...)
	}
}

// Warnf is used to write a warning into the trace with a formatted message.
// Min logLevel required for logging: LevelWarning(2)
func (lvl UplevelLogger) Warnf(context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelWarning {
		lvl.up.Warnf(context, function, format, a...)
	}
}

// Queryf is used to write a query into the trace with a formatted message.
// Min logLevel required for logging: LevelTrace(4)
func (lvl UplevelLogger) Queryf(context interface{}, function string, format string, a ...interface{}) {
	if lvl.l.level() >= LevelTrace {
		lvl.up.Queryf(context, function, format, a...)
	}
}

// DataKV is used to write a key/value pair into the trace.
// Min logLevel required for logging: LevelOutput(3)
func (lvl UplevelLogger) DataKV(context interface{}, function string, key string, value interface{}) {
	if lvl.l.level() >= LevelOutput {
		lvl.up.DataKV(context, function, key, value)
	}
}

// DataBlock is used to write a block of data into the trace.
// Min logLevel required for logging: LevelOutput(3)
func (lvl UplevelLogger) DataBlock(context interface{}, function string, block interface{}) {
	if lvl.l.level() >= LevelOutput {
		lvl.up.DataBlock(context, function, block)
	}
}

// DataString is used to write a string with CRLF each on their own line.
// Min logLevel required for logging: LevelOutput(3)
func (lvl UplevelLogger) DataString(context interface{}, function string, message string) {
	if lvl.l.level() >= LevelOutput {
		lvl.up.DataString(context, function, message)
	}
}

// DataTrace is used to write a block of data from an io.Stringer respecting each line.
// Min logLevel required for logging: LevelOutput(3)
func (lvl UplevelLogger) DataTrace(context interface{}, function string, formatters ...Formatter) {
	if lvl.l.level() >= LevelOutput {
		lvl.up.DataTrace(context, function, formatters...)
	}
}

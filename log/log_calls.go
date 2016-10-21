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

// Start is used for the entry into a function.
func Start(context interface{}, function string) {
	Up1.Start(context, function)
}

// Startf is used for the entry into a function with a formatted message.
func Startf(context interface{}, function string, format string, a ...interface{}) {
	Up1.Startf(context, function, format, a...)
}

// Complete is used for the exit of a function.
func Complete(context interface{}, function string) {
	Up1.Complete(context, function)
}

// Completef is used for the exit of a function with a formatted message.
func Completef(context interface{}, function string, format string, a ...interface{}) {
	Up1.Completef(context, function, format, a...)
}

// CompleteErr is used to write an error with complete into the trace.
func CompleteErr(err error, context interface{}, function string) {
	Up1.CompleteErr(err, context, function)
}

// CompleteErrf is used to write an error with complete into the trace with a formatted message.
func CompleteErrf(err error, context interface{}, function string, format string, a ...interface{}) {
	Up1.CompleteErrf(err, context, function, format, a...)
}

// Err is used to write an error into the trace.
func Err(err error, context interface{}, function string) {
	Up1.Err(err, context, function)
}

// Errf is used to write an error into the trace with a formatted message.
func Errf(err error, context interface{}, function string, format string, a ...interface{}) {
	Up1.Errf(err, context, function, format, a...)
}

// ErrFatal is used to write an error into the trace then terminate the program.
func ErrFatal(err error, context interface{}, function string) {
	Up1.ErrFatal(err, context, function)
}

// ErrFatalf is used to write an error into the trace with a formatted message then terminate the program.
func ErrFatalf(err error, context interface{}, function string, format string, a ...interface{}) {
	Up1.ErrFatalf(err, context, function, format, a...)
}

// ErrPanic is used to write an error into the trace then panic the program.
func ErrPanic(err error, context interface{}, function string) {
	Up1.ErrPanic(err, context, function)
}

// ErrPanicf is used to write an error into the trace with a formatted message then panic the program.
func ErrPanicf(err error, context interface{}, function string, format string, a ...interface{}) {
	Up1.ErrPanicf(err, context, function, format, a...)
}

// Tracef is used to write information into the trace with a formatted message.
func Tracef(context interface{}, function string, format string, a ...interface{}) {
	Up1.Tracef(context, function, format, a...)
}

// Warnf is used to write a warning into the trace with a formatted message.
func Warnf(context interface{}, function string, format string, a ...interface{}) {
	Up1.Warnf(context, function, format, a...)
}

// Queryf is used to write a query into the trace with a formatted message.
func Queryf(context interface{}, function string, format string, a ...interface{}) {
	Up1.Queryf(context, function, format, a...)
}

// DataKV is used to write a key/value pair into the trace.
func DataKV(context interface{}, function string, key string, value interface{}) {
	Up1.DataKV(context, function, key, value)
}

// DataBlock is used to write a block of data into the trace.
func DataBlock(context interface{}, function string, block interface{}) {
	Up1.DataBlock(context, function, block)
}

// DataString is used to write a string with CRLF each on their own line.
func DataString(context interface{}, function string, message string) {
	Up1.DataString(context, function, message)
}

// DataTrace is used to write a block of data from an io.Stringer respecting each line.
func DataTrace(context interface{}, function string, formatters ...Formatter) {
	Up1.DataTrace(context, function, formatters...)
}

// Splunk is used to write a log message in a splunk-able format.
func Splunk(m ...SplunkPair) {
	Up1.Splunk(m...)
}

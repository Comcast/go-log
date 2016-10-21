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
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

// Uplevel controls the stack frame level for file name, line number
// and function name.  It can be used to embed logging calls in helper
// functions that report the file name, line number and function name of
// the routine that calls the helper.
type Uplevel int

// Up1 is short for Uplevel(1).
var Up1 Uplevel = 1

// Start is used for the entry into a function.
func (lvl Uplevel) Start(context interface{}, function string) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevStart), "%s: %s[%d]: %s: %v: %s: Started:\n", dt, l.prefix, pid, file, context, funcName)
}

// Startf is used for the entry into a function with a formatted message.
func (lvl Uplevel) Startf(context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevStart), "%s: %s[%d]: %s: %v: %s: Started: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...))
}

// Complete is used for the exit of a function.
func (lvl Uplevel) Complete(context interface{}, function string) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevStart), "%s: %s[%d]: %s: %v: %s: Completed:\n", dt, l.prefix, pid, file, context, funcName)
}

// Completef is used for the exit of a function with a formatted message.
func (lvl Uplevel) Completef(context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevStart), "%s: %s[%d]: %s: %v: %s: Completed: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...))
}

// CompleteErr is used to write an error with complete into the trace.
func (lvl Uplevel) CompleteErr(err error, context interface{}, function string) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevError), "%s: %s[%d]: %s: %v: %s: Completed ERROR: %s", dt, l.prefix, pid, file, context, funcName, err)
}

// CompleteErrf is used to write an error with complete into the trace with a formatted message.
func (lvl Uplevel) CompleteErrf(err error, context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevError), "%s: %s[%d]: %s: %v: %s: Completed ERROR: %s: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...), err)
}

// Err is used to write an error into the trace.
func (lvl Uplevel) Err(err error, context interface{}, function string) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevError), "%s: %s[%d]: %s: %v: %s: ERROR: %s", dt, l.prefix, pid, file, context, funcName, err)
}

// Errf is used to write an error into the trace with a formatted message.
func (lvl Uplevel) Errf(err error, context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevError), "%s: %s[%d]: %s: %v: %s: ERROR: %s: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...), err)
}

// ErrFatal is used to write an error into the trace then terminate the program.
func (lvl Uplevel) ErrFatal(err error, context interface{}, function string) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevError), "%s: %s[%d]: %s: %v: %s: ERROR: %s", dt, l.prefix, pid, file, context, funcName, err)
	output(Dev.get(DevError), "%s: %s[%d]: %s: %v: %s: TERMINATING\n", dt, l.prefix, pid, file, context, funcName)
	Shutdown()
	os.Exit(1)
}

// ErrFatalf is used to write an error into the trace with a formatted message then terminate the program.
func (lvl Uplevel) ErrFatalf(err error, context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevError), "%s: %s[%d]: %s: %v: %s: ERROR: %s: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...), err)
	output(Dev.get(DevError), "%s: %s[%d]: %s: %v: %s: TERMINATING\n", dt, l.prefix, pid, file, context, funcName)
	Shutdown()
	os.Exit(1)
}

// ErrPanic is used to write an error into the trace then panic the program.
func (lvl Uplevel) ErrPanic(err error, context interface{}, function string) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevPanic), "%s: %s[%d]: %s: %v: %s: ERROR: %s", dt, l.prefix, pid, file, context, funcName, err)
	output(Dev.get(DevPanic), "%s: %s[%d]: %s: %v: %s: TERMINATING\n", dt, l.prefix, pid, file, context, funcName)
	Shutdown()
	panic("Terminating Program")
}

// ErrPanicf is used to write an error into the trace with a formatted message then panic the program.
func (lvl Uplevel) ErrPanicf(err error, context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevPanic), "%s: %s[%d]: %s: %v: %s: ERROR: %s: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...), err)
	output(Dev.get(DevPanic), "%s: %s[%d]: %s: %v: %s: TERMINATING\n", dt, l.prefix, pid, file, context, funcName)
	Shutdown()
	panic("Terminating Program")
}

// Tracef is used to write information into the trace with a formatted message.
func (lvl Uplevel) Tracef(context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevTrace), "%s: %s[%d]: %s: %v: %s: Trace: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...))
}

// Warnf is used to write a warning into the trace with a formatted message.
func (lvl Uplevel) Warnf(context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevWarning), "%s: %s[%d]: %s: %v: %s: Warning: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...))
}

// Queryf is used to write a query into the trace with a formatted message.
func (lvl Uplevel) Queryf(context interface{}, function string, format string, a ...interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevQuery), "%s: %s[%d]: %s: %v: %s: Query: %s", dt, l.prefix, pid, file, context, funcName, fmt.Sprintf(format, a...))
}

// DataKV is used to write a key/value pair into the trace.
func (lvl Uplevel) DataKV(context interface{}, function string, key string, value interface{}) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)
	output(Dev.get(DevData), "%s: %s[%d]: %s: %v: %s: DATA: %s: %v", dt, l.prefix, pid, file, context, funcName, key, value)
}

// DataBlock is used to write a block of data into the trace.
func (lvl Uplevel) DataBlock(context interface{}, function string, block interface{}) {
	if v, ok := block.(string); ok {
		(lvl + 1).DataString(context, function, v)
		return
	}

	d, err := json.MarshalIndent(block, "", "    ")
	if err != nil {
		d = []byte(err.Error())
	}

	(lvl + 1).DataString(context, function, string(d))
}

// DataString is used to write a string with CRLF each on their own line.
func (lvl Uplevel) DataString(context interface{}, function string, message string) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)

	if message == "" {
		output(Dev.get(DevData), "%s: %s[%d]: %s: %v: %s: DATA: %%!ds(MISSING)\n", dt, l.prefix, pid, file, context, funcName)
		return
	}

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s: %s[%d]: %s: %v: %s: DATA:\n", dt, l.prefix, pid, file, context, funcName)

	lines := bytes.Split([]byte(message), []byte{'\n'})
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		fmt.Fprintf(&buf, "\t%s\n", line)
	}

	output(Dev.get(DevData), buf.String())
}

// DataTrace is used to write a block of data from an io.Stringer respecting each line.
func (lvl Uplevel) DataTrace(context interface{}, function string, formatters ...Formatter) {
	dt, file, funcName, pid := dtFile(2+int(lvl), function)

	var lines [][]byte
	for _, f := range formatters {
		if f != nil {
			lines = append(lines, bytes.Split([]byte(f.Format()), []byte{'\n'})...)
		}
	}

	var buf bytes.Buffer

	fmt.Fprintf(&buf, "%s: %s[%d]: %s: %v: %s: DATA:\n", dt, l.prefix, pid, file, context, funcName)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		fmt.Fprintf(&buf, "\t%s\n", line)
	}

	message := buf.String()
	if message == "" {
		output(Dev.get(DevData), "\t%%!ds(MISSING)\n")
		return
	}

	output(Dev.get(DevData), message)
}

// splunkEncode encodes a value to be splunkable.
// If a value is a string that contains space character(s), that value will be
// encompassed within double quotes.
func splunkEncode(ifc interface{}) string {
	if v, ok := ifc.(string); ok && strings.Contains(v, " ") {
		return fmt.Sprintf("%q", v)
	}
	return fmt.Sprintf("%v", ifc)
}

// SplunkValue represents a slice of values to be logged in splunk.
type SplunkValue []interface{}

// String is a stringer function for the SplunkValue (which is a slice of SplunkPairs).
// Its main function is to encompass a list (empty, single member, or multiple members) within
// square brackets with ", " as a separator.
func (sl SplunkValue) String() string {
	var buf bytes.Buffer

	buf.WriteString("[")
	for i, v := range sl {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(splunkEncode(v))
	}
	buf.WriteString("]")

	return buf.String()
}

// SplunkPair represents the key/value pairs to be logged in splunk.
type SplunkPair struct {
	Key   string
	Value interface{}
}

// Splunk is used to write a log message in a splunk-able format.
func (lvl Uplevel) Splunk(m ...SplunkPair) {
	var buf bytes.Buffer

	for _, i := range m {
		buf.WriteString(" ")
		buf.WriteString(splunkEncode(i.Key))
		buf.WriteString("=")
		buf.WriteString(splunkEncode(i.Value))
	}

	var dateTime string
	if atomic.LoadInt32(&l.test) == 1 {
		dateTime = time.Date(2009, time.November, 10, 15, 0, 0, 0, time.UTC).UTC().Format(layout)
	} else {
		dateTime = time.Now().UTC().Format(layout)
	}

	output(Dev.get(DevSplunk), "%s:%s\n", dateTime, buf.String())
}

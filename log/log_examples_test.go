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

package log_test

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"testing"
	"time"

	"github.com/Comcast/go-log/log"
)

// ExampleStart provides a basic example for using the log package.
func ExampleStart() {
	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: buf})

	{
		log.Start("1234", "Basic")

		v, err := strconv.ParseInt("10", 10, 64)
		if err != nil {
			log.CompleteErr(err, "1234", "Basic")
			return
		}

		log.Completef("1234", "Basic", "Conv[%d]", v)
	}

	log.Shutdown()
	fmt.Println(buf.String())
	// Output:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Started:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Completed: Conv[10]
}

// ExampleErr provides an example of logging an error.
func ExampleErr() {
	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: buf})

	{
		log.Start("1234", "Error")

		v, err := strconv.ParseInt("1080980980980980980898908", 10, 64)
		if err != nil {
			log.CompleteErr(err, "1234", "Error")

			// Flush the output for testing.
			log.Shutdown()
			fmt.Println(buf.String())
			return
		}

		log.Completef("1234", "Error", "Conv[%d]", v)
	}

	log.Shutdown()
	fmt.Println(buf.String())
	// Output:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Error: Started:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Error: Completed ERROR: strconv.ParseInt: parsing "1080980980980980980898908": value out of range
}

// ExampleDataKV provides an example of logging K/V pair data.
func ExampleDataKV() {
	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: buf})

	{
		log.Start("1234", "Data_KV")

		log.DataKV("1234", "Data_KV", "Value 1", 1)
		log.DataKV("1234", "Data_KV", "Hex Value 2", 0x00000002)

		log.Complete("1234", "Data_KV")
	}

	log.Shutdown()
	fmt.Println(buf.String())
	// Output:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_KV: Started:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_KV: DATA: Value 1: 1
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_KV: DATA: Hex Value 2: 2
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_KV: Completed:
}

// ExampleDataBlock provides an example of logging a block of data.
func ExampleDataBlock() {
	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: buf})

	{
		log.Start("1234", "Data_Block")

		data := `Test Data with
2 lines`

		log.DataBlock("1234", "Data_Block", data)

		log.Complete("1234", "Data_Block")
	}

	log.Shutdown()
	fmt.Println(buf.String())
	// Output:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_Block: Started:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_Block: DATA:
	// 	Test Data with
	// 	2 lines
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_Block: Completed:
}

type Message []byte

// Format implements the Formatter interface to produce logging output
// for this slice of bytes.
func (m Message) Format() string {
	var buf log.SafeBuffer

	// Message Bytes:	(0x0000) EE 6E 11 00 00 00 3E EA DE 18 00 00 2D 00 00 00
	// 					(0x0010) 00 00 00 00 00 0A 16 C3 9B 00 00 00 00 00 00 00

	fmt.Fprintf(&buf, "Message Bytes:\t")

	rows := (len(m) / 16) + 1
	l := len(m)
	for row := 0; row < rows; row++ {
		var r []byte
		st := row * 16

		if row < (rows - 1) {
			r = m[st : st+16]
		} else {
			r = m[st : st+(l-st)]
		}

		if row > 0 {
			fmt.Fprintf(&buf, "\t\t\t")
		}

		fmt.Fprintf(&buf, "(0x%.4X)", st)

		for i := 0; i < len(r); i++ {
			fmt.Fprintf(&buf, " %.2X", r[i])
		}
		fmt.Fprintf(&buf, "\n")
	}

	return buf.String()
}

// ExampleDataTrace provides an example of logging from a fmt.Stringer.
func ExampleDataTrace() {
	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: buf})

	{
		log.Start("1234", "Data_String")

		b1 := []byte{0xEE, 0x6E, 0x11, 0x00, 0x00, 0x00, 0x3E, 0xEA, 0xDE, 0x18, 0x00, 0x00, 0x2D, 0x00, 0x00, 0x00, 0x3E, 0xEA, 0xDE}
		b2 := []byte{0xEE, 0x6E, 0x11, 0x00, 0x00, 0x00, 0x3E, 0xEA, 0xDE, 0x18, 0x00, 0x00, 0x2D, 0x00, 0x00, 0x00, 0x3E, 0xEA, 0xDE}

		log.DataTrace("1234", "Data_String", Message(b1), Message(b2))

		log.Complete("1234", "Data_String")
	}

	log.Shutdown()
	fmt.Println(buf.String())
	// Output:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_String: Started:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_String: DATA:
	// 	Message Bytes:	(0x0000) EE 6E 11 00 00 00 3E EA DE 18 00 00 2D 00 00 00
	// 				(0x0010) 3E EA DE
	// 	Message Bytes:	(0x0000) EE 6E 11 00 00 00 3E EA DE 18 00 00 2D 00 00 00
	// 				(0x0010) 3E EA DE
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Data_String: Completed:
}

// ExampleTracef provides an example of logging from a fmt.Stringer and also tests newline handling.
func ExampleTracef() {
	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: buf})

	{
		// Messages without a newline should have one added, and message
		// that have a newline should *not* have an extra one added.
		log.Tracef("1234", "Tracef", "%s: %s", "1234", "Basic")
		log.Tracef("1234", "Tracef", "%s: %s\n", "1234", "Basic")
		log.Tracef("1234", "Tracef", "%s: %s", "ABCD", "Basic")
		log.Tracef("1234", "Tracef", "%s: %s\n", "ABCD", "Basic")
	}

	log.Shutdown()
	fmt.Println(buf.String())
	// Output:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Tracef: Trace: 1234: Basic
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Tracef: Trace: 1234: Basic
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Tracef: Trace: ABCD: Basic
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Tracef: Trace: ABCD: Basic
}

// ExampleUplevel_Start provides an example of using the level up functionality.
func ExampleUplevel_Start() {
	// The following code would generate a log message with line
	// number 13 (the line number where handleErr is called) and function name
	// "example.someFunc".
	//
	//    1: package example
	//    2:
	//    3: func handleErr(err error, context interface{}, function string) {
	//    4:        ...
	//    5:        // log with caller's file, line, and function name.
	//    6:        log.Up1.Err(err, context, function)
	//    7:        ...
	//    8: }
	//    9:
	//   10: func someFunc() {
	//   11:        ...
	//   12:        if err := doSomething(); err != nil {
	//   13:                handleErr(err, context, "")
	//   14:        }
	//   15:        ...
	//   16: }

	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: buf})

	{
		levelUp := func(context interface{}, function string) {
			log.Up1.Tracef(context, function, "Test")
		}

		log.Start("1234", "Basic")
		levelUp("1234", "Basic")
		log.Complete("1234", "Basic")
	}

	log.Shutdown()
	fmt.Println(buf.String())

	// Output:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Started:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Trace: Test
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Completed:
}

// Example_multipleInit tests that when Init is called back to back,
// the log flushes and takes the change.
func Example_multipleInit() {
	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: buf})
	log.Start("1234", "Basic")
	log.Complete("1234", "Basic")

	// Now init the log with a buffer to disable logging
	log.InitTest("EXAMPLE", 10, log.DevWriter{Device: log.DevAll, Writer: ioutil.Discard})
	log.Start("1234", "SHOULD NOT BE DISPLAYED")
	log.Complete("1234", "SHOULD NOT BE DISPLAYED")

	log.Shutdown()
	fmt.Println(buf.String())
	// Output:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Started:
	// 2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Completed:
}

func BenchmarkTracef(b *testing.B) {
	log.InitTest("BENCHMARK", 10, log.DevWriter{Device: log.DevAll, Writer: ioutil.Discard})
	for i := 0; i < b.N; i++ {
		log.Tracef("context", "function", "This is a test %d this is a test %d this is a test %d", i, i, i)
	}
}

// ExampleSplunk provides an example of logging a message in a splunk-able format.
func ExampleSplunk() {
	// Init the log system using a buffer for testing.
	buf := new(log.SafeBuffer)
	log.InitTest("TestSplunk", 10, log.DevWriter{Device: log.DevAll, Writer: buf})

	sl1 := log.SplunkValue{1, 2, 3, 4}           // slice of ints.
	sl2 := log.SplunkValue{"123.123", "123.124"} // slice of strings.
	sl3 := log.SplunkValue{6, "123.123"}         // slice of mixed values.

	m := []log.SplunkPair{
		{Key: "Key1", Value: "Value1"},
		{Key: "RequestTime", Value: time.Date(2019, time.November, 10, 15, 0, 0, 0, time.UTC).UTC().Format("2006/01/02 15:04:05.000")},
		{Key: "MAC", Value: "010203040506"},
		{Key: "ResponseCode", Value: 0},
		{Key: "Slice", Value: sl1},
		{Key: "name1", Value: sl2},
		{Key: "name2", Value: sl3},
	}
	log.Splunk(m...)

	log.Splunk(log.SplunkPair{Key: "SecondKey", Value: "SecondValue"},
		log.SplunkPair{Key: "RequestTime", Value: time.Date(2019, time.November, 10, 15, 0, 0, 0, time.UTC).UTC().Format("2006/01/02 15:04:05.000")},
		log.SplunkPair{Key: "MAC", Value: "010203040507"},
		log.SplunkPair{Key: "ResponseCode", Value: 0},
		log.SplunkPair{Key: "Slice", Value: sl1},
		log.SplunkPair{Key: "name1", Value: sl2},
		log.SplunkPair{Key: "name2", Value: sl3})

	log.Shutdown()
	fmt.Println(buf.String())

	// Output:
	// 2009/11/10 15:00:00.000: Key1=Value1 RequestTime="2019/11/10 15:00:00.000" MAC=010203040506 ResponseCode=0 Slice=[1, 2, 3, 4] name1=[123.123, 123.124] name2=[6, 123.123]
	// 2009/11/10 15:00:00.000: SecondKey=SecondValue RequestTime="2019/11/10 15:00:00.000" MAC=010203040507 ResponseCode=0 Slice=[1, 2, 3, 4] name1=[123.123, 123.124] name2=[6, 123.123]
}

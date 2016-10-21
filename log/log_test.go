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
	"bytes"
	"errors"
	"math"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Comcast/go-log/log"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// failed is the Unicode codepoint for an X mark.
const failed = "\u2717"

// logdest implements io.Writer and is the log package destination.
var logdest bytes.Buffer

// resetLog can be called at the beginning of a test or example.
func resetLog() { logdest.Reset() }

// displayLog can be called at the end of a test or example.
// It only prints the log contents if the -test.v flag is set.
func displayLog() {
	if !testing.Verbose() {
		return
	}
	logdest.WriteTo(os.Stderr)
}

// TestDevices tests that we can set multiple devices
func TestDevices(t *testing.T) {
	t.Log("Given the need to test multiple devices.")
	{
		var device1 bytes.Buffer
		var device2 bytes.Buffer

		log.InitTest("LOG", 10)
		log.Dev.Trace(&device1)
		log.Dev.Warning(&device2)

		t.Log("When we write to two different trace types.")
		{
			log.Tracef("Device1", "TestDevices", "Hello")
			log.Warnf("Device2", "TestDevices", "Hello")

			log.Shutdown()

			got := device1.String()
			exp := "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: Device1: TestDevices: Trace: Hello\n"

			if got == exp {
				t.Log("\t\tShould log the expected trace line for Trace.", succeed)
			} else {
				t.Error("\t\tShould log the expected trace line for Trace.", failed)
			}

			got = device2.String()
			exp = "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: Device2: TestDevices: Warning: Hello\n"

			if got == exp {
				t.Log("\t\tShould log the expected trace line for Warning.", succeed)
			} else {
				t.Error("\t\tShould log the expected trace line for Warning.", failed)
			}
		}
	}
}

// TestDataStringEmpty tests we receive the properly formatted
// trace line when message is empty.
func TestDataStringEmpty(t *testing.T) {
	t.Log("Given the need to validate the trace line for DataString calls.")
	{
		cases := []struct {
			message  string
			expected string
		}{
			{"test", "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestDataStringEmpty: DATA:\n\ttest\n"},
			{"", "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestDataStringEmpty: DATA: %!ds(MISSING)\n"},
		}

		for _, tt := range cases {
			t.Logf("\tWhen we use the value of \"%s\".", tt.message)
			{
				resetLog()
				log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
				log.DataString("TEST", "TestDataStringEmpty", tt.message)
				log.Shutdown()

				got := logdest.String()
				if got == tt.expected {
					t.Log("\t\tShould log the expected trace line.", succeed)
				} else {
					t.Error("\t\tShould log the expected trace line.", failed, got)
				}
			}
		}
	}
}

type blockWriter struct {
	Writes int32
	count  int32
}

// Write will simulate long periods of blocking. This will allow us
// to test that the program does not block on log writes.
func (b *blockWriter) Write(p []byte) (int, error) {
	c := atomic.AddInt32(&b.count, 1)
	if c == 1 {
		time.Sleep(100 * time.Millisecond)
	}
	atomic.AddInt32(&b.Writes, 1)
	return len(p), nil
}

// TestLogBlocking will test the log will not cause blocking when write
// begin to block.
func TestLogBlocking(t *testing.T) {
	var bw blockWriter

	log.InitTest("LOG", 0, log.DevWriter{Device: log.DevAll, Writer: &bw})

	t.Log("Given the need to make sure the logging doesn't stop the program.")
	{
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			log.Tracef("TEST", "function", "Log: %d", 0)
			wg.Done()
		}()

		for i := 0; i < 6; i++ {
			log.Tracef("TEST", "function", "Log: %d", i)
			time.Sleep(25 * time.Millisecond)
		}

		wg.Wait()

		c := atomic.LoadInt32(&bw.Writes)
		if c == 4 {
			t.Log("\tShould have only 4 log writes.", succeed)
		} else {
			t.Error("\tShould have only 4 log writes.", failed, bw.Writes)
		}
	}

	log.Shutdown()
}

type blockWriter2 struct {
	count int32
}

// Write will simulate long periods of blocking. This will allow us
// to test that the program does not block on log writes.
func (b *blockWriter2) Write(p []byte) (int, error) {
	c := atomic.AddInt32(&b.count, 1)
	if c%10 == 0 {
		time.Sleep(50 * time.Millisecond)
	}
	return len(p), nil
}

// TestLogBlocking will test the log will not cause blocking when write
// begin to block.
func TestLogBlocking2(t *testing.T) {
	var bw blockWriter2

	log.InitTest("LOG", 0, log.DevWriter{Device: log.DevAll, Writer: &bw})

	t.Log("Given the need to make sure the logging doesn't stop the program.")
	{
		now := time.Now()

		for i := 0; i < 5; i++ {
			for counter := 0; counter < 20; counter++ {
				log.Tracef("TEST", "function", "Log: %d", counter)
			}
			time.Sleep(50 * time.Millisecond)
		}

		dur := time.Since(now)

		c := atomic.LoadInt32(&bw.count)
		if c == 50 {
			t.Log("\tShould have only 50 log writes.", succeed)
		} else {
			t.Error("\tShould have only 50 log writes.", bw.count, failed)
		}

		const max = 500 * time.Millisecond

		if dur < max {
			t.Log("\tShould take less than 500 milliseconds", succeed)
		} else {
			t.Error("\tShould take less than 500 milliseconds", max, dur, failed)
		}
	}

	log.Shutdown()
}

// TestLoggingLevels tests that each logging level is working.
func TestLoggingLevels(t *testing.T) {
	t.Log("Given the need to test different logging levels.")
	{
		data := []struct {
			n string
			l int
			v int
		}{
			{"level 4", 4, 19},
			{"level 3", 3, 13},
			{"level 2", 2, 6},
			{"level 1", 1, 5},
		}

		for _, d := range data {
			t.Logf("\tWhen we are at logging %s", d.n)
			{
				var buf bytes.Buffer
				log.InitTest("LOG", 0, log.DevWriter{Device: log.DevAll, Writer: &buf})

				f := func() int {
					return d.l
				}
				l := log.NewLogger("test", f)

				// Level 1 : Total 4
				l.CompleteErr(errors.New("E"), "A", "B")       // 1 line
				l.CompleteErrf(errors.New("E"), "A", "B", "C") // 1 line
				l.Err(errors.New("E"), "A", "B")               // 1 line
				l.Errf(errors.New("E"), "A", "B", "C")         // 1 line

				// Level 2 : Total 5
				l.Warnf("A", "B", "C") // 1 line

				// Level 3 : Total 12
				l.DataKV("A", "B", "C", "D")                       // 1 lines
				l.DataBlock("A", "B", "C")                         // 2 lines
				l.DataString("A", "B", "C")                        // 2 lines
				l.DataTrace("A", "B", Message([]byte{0xEE, 0xEF})) // 2 lines

				// Level 4 : Total 18
				l.Start("A", "B")          // 1 line
				l.Startf("A", "B", "C")    // 1 line
				l.Complete("A", "B")       // 1 line
				l.Completef("A", "B", "C") // 1 line
				l.Tracef("A", "B", "C")    // 1 line
				l.Queryf("A", "B", "C")    // 1 line

				log.Shutdown()

				got := strings.Split(buf.String(), "\n")

				if len(got) == d.v {
					t.Logf("\t\tShould see %d trace lines. %v", d.v, succeed)
				} else {
					t.Errorf("\t\tShould see %d trace lines. %v %d", d.v, failed, len(got))
				}
			}
		}
	}
}

// TestLoggingUpLevels tests that each logging uplevel is working.
func TestLoggingUpLevels(t *testing.T) {
	t.Log("Given the need to test different logging uplevels.")
	{
		data := []struct {
			n string
			l int
			v int
		}{
			{"level 4", 4, 19},
			{"level 3", 3, 13},
			{"level 2", 2, 6},
			{"level 1", 1, 5},
		}

		for _, d := range data {
			t.Logf("\tWhen we are at logging %s", d.n)
			{
				var buf bytes.Buffer
				log.InitTest("LOG", 0, log.DevWriter{Device: log.DevAll, Writer: &buf})

				f := func() int {
					return d.l
				}
				l := log.NewLogger("test", f)

				// Level 1 : Total 4
				l.Up1.CompleteErr(errors.New("E"), "A", "B")       // 1 line
				l.Up1.CompleteErrf(errors.New("E"), "A", "B", "C") // 1 line
				l.Up1.Err(errors.New("E"), "A", "B")               // 1 line
				l.Up1.Errf(errors.New("E"), "A", "B", "C")         // 1 line

				// Level 2 : Total 5
				l.Up1.Warnf("A", "B", "C") // 1 line

				// Level 3 : Total 12
				l.Up1.DataKV("A", "B", "C", "D")                       // 1 lines
				l.Up1.DataBlock("A", "B", "C")                         // 2 lines
				l.Up1.DataString("A", "B", "C")                        // 2 lines
				l.Up1.DataTrace("A", "B", Message([]byte{0xEE, 0xEF})) // 2 lines

				// Level 4 : Total 18
				l.Up1.Start("A", "B")          // 1 line
				l.Up1.Startf("A", "B", "C")    // 1 line
				l.Up1.Complete("A", "B")       // 1 line
				l.Up1.Completef("A", "B", "C") // 1 line
				l.Up1.Tracef("A", "B", "C")    // 1 line
				l.Up1.Queryf("A", "B", "C")    // 1 line

				log.Shutdown()

				got := strings.Split(buf.String(), "\n")

				if len(got) == d.v {
					t.Logf("\t\tShould see %d trace lines. %v", d.v, succeed)
				} else {
					t.Errorf("\t\tShould see %d trace lines. %v %d", d.v, failed, len(got))
				}
			}
		}
	}
}

type SomeFormatter struct{}

func (SomeFormatter) Format() string {
	return "42"
}

type EmptyFormatter struct{}

func (EmptyFormatter) Format() string {
	return ""
}

func TestLoggingFuncs(t *testing.T) {
	t.Log("Given the need to call all different logging calls.")
	{
		const context = "TEST"
		cases := []struct {
			expected string
			f        func()
		}{
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: foo: Started:\n", func() {
				log.Start(context, "foo")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: foo: Started: walrus[500]\n", func() {
				log.Startf(context, "foo", "walrus[%d]", 500)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bar: Completed:\n", func() {
				log.Complete(context, "bar")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bar: Completed: horse[3]\n", func() {
				log.Completef(context, "bar", "horse[%d]", 3)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: baz: Completed ERROR: A\n", func() {
				log.CompleteErr(errors.New("A"), context, "baz")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: baz: Completed ERROR: puppies[777]: B\n", func() {
				log.CompleteErrf(errors.New("B"), context, "baz", "puppies[%d]", 777)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: boo: ERROR: C\n", func() {
				log.Err(errors.New("C"), context, "boo")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bee: ERROR: ip[127.0.0.1]: D\n", func() {
				log.Errf(errors.New("D"), context, "bee", "ip[%s]", "127.0.0.1")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: faa: Trace: len[13]\n", func() {
				log.Tracef(context, "faa", "len[%d]", 13)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: fii: Warning: usage[99.900000]\n", func() {
				log.Warnf(context, "fii", "usage[%f]", 99.9)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: beer: Query: howmany[0]\n", func() {
				log.Queryf(context, "beer", "howmany[%d]", 0)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: oom: DATA: 2b: !2b\n", func() {
				log.DataKV(context, "oom", "2b", "!2b")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\tasdf I'm running out of ideas.\n", func() {
				log.DataBlock(context, "moo", "asdf I'm running out of ideas.")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\t5\n", func() {
				log.DataBlock(context, "moo", 5)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\tjson: unsupported value: NaN\n", func() {
				log.DataBlock(context, "moo", math.NaN())
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: boo: DATA:\n\tnot\n\tthat\n\tthey\n\twere\n\tvery\n\tgood to start with.\n", func() {
				log.DataString(context, "boo", "not\nthat\nthey\nwere\nvery\ngood to start with.")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: boo: DATA:\n", func() {
				log.DataString(context, "boo", "\n\n\n\n")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: aah: DATA:\n\t42\n", func() {
				log.DataTrace(context, "aah", SomeFormatter{})
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: aah: DATA:\n", func() {
				log.DataTrace(context, "aah", EmptyFormatter{})
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: aah: DATA:\n", func() {
				log.DataTrace(context, "aah", nil)
			}},
		}
		for _, tt := range cases {

			resetLog()
			log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
			tt.f()
			log.Shutdown()

			got := logdest.String()
			if got != tt.expected {
				t.Errorf("\t\tLog should match expected. %s %q", failed, got)
				continue
			}
			t.Log("\t\tLog should match expected.", succeed)
		}
	}
}

func TestErrPanic(t *testing.T) {
	const expected = "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: ERROR: A\n2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: TERMINATING\n"
	defer func() {
		if r := recover(); r == nil {
			t.Error("\t\tErrPanic should have panicked.", failed)
		} else {
			t.Log("\t\tErrPanic should have panicked.", succeed)
		}

		got := logdest.String()
		if got != expected {
			t.Errorf("\t\tLog should match expected. %s %q", failed, got)
			return
		}
		t.Log("\t\tLog should match expected.", succeed)
	}()

	resetLog()
	log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
	log.ErrPanic(errors.New("A"), "TEST", "TestErrPanic")
}

func TestErrPanicf(t *testing.T) {
	const expected = "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: ERROR: we're doomed -bender-: A\n2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: TERMINATING\n"
	defer func() {
		if r := recover(); r == nil {
			t.Error("\t\tErrPanicf should have panicked.", failed)
		} else {
			t.Log("\t\tErrPanicf should have panicked.", succeed)
		}

		got := logdest.String()
		if got != expected {
			t.Errorf("\t\tLog should match expected. %s %q", failed, got)
			return
		}
		t.Log("\t\tLog should match expected.", succeed)
	}()

	resetLog()
	log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
	log.ErrPanicf(errors.New("A"), "TEST", "TestErrPanic", "we're doomed -%s-", "bender")
}

func TestLoggerFuncs(t *testing.T) {
	t.Log("Given the need to call all different logging calls.")
	{
		const context = "TEST"
		cases := []struct {
			expected string
			f        func(*log.Logger)
		}{
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: foo: Started:\n", func(ll *log.Logger) {
				ll.Start(context, "foo")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: foo: Started: walrus[500]\n", func(ll *log.Logger) {
				ll.Startf(context, "foo", "walrus[%d]", 500)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bar: Completed:\n", func(ll *log.Logger) {
				ll.Complete(context, "bar")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bar: Completed: horse[3]\n", func(ll *log.Logger) {
				ll.Completef(context, "bar", "horse[%d]", 3)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: baz: Completed ERROR: A\n", func(ll *log.Logger) {
				ll.CompleteErr(errors.New("A"), context, "baz")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: baz: Completed ERROR: puppies[777]: B\n", func(ll *log.Logger) {
				ll.CompleteErrf(errors.New("B"), context, "baz", "puppies[%d]", 777)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: boo: ERROR: C\n", func(ll *log.Logger) {
				ll.Err(errors.New("C"), context, "boo")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bee: ERROR: ip[127.0.0.1]: D\n", func(ll *log.Logger) {
				ll.Errf(errors.New("D"), context, "bee", "ip[%s]", "127.0.0.1")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: faa: Trace: len[13]\n", func(ll *log.Logger) {
				ll.Tracef(context, "faa", "len[%d]", 13)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: fii: Warning: usage[99.900000]\n", func(ll *log.Logger) {
				ll.Warnf(context, "fii", "usage[%f]", 99.9)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: beer: Query: howmany[0]\n", func(ll *log.Logger) {
				ll.Queryf(context, "beer", "howmany[%d]", 0)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: oom: DATA: 2b: !2b\n", func(ll *log.Logger) {
				ll.DataKV(context, "oom", "2b", "!2b")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\tasdf I'm running out of ideas.\n", func(ll *log.Logger) {
				ll.DataBlock(context, "moo", "asdf I'm running out of ideas.")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\t5\n", func(ll *log.Logger) {
				ll.DataBlock(context, "moo", 5)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\tjson: unsupported value: NaN\n", func(ll *log.Logger) {
				ll.DataBlock(context, "moo", math.NaN())
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: boo: DATA:\n\tnot\n\tthat\n\tthey\n\twere\n\tvery\n\tgood to start with.\n", func(ll *log.Logger) {
				ll.DataString(context, "boo", "not\nthat\nthey\nwere\nvery\ngood to start with.")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: aah: DATA:\n\t42\n", func(ll *log.Logger) {
				ll.DataTrace(context, "aah", SomeFormatter{})
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: aah: DATA:\n", func(ll *log.Logger) {
				ll.DataTrace(context, "aah", EmptyFormatter{})
			}},
		}
		for _, tt := range cases {

			resetLog()
			log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
			ll := log.NewLogger("LOG", func() int { return log.LevelTrace })
			tt.f(ll)
			log.Shutdown()

			got := logdest.String()
			if got != tt.expected {
				t.Errorf("\t\tLog should match expected. %s %q", failed, got)
				continue
			}
			t.Log("\t\tLog should match expected.", succeed)
		}
	}
}

func TestLoggerErrPanic(t *testing.T) {
	const expected = "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: ERROR: A\n2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: TERMINATING\n"
	defer func() {
		if r := recover(); r == nil {
			t.Error("\t\tErrPanic should have panicked.", failed)
		} else {
			t.Log("\t\tErrPanic should have panicked.", succeed)
		}

		got := logdest.String()
		if got != expected {
			t.Errorf("\t\tLog should match expected. %s %q", failed, got)
			return
		}
		t.Log("\t\tLog should match expected.", succeed)
	}()

	resetLog()
	log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
	ll := log.NewLogger("LOG", func() int { return log.LevelTrace })
	ll.ErrPanic(errors.New("A"), "TEST", "TestErrPanic")
}

func TestLoggerErrPanicf(t *testing.T) {
	const expected = "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: ERROR: we're doomed -bender-: A\n2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: TERMINATING\n"
	defer func() {
		if r := recover(); r == nil {
			t.Error("\t\tErrPanicf should have panicked.", failed)
		} else {
			t.Log("\t\tErrPanicf should have panicked.", succeed)
		}

		got := logdest.String()
		if got != expected {
			t.Errorf("\t\tLog should match expected. %s %q", failed, got)
			return
		}
		t.Log("\t\tLog should match expected.", succeed)
	}()

	resetLog()
	log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
	ll := log.NewLogger("LOG", func() int { return log.LevelTrace })
	ll.ErrPanicf(errors.New("A"), "TEST", "TestErrPanic", "we're doomed -%s-", "bender")
}

func TestUpLoggerFuncs(t *testing.T) {
	t.Log("Given the need to call all different logging calls.")
	{
		const context = "TEST"
		cases := []struct {
			expected string
			f        func(*log.UplevelLogger)
		}{
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: foo: Started:\n", func(ll *log.UplevelLogger) {
				ll.Start(context, "foo")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: foo: Started: walrus[500]\n", func(ll *log.UplevelLogger) {
				ll.Startf(context, "foo", "walrus[%d]", 500)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bar: Completed:\n", func(ll *log.UplevelLogger) {
				ll.Complete(context, "bar")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bar: Completed: horse[3]\n", func(ll *log.UplevelLogger) {
				ll.Completef(context, "bar", "horse[%d]", 3)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: baz: Completed ERROR: A\n", func(ll *log.UplevelLogger) {
				ll.CompleteErr(errors.New("A"), context, "baz")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: baz: Completed ERROR: puppies[777]: B\n", func(ll *log.UplevelLogger) {
				ll.CompleteErrf(errors.New("B"), context, "baz", "puppies[%d]", 777)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: boo: ERROR: C\n", func(ll *log.UplevelLogger) {
				ll.Err(errors.New("C"), context, "boo")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: bee: ERROR: ip[127.0.0.1]: D\n", func(ll *log.UplevelLogger) {
				ll.Errf(errors.New("D"), context, "bee", "ip[%s]", "127.0.0.1")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: faa: Trace: len[13]\n", func(ll *log.UplevelLogger) {
				ll.Tracef(context, "faa", "len[%d]", 13)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: fii: Warning: usage[99.900000]\n", func(ll *log.UplevelLogger) {
				ll.Warnf(context, "fii", "usage[%f]", 99.9)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: beer: Query: howmany[0]\n", func(ll *log.UplevelLogger) {
				ll.Queryf(context, "beer", "howmany[%d]", 0)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: oom: DATA: 2b: !2b\n", func(ll *log.UplevelLogger) {
				ll.DataKV(context, "oom", "2b", "!2b")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\tasdf I'm running out of ideas.\n", func(ll *log.UplevelLogger) {
				ll.DataBlock(context, "moo", "asdf I'm running out of ideas.")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\t5\n", func(ll *log.UplevelLogger) {
				ll.DataBlock(context, "moo", 5)
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: moo: DATA:\n\tjson: unsupported value: NaN\n", func(ll *log.UplevelLogger) {
				ll.DataBlock(context, "moo", math.NaN())
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: boo: DATA:\n\tnot\n\tthat\n\tthey\n\twere\n\tvery\n\tgood to start with.\n", func(ll *log.UplevelLogger) {
				ll.DataString(context, "boo", "not\nthat\nthey\nwere\nvery\ngood to start with.")
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: aah: DATA:\n\t42\n", func(ll *log.UplevelLogger) {
				ll.DataTrace(context, "aah", SomeFormatter{})
			}},
			{"2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: aah: DATA:\n", func(ll *log.UplevelLogger) {
				ll.DataTrace(context, "aah", EmptyFormatter{})
			}},
		}
		for _, tt := range cases {

			resetLog()
			log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
			ll := log.NewLogger("LOG", func() int { return log.LevelTrace })
			tt.f(&ll.Up1)
			log.Shutdown()

			got := logdest.String()
			if got != tt.expected {
				t.Errorf("\t\tLog should match expected. %s %q", failed, got)
				continue
			}
			t.Log("\t\tLog should match expected.", succeed)
		}
	}
}

func TestUpLoggerErrPanic(t *testing.T) {
	const expected = "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: ERROR: A\n2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: TERMINATING\n"
	defer func() {
		if r := recover(); r == nil {
			t.Error("\t\tErrPanic should have panicked.", failed)
		} else {
			t.Log("\t\tErrPanic should have panicked.", succeed)
		}

		got := logdest.String()
		if got != expected {
			t.Errorf("\t\tLog should match expected. %s %q", failed, got)
			return
		}
		t.Log("\t\tLog should match expected.", succeed)
	}()

	resetLog()
	log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
	ll := log.NewLogger("LOG", func() int { return log.LevelTrace })
	ll.Up1.ErrPanic(errors.New("A"), "TEST", "TestErrPanic")
}

func TestUpLoggerErrPanicf(t *testing.T) {
	const expected = "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: ERROR: we're doomed -bender-: A\n2009/11/10 15:00:00.000: LOG[69910]: file.go#512: TEST: TestErrPanic: TERMINATING\n"
	defer func() {
		if r := recover(); r == nil {
			t.Error("\t\tErrPanicf should have panicked.", failed)
		} else {
			t.Log("\t\tErrPanicf should have panicked.", succeed)
		}

		got := logdest.String()
		if got != expected {
			t.Errorf("\t\tLog should match expected. %s %q", failed, got)
			return
		}
		t.Log("\t\tLog should match expected.", succeed)
	}()

	resetLog()
	log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &logdest})
	ll := log.NewLogger("LOG", func() int { return log.LevelTrace })
	ll.Up1.ErrPanicf(errors.New("A"), "TEST", "TestErrPanic", "we're doomed -%s-", "bender")
}

func TestDoubleInit(t *testing.T) {
	log.InitTest("TEST", 0, log.DevWriter{Device: log.DevAll, Writer: new(bytes.Buffer)})

	var buf bytes.Buffer
	log.InitTest("LOG", 10, log.DevWriter{Device: log.DevAll, Writer: &buf})

	log.Warnf("ctx", "ExampleLog", "Hola, mundo")

	log.Shutdown()

	expected := "2009/11/10 15:00:00.000: LOG[69910]: file.go#512: ctx: ExampleLog: Warning: Hola, mundo\n"
	if got := buf.String(); got != expected {
		t.Errorf("Got:\n%q\nExpected:\n%q", got, expected)
	}
}

// TestLineNumber will ensure that the line numbers logged are correct.
func TestLineNumbers(t *testing.T) {
	context := "TestLineNumbers"
	str := "dummy string"
	dummyErr := errors.New("dummy error")

	// Will not be performing panic-related tests

	var buf bytes.Buffer
	log.Init(context, 0, log.DevWriter{Device: log.DevAll, Writer: &buf})
	defer log.Shutdown()

	logger := log.NewLogger("logger", func() int { return log.LevelTrace }) // highest level to display all messages

	_, _, thisLineNum, _ := runtime.Caller(0)
	lineDiff := 4

	thisLineNum += lineDiff
	log.Complete(context, str)
	testLineNumber(t, "log.Complete", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.CompleteErr(dummyErr, context, str)
	testLineNumber(t, "log.CompleteErr", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.CompleteErrf(dummyErr, context, str, str)
	testLineNumber(t, "log.CompleteErrf", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.Completef(context, str, str)
	testLineNumber(t, "log.Completef", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.DataBlock(context, str, nil)
	testLineNumber(t, "log.DataBlock", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.DataKV(context, str, str, nil)
	testLineNumber(t, "log.DataKV", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.DataString(context, str, str)
	testLineNumber(t, "log.DataString", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.DataTrace(context, str, nil)
	testLineNumber(t, "log.DataTrace", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.Err(dummyErr, context, str)
	testLineNumber(t, "log.Err", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.Errf(dummyErr, context, str, str)
	testLineNumber(t, "log.Errf", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.Queryf(context, str, str)
	testLineNumber(t, "log.Queryf", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.Start(context, str)
	testLineNumber(t, "log.Start", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.Startf(context, str, str)
	testLineNumber(t, "log.Startf", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.Tracef(context, str, str)
	testLineNumber(t, "log.Tracef", &buf, thisLineNum)

	thisLineNum += lineDiff
	log.Warnf(context, str, str)
	testLineNumber(t, "log.Warnf", &buf, thisLineNum)
	//---------------------
	thisLineNum += lineDiff
	logger.Complete(context, str)
	testLineNumber(t, "logger.Complete", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.CompleteErr(dummyErr, context, str)
	testLineNumber(t, "logger.CompleteErr", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.CompleteErrf(dummyErr, context, str, str)
	testLineNumber(t, "logger.CompleteErrf", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.Completef(context, str, str)
	testLineNumber(t, "logger.Completef", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.DataBlock(context, str, nil)
	testLineNumber(t, "logger.DataBlock", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.DataKV(context, str, str, nil)
	testLineNumber(t, "logger.DataKV", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.DataString(context, str, str)
	testLineNumber(t, "logger.DataString", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.DataTrace(context, str, nil)
	testLineNumber(t, "logger.DataTrace", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.Err(dummyErr, context, str)
	testLineNumber(t, "logger.Err", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.Errf(dummyErr, context, str, str)
	testLineNumber(t, "logger.Errf", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.Queryf(context, str, str)
	testLineNumber(t, "logger.Queryf", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.Start(context, str)
	testLineNumber(t, "logger.Start", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.Startf(context, str, str)
	testLineNumber(t, "logger.Startf", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.Tracef(context, str, str)
	testLineNumber(t, "logger.Tracef", &buf, thisLineNum)

	thisLineNum += lineDiff
	logger.Warnf(context, str, str)
	testLineNumber(t, "logger.Warnf", &buf, thisLineNum)

	thisLineNum += lineDiff
	testLoggerUp1(t, logger, &buf, thisLineNum)
}

// testLoggerUp1 ensures that Up1 calls produce the correct line number.
func testLoggerUp1(t *testing.T, logger *log.Logger, buf *bytes.Buffer, expectedLineNumber int) {
	context := "testLoggerUp1"
	str := "dummy string"
	dummyErr := errors.New("dummy error")

	// Will not be performing panic-related tests

	logger.Up1.Complete(context, str)
	testLineNumber(t, "logger.Up1.Complete", buf, expectedLineNumber)

	logger.Up1.CompleteErr(dummyErr, context, str)
	testLineNumber(t, "logger.Up1.CompleteErr", buf, expectedLineNumber)

	logger.Up1.CompleteErrf(dummyErr, context, str, str)
	testLineNumber(t, "logger.Up1.CompleteErrf", buf, expectedLineNumber)

	logger.Up1.Completef(context, str, str)
	testLineNumber(t, "logger.Up1.Completef", buf, expectedLineNumber)

	logger.Up1.DataBlock(context, str, nil)
	testLineNumber(t, "logger.Up1.DataBlock", buf, expectedLineNumber)

	logger.Up1.DataKV(context, str, str, nil)
	testLineNumber(t, "logger.Up1.DataKV", buf, expectedLineNumber)

	logger.Up1.DataString(context, str, str)
	testLineNumber(t, "logger.Up1.DataString", buf, expectedLineNumber)

	logger.Up1.DataTrace(context, str, nil)
	testLineNumber(t, "logger.Up1.DataTrace", buf, expectedLineNumber)

	logger.Up1.Err(dummyErr, context, str)
	testLineNumber(t, "logger.Up1.Err", buf, expectedLineNumber)

	logger.Up1.Errf(dummyErr, context, str, str)
	testLineNumber(t, "logger.Up1.Errf", buf, expectedLineNumber)

	logger.Up1.Queryf(context, str, str)
	testLineNumber(t, "logger.Up1.Queryf", buf, expectedLineNumber)

	logger.Up1.Start(context, str)
	testLineNumber(t, "logger.Up1.Start", buf, expectedLineNumber)

	logger.Up1.Startf(context, str, str)
	testLineNumber(t, "logger.Up1.Startf", buf, expectedLineNumber)

	logger.Up1.Tracef(context, str, str)
	testLineNumber(t, "logger.Up1.Tracef", buf, expectedLineNumber)

	logger.Up1.Warnf(context, str, str)
	testLineNumber(t, "logger.Up1.Warnf", buf, expectedLineNumber)
}

// testLineNumber processes the logging line, extracts the line number and compares it against what
// is expected.
func testLineNumber(t *testing.T, testCall string, buf *bytes.Buffer, expectedLineNumber int) {
	// sleep a little before reading to make sure the string gets pushed in the buffer.
	time.Sleep(50 * time.Millisecond)

	str := buf.String()
	buf.Reset() // done with the buffer, clean it

	// Line number follows the pound sign
	re := regexp.MustCompile(`#(\d*)`)

	list := re.FindAllStringSubmatch(str, 1)
	if len(list) == 0 || len(list[0]) < 2 {
		t.Errorf("%s: Failed to find line number in output: %s", testCall, str)
		return
	}
	numStr := list[0][1]
	n, err := strconv.Atoi(numStr)
	if err != nil {
		t.Errorf("%s: Bad number format: %s, err: %v", testCall, numStr, err)
	}

	if n != expectedLineNumber {
		t.Errorf("%s: Str: %q: Expected line number %d, got %d", str, testCall, expectedLineNumber, n)
	}
}

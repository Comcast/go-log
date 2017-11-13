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
	"os"
	"regexp"
	"testing"
	"time"
)

func TestDtFile(t *testing.T) {
	const calldepth = 20
	t.Log("Given the need to get date time and file.")
	{
		expectedFuncName := "TestDtFile"
		dateTime, file, funcName, pid := dtFile(2, expectedFuncName)

		// At time of writing this function will return "testing.go#485". But adding
		// test might change the line number may change the second part and there's
		// no guarantee that go will always run tests inside testing.go. So we use a
		// regex to verify as close as possible.
		match, err := regexp.Match("[a-zA-Z0-9].go#\\d+", ([]byte)(file))
		if err != nil {
			t.Error("the regex is broken, please fix the test.", err)
			return
		}
		if !match {
			t.Error("\tfile should match with expected regex. ", failed)
		} else {
			t.Log("\tfile should match with expected regex. ", succeed)
		}

		// verify the process id
		if pid != os.Getpid() {
			t.Error("\tProcess ID should match.", failed)
		} else {
			t.Log("\tProcess ID should match.", succeed)
		}

		//verify funcname
		if funcName != expectedFuncName {
			t.Error("\tfuncName should match.", failed)
		} else {
			t.Log("\tfuncName should match.", succeed)
		}

		// verify that dateTime is roughly now. It's the best we can do.
		dt, err := time.Parse(layout, dateTime)
		if err != nil {
			t.Error("\tdateTime should be parsable.", failed)
		} else {
			t.Log("\tdateTime should be parsable.", succeed)
		}

		// If they're more then 5 seconds apart from each other then we either
		// have a performance problem on go-log, something terrible happened on
		// the machine performing the test, or dtFile is just plain wrong.
		timeDiff := time.Now().Sub(dt)
		if timeDiff > 5*time.Second {
			t.Error("\tNow-dateTime shoud be less then 5 second away.", failed)
		} else {
			t.Log("\tNow-dateTime shoud be less then 5 second away.", succeed)
		}
	}
	t.Log("Given a way too big caller depth")
	{
		expectedFuncName := "TestDtFile"
		dateTime, file, funcName, pid := dtFile(calldepth, expectedFuncName)

		// with a broken caller depth the filename returned should be unknown
		// and line number is zero
		if file != "unknown.go#0:" {
			t.Error("\tfile should match \"unknown.go#0:\". ", failed)
		} else {
			t.Log("\tfile should match \"unknown.go#0:\". ", succeed)
		}

		// verify the process id
		if pid != os.Getpid() {
			t.Error("\tProcess ID should match.", failed)
		} else {
			t.Log("\tProcess ID should match.", succeed)
		}

		// verify funcname, it should be unknown because the caller depth is too
		// big
		if funcName != "missing" {
			t.Error("\tfuncName should match.", failed)
		} else {
			t.Log("\tfuncName should match.", succeed)
		}

		// verify that dateTime is roughly now. It's the best we can do.
		dt, err := time.Parse(layout, dateTime)
		if err != nil {
			t.Error("\tdateTime should be parsable.", failed)
		} else {
			t.Log("\tdateTime should be parsable.", succeed)
		}

		// If they're more then 5 seconds apart from each other then we either
		// have a performance problem on go-log, something terrible happened on
		// the machine performing the test, or dtFile is just plain wrong.
		timeDiff := time.Now().Sub(dt)
		if timeDiff > 5*time.Second {
			t.Error("\tNow-dateTime shoud be less then 5 second away.", failed)
		} else {
			t.Log("\tNow-dateTime shoud be less then 5 second away.", succeed)
		}
	}

	t.Log("Given no funcName.")
	{
		// this is the actual function that is using this test.
		expectedFuncName := "testing.tRunner"
		dateTime, file, funcName, pid := dtFile(2, "")

		// At time of writing this function will return "testing.go#485". But adding
		// test might change the line number may change the second part and there's
		// no guarantee that go will always run tests inside testing.go. So we use a
		// regex to verify as close as possible.
		match, err := regexp.Match("[a-zA-Z0-9].go#\\d+", ([]byte)(file))
		if err != nil {
			t.Error("the regex is broken, please fix the test.", err)
			return
		}
		if !match {
			t.Error("\tfile should match with expected regex. ", failed)
		} else {
			t.Log("\tfile should match with expected regex. ", succeed)
		}

		// verify the process id
		if pid != os.Getpid() {
			t.Error("\tProcess ID should match.", failed)
		} else {
			t.Log("\tProcess ID should match.", succeed)
		}

		//verify funcname
		if funcName != expectedFuncName {
			t.Error("\tfuncName should match.", failed)
		} else {
			t.Log("\tfuncName should match.", succeed)
		}

		// verify that dateTime is roughly now. It's the best we can do.
		dt, err := time.Parse(layout, dateTime)
		if err != nil {
			t.Error("\tdateTime should be parsable.", failed)
		} else {
			t.Log("\tdateTime should be parsable.", succeed)
		}

		// If they're more then 5 seconds apart from each other then we either
		// have a performance problem on go-log, something terrible happened on
		// the machine performing the test, or dtFile is just plain wrong.
		timeDiff := time.Now().Sub(dt)
		if timeDiff > 5*time.Second {
			t.Error("\tNow-dateTime shoud be less then 5 second away.", failed)
		} else {
			t.Log("\tNow-dateTime shoud be less then 5 second away.", succeed)
		}
	}

	t.Log("Given a way too big caller depth and no funcName")
	{
		// here we are testing wether "missing" takes precedence over not giving
		// a funcName.
		expectedFuncName := "missing"
		dateTime, file, funcName, pid := dtFile(calldepth, "")

		// with a broken caller depth the filename returned should be unknown
		// and line number is zero
		if file != "unknown.go#0:" {
			t.Error("\tfile should match \"unknown.go#0:\". ", failed)
		} else {
			t.Log("\tfile should match \"unknown.go#0:\". ", succeed)
		}

		// verify the process id
		if pid != os.Getpid() {
			t.Error("\tProcess ID should match.", failed)
		} else {
			t.Log("\tProcess ID should match.", succeed)
		}

		// verify funcname, it should be unknown because the caller depth is too
		// big
		if funcName != expectedFuncName {
			t.Error("\tfuncName should match.", failed)
		} else {
			t.Log("\tfuncName should match.", succeed)
		}

		// verify that dateTime is roughly now. It's the best we can do.
		dt, err := time.Parse(layout, dateTime)
		if err != nil {
			t.Error("\tdateTime should be parsable.", failed)
		} else {
			t.Log("\tdateTime should be parsable.", succeed)
		}

		// If they're more then 5 seconds apart from each other then we either
		// have a performance problem on go-log, something terrible happened on
		// the machine performing the test, or dtFile is just plain wrong.
		timeDiff := time.Now().Sub(dt)
		if timeDiff > 5*time.Second {
			t.Error("\tNow-dateTime shoud be less then 5 second away.", failed)
		} else {
			t.Log("\tNow-dateTime shoud be less then 5 second away.", succeed)
		}
	}
}

func TestOutput(t *testing.T) {
	t.Log("Given no format passed to output.")
	{
		var buf SafeBuffer
		Init("TEST", 0, DevWriter{
			Device: DevAll,
			Writer: &buf,
		})

		// We expect this to generate some message because format is empty.
		output(&buf, "")

		// don't defer the shutdown because we need a clean start for the next
		// part of the test.
		Shutdown()

		if buf.String() != emptyMessage {
			t.Error("\tempty format should generate error message.", failed)
		}
		t.Log("\tempty format should generate error message.", succeed)

	}
	t.Log("Given no format passed to output.")
	{
		var buf bytes.Buffer
		Init("TEST", 0, DevWriter{
			Device: DevAll,
			Writer: &buf,
		})

		// Shutdown right now to trigger the no write on `output`.
		Shutdown()

		// We expect this to generate some message because format is empty.
		output(&buf, "")

		if buf.String() != "" {
			t.Error("\tempty format should contain nothing.", failed)
		}
		t.Log("\tempty format should contain nothing.", succeed)
	}
}

// This test doesn't actually work because if you pass nil to output another
// goroutine panics, not this one. I need a second opinion on this. However, I
// have added a check for nil in output to prevent it from panicking.
func TestOutputNilWriter(t *testing.T) {
	t.Log("Given no writer passed to output.")
	{
		var buf bytes.Buffer
		Init("TEST", 0, DevWriter{
			Device: DevAll,
			Writer: &buf,
		})

		defer Shutdown()

		defer func() {
			if r := recover(); r != nil {
				t.Error("\tGiving a nil writer to output should not panic.", failed)
			}
			t.Log("\tGiving a nil writer to output should not panic.", succeed)
		}()

		// Should not panic if writer is nil
		output(nil, "Asdf %d", 2)
	}
}

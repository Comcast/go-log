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
	"io"
	"os"
	"testing"
)

// succeed is the Unicode codepoint for a check mark.
const succeed = "\u2713"

// failed is the Unicode codepoint for an X mark.
const failed = "\u2717"

func TestDevAPI(t *testing.T) {
	// Setup the logger for this test.
	Init("TEST", 0, DevWriter{})
	defer Shutdown()

	possibleDevices := [...]struct {
		Device  int8
		SetFunc func(io.Writer)
	}{
		{DevStart, Dev.Start},
		{DevError, Dev.Error},
		{DevPanic, Dev.Panic},
		{DevTrace, Dev.Trace},
		{DevWarning, Dev.Warning},
		{DevQuery, Dev.Query},
		{DevData, Dev.Data},
		{DevSplunk, Dev.Splunk},
	}

	t.Log("Given the need to set all devices.")
	{

		// Set all writers to stdout
		Dev.All(os.Stdout)

		// Test that the alldevice was indeed set to every single device.
		for _, d := range possibleDevices {
			if Dev.get(d.Device) != os.Stdout {
				t.Errorf("\tDevice %d should be Stdout. %s", d.Device, failed)
				continue
			}
			t.Logf("\tDevice %d should be Stdout. %s", d.Device, succeed)
		}
	}

	// just set everything to nil and verify that indeed they are nil.
	Dev.All(nil)
	for _, d := range possibleDevices {
		if Dev.get(d.Device) != nil {
			t.Errorf("\tDevice %d should be nil. %s", d.Device, failed)
			continue
		}
	}

	t.Log("Given the need to set each device individually.")
	{
		for _, d := range possibleDevices {
			d.SetFunc(os.Stdout)
			if Dev.get(d.Device) != os.Stdout {
				t.Errorf("\tDevice %d should be Stdout. %s", d.Device, failed)
			}
			t.Logf("\tDevice %d should be Stdout. %s", d.Device, succeed)
		}
	}
}

func TestInitOnlyOneDevice(t *testing.T) {
	// Initialize the library with a single device.
	Init("TEST", 0, DevWriter{
		Device: DevStart,
		// Here we use stdin because the library uses stdout or stderr by
		// default. With stdin we can verify that only 1 device was set and make
		// sure it wasn't the default.
		Writer: os.Stdin,
	})

	// When we're done the test we clean up.
	defer Shutdown()

	nilDevice := [...]int8{DevError, DevPanic, DevTrace, DevWarning,
		DevQuery, DevData, DevSplunk}

	if Dev.get(DevStart) != os.Stdin {
		t.Error("\tDevice DevStart should be Stdout.", failed)
		return
	}
	t.Log("\tDevice DevStart should be Stdout.", succeed)

	for _, d := range nilDevice {
		// verify that only DevStart was set to stdin
		if Dev.get(d) == os.Stdin {
			t.Errorf("\tDevice %d should not be stdin. %s", d, failed)
			continue
		}
		t.Logf("\tDevice %d should not be stdin. %s", d, succeed)
	}
}

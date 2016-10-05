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

import "io"

// Set of constants that represent different trace lines
// types. Used to map different devices to the types.
const (
	// DevAll will update all devices at the time it is applied.
	DevAll int8 = iota

	DevStart
	DevError
	DevPanic
	DevTrace
	DevWarning
	DevQuery
	DevData
	DevSplunk
)

// DevWriter can be used in Init to change the default
// writers for use.
type DevWriter struct {
	Device int8
	Writer io.Writer
}

// dev provides a name space for the device methods.
type dev struct{}

// Dev provides access to the set of device methods. The goal of this method
// set is to allow the library user to redirect certain method, like Err or
// Warning, to different devices, like StdErr or StdOut.
var Dev dev

// get returns the device for the specified type.
func (dev) get(d int8) io.Writer {
	var w io.Writer

	l.destMu.RLock()
	{
		w = l.dest[d]
	}
	l.destMu.RUnlock()

	return w
}

// All sets all destinations to the specified device.
func (dev) All(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevStart] = w
		l.dest[DevError] = w
		l.dest[DevPanic] = w
		l.dest[DevTrace] = w
		l.dest[DevWarning] = w
		l.dest[DevQuery] = w
		l.dest[DevData] = w
		l.dest[DevSplunk] = w
	}
	l.destMu.Unlock()
}

// Start sets the Start and Complete functions device.
func (dev) Start(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevStart] = w
	}
	l.destMu.Unlock()
}

// Error sets the Error functions device.
func (dev) Error(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevError] = w
	}
	l.destMu.Unlock()
}

// Panic sets the panic functions device.
func (dev) Panic(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevPanic] = w
	}
	l.destMu.Unlock()
}

// Trace sets the trace functions device.
func (dev) Trace(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevTrace] = w
	}
	l.destMu.Unlock()
}

// Warning sets the warning functions device.
func (dev) Warning(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevWarning] = w
	}
	l.destMu.Unlock()
}

// Query sets the query functions device.
func (dev) Query(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevQuery] = w
	}
	l.destMu.Unlock()
}

// Data sets the data functions device.
func (dev) Data(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevData] = w
	}
	l.destMu.Unlock()
}

// Splunk sets the splunk functions device.
func (dev) Splunk(w io.Writer) {
	l.destMu.Lock()
	{
		l.dest[DevSplunk] = w
	}
	l.destMu.Unlock()
}

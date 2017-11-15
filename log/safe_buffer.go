/**
* Copyright 2017 Comcast Cable Communications Management, LLC
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
	"io"
	"sync"
)

// SafeBuffer is an object that can be used to safely use bytes.Buffer.
// It uses a mutex to protect the buffer and wraps few methods that can be used during
// various test cases.
type SafeBuffer struct {
	mu sync.Mutex // Mutex to safeguard the buffer
	b  bytes.Buffer
}

// Write is a wrapper to safely call bytes.Buffer's Write.
func (b *SafeBuffer) Write(ab []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.b.Write(ab)
}

// WriteTo is a wrapper to safely call bytes.Buffer's WriteTo.
func (b *SafeBuffer) WriteTo(w io.Writer) (int64, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.b.WriteTo(w)
}

// String is a wrapper to safely call bytes.Buffer's String.
func (b *SafeBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.b.String()
}

// Reset is a wrapper to safely call bytes.Buffer's Reset.
func (b *SafeBuffer) Reset() {
	b.mu.Lock()
	b.b.Reset()
	b.mu.Unlock()
}

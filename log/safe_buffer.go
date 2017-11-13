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
	mu sync.Mutex // Mutext to safeguard the buffer
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

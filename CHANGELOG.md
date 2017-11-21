## Changelog

### 11-21-2017

- [Issue #8](https://github.com/Comcast/go-log/issues/8) - Update timestamp nanosecond precision to 9 digits.

### 11-13-2017
- [Issue #5](https://github.com/Comcast/go-log/issues/5) - Exports SafeBuffer so that (external) tests can use it.

### 11-03-2017
- [Issue #3](https://github.com/Comcast/go-log/issues/3) - Emits logs in bulk for better performance
  - Introduced bulk logs timer (defaults to 1s)
  - Made stall timeout variable

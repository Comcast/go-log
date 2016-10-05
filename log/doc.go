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

// Package log is an important part of the application and having a consistent
// logging mechanism and structure is mandatory. With several teams writing
// different components that talk to each other, being able to read each others
// logs could be the difference between finding bugs quickly or wasting hours.
//
// Loggers
//
// With the log package we have the ability to create custom loggers that can
// be configured to write to one or many devices. This not only simplifies things,
// but will keep each log trace in correct sequence.
//
// Logging levels
//
// This package includes logging levels. Not everything needs to be logged all the time,
// so logging can be performed on a module/package basis.
//
// Tracing Formats
//
// There are two types of tracing lines we need to log. One is a trace line that
// describes where the program is, what it is doing and any data associated with
// that trace. The second is formatted data such as a JSON document or binary
// dump of data. Each serve a different purpose but they both exists within the
// same scope of space and time.
//
// The format of each trace line needs to be consistent and helpful or else the
// logging will just be noise and ultimately useless.
//
//     YYYY/MM/DD HH:MM:SS.ZZZ: APP[PID]: file.go#LN: Context: Func: Tag: Var[value]: Messages
//
// Here is a breakdown of each section and a sample value:
//
//     YYYY/MM/DD       Date of the trace log line in UTC.
//                      Ex. 2015/03/23
//
//     HH:MM:SS.ZZZZZZ  Time of the trace log line with microsecond in UTC.
//                      Ex. 14:02:42.123
//
//     APP              The application or service name. Set on Init.
//
//     PID              The process id for the running program.
//
//
//     file.go#LN       The name of the source code file and line number.
//                      Ex. main.go#15
//
//     Context:         Any context that is passed to the logging function.
//
//     Func:            The name of the function writing the trace log.
//
//     Tag:             The tag of trace line.
//       Started:           Start of a function/method call.
//       Completed:         Return of a function/method call.
//       Completed ERROR:   Return of a function/method call with error.
//       ERROR:             Error trace.
//       TERMINATING:       Termination of the application.
//       Trace:             All messages.
//       Warning:           Warning trace.
//       Query:             Any query that can be copied/pasted and run.
//       DATA:              A dump of data
//
//     Var[value]:      Optional, Data values, parameters or return values.
//                      Ex. ID[1234]
//
//     Messages:        Optional, Any extended information with proper grammar.
//                      Ex. Waiting on SMSC to acknowledge request.
//
// Here are examples of how trace lines would show in the log:
//
//     2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Started:
//     2009/11/10 15:00:00.000: EXAMPLE[69910]: file.go#512: 1234: Basic: Completed: Conv[10]
//
// API Documentation and Examples
//
// The API for the log package is focused on initializing the logger and then
// provides function abstractions for the different tags we have defined.
//
package log

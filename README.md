[![Xybor founder](https://img.shields.io/badge/xybor-huykingsofm-red)](https://github.com/huykingsofm)
[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/xylog.svg)](https://pkg.go.dev/github.com/xybor-x/xylog)
[![GitHub Repo stars](https://img.shields.io/github/stars/xybor-x/xylog?color=yellow)](https://github.com/xybor-x/xylog)
[![GitHub top language](https://img.shields.io/github/languages/top/xybor-x/xylog?color=lightblue)](https://go.dev/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor-x/xylog)](https://go.dev/blog/go1.18)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/xybor-x/xylog?include_prereleases)](https://github.com/xybor-x/xylog/releases/latest)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/e9146e2ccfd745a48a42fc723918d4cb)](https://www.codacy.com/gh/xybor-x/xylog/dashboard?utm_source=github.com&utm_medium=referral&utm_content=xybor-x/xylog&utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/e9146e2ccfd745a48a42fc723918d4cb)](https://www.codacy.com/gh/xybor-x/xylog/dashboard?utm_source=github.com&utm_medium=referral&utm_content=xybor-x/xylog&utm_campaign=Badge_Coverage)
[![Go Report](https://goreportcard.com/badge/github.com/xybor-x/xylog)](https://goreportcard.com/report/github.com/xybor-x/xylog)

# Introduction

Package xylog is a logging module based on the design of python logging

# Feature

The basic structs defined by the module, together with their functions, are
listed below:

1.  `Logger` is directly used by application code. It creates `LogRecord` and
    sends to `Handlers`.
2.  `Handler` converts `LogRecord` (created by `Logger`) to logging messages and
    sends to `Emitters`.
3.  `Emitter` writes logging messages (created by `Handler`) to appropriate
    destination.
4.  `Filter` is used by `Logger` and `Handler` to determine which `LogRecord`
    should be logged.
5.  `Formatter` is used by `Handler` to specify how a `LogRecord` is converted
    to the logging message.

## Logger

`Logger` should NEVER be instantiated directly. Instead, it will be created
through the function `GetLogger(name)`. Multiple calls to `GetLogger()` with the
same name will always return the same `Logger` object.

`Logger` names are dot-separated hierarchical names, such as "a", "a.b", "a.b.c"
or similar. For "a.b.c", its parents are "a" and "a.b".

When a `LogRecord` passes through a `Logger`, it will be handled by all
`Handlers` of `Logger` itself and `Logger`'s parents.

You can logs a message by using one of built-in logging methods, such as:

```golang
func Log(level int, a ...any)
func Logf(level int, msg string, a ...any)

func Debug(a ...any)
func Debugf(msg string, a ...any)
```

### EventLogger

`EventLogger` is a `Logger` wrapper supporting to compose logging message by
key-value fields.

Use `Event` method of `Logger` to create a `EventLogger`. You must create a
unique `EventLogger` each time you want to log.

If you call `Logger.AddField`, every `EventLogger` created by that `Logger`
always contains the added fields without adding again.

`EventLogger` provides structured format (default) and JSON format for logging
message.

## Logging level

The numeric values of logging levels are given in the following table. These are
primarily of interest if you want to define your own levels, and need them to
have specific values relative to the predefined levels. If you define a level
with the same numeric value, it overwrites the predefined value.

| Level        | Numeric value |
| ------------ | ------------- |
| CRITICAL     | 50            |
| ERROR/FATAL  | 40            |
| WARN/WARNING | 30            |
| INFO         | 20            |
| DEBUG        | 10            |
| NOTSET       | 0             |

## Handler

`Handler` handles `LogRecord` and converts it to logging message.

Like `Logger`, `Handler` can also be determined by its name. But the name is not
hierarchical.

To get an `Handler`, call `GetHandler` with its name. If the `Handler` doesn't
yet existed, create a new one. Many calls to `GetHandler` with the same name
will always give the same `Handler`. An exception is empty name which represents
for anonymous `Handlers`.

## Emitter

`Emitter` writes log messages to specified destination.

`StreamEmitter` can be used to print logging message into `stdout` or `stderr`.

`FileEmitter` can be used to write logging message to files. Library provides a
capability of rotating log with limited size or time.

## Formatter

`Formatter` converts a `LogRecord` to text.

Attributes of `LogRecord` are called macro. Macros' value is filled when the
`LogRecord` is initiated. Using macros is the easy way to construct a logging
message with dynamic and complex values.

`TextFormatter` is a simple `Formatter` which uses macros and format string to
format the message.

`JSONFormatter` allows to create a logging message of JSON format.

`StructureFormatter` allows to create a logging message with format of
`key=value`.

| MACRO             | DESCRIPTION                                                                                                                                      |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------------------ |
| `asctime`         | Textual time when the LogRecord was created.                                                                                                     |
| `created`         | Time when the LogRecord was created (time.Now().Unix() return value).                                                                            |
| `filename`        | Filename portion of pathname.                                                                                                                    |
| `funcname`        | Function name logged the record.                                                                                                                 |
| `levelname`       | Text logging level for the message ("DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL").                                                            |
| `levelno`         | Numeric logging level for the message (DEBUG, INFO, WARNING, ERROR, CRITICAL).                                                                   |
| `lineno`          | Source line number where the logging call was issued.                                                                                            |
| `message`         | The logging message.                                                                                                                             |
| `module`          | The module called log method.                                                                                                                    |
| `msecs`           | Millisecond portion of the creation time.                                                                                                        |
| `name`            | Name of the logger.                                                                                                                              |
| `pathname`        | Full pathname of the source file where the logging call was issued.                                                                              |
| `process`         | Process ID.                                                                                                                                      |
| `relativeCreated` | Time in milliseconds when the LogRecord was created, relative to the time the logging module was loaded (typically at application startup time). |

## Filter

`Filter` instances are used to perform arbitrary filtering of `LogRecord`.

A `Filter` struct needs to define `Format(LogRecord)` method, which return true
if it allows to log the `LogRecord`, and vice versa.

`Filter` can be used in both `Handler` and `Logger`.

# Benchmark

CPU: Intel(R) Xeon(R) Platinum 8272CL CPU @ 2.60GHz

| op name                | time per op |
| ---------------------- | ----------- |
| GetSameLogger          | 199ns       |
| GetRandomLogger        | 325ns       |
| GetSameHandler         | 5ns         |
| GetRandomHandler       | 30ns        |
| TextFormatter          | 632ns       |
| JSONFormatter          | 3032ns      |
| LogWithoutHandler      | 36ns        |
| EventLogWithoutHandler | 246ns       |
| LogWithOneHandler      | 3604ns      |
| LogWith100Handler      | 101268ns    |
| LogWithStream          | 7389ns      |
| LogWithFile            | 11372ns     |
| LogWithRotateFile      | 16234ns     |

# Example

See more examples [here](./example_test.go).

## Simple usage

```golang
var handler = xylog.GetHandler("")
handler.AddEmitter(xylog.StderrEmitter)
handler.SetFormatter(xylog.NewTextFormmater("%(level)s %(message)s"))

var logger = xylog.GetLogger("example.service")
logger.AddHandler(handler)
logger.SetLevel(xylog.DEBUG)

logger.Debug("foo")

// Output:
// DEBUG foo
```

## Rotating File Emitter

```golang
// Create a rotating emitter which rotates to another files if current file
// size is over than 30 bytes. Backup maximum of two log files.
var emitter = xylog.NewSizeRotatingFileEmitter("example.log", 30, 2)

var handler = xylog.GetHandler("")
handler.AddEmitter(emitter)
handler.SetFormatter(xylog.NewTextFormatter("%(message)s"))

var logger = xylog.GetLogger("example_file_emitter")
logger.SetLevel(xylog.DEBUG)
logger.AddHandler(handler)

for i := 0; i < 20; i++ {
    // logger will write 80 bytes (including newlines).
    logger.Debug("foo")
}

if _, err := os.Stat("example.log"); err == nil {
    fmt.Println("Created example.log")
}

if _, err := os.Stat("example.log.1"); err == nil {
    fmt.Println("Created example.log.1")
}

if _, err := os.Stat("example.log.2"); err == nil {
    fmt.Println("Created example.log.2")
}

// Output:
// Created example.log
// Created example.log.1
// Created example.log.2
```

## Get the existed Handler

```golang
var handler = xylog.GetHandler("bar")
handler.AddEmitter(xylog.StderrEmitter)
handler.SetFormatter(xylog.NewTextFormatter("%(message)s"))

var handler = xylog.GetHandler("bar")
var logger = xylog.GetLogger("example")
logger.AddHandler(handler)
logger.Critical("foo foo")

// Output:
// CRITICAL foo foo
```

## Formatters

```golang
var formatter = xylog.NewStructureFormatter().
    AddField("module", "name").
    AddField("level", "levelname").
    AddField("", "message")
var handler = xylog.GetHandler("")
handler.AddEmitter(xylog.StdoutEmitter)
handler.SetFormatter(formatter)

var logger = xylog.GetLogger("example.StructureFormatter")
logger.AddHandler(handler)
logger.SetLevel(xylog.DEBUG)
logger.Event("create").Field("employee", "david").Debug()

// Output:
// module=example.StructureFormatter level=DEBUG event=create employee=david
```

## Event Logger

```golang
var handler = xylog.GetHandler("")
handler.AddEmitter(xylog.StdoutEmitter)
handler.SetFormatter(xylog.NewTextFormatter(
    "module=%(name)s level=%(levelname)s %(message)s"))

var logger = xylog.GetLogger("example")
logger.AddHandler(handler)
logger.SetLevel(xylog.DEBUG)
logger.AddField("boss", "foo")

logger.Event("create").Field("product", 1235).Debug()
logger.Event("use").Field("product", "bar").JSON().Debug()

// Output:
// module=example level=DEBUG boss=foo event=create product=1235
// module=example level=DEBUG {"boss":"foo","event":"use","product":"bar"}
```

## Filter definition

```golang
// LoggerNameFilter only logs out records belongs to a specified logger.
type LoggerNameFilter struct {
    name string
}

func (f *LoggerNameFilter) Filter(r xylog.LogRecord) bool {
    return f.name == r.name
}

var handler = xylog.GetHandler("")
handler.AddEmitter(xylog.StderrEmitter)

var logger = xylog.GetLogger("example.filter")
logger.AddHandler(handler)
logger.SetLevel(xylog.DEBUG)
logger.AddFilter(&LoggerNameFilter{"exampel.filter.chat"})

xylog.GetLogger("example.filter.auth").Debug("auth foo")
xylog.GetLogger("example.filter.chat").Debug("chat foo")

// Output:
// chat foo
```

## Root logger

```golang
// A simple program with only one application area could use directly the root
// logger.
var handler = xylog.GetHandler("")
handler.AddEmitter(xylog.StdoutEmitter)

xylog.SetLevel(xylog.DEBUG)
xylog.AddHandler(handler)

xylog.Debug("bar")

// Output:
// bar
```

## Standard xybor logger

```golang
// If a logger is named with the prefix of "xybor.", it belongs to the xybor
// logger.
// This logger has the following properties:
//  - Logger and handler level are WARNING.
//  - Logging messages contain time, level name, logger name.
//  - Output is stderr.
var logger = xylog.GetLogger("xybor.foo")
logger.Debug("start")

// Output:
// time=[time] level=DEBUG module=xybor.foo message=start
```

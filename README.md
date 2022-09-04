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

1.  Loggers expose the interface that application code directly uses.
2.  Handlers convert log records (created by loggers) to log messages, then send
    them to the Emitter.
3.  Emitters write log messages to appropriate destination.
4.  Filters provide a finer grained facility for determining which log
    records to output.
5.  Formatters specify the layout of log records in the final output.

## Logger

Loggers should NEVER be instantiated directly, but always through the
module-level function `xylog.GetLogger(name)`. Multiple calls to `GetLogger()`
with the same name will always return a reference to the same Logger object.

You can logs a message by using one of built-in logging methods, such as:

```golang
func Log(level int, a ...any)
func Logf(level int, msg string, a ...any)

func Debug(a ...any)
func Debugf(msg string, a ...any)
```

To add `Handler` or `Filter` instances to the `logger`, call `AddHandler` or
`AddFilter` methods.

To adjust the level, using `SetLevel` method.

### EventLogger

`EventLogger` is a logger wrapper supporting to compose logging message by
key-value fields.

Use `Event` method of `Logger` to create a `EventLogger`. You need to create a
unique `EventLogger` each time you want to log.

If you call `Logger.AddField`, every `EventLogger` created by that `Logger` will
have the added fields without adding again.

`EventLogger` also provides JSON format for logging message (only the message
macro). You can use this format even if you do not set formatter as a
`JSONFormatter`.

## Logging level

The numeric values of logging levels are given in the following table. These are
primarily of interest if you want to define your own levels, and need them to
have specific values relative to the predefined levels. If you define a level
with the same numeric value, it overwrites the predefined value; the predefined
name is lost.

| Level        | Numeric value |
| ------------ | ------------- |
| CRITICAL     | 50            |
| ERROR/FATAL  | 40            |
| WARN/WARNING | 30            |
| INFO         | 20            |
| DEBUG        | 10            |
| NOTSET       | 0             |

## Handler

`Handler` handles logging events. A `Handler` need to be instantiated with an
`Emitter` instance rather than creating directly.

Any `Handler` with a not-empty name will be associated with its name. Calling
`NewHandler` twice with the same name will cause a panic. If you want to create
an anonymous `Handler`, call this function with an empty name.

To get an existed `Handler`, call `GetHandler` with its name. 

`Handler` can use `SetFormatter` method to format the logging message.

Like `Logger`, `Handler` is also able to call `AddFilter`.

## Emitter

`Emitter` instances write log messages to specified destination.

`StreamEmitter` can be used to print logging message into stdout or stderr.

`FileEmitter` can be used to write logging message to files. It can rotate to
log into another file if the file exceed the limit size or time.

## Formatter

`Formatter` instances are used to convert a `LogRecord` to text.

`Formatter` need to know how a `LogRecord` is constructed. They are responsible
for converting a `LogRecord` to a string which can be interpreted by either a
human or an external system.

`TextFormatter` is a simple `Formatter` which uses logging macros and format
string to format the message.

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
var handler = xylog.NewHandler("xybor", xylog.StdoutEmitter)
handler.SetFormatter(xylog.NewTextFormmater("%(level)s %(message)s"))

var logger = xylog.GetLogger("xybor.service")
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
var handler = xylog.NewHandler("", emitter)
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
// Get the handler of the first example.
var handler = xylog.GetHandler("xybor")
var logger = xylog.GetLogger("xybor.example")
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
var handler = xylog.NewHandler("", xylog.NewStreamEmitter(os.Stdout))
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
var handler = xylog.NewHandler("", xylog.NewStreamEmitter(os.Stdout))
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

// Get the logger of the first example.
var logger = xylog.GetLogger("xybor.service")
logger.AddFilter(&LoggerNameFilter{"xybor.service.chat"})

logger.Debug("foo")
xylog.GetLogger("xybor.service.auth").Debug("auth foo")
xylog.GetLogger("xybor.service.chat").Debug("chat foo")

// Output:
// chat foo
```

## Root logger

```golang
// A simple program with only one application area could use directly the root
// logger.
var handler = xylog.NewHandler("", xylog.StdoutEmitter)
xylog.SetLevel(xylog.DEBUG)
xylog.AddHandler(handler)

xylog.Debug("bar")

// Output:
// bar
```

## Xybor standard format

```golang
// All loggers need to be started with the prefix "xybor."
var logger = xylog.GetLogger("xybor.foo")

func main() {
    logger.Debug(1)
}

// Output:
// time=[time] source=example.go.main:13 level=DEBUG module=foo message=1
```

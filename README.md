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

Package xylog is designated for [structured logging](#structured-logging), [high performance](#benchmark), and
readable syntax.

The library is combined by
[python logging](https://docs.python.org/3/library/logging.html) design and
[zap](https://github.com/uber-go/zap) encoding approach.

# Quick start

`Logger` is directly used by application code. `Logger` names are dot-separated
hierarchical names, such as "a", "a.b", "a.b.c" or similar. For "a.b.c", its
parents are "a" and "a.b".

A `Logger` is obtained by `GetLogger` method. If the `Logger` with that name
hasn't existed before, the method will create a new one.

```golang
var logger = xylog.GetLogger("example")
defer logger.Flush()
```

`Handler` is responsible for creating the logging message. Like `Logger`,
`Handler` is also determined by its name, however, the name is not hierarchical.
Every call of `GetHandler` with the same name will give the same `Handler`.

_Exception: `Handlers` with the empty name are always different._

```golang
var handler = xylog.GetHandler("handler")
```

`Emitter` writes logging message to the specified output. Currently, only the
`StreamEmitter` is supported. You can use any \``in this`Emitter\`.
Writer

```golang
var emitter = xylog.NewStreamEmitter(os.Stdout)
```

When a logging method is called, the `Logger` creates a `LogRecord` and sends it
to underlying `Handlers`. `Handlers` converts `LogRecord` to text and sends it
to `Emitters`.

```golang
handler.AddEmitter(emitter)
logger.AddHandler(handler)
```

After preparing `Logger`, `Handler`, and `Emitter`, you can log the first
message.

```golang
logger.Debug("foo") // This message is blocked by Logger preferred level.
logger.Warning("bar")

// Output:
// message=bar
```

# Logging level

Both `Logger` and `Handler` has its own preferred level. If a logging level is
lower than the preferred one, the message will not be logged.

By default:

-   `Handler` logs all logging level.
-   `Logger`'s preferred level depends on its parents.

You can set the new preferred level of `Logger` and `Handler`.

```golang
logger.SetLevel(xylog.DEBUG)

logger.Debug("foo")
logger.Warning("bar")

// Output:
// message=foo
// message=bar
```

In the following example, however, the first message with DEBUG can bypass the
`Logger`, but will be prevented by `Handler`.

```golang
logger.SetLevel(xylog.DEBUG)
handler.SetLevel(xylog.INFO)

logger.Debug("foo")
logger.Warning("bar")

// Output:
// message=bar
```

The numeric values of logging levels are given in the following table. If you
define a level with the same numeric value, it overwrites the predefined value.

| Level        | Numeric value |
| ------------ | ------------- |
| CRITICAL     | 50            |
| ERROR/FATAL  | 40            |
| WARN/WARNING | 30            |
| INFO         | 20            |
| DEBUG        | 10            |
| NOTSET       | 0             |

# Structured logging

If the logging message has more than one field, `EventLogger` can help.

```golang
logger.Event("add-user").Field("name", "david").Field("email", "david@dad.com").Info()

// Output:
// event=add-user name=david email=david@dad.com
```

You also add a field to `Logger` or `Handler` permanently. All logging messages
will always include permanent fields.

```golang
logger.AddField("host", "localhost")
handler.AddField("port", 3333)

logger.Info("start server")

// Output:
// host=localhost port=3333 message="start server"
```

_NOTE: Fixed fields added to `Handler` will log faster than the oneadded to `Logger`_

`Handler` can support different encoding types. By default, it is
`TextEncoding`.

You can log the message with JSON format too.

```golang
import "github.com/xybor-x/xylog/encoding"

handler.SetEncoding(encoding.NewJSONEncoding())

logger.Warning("this is a message")
logger.Event("failed").Field("id", 1).Error()

// Output:
// {"message":"this is a message"}
// {"event":"failed","id": 1}
```

# Macros

You can log special fields whose values change every time you log. These fields
called macros.

Only the `Handler` can add macros.

```golang
handler.AddMacro("level", "levelname")

logger.Warning("this is a warning message")

// Output:
// level=WARNING message="this is a warning message"
```

The following table shows supported macros.

| MACRO             | DESCRIPTION                                                                                             |
| ----------------- | ------------------------------------------------------------------------------------------------------- |
| `asctime`         | Textual time when the LogRecord was created.                                                            |
| `created`         | Time when the LogRecord was created (time.Now().Unix() return value).                                   |
| `filename`\*      | Filename portion of pathname.                                                                           |
| `funcname`\*      | Function name logged the record.                                                                        |
| `levelname`       | Text logging level for the message ("DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL").                   |
| `levelno`         | Numeric logging level for the message (DEBUG, INFO, WARNING, ERROR, CRITICAL).                          |
| `lineno`\*        | Source line number where the logging call was issued.                                                   |
| `module`\*        | The module called log method.                                                                           |
| `msecs`           | Millisecond portion of the creation time.                                                               |
| `name`            | Name of the logger.                                                                                     |
| `pathname`        | Full pathname of the source file where the logging call was issued.                                     |
| `process`         | Process ID.                                                                                             |
| `relativeCreated` | Time in milliseconds between the time LogRecord was created and the time the logging module was loaded. |

_\* These are macros that are only available if `xylog.SetFindCaller` is called with `true`._

# Filter

`Filter` instances are used to perform arbitrary filtering of `LogRecord`.

A `Filter` struct needs to define `Format(LogRecord)` method, which returns true
if it allows to log the `LogRecord`, and vice versa.

`Filter` can be used in both `Handler` and `Logger`.

```golang
type NameFilter struct {
    name string
}

func (f *NameFilter) Filter(record xylog.LogRecord) bool {
    return f.name == record.Name
}

handler.AddFilter(&NameFilter{"example.user"})

var userLogger = xylog.GetLogger("example.user")
var serviceLogger = xylog.GetLogger("example.service")

userLogger.Warning("this is the user logger")
serviceLogger.Warning("this is the service logger")

// Output:
// message="this is the user logger"
```

# Hierarchical logger

As the first section mentioned, the `Logger`'s name is hierarchical. With this
feature, you can setup a common `Logger` with a specified configuration and uses
in different application zones.

```golang
// common/setup.go
func init() {
    var emitter = xylog.NewStreamEmitter(os.Stderr)
    var handler = xylog.GetHandler("")
    handler.AddEmitter(emitter)
    handler.SetEncoding(encoding.NewJSONEncoding())
    handler.AddMacro("time", "asctime")
    handler.AddMacro("level", "levelname")

    var logger = xylog.GetLogger("parent")
    logger.AddHandler(handler)
    logger.SetLevel(xylog.WARNING)
}
```

```golang
// user/foo.go
import _ "common"

var logger = xylog.GetLogger("parent.user")
defer logger.Flush()
logger.SetLevel(xylog.INFO)
logger.AddField("module", "user")

logger.Info("this is user module")
logger.Debug("this is a not logged message")

// Output:
// time=[time] level=INFO module=user message="this is user module"
```

```golang
// service/bar.go
import _ "common"

var logger = xylog.GetLogger("parent.service")
defer logger.Flush()
logger.AddField("module", "service")
logger.AddField("service", "bar")

logger.Warning("this is service module")

// Output:
// time=[time] level=INFO module=service service=bar message="this is service module"
```

# Benchmark

CPU: AMD Ryzen 7 5800H (3.2Ghz)

These following benchmarks are measured with all [macros](#macros) and
`SetFindCaller` is called with `false`.

| op             | time per op | alloc per op |
| -------------- | ----------: | -----------: |
| Disable        |        51ns |  0 allocs/op |
| WithoutHandler |       247ns |  3 allocs/op |
| TextEncoding   |       653ns | 14 allocs/op |
| JSONEncoding   |       631ns | 14 allocs/op |

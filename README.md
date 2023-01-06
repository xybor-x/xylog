[![xybor founder](https://img.shields.io/badge/xybor-huykingsofm-red)](https://github.com/huykingsofm)
[![Go Reference](https://pkg.go.dev/badge/github.com/xybor-x/xylog.svg)](https://pkg.go.dev/github.com/xybor-x/xylog)
[![GitHub Repo stars](https://img.shields.io/github/stars/xybor-x/xylog?color=yellow)](https://github.com/xybor-x/xylog)
[![GitHub top language](https://img.shields.io/github/languages/top/xybor-x/xylog?color=lightblue)](https://go.dev/)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/xybor-x/xylog)](https://go.dev/blog/go1.18)
[![GitHub release (release name instead of tag name)](https://img.shields.io/github/v/release/xybor-x/xylog?include_prereleases)](https://github.com/xybor-x/xylog/releases/latest)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/e9146e2ccfd745a48a42fc723918d4cb)](https://www.codacy.com/gh/xybor-x/xylog/dashboard?utm_source=github.com&utm_medium=referral&utm_content=xybor-x/xylog&utm_campaign=Badge_Grade)
[![Codacy Badge](https://app.codacy.com/project/badge/Coverage/e9146e2ccfd745a48a42fc723918d4cb)](https://www.codacy.com/gh/xybor-x/xylog/dashboard?utm_source=github.com&utm_medium=referral&utm_content=xybor-x/xylog&utm_campaign=Badge_Coverage)
[![Go Report](https://goreportcard.com/badge/github.com/xybor-x/xylog)](https://goreportcard.com/report/github.com/xybor-x/xylog)

# Introduction

Package xylog is designed for [leveled](#logging-level) and
[structured](#structured-logging) logging, [dynamic fields](#macros),
[high performance](#benchmark), [zone management](#hierarchical-logger), simple
configuration, and readable syntax.

The library is combined by
[python logging](https://docs.python.org/3/library/logging.html) design and
[zap](https://github.com/uber-go/zap) encoding approach.

# Quick start

You can easily configure a logger with `SimpleConfig`.

There are some fields you can modify with this way (note that all fields are
optional):

-   `Name` is the name of Logger. It can be used later with `GetLogger`
    function. Default to an empty name (the root logger).

-   `Encoding` to format the output. Default to `TextEncoding`.

-   `Filename` specifies that `Logger` will write the output to a file. Do NOT
    use together with `Writer`.

-   `Filemode` specifies the mode to open file. Default to `APPEND` \| `CREATE`
    \| `WRONLY`.

-   `Fileperm` specifies the permission when creating the file. Default to 0666.

-   `Level` specifies the logging level. Default to `WARNING`.

-   `TimeLayout` when format the time string. Default to `RFC3339Nano`.

-   `Writer` specifies that Logger will write the output to a file. Do NOT use
    together with `Filename`.

```golang
var config = &xylog.SimpleConfig{
    Name:   "simple-logger",
    Level:  xylog.DEBUG,
    Writer: os.Stdout,
}

var logger, err = config.AddMacro("level", "levelname").Apply()
if err != nil {
    fmt.Println("An error occurred:", err)
    os.Exit(1)
}
defer xylog.Flush()

logger.Debug("logging message")
logger.Event("create-user").Field("username", "foo").
    Field("email", "bar@buzz.com").Field("Age", 25).Info()

// Output:
// level=DEBUG messsage="logging message"
// level=INFO event=create-user username=foo email=bar@buzz.com Age=25
```

# Full configuration

`Logger` is directly used by application code. `Logger` names are dot-separated
hierarchical names, such as "a", "a.b", "a.b.c" or similar. For "a.b.c", its
parents are "a" and "a.b".

A `Logger` is obtained using `GetLogger` method. If the `Logger` with that name
didn't exist before, the method will create a new one. The `Logger` with empty
name is the root one.

```golang
var logger = xylog.GetLogger("example")
defer xylog.Flush()
```

`Handler` is responsible for generating logging messages. Like `Logger`,
`Handler` is also identified by its name, however, the name is not hierarchical.
Every `GetHandler` call with the same name gives the same `Handler`.

_Exception: `Handlers` with the empty names are always different._

```golang
var handler = xylog.GetHandler("handler")
```

`Emitter` writes logging messages to the specified output. Currently, only
`StreamEmitter` is supported. You can use any `Writer` in this `Emitter` type.

```golang
var emitter = xylog.NewStreamEmitter(os.Stdout)
```

When a logging method is called, the `Logger` creates a `LogRecord` and sends it
to underlying `Handlers`. `Handlers` convert `LogRecord` to text and send it
to `Emitters`.

A `Logger` can have multiple `Handlers`, and a `Handler` can have multiple
`Emitters`.

```golang
handler.AddEmitter(emitter)
logger.AddHandler(handler)
```

After preparing `Logger`, `Handler`, and `Emitter`, you can log the first
messages.

```golang
logger.Debug("foo") // This message is blocked by Logger's preferred level.
logger.Warning("bar")

// Output:
// message=bar
```

# Logging level

Both `Logger` and `Handler` has its own preferred level. If a logging level is
lower than the preferred one, the message will not be logged.

By default:

-   `Handler`'s preferred level is `NOTSET` (it logs all logging levels).
-   `Logger`'s preferred level depends on its parents. When a `Logger` is newly
    created, its preferred level is the nearest parent's one. The root logger's
    preferred level is `WARNING`.

You can set a new preferred level for both `Logger` and `Handler`.

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

_NOTE: Fixed fields added to `Handler` will log faster than the one added to `Logger`_

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

`Filter` can be used by `Handlers` and `Loggers` for more sophisticated
filtering than is provided by levels.

A `Filter` instance needs to define `Filter(LogRecord)` method, which returns
`true` if it allows logging the `LogRecord`, and vice versa.

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
defer xylog.Flush()
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
defer xylog.Flush()
logger.AddField("module", "service")
logger.AddField("service", "bar")

logger.Warning("this is service module")

// Output:
// time=[time] level=INFO module=service service=bar message="this is service module"
```

# Benchmark

CPU: AMD Ryzen 7 5800H (3.2Ghz)

These benchmark of xylog are measured with `SetFindCaller` is called with
`false`.

_NOTE: The benchmarks are run on a different CPU from the [origin](https://github.com/uber-go/zap/blob/master/README.md#performance), so the benchmark values may be different too._

Log a message and 10 fields:

| Package             |        Time | Time % to zap | Objects Allocated |
| :------------------ | ----------: | ------------: | ----------------: |
| :zap: zap           |  1707 ns/op |           +0% |       5 allocs/op |
| :zap: zap (sugared) |  2043 ns/op |          +20% |      10 allocs/op |
| zerolog             |   884 ns/op |          -48% |       1 allocs/op |
| go-kit              |  6255 ns/op |         +266% |      58 allocs/op |
| logrus              |  8384 ns/op |         +391% |      80 allocs/op |
| apex/log            | 22707 ns/op |        +1230% |      65 allocs/op |
| log15               | 25461 ns/op |        +1391% |      75 allocs/op |
| :rocket: xylog      |  3518 ns/op |         +106% |      77 allocs/op |

Log a message with a logger that already has 10 fields of context:

| Package             |        Time | Time % to zap | Objects Allocated |
| :------------------ | ----------: | ------------: | ----------------: |
| :zap: zap           |   140 ns/op |           +0% |       0 allocs/op |
| :zap: zap (sugared) |   181 ns/op |          +29% |       1 allocs/op |
| zerolog             |    89 ns/op |          -36% |       0 allocs/op |
| go-kit              |  5963 ns/op |        +4159% |      57 allocs/op |
| logrus              |  6590 ns/op |        +4607% |      69 allocs/op |
| apex/log            | 21777 ns/op |       +15455% |      54 allocs/op |
| log15               | 15124 ns/op |       +10702% |      71 allocs/op |
| :rocket: xylog      |   416 ns/op |         +197% |       6 allocs/op |

Log a static string, without any context or `printf`-style templating:

| Package             |       Time | Time % to zap | Objects Allocated |
| :------------------ | ---------: | ------------: | ----------------: |
| :zap: zap           |  154 ns/op |           +0% |       0 allocs/op |
| :zap: zap (sugared) |  195 ns/op |          +27% |       1 allocs/op |
| zerolog             |   87 ns/op |          -44% |       0 allocs/op |
| go-kit              |  382 ns/op |         +148% |      10 allocs/op |
| logrus              | 1008 ns/op |         +554% |      24 allocs/op |
| apex/log            | 1744 ns/op |        +1032% |       6 allocs/op |
| log15               | 4246 ns/op |        +2657% |      21 allocs/op |
| :rocket: xylog      |  447 ns/op |         +190% |       6 allocs/op |

# v0.4.1 (Jan 6, 2023)

-   Change name of NewStreamEmitter to NewBufferEmitter.
-   Change name of NewDefaultEmitter to NewStreamEmitter.
-   Remove default logger of xybor.

# v0.4.0 (Jan 6, 2023)

-   Allow determine the buffer size in emitter.
-   Support simple configuration.
-   Increase at least 50% performance by using EventLogger pool.
-   Modify benchmark: use JSONEncoding instead of TextEncoding; use io.Discard
    instead of devnull.

# v0.3.0 (Sep 17, 2022)

-   Fix misspelling.
-   Remove Logger.Flush method. Flush now is a global function which flushes all
    emitters.

# v0.2.0 (Sep 17, 2022)

-   Handler now has many Emitters (instead only one).
-   Remove root logger methods.
-   Remove Rotating Emitters.
-   Add more get methods to Logger and Handler.
-   Public some functions, methods.
-   Use lock more exactly.
-   Refactor unittests.
-   Improve performance by using buffer, pool, and reducing unnecessary works.
-   Remove the old TextFormatter. StructuredFormatter is renamed to
    TextFormatter.
-   EventLogger is adapted with the formatter. EventLogger.JSON is removed.
-   Remove `message` macro.
-   Integrate Formatter to Handler and rename it to Encoding.
-   Use encoding approach of zap.
-   Add benchmark comparison.

# v0.1.0 (Sep 04, 2022)

-   Add JSONFormatter and StructureFormatter.
-   Change extra, message, and log design.
-   It is possible to log stack trace.

# v0.0.1 (Sep 02, 2022)

-   Migrated from the [origin project](https://github.com/xybor/xyplatform).

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

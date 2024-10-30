## Abstract 
Project started by Arpit Bhayani. 
It is an in-memory realtime database with Redis like commands, and SQL Like reactivity.

## Index
- [Day 1](#day-1-3rd-october-2024)
- [Day 2](#day-2-5th-october-2024)

## Day 1 (3rd October 2024):
`main.go`: 
it is the starting point of an go application.
I can see an init() function, which implies, it will be called even before the main function, where we are configuring the flag variables for our command line.


We are using built-in `flag` library which golang provides for configuring CLI flags for our application.

Some of the functions are:
`flag.StringVar(), flag.IntVar(), flag.BoolVar()` for configuring diff. type of flags + initial config

`Config` Struct: it uses `mapstructure` for directly unmarshalling of configuration data, which is in JSON, YAML format to struct.

It is in-memory database which uses Append only file for persistence, and user have to authenticate with username and password to register to the dice client.

3 main sections: 
1. server

diff. fields:   
- `Addr`: server IP address
- `Port`: dice server port (which is 7063 by default)
- `KeepAlive`: keep alive duration, for how long connection should stay alive
- `Timeout`: connection timeout
- `MaxConn`: max. no. of connections to the server
- `ShardCronFrequency`: frequency of shard cron jobs
- `MultiplexerPollTimeout`: timeout for multiplexer polling
- `MaxClients`: max. no. of clients
- `MaxMemory`: max. memory usage
- `EvictionPolicy`: policy for evicting data when memory limit is reached
- `EvictionRatio`: ratio at which starting evicting data
- `KeysLimit`: max. no. of keys
- `AOF`: append only file for persistence (Write ahead log)
- `PersistenceEnabled`: whether persistence enabled
- `WriteAOFOnCleanup`: Whether to write AOF on cleanup
- `LFULogFactor`: Log Factor for LFU (Least frequently used)
- `LogLevel`: Logging Level (like error, success, warning)
- `PrintPrettyLogs`: Whether to print pretty logs
- `EnableMultiThreading`: Whether to enable multithreading
- `StopMapInitialize`: Initial size of store map

2. auth

- `Username` and `Password` for authentication

3. network

- `IOBufferLength`: length of I/O buffer
- `IOBufferLengthMAX`: max. length of I/O buffer


## Day 2 (5th October, 2024): 
Understanding dicedb internals: logger functionality

- We have diff. levels of log: 
1. debug
2. warn
3. info 
4. error

Here, we are using `zerolog` library for creating new logger instance, with timestamps.

Creating new instance of zerolog:

1. import zerolog library: 
```bash
go get -u github.com/rs/zerolog
```

> Side Note: -u flag not only import the library; but also updates the dependencies, the library depend upon, to latest minor/patch version.

2. attach IO writer to it, by creating new console writer

```golang
var writer io.Writer = zerolog.ConsoleWriter({Out: os.Stderr });
```

3. Creating New zerolog logger instance


```golang
logger:= zerolog.New(writer).Level(level); // level can be "info" | "error" | "debug" | "warn"
```

4. Additionally, we can attach timestamps to each log that gets created using `TimeStamp()` method

```golang
logger:= logger.With().TimeStamp().Logger();
```

Finally, what they are doing is, they are attaching zerolog to slog instance (WTF: Need to deep dive on this, maybe it's similar to morgan and winston setup in node.js)

```golang

finalLogger:= slog.New(newZeroLogHandler(&logger));

struct ZeroLogHandler  {
    logger *zerolog.Logger
}

func newZeroLogHandler (logger* zerolog.Logger) *ZeroLogHandler {
return &ZeroLogHandler {
    logger: logger;
}
}
```

Conclusion: 
- zerolog provides pretty logs, logs with timestamps. 
- But, it does not handles levels exclusively and need extra layer of `slog` which maintains diff. types of log levels, and subsequently we can print those pretty logs with timestamps to the console.

Ques: Why use both `slog` and `zerolog`?

Ans. slog (Go's standard logging package)

Advantages:

- Part of the Go standard library (as of Go 1.21)
- Provides a standard interface for logging that can be used across different libraries and projects
- Structured logging support
- Leveled logging
- Contextual logging with WithGroup and WithAttrs

Disadvantages:

- Relatively new, so it may lack some advanced features or optimizations
- May not be as performant as some third-party logging libraries

zerolog
Advantages:

- High-performance, zero-allocation JSON logger
    - Zerolog is designed to be extremely fast and efficient.
    - It uses a zero-allocation approach, meaning it minimizes memory allocations during logging operations.
    - This results in very low overhead, making it suitable for high-performance applications.
- Structured logging with a fluent API
    - Zerolog provides a fluent (chainable) API for adding structured data to log entries.
    ```golang
    log.Info().
    Str("event", "user_login").
    Int("user_id", 123).
    Msg("User logged in successfully")

    ```
    - This creates easily parseable JSON log entries.
- Highly customizable
    - Zerolog allows extensive customization of log output.
    - You can customize timestamp formats, level names, and add custom fields globally.
    - It supports custom formatters and writers for complete control over log output.

Built-in pretty printing for development
Contextual logging

Disadvantages:

- Not part of the standard library, requiring an external dependency
- May have a steeper learning curve compared to simpler logging libraries

Benefits of the Integration

- Standard Interface: By using slog as the primary interface, the code remains compatible with the Go standard library and other libraries that use slog.
- Performance: The integration allows leveraging zerolog's high-performance logging backend while using slog's interface.
- Flexibility: It's easier to switch between different logging backends (e.g., zerolog, zap) while maintaining the same slog interface in the application code.
- Feature Combination: This approach combines slog's standard interface and structured logging capabilities with zerolog's performance and additional features like pretty printing.
- Future-Proofing: As slog evolves, the application can easily adapt to new features while still benefiting from zerolog's optimizations.

Potential Drawbacks

- Complexity: Integrating two logging systems adds some complexity to the codebase.
- Overhead: There might be a small performance overhead due to the translation layer between slog and zerolog.
- Maintenance: Keeping the integration up-to-date with both slog and zerolog updates might require additional effort.

## Abstract 
Project started by Arpit Bhayani. 
It is an in-memory realtime database with Redis like commands, and SQL Like reactivity.

## Index
- [Main Config](#main-configuration)
- [DiceDB Internals: Logger](#dicedb-internals-logger)
- [WebSocket Server Implementation](#websocket-server-implementation)
- [Implementation of Append Only FIle](#implementation-of-append-only-file-aof-for-disk-persistence-in-memory-based-store)
- [Swiss Table implementation](#swiss-table-implementation)

## Main Configuration
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


## DiceDB Internals: Logger: 
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

#### Ques: Why use both `slog` and `zerolog`?

Ans. 
##### slog (Go's standard logging package)

##### Advantages:

- Part of the Go standard library (as of Go 1.21)
- Provides a standard interface for logging that can be used across different libraries and projects
- Structured logging support
- Leveled logging
- Contextual logging with WithGroup and WithAttrs

##### Disadvantages:

- Relatively new, so it may lack some advanced features or optimizations
- May not be as performant as some third-party logging libraries

##### zerolog
##### Advantages:

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

- Built-in pretty printing for development
- Contextual logging

##### Disadvantages:

- Not part of the standard library, requiring an external dependency
- May have a steeper learning curve compared to simpler logging libraries

##### Benefits of the Integration

- **Standard Interface**: By using slog as the primary interface, the code remains compatible with the Go standard library and other libraries that use slog.
- **Performance**: The integration allows leveraging zerolog's high-performance logging backend while using slog's interface.
- **Flexibility**: It's easier to switch between different logging backends (e.g., zerolog, zap) while maintaining the same slog interface in the application code.
- **Feature Combination**: This approach combines slog's standard interface and structured logging capabilities with zerolog's performance and additional features like pretty printing.
- **Future-Proofing**: As slog evolves, the application can easily adapt to new features while still benefiting from zerolog's optimizations.

Potential Drawbacks

- **Complexity**: Integrating two logging systems adds some complexity to the codebase.
- **Overhead**: There might be a small performance overhead due to the translation layer between slog and zerolog.
- **Maintenance**: Keeping the integration up-to-date with both slog and zerolog updates might require additional effort.


## Websocket server implementation

- Use of `gorilla/websocket` library
```golang


type WsServer struct {
    upgrader websocket.Upgrader
    server *http.Server
}

func NewWsServer() *WsServer {
    // create new HTTP request multiplexer
    mux := http.NewServeMux();
    srv := &http.Server{
        Addr : fmt.Sprintf(":%d", 3000)
        Handler : mux
    }

    upgrader = websocket.Upgrader{
        //CORS configuration: allowing all requests
        CheckOrigin: func (r* http.Request) bool {
            return true;
        }
    }

    wsServer := &WsServer {
        server : srv,
        upgrader: upgrader,
    }

    mux.HandleFunc("/", WsHandler);
    mux.HandleFunc("/health", func (w http.ResponseWriter, r* http.Request) {
       if _, err :=  w.Write([]byte("OK")); err != nil {
        return;
       }
    })
    return wsServer;
}


func (s* WsServer) Run(ctx context.Context) error {
    var wg sync.WaitGroup
    var err error
    wsCtx,cancelWebsocket = context.WithCancel(ctx)
    // cancel the derived context when server stops running: when error occured
    defer cancelWebsocket();

    // HANDLING GRACEFUL SHUTDOWN OF WEBSOCKET CONNECTION
    wg.Add(1);
    go func() {
        // decrement wait group counter when goroutine finish running
        defer wg.Done();
        // parent context cancelled
        <-ctx.Done()
        // so, shutdown ws server
       err =  s.server.Shutdown(wsCtx)
    }()
    
    wg.Add(1);
    go func() {
        defer wg.Done()
        err = s.server.ListenAndServe()
    }();

    wg.Wait()
    return err
}

// Websocket handler
func (s* WsServer) WsHandler(w http.ResponseWriter, r* http.Request) {
    // update http to websocket
   conn, err:= s.upgrader.Upgrade(w,r,nil)
   if err != nil {
    return;
   }

// close ws connection
   defer func() {
    conn.WriteMessage(websocket.CloseMessage)
    conn.Close()
   }()

   for {
    // read from connection
    _, message, err := conn.ReadMessage()
    if err != nil {
        break
    }

// echo back the text message to client
    if err:= conn.WriteMessage(websocket.TextMessage,message); err != nil {
        break
    }
   }
}

```


#### Handling Graceful shutdown of websocket server
> GRACEFUL SHUTDOWN: tries to complete all the in-flight requests before shutting down.
1. When the parent context is cancelled, we shutdown the websocket server, waiting for it in a separate goroutine, which gets created when we run ws server.

```golang

go func(){
    defer wg.Done(); // decrement the wait group counter
<- ctx.Done() // blocking the gorouting till parent context is not cancelled
server.Shutdown(wsCtx); // shutdown the server
}()

```

#### Parent-child context relationship

Context represents lifecycle of a particular operation or entire application


Question arises as why we are shutting down websocket server on the basis of child context and not the parent context?

Here's the reasoning behind this choice:

1. **Parent Context Lifespan**: The parent ctx context represents the overall lifetime of the application or a specific operation. It's possible that this context may be canceled for reasons unrelated to the WebSocket server, such as a timeout or a signal received at a higher level.
2. **WebSocket Server Lifecycle**: The WebSocket server has its own lifecycle, which should be managed independently from the parent context. If the parent context is canceled, it doesn't necessarily mean that the WebSocket server should be immediately shut down.
3. **Shutdown Control**: By using the websocketCtx context, you can control the timeout and cancellation of the WebSocket server's shutdown process separately from the parent context. This allows you to ensure a graceful shutdown of the WebSocket server, even if the parent context is canceled prematurely.
4. **Error Handling**: If an error occurs during the shutdown process, you can capture and handle it more easily when you're using a dedicated context for the shutdown, rather than relying on the parent context.


## Implementation of Append Only FIle (AOF) for disk persistence in memory based store.

For persisting the data to disk, we flush all the write/update operations to disk to an Append only file.

```golang
    type AOF struct {
        file *os.file
        writer *bufio.Writer
        mutex sync.Mutex 
        path string // file path
    }

    func NewAOF(path string) (*AOF, error) {
        // open file in append mode, write only mode, and if not created, create the file with user permission set(6: owner(rw-)4: group(r--)4: others(r--))
       f, err:=  os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, fs.FileMode(0644))
       if err != nil{
        return nil, err
       }


        return &AOF {
            file: f,
            writer: bufio.NewWriter(f),
            path: path,
        }
    }


```

#### How to efficiently write to disk
We need to flush the commands to disk to maintaining durability in the database, but write to disk is very slow operation, via write() system call to the given file descriptor.

We can solve this issue by:
i. writing data to memory buffer
ii. flush data from memory to OS buffer/cache
iii. flush data from OS buffer/cache to file for disk persistance, by synchronizing file with the cache data.

```golang
// write operations to file: atomic operation via mutex
    func (a *AOF) Write(data string) error {
        a.mutex.Lock()
        defer a.mutex.Unlock()

// write data to memory buffer
        if _, err := a.writer.WriteString(data + "\n"); err != nil {
            return err
        }
        // flush data to OS buffer/cache
        if err := a.writer.Flush(); err != nil {
            return err
        }
        // CALLS fsync() system call to flush OS cache to disk
        return a.file.Sync()
    }
```

#### Efficiently Closing the file 

- Atomic operation via mutex
- flush the remaining buffer to OS buffer
- close the file
```golang

func (a* AOF) Close() error {
    // acquiring lock
    a.mutex.Lock()
    defer a.mutex.Unlock()
    // CRITICAL SECTION
    // flush the memory buffer to OS cache
    if err := a.writer.Flush(); err != nil {
        return err
    }
// close the file via it's file descriptor
    return a.file.Close()
}
```


## Swiss Table Implementation

### Abstract
Swiss table is a hash table implementation which is highly efficient in lookups and insertion operations.
It stores it's elements in contiguous manner, and does not do chaining of elements in case of collision in the insertion operations.
It is the variation of open addressing hash table.

Key features of Swiss Table:
1. **Quadratic Probing**: in case of collisions in the hash values, swiss tables uses quadratic probing.

index = (hash(key) + i^2)%array_size; where i is the probe number starting from 0.

2. **Tombstone Marker**:it uses tombstone markers to make the elements as deleted. Instead of being left empty, the places are marked as tombstones. 

#### Working of Tombstones with diff. operations: 
1. **Lookups**: When during lookups, if tombstones are found, it will continue probing. The logic behind tombstone is that, if the element is removed and if that element caused other previous hashes to be probed due to collisions then, we need to mark is as pseudo-filled, so that if there is lookup for the probed element then it could be possible. If we do not add tombstone here, the traversal will stop immediately as soon as it find empty location, causing lookup failure.


WITH TOMBSTONES: 
```
Initial state (size=7):
[A, B, C, _, _, _, _]
hash(A) = 0
hash(B) = 0 (collided, probed to 1)
hash(C) = 1 (collided, probed to 2)

After deleting B:
[A, T, C, _, _, _, _]  (T = tombstone)

Now when looking up C:
1. hash(C) = 1
2. Position 1 is marked as tombstone
3. Continue probing to find C at position 2 
```


WITHOUT TOMBSTONES: 

```
[A, _, C, _, _, _, _]

Looking up C:
1. hash(C) = 1
2. Position 1 is empty
3. Stop searching (WRONG - we would never find C!)
```

2. **Insertion**: During insertion operation, tombstone will be replaced with new elements.
3 **Deletion**:  We use tombstones instead of empty markers to maintain probe chains.

So tombstones aren't about blocking positions - they're about maintaining the "probe history" so we can still find elements that probed past deleted entries. New insertions are free to reuse tombstone positions.

---

DiceDB uses cockroachdb's swiss table implementation in golang.

In this, the hash function returns a 64 bit hash value. It consists of 2 parts:
1. H1, 57 bit hash value to identify the element index within the table itself.
2. H2, 7 bit used to store metadata of this element: empty, deleted or full.

```golang
import "github.com/cockroachdb/swiss";

// Key is of comparable type: like Struct* not allowed
type SwissTable[K comparable, V any] struct {
	M *swiss.Map[K, V]
}

func (t *SwissTable[K, V]) Put(key K, value V) {
	t.M.Put(key, value)
}

func (t *SwissTable[K, V]) Get(key K) (V, bool) {
	return t.M.Get(key)
}

func (t *SwissTable[K, V]) Delete(key K) {
	t.M.Delete(key)
}

func (t *SwissTable[K, V]) Len() int {
	return t.M.Len()
}

func (t *SwissTable[K, V]) All(f func(k K, obj V) bool) {
	t.M.All(f)
}
```

- It creates a generic swiss table wrapper for standard function implementations.



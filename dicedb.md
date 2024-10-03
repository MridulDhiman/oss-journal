Project started by Arpit Bhayani. 
It is an in-memory realtime database with Redis like commands, and SQL Like reactivity.

## Day 1 (3rd October 2024):
`main.go`: 
it is the starting point of an go application.
I can see an init() function, which implies, it will be called even before the main function, where we are configuring the flag variables for our command line.


We are using built-in `flag` library which golang provides for configuring CLI flags for our application.

Some of the functions are:
`flag.StringVar(), flag.IntVar(), flag.BoolVar()` for configuring diff. type of flags + initial config

`Config` Struct: it uses `mapstructure` for directly unmarshalling of configuration data, which is in JSON, YAML format to struct.
3 main sections: 
1. server

diff. fields:
- Addr: server IP address
- Port: dice server port (which is 7063 by default)
- KeepAlive: keep alive duration, for how long connection should stay alive
- Timeout: connection timeout
- MaxConn: max. no. of connections to the server
- ShardCronFrequency: frequency of shard cron jobs
- MultiplexerPollTimeout: timeout for multiplexer polling
- MaxClients: max. no. of clients
- MaxMemory: max. memory usage
- EvictionPolicy: policy for evicting data when memory limit is reached
- EvictionRatio: ratio at which starting evicting data
- KeysLimit: max. no. of keys
- AOF file: append only file for persistence (Write ahead log)
- PersistenceEnabled: whether persistence enabled
- WriteAOFOnCleanup: Whether to write AOF on cleanup
- LFULogFactor: Log Factor for LFU (Least frequently used)
- LogLevel: Logging Level (like error, success, warning)
- PrintPrettyLogs: Whether to print pretty logs
- EnableMultiThreading: Whether to enable multithreading
- StopMapInitialize: Initial size of store map

2. auth

- Username and Password for authentication

3. network

- IOBufferLength: length of I/O buffer
- IOBufferLengthMAX: max. length of I/O buffer
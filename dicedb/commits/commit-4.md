## Abstract

EPOLL is a linux kernel system call for a scalable I/O event notification mechanism.
It is used to monitor multiple file descriptors to see if I/O is possible on any of them. It is used to replace older POSIX system calls (select(2) and poll(2)), by providing better performance in case if there are large no. of watched file descriptors. Older system calls work in O(N) time, while epoll works in O(1) time. It is basically I/O multiplexer, where all the events are registered and it keeps watching/polling for file descriptors, in case a particular event occurs, to notify the system regarding that, so that proper action can be taken, via default or custom handlers.

I/O bound tasks include: reading from file, waiting on network for new request, performing DNS lookup etc., where the `epoll` will keep watching.

## Commit(d3da078ec2b7e5802bbf901caf58e1f7489d5fcf): Supporting concurrent clients via EPOLL

In this commit, we implement an asynchronous TCP server (which does block the main thread), using Linux's epoll mechanism to handle multiple clients asychronously. It can handle 20,000 concurrent clients using event driven I/O instead of creating new thread per connection.

Here, we will be monitoring network socket file descriptor for data from the redis client via TCP connection with RESP serialization.

### server/async_tcp.go: package server

#### RunAsyncTCPServer()

- create slice of epoll events for handling 20,000 max. concurrent clients
- create a socket using configurations like:
    - address family: IPv4(AF_INET)(in our case) and other options are: IPv6(AF_INET6), unix based socket (AF_UNIX)
    - socket Type: we will creating a TCP socket(SOCK_STREAM)
    - socket mode: non blocking(O_NONBLOCK)/blocking: essential for asynchronous I/O.
    - protocol no.: 0 for TCP protocol
- so, we will create a TCP socket bound to IPv4 address family in non-blocking mode, and return socket's file descriptor.
- Now, we will bind socket(using it's file descriptor) to dice server's IP and port.
- start listening for incoming connections. 
- once, socket file descriptor start listening for incoming connections (max. upto max_client), we need to register the socket file descriptor to be monitored for EPOLLIN(read) events. Whenever the redis client sends new command, it gets stored in kernel network buffer, but our dice server don't know about it. So, as we registered EPOLLIN event for our TCP server's socket file descriptor, we would be notified by EPOLL, in case of new command, so that we can parse the payload and process the request.
- wait for EPOLLIN events to epoll file descriptor using `epoll_wait()` in a loop.
- We will verify whether the current event is for socket server, if yes then we will accept the incoming connection from server and create new epoll instance for listening for new data on this TCP connection, by registering socket client for epoll monitoring. 
- For writing and reading to/from TCP connection, we use write() and read() system calls on file descriptors, and for that use io.ReadWriter interface instead of net.Conn that we were previously using in case of Synchronous TCP Server.
- in case of error while reading from connection, close the connection by stop listening to socket client, decrement the no. of connected clients.







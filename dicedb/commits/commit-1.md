### 1st commit(dd524f174bd6de83ea8a27526f8f1d50436d3b00): setup synchronous TCP server and CLI flags

* Setup `host` and `port` flags, which default values of `0.0.0.0` (allowing all network interfaces) and `7379`
* start listening on socket (host + port), it would be an blocking call(blocking, as in blocking the main thread, listening for connection to be initiated by the client and server accepting it via `conn.Accept()`) and establish connection
* once connection is established, we will keep listening for new message to be sent from the client, and we would store that inside bytes buffer of 512 bytes. 
* If there is no error, we will echo back the message to client, acknowledging that we received the message successfully.
* else if there is some sort of error in case, then that means user have been disconnected intentionally/unintentionally and we have to close the connections and decrease the no. of concurrent clients by 1.
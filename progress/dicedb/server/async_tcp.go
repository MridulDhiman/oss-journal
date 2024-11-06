//go:build linux
// +build linux

// async_tcp.go
package server

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/MridulDhiman/dice/config"
	"github.com/MridulDhiman/dice/core"
)

var conn_clients int = 0;
var cronFrequency = 1 * time.Second
var lastCronExecutionTime = time.Now()

func RunAsyncTCPServer() error {
	log.Println("Starting asynchronous TCP server in: ", config.Host, config.Port)
	// max. no. of concurrent clients
	max_clients := 20000

	// use Linux based EPOLL to handle concurrent clients, creating EPOLL events, each of which will be monitoring network socket, waiting for new command to be read from the client
	var events []syscall.EpollEvent = make([]syscall.EpollEvent, max_clients)

	// create new TCP socket in non-blocking mode
	socketFD, err := syscall.Socket(syscall.AF_INET, syscall.O_NONBLOCK|syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}

	// cleanup the socket file descriptor after server being terminated...
	defer syscall.Close(socketFD)

	// set the file descriptor to be non blocking mode
	if err := syscall.SetNonblock(socketFD, true); err != nil {
		return err
	}

	// parse IP from Host string
	ipv4 := net.ParseIP(config.Host)
	if err := syscall.Bind(socketFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ipv4[0], ipv4[1], ipv4[2], ipv4[3]},
	}); err != nil {
		log.Println("Could not bind IP and port to socket")
		return err
	}

	// start listening to incoming connections
	if err := syscall.Listen(socketFD, max_clients); err != nil {
		log.Println("could not listen to incoming connections")
		return err
	}



	// create new EPOLL instance
	epollFD, err := syscall.EpollCreate1(0);


	// configure socket to be monitored/polled for the read events
	var socketServerEvent syscall.EpollEvent = syscall.EpollEvent{
		Fd:     int32(socketFD),
		Events: syscall.EPOLLIN,
	}

	// modify the epoll file descriptor to monitor server socket for read events
	// add new file descriptor to EPOLL: EPOLL_CTL_ADD
	if err:= syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, socketFD, &socketServerEvent); err != nil {
		return err
	}

	// waiting for events in loop
	for {
		
		// after 1 second delete the expired keys
		if time.Now().After(lastCronExecutionTime.Add(cronFrequency)) {
			core.DeleteExpiredKeys()
			lastCronExecutionTime = time.Now() // update the last cron execution time to current time
		}


		// n ready file descriptors to which read is available
		nevents, err:= syscall.EpollWait(epollFD, events[:], -1) // -1 timeout means wait indefinitely 
		if err != nil {
			continue
		}

	
		for i:= 0;i<nevents;i++ {
			if int(events[i].Fd) == socketFD {
				// accept the incoming connection
				fd, _, err:= syscall.Accept(socketFD)
				if err != nil {
					continue;
				}

				// increment the no. of clients
				conn_clients++;
				// configure new connection as non-blocking
				syscall.SetNonblock(fd, true);

				// add this TCP connection to be monitored
				var socketClientEvent syscall.EpollEvent = syscall.EpollEvent{
					Fd: int32(fd),
					Events: syscall.EPOLLIN,
				}

				if err:= syscall.EpollCtl(epollFD, syscall.EPOLL_CTL_ADD, fd, &socketClientEvent); err != nil {
					log.Fatal(err)
				}

			} else {
				// read the data, parse the payload and response with appt. reply
				comm:= core.FDComm{Fd: int(events[i].Fd)}
				cmds, err:= readFromConn(comm)
				if err != nil {
					// close the connection
					syscall.Close(int(events[i].Fd))
					// decrement the no. of clients
					conn_clients--
					continue;
				}
				 writeToConn(cmds, comm); 
			}
		}

	}
}

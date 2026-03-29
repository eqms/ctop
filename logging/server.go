package logging

import (
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	socketPath = "./ctop.sock"
	socketAddr = "127.0.0.1:9000"
)

var server struct {
	wg sync.WaitGroup
	ln net.Listener
}

func getListener() net.Listener {
	var ln net.Listener
	var err error
	if debugModeTCP() {
		ln, err = net.Listen("tcp", socketAddr)
	} else {
		ln, err = net.Listen("unix", socketPath)
	}
	if err != nil {
		panic(err)
	}
	return ln
}

func StartServer() {
	server.ln = getListener()

	go func() {
		for {
			conn, err := server.ln.Accept()
			if err != nil {
				// Check if the error is a timeout (Temporary is deprecated since Go 1.18)
				if nErr, ok := err.(net.Error); ok && nErr.Timeout() {
					continue
				}
				return
			}
			go handler(conn)
		}
	}()

	Log.Notice("logging server started")
}

func StopServer() {
	server.wg.Wait()
	if server.ln != nil {
		_ = server.ln.Close()
	}
}

func handler(wc io.WriteCloser) {
	server.wg.Add(1)
	defer server.wg.Done()
	defer func() { _ = wc.Close() }()

	ch := Log.tail()
	defer Log.untail(ch)

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				_, _ = wc.Write([]byte("bye\n"))
				return
			}
			_, _ = fmt.Fprintf(wc, "%s\n", msg)
		case <-Log.done:
			_, _ = wc.Write([]byte("bye\n"))
			return
		}
	}
}

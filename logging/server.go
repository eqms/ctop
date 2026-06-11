package logging

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
)

const socketAddr = "127.0.0.1:9000"

var server struct {
	wg sync.WaitGroup
	ln net.Listener
}

// socketPath returns a per-user private location for the debug socket
// instead of the CWD, so other users sharing a directory cannot attach.
func socketPath() string {
	dir := os.Getenv("XDG_RUNTIME_DIR")
	if dir == "" {
		dir = filepath.Join(os.TempDir(), fmt.Sprintf("ctop-%d", os.Getuid()))
		if err := os.MkdirAll(dir, 0o700); err != nil {
			panic(err)
		}
		if err := os.Chmod(dir, 0o700); err != nil {
			panic(err)
		}
	}
	return filepath.Join(dir, "ctop.sock")
}

func getListener() net.Listener {
	var ln net.Listener
	var err error
	if debugModeTCP() {
		ln, err = net.Listen("tcp", socketAddr)
	} else {
		path := socketPath()
		// remove stale socket from a previous unclean shutdown
		_ = os.Remove(path)
		ln, err = net.Listen("unix", path)
		if err == nil {
			err = os.Chmod(path, 0o600)
		}
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

	Log.Notice("logging server started on " + server.ln.Addr().String())
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

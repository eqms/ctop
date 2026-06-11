# Debug Mode

`ctop` comes with a built-in logging facility and local socket server to simplify debugging at run time.

## Quick Start

If running `ctop` via Docker, debug logging can be most easily enabled as below:
```bash
docker run -ti --rm \
           --name=ctop \
           -e CTOP_DEBUG=1 \
           -e CTOP_DEBUG_TCP=1 \
           -p 9000:9000 \
           -v /var/run/docker.sock:/var/run/docker.sock \
           ghcr.io/eqms/ctop:latest
```

Log messages can be followed by connecting to the default listen address:
```bash
curl -s localhost:9000
```

example output:
```
2024/01/15 15:06:43 INFO logger initialized
2024/01/15 15:06:43 INFO loaded config param key=filterStr value=""
2024/01/15 15:06:43 INFO loaded config param key=sortField value=state
2024/01/15 15:06:43 INFO loaded config switch key=sortReversed value=false
2024/01/15 15:06:43 INFO loaded config switch key=allContainers value=true
2024/01/15 15:06:43 INFO loaded config switch key=enableHeader value=true
2024/01/15 15:06:43 INFO collector started container=7120f83ca...
...
```

## Unix Socket

Debug mode is enabled via the `CTOP_DEBUG` environment variable:

```bash
CTOP_DEBUG=1 ./ctop
```

The socket is created in a per-user private directory: `$XDG_RUNTIME_DIR/ctop.sock` if `XDG_RUNTIME_DIR` is set, otherwise `$TMPDIR/ctop-<uid>/ctop.sock` (mode `0600`).

While `ctop` is running, you can connect to the logging socket via socat or similar tools:
```bash
socat unix-connect:"${XDG_RUNTIME_DIR:-$TMPDIR/ctop-$(id -u)}/ctop.sock" stdio
```

## TCP Logging Socket

In lieu of using a local unix socket, TCP logging can be enabled via the `CTOP_DEBUG_TCP` environment variable:

```bash
CTOP_DEBUG=1 CTOP_DEBUG_TCP=1 ./ctop
```

A TCP listener for streaming log messages will be started on the default listen address (`127.0.0.1:9000`).

**Note:** the TCP listener has no authentication — any local process can connect and read the log stream. Only enable it on trusted single-user hosts; prefer the Unix socket otherwise.

## Log to file

You can also log to a file by specifying `CTOP_DEBUG_FILE=/path/to/ctop.log` environment variable:
```sh
CTOP_DEBUG=1 CTOP_DEBUG_FILE=ctop.log ./ctop
```

This is useful for GoLand to see logs right in debug panel: 
* Edit Run configuration 
* Go to Logs tab
* Specify this log file in "Log file to be shown in console".
Then during debugging you'll see the log tab in debug panel:

![Debug in GoLand](img/goland_debug.png)

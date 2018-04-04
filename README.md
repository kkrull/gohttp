# HTTP Server for Go

This is an HTTP server for Go.


## Requirements

This is being developed on Go 1.10, 64-bit.


## Running

```bash
go run /path/to/github.com/kkrull/gohttp/gohttp.go -p <port> -d <content root directory>
```

Note: It may be necessary to run `go build` the `gohttp` executable separately, if a
[graceful exit code from `SIGTERM`](https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe)
is desired.


## Testing

```bash
$ go get -t #Download dependencies, including those used by tests
$ ginkgo watch #Watch source/test files for changes
```


## TCP Traffic

When the server just listens, accepts, and closes a connection.

```bash
$ ./gohttp -p 1234 -d ... #Server
$ netcat -vz -4 localhost 1234 #Client
```

Packet sniffing

```bash
$ tcpdump -D #Show interfaces; find localhost
$ tcpdump -i <interface> -s 0 -w gohttp--netcat-4.pcap #Capture
$ tcpdump -4 <file> #View file
```

There's also `curl --trace <hex as ascii dump file> ...`

# HTTP Server for Go [![Build Status](https://travis-ci.org/kkrull/gohttp.svg?branch=master)](https://travis-ci.org/kkrull/gohttp)

This is an HTTP server for Go.


## Requirements

This is being developed on Go 1.10, 64-bit.


## Installation

Install Go 1.10 with [their installer](https://golang.org/doc/install), or with `brew install go` if you use homebrew.


### Set up Go environment

Go recommends following a few conventions on setting up your environment

* Set `GOPATH`.  You can set it to the value from `go env GOPATH`, if you're not sure of the conventional path
  on your system.
* Add Go binaries to your system path.  I recommend putting the following into your startup scripts (`.bash_profile` et al)

```bash
go version >/dev/null 2>&1
if (( $? == 0 ))
then
  export PATH="$PATH:$(go env GOPATH)/bin"
fi
```


### Install support tools

This installs [goimports](https://godoc.org/golang.org/x/tools/cmd/goimports) for formatting and organizing imports
and [ginkgo](http://onsi.github.io/ginkgo/) for spec-style testing.

```bash
$ go get github.com/onsi/ginkgo/ginkgo
$ go get golang.org/x/tools/cmd/goimports
```

When you are done, `GOPATH/bin` should contain `ginkgo` and `goimports`.  
`which ginkgo` and `which goimports` should then work if the binaries are installed and present in your `PATH`.


### Clone this repository

```bash
$ cd $(go env GOPATH)
$ mkdir src
$ cd src
$ git clone git@github.com:kkrull/gohttp.git
```

### Set up `pre-push` hook

Set up a Git hook to double check that code is formatted and imports are sorted/curated before pushing.

```bash
$ cp bin/pre-push .git/hooks/pre-push
```


## Running

From the path where you cloned this repository:

```bash
$ go get -t -u -v
$ go build
$ ./gohttp -p <port> -d <content root directory>
```

Note that if you build and run this with `go run`, it will not
[handle `SIGTERM` from Ctrl+C](https://stackoverflow.com/questions/11268943/is-it-possible-to-capture-a-ctrlc-signal-and-run-a-cleanup-function-in-a-defe)
correctly.

When you want to exit the server, press `Ctrl+C`.


## Testing

```bash
$ go get -t #Download dependencies, including those used by tests
$ ginkgo -r #Run tests in all packages
```

Continuous Integration happens on [Travis CI](https://travis-ci.org/kkrull/gohttp).
See `.travis.yml` in this repository for details in the CI environment and how the tests are run.

Additional testing is performed by a version of `cob_spec` that has been configured to start/stop this server.
This version of `cob_spec` can be found [on GitHub](https://github.com/kkrull/cob_spec).


## Support scripts

A few steps of the development process are being automated, as the project takes shape.
These are located in the `bin/` directory:

* `bin/build-and-start.sh`: Re-builds the local binary `gohttp` and runs it.  Pass it the same options you would if you
  were running `gohttp` directly.
* `bin/update-dependencies.sh`: Updates all Go libraries in your `GOPATH` and runs tests to make sure everything still
  works.  *Note that this repository's current branch must have an upstream branch, for this to work.*


## Developer notes

### TCP Traffic

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

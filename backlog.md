# Developer backlog

## Content-type missing on empty file

It should probably be text/plain, but it's worth double-checking to see if there's a specification for this.

```shell
$ curl -4 -v 'http://localhost:1234/cat-form/.gitkeep' ; echo ''
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 1234 (#0)
> GET /cat-form/.gitkeep HTTP/1.1
> Host: localhost:1234
> User-Agent: curl/7.54.0
> Accept: */*
> 
< HTTP/1.1 200 OK
< Content-Type: 
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
```


## Handling requests

Missing ways to cause I/O errors with

* `io.Writer`
* `net.TCPConn`
* `os.File`


Other items:

* Is `Target` needed in each controller?  Is it just used to route to the right controller?


## Naming

- `Controller` might be better referred to as a `Resource` that supports (HTTP) methods.
- `package.Route` might be better called `MethodRoute` and be dedicated to de-muxing HTTP method **only**
- `Router.Route` might be better called `ResourceRoute` and be dedicated to de-muxing over Target/Resource **only**


## Request parsing

Denial of Service: Should allow a request line of at least 8,000 octets.  On the flip side, prevent a
Denial of Service attack with a request line that (almost) never ends (has no CR).
In this case, it should only read up to 8,000 octets and give up if it hasn't seen CR yet.
See RFC 7230 Section 3.1.1 for details

HTTP version

* Given HTTP/1.0 -- it could respond 426 Upgrade Required with Upgrade: HTTP/1.1 and Connection: Upgrade
* Given HTTP/1.2+ -- should it respond with 501 Not Implemented?
  RFC 7231 seems to suggest that it's only meant for an unsupported _method_.
* Given HTTP/2+ -- it could respond 505 HTTP Version Not Supported


## Linting

Enable more linters and drive out the issues.

Here are all the possibilities:

```
    "deadcode",
    "dupl",
    "errcheck",
    "gas",
    "goconst",
    "gocyclo",
    "gofmt",
    "goimports",
    "golint",
    "gosimple",
    "gotype",
    "gotypex",
    "ineffassign",
    "interfacer",
    "lll",
    "maligned",
    "megacheck",
    "misspell",
    "nakedret",
    "safesql",
    "staticcheck",
    "structcheck",
    "test",
    "testify",
    "unconvert",
    "unparam",
    "unused",
    "varcheck",
    "vet",
    "vetshadow"
```

Start with `errcheck`.

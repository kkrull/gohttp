# Developer backlog

## Universal HEAD requests

Any time GET is supported, HEAD should be too.  
Come up with some general mechanism to have a default resource method that calls `GetResource#Get` and ignores the message body.


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


## Header-only responses

A number of responses don't have or need any message body.
Come up with a way to consistently set `Content-Length` to 0.


## Handling requests

Missing ways to cause I/O errors with

* `io.Writer`
* `net.TCPConn`
* `os.File`


## DRYness

* Invoking requests -- buffering and parsing the response message
* Building up request messages -- a builder would come in handy
* Constants for commonly-used content types


## Routing

- `package.Route` might be better called `MethodRoute` and be dedicated to de-muxing HTTP method **only**
- `Router.Route` might be better called `ResourceRoute` and be dedicated to de-muxing over Target/Resource **only**
- The route configuration is getting big enough that it would be handy to configure from a file instead of modifying code all the time.
- It would be more clear which paths go to which routes/resources if each constructor takes the path(s) it serves are parameters.


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


## Cookies

* Sign cookies -- something like sha256
* Encrypt cookies -- with a strong encryption algorithm and salt

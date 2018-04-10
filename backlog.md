# Developer backlog

## Handling requests

* If there is an error handling a request, it should respond 500 Internal Server Error, if at all possible.
* Missing ways to cause I/O errors with
  * `io.Writer`
  * `net.TCPConn`
* HTTP version
  * Given HTTP/1.0 -- it could respond 426 Upgrade Required with Upgrade: HTTP/1.1 and Connection: Upgrade
  * Given HTTP/1.2+ -- should it respond with 501 Not Implemented?
    RFC 7231 seems to suggest that it's only meant for an unsupported _method_.
  * Given HTTP/2+ -- it could respond 505 HTTP Version Not Supported
* Hard-coded responses can be turned into XyzResponse types
  * 400 Bad Request
  * 500 Internal Server Error


## Request parsing

RFC7230 Section 3.1.1

* Should allow a request line of at least 8,000 octets.  On the flip side, prevent a
  Denial of Service attack with a request line that (almost) never ends (has no CR).
  In this case, it should only read up to 8,000 octets and give up if it hasn't seen CR yet.


## General

* Is the current request handling, parsing, and error handling getting large enough to benefit from being in a separate
  `handler/filesystem` package?

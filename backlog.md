# Developer backlog


## Response message construction

It might be nice for `fs.DirectoryListing#WriteTo` to accept `msg.ResponseBuilder`

* `#WithStatus(status, reason)`
* `#WithHeader(name, value)`
* `#WithContentLength(int)`
* `#WriteToBody(string|[]byte)`
* `#Send(io.Writer)`


## Handling requests

Missing ways to cause I/O errors with

* `io.Writer`
* `net.TCPConn`
* `os.File`


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


## Tests

Some types have been extracted after tests have been written on an outer layer.  If there is going to be a lot more
development on these tests, it may make sense to refactor some of these tests to

* test the delegation to the recently-extracted type
* move / refactor the existing tests that apply to the recently-extracted type, down to that level

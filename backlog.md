# Developer backlog


## HEAD requests

* `http.Server` -- ok
  * parses the request with `RequestParser`
  * runs it with `Request#Handle(io.Writer)`
* `http.RequestParser`
  * reads the request from the `bufio.Reader`
  * returns a `fs.GetRequest`
  * returns a `servererror.NotImplemented`
* `GetRequest`
  * Checks the existence and type of `Target`
  * Constructs a `getresponse` aka `http.Response` with the contents of the filesystem
  * Writes the response back to the client with `getresponse#WriteTo(io.Writer)`
* `fs.DirectoryListing` -- ok
  * Writes the HTTP response to the client as an HTML page
* The `msg` package -- ok
  * Handles details of how to format the response as a string


### Improvements

It might be nice for `fs.DirectoryListing#WriteTo` to accept `msg.ResponseBuilder`

* `#WithStatus(status, reason)`
* `#WithHeader(name, value)`
* `#WithContentLength(int)`
* `#WriteToBody(string|[]byte)`
* `#Send(io.Writer)`


Why to `http.Request` and `http.Response` look exactly the same?  That seems odd.


## Handling requests

* Missing ways to cause I/O errors with
  * `io.Writer`
  * `net.TCPConn`
  * `os.File`
* Hard-coded responses can be turned into XyzResponse types
  * 400 Bad Request
  * 500 Internal Server Error


## Request parsing

* Denial of Service: Should allow a request line of at least 8,000 octets.  On the flip side, prevent a
  Denial of Service attack with a request line that (almost) never ends (has no CR).
  In this case, it should only read up to 8,000 octets and give up if it hasn't seen CR yet.
  See RFC 7230 Section 3.1.1 for details
* HTTP version
  * Given HTTP/1.0 -- it could respond 426 Upgrade Required with Upgrade: HTTP/1.1 and Connection: Upgrade
  * Given HTTP/1.2+ -- should it respond with 501 Not Implemented?
    RFC 7231 seems to suggest that it's only meant for an unsupported _method_.
  * Given HTTP/2+ -- it could respond 505 HTTP Version Not Supported


## Packaging

* `http.RequestParser` shouldn't depend on `fs.GetRequest`.  It may be better to configure `RequestParser` with
  what it needs to handle the request, once it has been parsed.
* Is the current request handling, parsing, and error handling getting large enough to benefit from being in a separate
  `handler/filesystem` package?


## Tests

* Some types have been extracted after tests have been written on an outer layer.  If there is going to be a lot more
  development on these tests, it may make sense to refactor some of these tests to
  * test the delegation to the recently-extracted type
  * move / refactor the existing tests that apply to the recently-extracted type, down to that level

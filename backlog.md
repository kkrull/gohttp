# Developer backlog


## Request parsing

RFC7230 Section 3.1.1

- Should allow a request line of at least 8,000 octets.  On the flip side, prevent a
  Denial of Service attack with a request line that (almost) never ends (has no CR).
  In this case, it should only read up to 8,000 octets and give up if it hasn't seen CR yet.
- I/O errors when reading from the socket.
  Need to invest in a mock bufio.Reader that returns errors.

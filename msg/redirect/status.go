package redirect

import "github.com/kkrull/gohttp/msg"

var (
	FoundStatus = msg.Status{Code: 302, Reason: "Found"}
)

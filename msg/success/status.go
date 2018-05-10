package success

import "github.com/kkrull/gohttp/msg"

var (
	OKStatus             = msg.Status{200, "OK"}
	PartialContentStatus = msg.Status{206, "Partial Content"}
)

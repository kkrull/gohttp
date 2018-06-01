package success

import "github.com/kkrull/gohttp/msg"

var (
	OKStatus             = msg.Status{200, "OK"}
	CreatedStatus        = msg.Status{201, "Created"}
	NoContentStatus      = msg.Status{204, "No Content"}
	PartialContentStatus = msg.Status{206, "Partial Content"}
)

package servererror

import "github.com/kkrull/gohttp/msg"

var (
	InternalServerErrorStatus = msg.Status{500, "Internal Server Error"}
	NotImplementedStatus      = msg.Status{501, "Not Implemented"}
)

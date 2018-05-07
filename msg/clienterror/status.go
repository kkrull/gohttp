package clienterror

import "github.com/kkrull/gohttp/msg"

var (
	BadRequestStatus       = msg.Status{400, "Bad Request"}
	NotFoundStatus         = msg.Status{404, "Not Found"}
	MethodNotAllowedStatus = msg.Status{405, "Method Not Allowed"}
)

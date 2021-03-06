package clienterror

import "github.com/kkrull/gohttp/msg"

var (
	BadRequestStatus          = msg.Status{400, "Bad Request"}
	UnauthorizedStatus        = msg.Status{401, "Unauthorized"}
	ForbiddenStatus           = msg.Status{403, "Forbidden"}
	NotFoundStatus            = msg.Status{404, "Not Found"}
	MethodNotAllowedStatus    = msg.Status{405, "Method Not Allowed"}
	ConflictStatus            = msg.Status{409, "Conflict"}
	PreconditionFailedStatus  = msg.Status{412, "Precondition Failed"}
	RangeNotSatisfiableStatus = msg.Status{416, "Range Not Satisfiable"}
)

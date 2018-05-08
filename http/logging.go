package http

// A null object for RequestLogger that does nothing
type noLogger struct{}

func (noLogger) Parsed(message RequestMessage) {}

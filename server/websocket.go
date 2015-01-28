package server

type (
	WebSocketHandler interface {
		Init(*Server) error
		Destroy()
		Process()
	}
)

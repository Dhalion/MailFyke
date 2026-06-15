package smtp

type Server struct {
	// TODO: implement SMTP server with go-smtp
}

func New() *Server {
	return &Server{}
}

func (s *Server) Start(addr string) error {
	// TODO: start SMTP listener with AUTH LOGIN/PLAIN
	return nil
}

func (s *Server) Stop() error {
	// TODO: graceful shutdown
	return nil
}

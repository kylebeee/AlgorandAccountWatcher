package server

func (s *Server) routes() {
	s.GET("/add/:address", s.handleAddToWatchlist())
	s.GET("/list", s.handleGetWatchlist())
}

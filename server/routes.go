package server

import "github.com/gin-gonic/gin"

func (s *Server) routes(router *gin.Engine) {
	// adding addresses should really be a POST request but lets keep it easy to test from a browser
	router.GET("/add/:address", s.handleAddToWatchlist())
	router.GET("/list", s.handleGetWatchlist())
}

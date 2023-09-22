package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

type Server struct {
	*http.Server
	LocalTime *time.Location
	Sentry    *sentry.Client
	Algod     *algod.Client
	LastBlock uint64
	WatchList *WatchList[string]
}

func (s *Server) Close() {
	if s.Sentry != nil {
		sentry.Flush(2 * time.Second)
	}

	err := s.Shutdown(context.Background())
	if err != nil {
		log.Fatalf("[ERROR][Server.Close] %v\n", err)
	}
	log.Println("[INFO][Server.Close] server shutdown")
}

// middleware we add so rest api can be called from any domain via javascript
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

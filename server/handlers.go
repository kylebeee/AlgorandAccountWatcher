package server

import (
	"fmt"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func (s *Server) handleAddToWatchlist() gin.HandlerFunc {
	const operation = "handleAddToWatchlist"
	var hasSentry = s.Sentry != nil

	type request struct {
		Address string `uri:"address" binding:"required"`
	}

	type response struct {
		Ok     bool   `json:"ok"`
		Result string `json:"result,omitempty"`
		Error  string `json:"error,omitempty"`
	}

	return func(c *gin.Context) {
		var (
			err  error
			req  request
			resp response
			hub  *sentry.Hub
		)

		if hasSentry {
			hub = sentrygin.GetHubFromContext(c)
		}

		err = c.ShouldBindUri(&req)
		if err != nil {
			if hasSentry {
				hub.CaptureException(err)
			} else {
				fmt.Printf("[ERROR][%v] %v\n", operation, err)
			}
			resp.Error = "an address is required"
			c.JSON(400, resp)
			return
		}

		if len(req.Address) != 58 {
			if hasSentry {
				hub.CaptureException(err)
			} else {
				fmt.Printf("[ERROR][%v] %v\n", operation, err)
			}
			resp.Error = "invalid address"
			c.JSON(400, resp)
			return
		}

		err = s.WatchList.AddToQueue([]string{req.Address})
		if err != nil {
			if hasSentry {
				hub.CaptureException(err)
			} else {
				fmt.Printf("[ERROR][%v] %v\n", operation, err)
			}
			resp.Error = "an address is required"
			c.JSON(400, resp)
			return
		}

		resp.Ok = true
		resp.Result = "address added to watchlist"
		c.JSON(200, resp)
	}
}

func (s *Server) handleGetWatchlist() gin.HandlerFunc {

	type response struct {
		Ok     bool                           `json:"ok"`
		Result map[string]SlimmedAccountState `json:"result,omitempty"`
		Error  string                         `json:"error,omitempty"`
	}

	return func(c *gin.Context) {
		resp := response{
			Result: map[string]SlimmedAccountState{},
		}
		s.WatchList.Subs.Range(func(key, value interface{}) bool {
			resp.Result[key.(string)] = value.(SlimmedAccountState)
			return true
		})

		resp.Ok = true
		c.JSON(200, resp)
	}
}

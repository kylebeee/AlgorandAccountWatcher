package server

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/algorand/go-algorand-sdk/v2/client/v2/algod"
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

// New is our server constructor function
// this function takes in paramaters from the command line
// and returns a server struct ready to run
func New(mainnet bool, sentryDSN string, tz string, port string) *Server {
	var (
		err error
		env string
	)

	if mainnet {
		env = "mainnet"
	} else {
		env = "testnet"
	}

	// s := &Server{
	// 	Engine: gin.Default(),
	// }

	router := gin.Default()
	router.Use(CORSMiddleware())

	if sentryDSN != "" {
		hub := sentry.CurrentHub()

		client, err := sentry.NewClient(sentry.ClientOptions{
			Dsn:         sentryDSN,
			Environment: env,
		})
		if err != nil {
			log.Fatalln(err)
		}

		hub.BindClient(client)
		router.Use(sentrygin.New(sentrygin.Options{
			Repanic: true,
		}))
	}

	s := &Server{
		Server: &http.Server{
			Addr:    ":" + port,
			Handler: router,
		},
	}

	s.routes(router)

	s.LocalTime, err = time.LoadLocation(tz)
	if err != nil {
		log.Fatalln(err)
	}

	s.Algod, err = algod.MakeClient(fmt.Sprintf("https://%v-api.algonode.cloud", env), "")
	if err != nil {
		log.Fatalln(err)
	}

	s.WatchList = NewWatchList[string](
		[]Processor[string]{
			func(address string) error {
				return s.ProcessAddress(address)
			},
		},
		true,
	)
	go s.WatchList.Start()

	go func() {
		interval := time.NewTicker(time.Minute)
		for {
			<-interval.C
			s.WatchList.Subs.Range(func(key, value interface{}) bool {
				address := key.(string)
				err := s.ProcessAddress(address)
				if err != nil {
					fmt.Println(err)
				}
				return true
			})
		}
	}()

	return s
}

func NewWatchList[L ListType](processors []Processor[L], ignoreFailures bool) *WatchList[L] {
	return &WatchList[L]{
		Subs:       &sync.Map{},
		Processors: processors,
		ListChan:   make(chan []L, 100),
	}
}

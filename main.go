package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kylebeee/AlgorandAccountWatcher/server"
	"github.com/kylebeee/AlgorandAccountWatcher/utils"
)

func main() {

	// flags to add
	mainnet := flag.Bool("m", false, "whether to run on mainnet or not")
	sentryDSN := flag.String("sentry", "", "sentry dsn for where to send error logs")
	timezone := flag.String("tz", "Local", "timezone to use for logging")
	port := flag.String("p", "8080", "port to run server on")
	flag.Parse()

	// start server
	s := server.New(*mainnet, *sentryDSN, *timezone, *port)

	// watch for close signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGTERM, os.Interrupt, syscall.SIGABRT, syscall.SIGTSTP)
	go func() {
		<-sigc
		fmt.Printf("[SERVER][%s] Shutting Down\n", time.Now().Format(utils.TimeFormat))
		s.Close()
	}()

	s.ListenAndServe()
}

# Algorand Account Watcher
Algorand account watcher service example.

## Getting Started
This is a simple example of an account watching service http server & rest api. It utilizes a few simple queues, channels and goroutines to handle the following:
- supports adding new addresses to the watchlist
- checks the address list every 60 seconds for balance changes
- logs changes in the balances of listed addresses
- supports listing all addresses being watched & their latest balances

### Prerequisites
You'll need go 1.21 installed on your local machine. See https://go.dev/ for installation instructions or use Homebrew on MacOS.

### Running the Server

This Rest API supports several optional flags:
<br><br> -m: whether to run on mainnet or not (defaults to testnet)
<br> -sentryDSN: the sentry DSN to use for error reporting
<br> -tz: the timezone to use (defaults to Local)
<br> -port: the port to listen on (defaults to 8080)

You can compile and run the server with the following command
( just make sure you're in the root directory of the project ):
```bash
go run ./main.go
```

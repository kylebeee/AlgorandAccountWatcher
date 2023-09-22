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
<br><br> -m : whether to run on mainnet or not (defaults to testnet)
<br> -sentry : the sentry DSN to use for error reporting
<br> -tz : the timezone to use (defaults to Local)
<br> -p : the port to listen on (defaults to 8080)

You can compile and run the server with the following command
( just make sure you're in the root directory of the project ):
```bash
go run ./main.go
```

### Endpoints

The server only has two endpoints that can be called

```
GET /add/{address}
```
*We really should be doing a POST request here but lets keep it easy to test in the browser*

This endpoint will add the address to the watchlist and start checking it for balance changes. It will return a 200 status code if successful or a 400 status code if the address is invalid.

```
GET /list
```

This endpoint will return a json object containing all the addresses being watched and their latest balances. It will return a 200 status code if successful or a 500 status code if there was an error.

Both endpoints return the following structure
```json
{
    "ok": "boolean",
    "results": "any",
    "error": "string"
}
```
The result being either a confirmation message or the list of addresses and their balances. results or error will be missing depending on the status of the request.
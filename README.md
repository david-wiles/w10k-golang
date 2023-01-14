# W10k Challenge (Golang)

This is a simple web server to see if we can handle 10k concurrent websocket connections, and what limits the server
would have in different situations. 

You can build either executable with go build <program>/main.go.

## Broadcast

The server only does a couple things:

* Prints messages it receives
* Sends the current time to all websockets at the interval defined by `PING_INTERVAL`

## Client2Client

This server implements a basic chat service. The first 36 characters of a text message specify the destination UUID, 
and the rest of the message will be relayed to that websocket.

[k6](https://k6.io/docs/) is a good tool for load testing servers with virtual users. See 
[w10k-k6-clients](https://github.com/david-wiles/w10k-k6-clients) for the test files.

This implementation is based on the [gorilla/websocket chat example](https://github.com/gorilla/websocket/tree/76ecc29eff79f0cedf70c530605e486fc32131d1/examples/chat).

You can deploy this project to a DigitalOcean droplet with `./deploy.sh <domain>`. You will need to add your 
DigitalOcean token and private key to tf/terraform.tfvars as do_token and pvt_key, respectively.

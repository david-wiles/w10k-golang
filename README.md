# W10k Challenge (Golang)

This is a simple web server to see if we can handle 10k concurrent websocket connections, and what limits the server
would have in different situations. The server only does a couple things:

* Prints messages it receives
* Sends the current time to all websockets at the interval defined by `PING_INTERVAL`

Critically, this server is only given 0.5 CPUs and 1Gi of memory.

You can create the image with 

```
docker build -t w10k-go:v1 .
```

and deploy it using

```
docker run --cpus="0.5" --memory="1Gi" --env PING_INTERVAL=10s -p 8080:8080 w10k-go:v1
```

[k6](https://k6.io/docs/) is a good tool for load testing servers with virtual users. See 
[w10k-k6-clients](https://github.com/david-wiles/w10k-k6-clients) for the test files.

This implementation is based on the [gorilla/websocket chat example](https://github.com/gorilla/websocket/tree/76ecc29eff79f0cedf70c530605e486fc32131d1/examples/chat). 

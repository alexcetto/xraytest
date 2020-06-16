# xraytest

Either clone or use: 

`go get github.com/alexcetto/xraytest` and `cd $GOPATH/src/github.com/alexcetto/xraytest`

Run the program with `go run .`

Send a request with curl `curl -X GET localhost:8067`

Observe if you get a big number of `write udp 127.0.0.1:55750->127.0.0.1:2000: write: message too long`

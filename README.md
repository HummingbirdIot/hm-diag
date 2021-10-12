# Helium Diagnostic 

## build

build for your platform, get executable binary file `hm-diag`.

eg:

linux amd64
```
env GOOS=linux GOARCH=amd64 go build
```

linux arm64
```
env GOOS=linux GOARCH=arm64 go build
```

## usage

```
$ hm-diag -h

Helium Diagnostic
Usage: [options] [get | server] 

Subcommand:
  get
        get info to stdout
  server
        run http server, can omit it
Options:
  -m string
        miner http url (default "http://127.0.0.1:4467")
  -p string
        server listening port (default "8090")
```

## example

run server

```bash
$ hm-diag


2021/10/11 22:40:06 server listening on port 8090
```

then get info by http or browser:

```bash
$ curl localhost:8090

{
  "fetch_time": "2021-10-11T22:41:17.825123+08:00",
  "info_height": 1049597,
  "info_region": "US915",
  "peer_addr": "/p2p/11SLTwMavXGwG9T7hogA1xErwUpLZYKQJfYYVvR4s4Kk8hyGGLb",
  "peer_book": [
		// ... omit
  ]
 }

```

Or get the information directly from the command line:

```bash
$ hm-diag get

{
  "fetch_time": "2021-10-11T22:38:33.765549+08:00",
  "info_height": 1049595,
  "info_region": "US915",
  "peer_addr": "/p2p/11SLTwMavXGwG9T7hogA1xErwUpLZYKQJfYYVvR4s4Kk8hyGGLb",
	"peer_book": [
		// ... omit
  ]t
 }
```
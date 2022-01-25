all: server client
  
server: cmd/server/main.go
	go build github.com/bigwhite/tcp-server-demo1/cmd/server
client: cmd/client/main.go
	go build github.com/bigwhite/tcp-server-demo1/cmd/client

clean:
	rm -fr ./server
	rm -fr ./client

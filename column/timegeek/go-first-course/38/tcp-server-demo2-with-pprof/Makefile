all: server client
  
server: cmd/server/main.go
	go build -race github.com/bigwhite/tcp-server-demo2/cmd/server
client: cmd/client/main.go
	go build -race github.com/bigwhite/tcp-server-demo2/cmd/client

clean:
	rm -fr ./server
	rm -fr ./client

From ubuntu:14.04

RUN apt-get update \
      && apt-get install -y software-properties-common \
      && add-apt-repository ppa:gophers/archive \
      && apt-get update \
      && apt-get install -y golang-1.9-go \
                            git \
      && rm -rf /var/lib/apt/lists/*

ENV GOPATH /root/go
ENV GOROOT /usr/lib/go-1.9
ENV PATH="/usr/lib/go-1.9/bin:${PATH}"

COPY ./httpserver.go /root/httpserver.go
RUN go build -o /root/httpd /root/httpserver.go \
      && chmod +x /root/httpd

WORKDIR /root
ENTRYPOINT ["/root/httpd"]

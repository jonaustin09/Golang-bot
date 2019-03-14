FROM golang:1.12

# Install protoc
RUN curl -s -L https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip > protoc-3.7.0-linux-x86_64.zip
RUN apt update  && apt install unzip
RUN unzip protoc-3.7.0-linux-x86_64.zip -d protoc-3.7.0-linux-x86_64
RUN mv protoc-3.7.0-linux-x86_64/bin/protoc /usr/local/sbin
RUN mv protoc-3.7.0-linux-x86_64/include /usr/local/
RUN rm -rf protoc-3.7.0-linux-x86_64*

RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/pressly/goose/cmd/goose
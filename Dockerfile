FROM golang:1.12

# Install protoc, apply migrations, populate default env file
RUN curl -s -L https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip > protoc-3.7.0-linux-x86_64.zip &&\
    apt update  &&\
    apt install unzip &&\
    unzip protoc-3.7.0-linux-x86_64.zip -d protoc-3.7.0-linux-x86_64 &&\
    mv protoc-3.7.0-linux-x86_64/bin/protoc /usr/local/sbin &&\
    mv protoc-3.7.0-linux-x86_64/include /usr/local/ &&\
    rm -rf protoc-3.7.0-linux-x86_64* &&\
    go get github.com/golang/protobuf/protoc-gen-go &&\
    go get github.com/pressly/goose/cmd/goose

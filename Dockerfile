FROM golang:1.12

RUN go get github.com/jinzhu/gorm
RUN go get github.com/sirupsen/logrus
RUN go get gopkg.in/tucnak/telebot.v2
RUN go get github.com/satori/go.uuid
RUN go get github.com/mattn/go-sqlite3
RUN go get github.com/joho/godotenv

# Install protoc
RUN curl -s -L https://github.com/protocolbuffers/protobuf/releases/download/v3.7.0/protoc-3.7.0-linux-x86_64.zip > protoc-3.7.0-linux-x86_64.zip
RUN apt update  && apt install unzip
RUN unzip protoc-3.7.0-linux-x86_64.zip -d protoc-3.7.0-linux-x86_64
RUN mv protoc-3.7.0-linux-x86_64/bin/protoc /usr/local/sbin
RUN mv protoc-3.7.0-linux-x86_64/include /usr/local/
RUN rm -rf protoc-3.7.0-linux-x86_64*

RUN go get google.golang.org/grpc
RUN go get github.com/golang/protobuf/protoc-gen-go
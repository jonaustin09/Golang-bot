FROM golang:1.12

RUN go get github.com/jinzhu/gorm
RUN go get github.com/sirupsen/logrus
RUN go get gopkg.in/tucnak/telebot.v2
RUN go get github.com/satori/go.uuid
RUN go get github.com/mattn/go-sqlite3
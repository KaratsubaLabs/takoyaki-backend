FROM golang:1.16-alpine
WORKDIR /app

RUN apk update
RUN apk add mkpasswd

COPY go.mod go.sum ./
COPY .env ./

RUN go mod download

COPY *.go api cli db util vps ./

RUN go build -o /takoyaki github.com/KaratsubaLabs/takoyaki-backend
CMD [ "/takoyaki", "server" ]

FROM golang:1.16-alpine
WORKDIR /app

RUN apk update
RUN apk add mkpasswd

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./
COPY .env ./

RUN go build -o /takoyaki
CMD [ "/takoyaki", "server" ]

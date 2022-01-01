FROM golang:1.16-alpine
WORKDIR /app

RUN apk update
RUN apk add mkpasswd

# don't do this - copy only the subset of files you need (or use dockerignore)
COPY . ./

RUN go mod download

RUN go build -o /takoyaki github.com/KaratsubaLabs/takoyaki-backend
CMD [ "/takoyaki", "server" ]

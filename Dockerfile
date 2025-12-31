FROM golang:1.23-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
WORKDIR /usr/src/app/src

RUN go build -v -o /usr/local/bin/app

CMD ["app"]
FROM golang:1.24.4-alpine

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

ENV RAPID_API_KEY=$RAPID_API_KEY
ENV RAPID_API_HOST=$RAPID_API_HOST
ENV PORT=$PORT

COPY . .

RUN go build -v -o /usr/src/app/wellness

CMD ["/usr/src/app/wellness"]
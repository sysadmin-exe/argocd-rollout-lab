
FROM golang:1.23-alpine AS base
WORKDIR /my-app
COPY go.mod ./
RUN go mod download
COPY . .

RUN go build -o app

EXPOSE 8000

CMD ["/my-app/app"]

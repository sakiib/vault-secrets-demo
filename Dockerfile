FROM golang:latest as build
WORKDIR /app
COPY . ./
RUN go build -o bin/app main.go
ENTRYPOINT ["./bin/app"]
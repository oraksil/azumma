FROM golang:1.14.3 AS builder
WORKDIR /build
COPY . /build
RUN GOOS=linux GOARCH=386 go build -o app cmd/app.go

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /deploy
COPY --from=builder /build/app .
CMD ["./app"]  

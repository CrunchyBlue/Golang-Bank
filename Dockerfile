FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

FROM alpine:3.16 
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY db/migrations ./migrations
COPY app.env .
COPY start.sh .

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]

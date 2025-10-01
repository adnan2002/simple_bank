FROM golang:1.25.1-alpine3.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY Makefile .
COPY db/migration ./db/migration
COPY wait-for.sh .
COPY start.sh .

RUN apk add --no-cache make curl tar netcat-openbsd
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz \
 | tar xvz && mv migrate /usr/bin/

RUN chmod +x wait-for.sh start.sh

EXPOSE 8080
CMD ["/app/start.sh"]
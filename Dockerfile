FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o web-server .

FROM alpine:latest

RUN addgroup -S app && adduser -S -G app -h /home/app app

RUN mkdir -p /home/app/var/www

COPY www/ /home/app/var/www/
COPY --from=builder /build/web-server /home/app/web-server

RUN chown -R app:app /home/app

USER app
WORKDIR /home/app

EXPOSE 8080

ENTRYPOINT ["./web-server"]

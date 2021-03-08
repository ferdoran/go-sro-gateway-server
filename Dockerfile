FROM golang:1.15.3 as builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go clean --modcache
RUN GOOS=linux CGO_ENABLED=0 go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /etc/go-sro-gateway-server/
COPY --from=builder /app/main .

EXPOSE 15779
CMD ["./main"]
FROM golang:1.22.3-alpine3.19 as builder
RUN apk update
RUN apk add -U --no-cache ca-certificates && update-ca-certificates

WORKDIR /GoRateLimiterFC/cmd/server
COPY . /GoRateLimiterFC
EXPOSE 8080
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o GoRateLimiterFC

FROM scratch
WORKDIR /GoRateLimiterFC/cmd/server
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /GoRateLimiterFC/cmd/server .
EXPOSE 8080
ENTRYPOINT [ "./GoRateLimiterFC" ]
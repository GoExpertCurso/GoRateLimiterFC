FROM golang:latest

WORKDIR /app


COPY . /GoRateLimiterFC
CMD ["tail", "-f", "/dev/null"]
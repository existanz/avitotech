FROM golang:1.23

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /build ./cmd/api/main.go \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]
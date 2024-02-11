FROM golang:1.20

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY ./cmd/main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-go

EXPOSE 8080

CMD ["/docker-go"]

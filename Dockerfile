FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -C cmd -o tma_dashboard

CMD ["./cmd/tma_dashboard"]

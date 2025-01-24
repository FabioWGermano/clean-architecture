FROM golang:1.23.3

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/ordersystem

RUN go build -o /main main.go wire_gen.go

CMD ["/main"]
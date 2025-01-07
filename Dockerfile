FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum .env ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/ordersystem

# Instale o netcat
# RUN apt-get update && apt-get install -y netcat-openbsd

RUN go build -o /main main.go wire_gen.go

CMD ["/main"]
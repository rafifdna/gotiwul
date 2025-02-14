FROM golang:1.21-alpine

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o gotiwul ./cmd/gotiwul

EXPOSE 80 443

CMD ["./gotiwul"]
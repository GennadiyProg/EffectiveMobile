FROM golang:1.25.3

WORKDIR /app
COPY . .
RUN go build -o main .

CMD ["./main"]
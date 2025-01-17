FROM golang:1.23

WORKDIR /

COPY . .

RUN go build -o main ./cmd

EXPOSE 8080

CMD ["./main"]
FROM golang:1.21-alpine


WORKDIR /go/src/Wallet


COPY . .


RUN go build -o mainFile ./cmd/main.go


EXPOSE 8080


CMD ["./mainFile"]
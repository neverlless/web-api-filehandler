FROM golang:1.22.0 AS builder

WORKDIR /src

COPY . .

RUN go get -v . \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o web-api-filehandler

FROM alpine:3.19.1

WORKDIR /opt

COPY --from=builder /src/web-api-filehandler .

RUN chmod +x web-api-filehandler

CMD ["./web-api-filehandler"]

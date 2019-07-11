FROM golang:1.11.4-stretch as builder

WORKDIR /root/go/okcoin/alert-sender

COPY . ./

ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build -o /root/app main.go

FROM alpine:3.8

WORKDIR /root

COPY --from=builder /root/app .

ENTRYPOINT ["/root/app"]
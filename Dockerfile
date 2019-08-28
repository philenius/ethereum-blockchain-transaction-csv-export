FROM golang:1.12.9-stretch AS builder

WORKDIR /go/src/ethereum-blockchain-transaction-csv-export/

COPY . .

RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ethereum-blockchain-transaction-csv-export

FROM alpine:3.10.2

WORKDIR /app
ENV LOGXI=*=INF

COPY --from=builder /go/src/ethereum-blockchain-transaction-csv-export/ethereum-blockchain-transaction-csv-export .

ENTRYPOINT [ "./ethereum-blockchain-transaction-csv-export" ]
CMD [ "-start", "0", "-count", "100" ]

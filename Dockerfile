FROM golang:1.12.9-stretch AS builder

WORKDIR /go/src/ethereum-blockchain-transaction-csv-export/

COPY . .

RUN go get -d -v ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ethereum-blockchain-transaction-csv-export

FROM alpine:3.10.2

WORKDIR /app

ENV LOGXI=*=INF \
    BLOCK_COUNT=1000 \
    HOSTNAME=127.0.0.1 \
    PORT=8545 \
    START_BLOCK_HEIGHT=0 \
    STATS_INTERVAL_IN_SECONDS=5 \
    WORKER_COUNT_FOR_BLOCKS=10 \
    WORKER_COUNT_FOR_TRANSACTIONS=20

COPY --from=builder /go/src/ethereum-blockchain-transaction-csv-export/ethereum-blockchain-transaction-csv-export .

ENTRYPOINT [ "./ethereum-blockchain-transaction-csv-export" ]

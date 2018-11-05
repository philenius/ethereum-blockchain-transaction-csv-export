# Ethereum Blockchain Transaction CSV Export

This application extracts transactions from the Ethereum blockchain and exports the data as a CSV file. It uses the RPC API of [go-ethereum / geth](https://github.com/ethereum/go-ethereum), the Golang implementation fo the Ethereum protocol.

## Run
* You'll need an Ethereum node in order to execute API requests. If you don't have your own net yet, then you can start a local Ethereum node as a Docker container with this command:
  ```bash
  docker run -d --name ethereum-node -v /home/user/ethereum:/root -p 8545:8545 -p 30303:30303 \
         ethereum/client-go --rpc --rpcapi "eth,net,web3" --rpcaddr 0.0.0.0 --syncmode "fast"
  ```

* You'll also need a local installation of Golang to run this application:
  ```bash
  go get ./...
  go build
  LOGXI=* ./ethereum-blockchain-transaction-csv-export -start 46147 -count 10
  ```
  This demo should fetch the blocks 46147 through 46157 including the first ever made transaction on the Ethereum blockchain.  

:warning: It might take some time, until your Ethereum node synchronized enough blocks to respond to API requests for those blocks. You can also set the flag `--start` to `0`. But, keep in mind that about the first 46.000 blocks of the Ethereum blockchain contain any transactions.

**Another example:**  

This command exports all transactions contained in the first 4 million blocks (estimated time to completion: 11 hours):

```bash
LOGXI=WRN nohup ./ethereum-blockchain-transaction-csv-export \
    -count 4000000 \
    -blockConcurr 20 \
    -txConcurr 30 \
    &
```

## Configuration
This application can be configured using command line flags. Type `--help` to find out more about the available configurations:
```bash
./ethereum-blockchain-transaction-csv-export --help
Usage of ./ethereum-blockchain-transaction-csv-export:
  -blockConcurr int
        The count of concurrent workers for fetching block. (default 10)
  -count int
        The total amount of blocks to fetch. (default 1000)
  -host string
        The hostname / IP address of the Ethereum node. (default "127.0.0.1")
  -port int
        The port number of the Ethereum node. (default 8545)
  -start int
        The height / number of the block where to start.
  -statsIntervalSec int
        The invertal in seconds to display stats. (default 5)
  -txConcurr int
        The count of concurrent workers for fetching transactions. (default 20)
```

# Ethereum Blockchain Transaction CSV Export

This application extracts transactions from the Ethereum blockchain and exports the data as a CSV file. It uses the RPC API of [go-ethereum / geth](https://github.com/ethereum/go-ethereum), the Golang implementation fo the Ethereum protocol.

## Run
* You'll need an Ethereum node for executing API requests. If you don't have your own net yet, then you can start a local Ethereum node as a Docker container with this command:
  ```bash
  docker run -d --name ethereum-node -v /home/user/ethereum:/root -p 8545:8545 -p 30303:30303 \
         ethereum/client-go --rpc --rpcapi "eth,net,web3" --rpcaddr 0.0.0.0 --syncmode "fast"
  ```

* You'll also need a local installation of Golang for running this application:
  ```bash
  go get ./...
  go build
  LOGXI=* ./ethereum-blockchain-transaction-csv-export -start 46147 -count 10
  ```
  This demo should fetch the blocks 46147 through 46157 inlcuding one the first ever made transaction on the Ethereum blockchain.  

:warning: It might take some time, until your Ethereum node synchronized enough blocks to respond to API requests for those blocks. You can also set the flag `--start` to `0`. But, keep in mind that about the first 46.000 blocks of the Ethereum blockchain contain any transactions.

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
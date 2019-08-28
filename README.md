# Ethereum Blockchain Transaction CSV Export

This application extracts transactions from the Ethereum blockchain and exports the data as a CSV file. It uses the RPC API of [go-ethereum / geth](https://github.com/ethereum/go-ethereum), the Golang implementation fo the Ethereum protocol.

## Run
* You'll need an Ethereum node in order to execute API requests. If you don't have your own net yet, then you can start a local Ethereum node as a Docker container with this command:
   ```bash
   docker run \
       --name ethereum-node \
       -v /home/user/ethereum:/root \
       -p 8545:8545 \
       -p 30303:30303 \
       ethereum/client-go \
       --rpc \
       --rpcapi "eth,net,web3" \
       --rpcaddr 0.0.0.0 \
       --syncmode "fast"
   ```

* This application is dockerized. Use these commands to build and run the application:
   ```bash
   docker build -t ethereum-blockchain-transaction-csv-export .

   # create directory for log and CSV files
   mkdir ./output/

   docker run \
      --network host \
      -v $(pwd)/output:/output \
      -e LOGXI="*" \
      -e START_BLOCK_HEIGHT="46147" \
      -e BLOCK_COUNT="10" \
      ethereum-blockchain-transaction-csv-export
   ```

   This demo should fetch the blocks 46147 through 46157. These blocks contain only one transaction, the first transaction that was ever made on the Ethereum blockchain.

* The resulting CSV export is stored in the file `./output/geth_tx_export_<yyyy-MM-dd-HH-mm-ss>.csv`. The CSV file contains the following columns:
  
  | Column            | Description                                                  |
  | ----------------- | ------------------------------------------------------------ |
  | tx_hash           | 32 Bytes hash of the transaction.     |
  | tx_nonce          | The number of transactions made by the sender prior to this one. |
  | tx_block_hash     | 32 Bytes hash of the block where this transaction was in. |
  | tx_block_number   | Block number where this transaction was in. `null` when its pending. |
  | tx_index          | Integer of the transaction's index position in the block. `null` when its pending. |
  | tx_from           | 20 Bytes address of the sender |
  | tx_to             | 20 Bytes address of the receiver. `null` when its a contract creation transaction. |
  | tx_value          | Value transferred in Wei. |
  | tx_gas            | Gas provided by the sender. |
  | tx_gas_price      | Gas price provided by the sender in Wei. |
  | tx_input          | The data send along with the transaction. |
  | tx_timestamp      | Integer of the unix timestamp when the transaction was sent. |


:warning: It might take some time, until your Ethereum node synchronized enough blocks to respond to API requests for those blocks. You can also set the environment variable `START_BLOCK_HEIGHT` to `0`. But, keep in mind that about the first 46.000 blocks of the Ethereum blockchain don't contain any transactions.

**Another example:**  

This command exports all transactions contained in the first 4 million blocks (this process might take several hours):

```bash
docker run \
  --network host \
  -v $(pwd)/output:/output \
  -e LOGXI="*=INF" \
  -e START_BLOCK_HEIGHT="0" \
  -e BLOCK_COUNT="4000000" \
  -e WORKER_COUNT_FOR_BLOCKS="20" \
  -e WORKER_COUNT_FOR_TRANSACTIONS="30" \
  ethereum-blockchain-transaction-csv-export
```

## Configuration

This application is configured via the environment. The following environment variables can be used:

```
BLOCK_COUNT
  [description] The total amount of blocks to fetch.
  [type]        Integer
  [default]     1000
  [required]    
HOSTNAME
  [description] The hostname / IP address of the Ethereum node.
  [type]        String
  [default]     127.0.0.1
  [required]    
PORT
  [description] The port number of the Ethereum node.
  [type]        Integer
  [default]     8545
  [required]    
START_BLOCK_HEIGHT
  [description] The height / number of the first block to fetch.
  [type]        Integer
  [default]     0
  [required]    
STATS_INTERVAL_IN_SECONDS
  [description] The invertal in seconds to display stats.
  [type]        Integer
  [default]     5
  [required]    
WORKER_COUNT_FOR_BLOCKS
  [description] The count of concurrent workers for fetching blocks.
  [type]        Integer
  [default]     10
  [required]    
WORKER_COUNT_FOR_TRANSACTIONS
  [description] The count of concurrent workers for fetching transactions.
  [type]        Integer
  [default]     20
  [required]
```

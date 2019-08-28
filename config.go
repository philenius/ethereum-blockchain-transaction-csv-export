package main

type config struct {
	BlockCount                 int    `default:"1000" desc:"The total amount of blocks to fetch." split_words:"true"`
	Hostname                   string `default:"127.0.0.1" desc:"The hostname / IP address of the Ethereum node." split_words:"true"`
	Port                       int    `default:"8545" desc:"The port number of the Ethereum node." split_words:"true"`
	StartBlockHeight           int    `default:"0" desc:"The height / number of the first block to fetch." split_words:"true"`
	StatsIntervalInSeconds     int    `default:"5" desc:"The invertal in seconds to display stats." split_words:"true"`
	WorkerCountForBlocks       int    `default:"10" desc:"The count of concurrent workers for fetching blocks." split_words:"true"`
	WorkerCountForTransactions int    `default:"20" desc:"The count of concurrent workers for fetching transactions." split_words:"true"`
}

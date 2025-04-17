# About

This is a simple Proof of Work (PoW) blockchain implementation with basic functionalities. The blockchain includes three main features: transaction handling, mining, and retrieving the blockchain.

# Features

1. Transaction

- Allows users to create transactions.
- Triggered via the API endpoint /transaction.
- Transactions are temporarily stored in the mempool until they are added to a block through mining.

2. Mining

- The mining process records transactions from the mempool into a new block.
- The time required for mining +- 5 minutes.
- Once mined, the block is added to the blockchain.
- miner will get 100 reward if successfully mining the block
- reward wll decrease 50% every 1000 block created

3. Get Blockchain

- Retrieves the entire blockchain.
- Allows users to see the number of blocks created so far.

# API

| Endpoint            | Method | Description                                    |
| ------------------- | ------ | ---------------------------------------------- |
| `/transaction`      | POST   | Add a new transaction to the mempool.          |
| `/mine/:address`    | GET    | Mine a new block and add it to the blockchain. |
| `/blockchain`       | GET    | Retrieve the entire blockchain.                |
| `/balance/:address` | GET    | Retrieve balance by address.                   |

# How it works

Users create transactions that are temporarily stored in the mempool.

Mining processes these transactions, groups them into a block, and adds the block to the blockchain.

Users can retrieve the blockchain to check the number of blocks and stored transactions.

# Project structure

📁proof-of-work/
├── 📁 blockchain/
│ ├── block.go
│ ├── blockchain.go
│ ├── pow.go
│ ├── transaction.go
│ └── wallet.go
│
├── 📁 server/
│ ├── handler.go
│ └── router.go
│
├── .env
├── go.mod
├── main.go
└── README.md

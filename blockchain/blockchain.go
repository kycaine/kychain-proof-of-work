package blockchain

import "fmt"

type Blockchain struct {
	Blocks  []*Block
	Mempool []Transaction
}

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock([]Transaction{}, "0")
	return &Blockchain{
		Blocks: []*Block{genesisBlock},
	}
}

func (bc *Blockchain) AddTransaction(sender, recipient string, amount int) {
	tx := Transaction{Sender: sender, Recipient: recipient, Amount: amount}
	bc.Mempool = append(bc.Mempool, tx)
	fmt.Println("Transaction added to mempool!")
}

func (bc *Blockchain) MineBlock() *Block {
	if len(bc.Mempool) == 0 {
		fmt.Println("Mempool is empty, there are no transactions to mine.")
		return nil
	}

	previousBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(bc.Mempool, previousBlock.Hash)

	newBlock.Index = len(bc.Blocks)
	bc.Blocks = append(bc.Blocks, newBlock)
	bc.Mempool = []Transaction{}

	fmt.Println("Block mined successfully:", newBlock.Hash)
	return newBlock
}

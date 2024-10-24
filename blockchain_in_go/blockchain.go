package main

type Blockchain struct {
	blocks []*Block
}

func (b *Blockchain) AddBlock(data string) {
	prevBlock := b.blocks[len(b.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	b.blocks = append(b.blocks, newBlock)
}

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock("genesis block", []byte{})
	blockchain := &Blockchain{[]*Block{genesisBlock}}
	return blockchain
}

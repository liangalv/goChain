package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func init() {
	log.SetPrefix("goChain Block: ")
}

func main() {
}

type BlockChain struct {
	memPool []*Transaction
	chain   []*Block
}

func NewBlockChain() *BlockChain {
	bc := new(BlockChain)
	//Generates an Empty Block as the prevHash of the Genesis Block
	bc.CreateBlock(0, (&Block{}).Hash(), make([]*Transaction, 0, 1000))
	return bc
}

func (bc *BlockChain) CreateBlock(nonce int, prevHash [32]byte, trans []*Transaction) *Block {
	b := &Block{
		nonce:        nonce,
		prevHash:     prevHash,
		timestamp:    time.Now().UnixNano(),
		transactions: trans,
	}
	bc.chain = append(bc.chain, b)
	return b
}

// Helper Methods
func (bc *BlockChain) LastBlock() *Block {
	return bc.chain[len(bc.chain)-1]
}

func (bc *BlockChain) hashPrevBlock() [32]byte {
	return bc.LastBlock().Hash()
}

// Formatting methods
func (bc *BlockChain) String() string {
	var sb strings.Builder
	div := strings.Repeat("=", 25)
	for i, block := range bc.chain {
		sb.WriteString(fmt.Sprintf("%s Block:%d %s", div, i, div))
		sb.WriteString(fmt.Sprintf(block.String()))
	}
	return sb.String()
}
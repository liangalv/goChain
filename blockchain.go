package main

import (
	"fmt"
	"log"
	"strings"
	"time"
	"golang.org/x/crypto/sha3"
)

func init() {
	log.SetPrefix("goChain Block: ")
}

func main() {
	bc := NewBlockChain()
	fmt.Println(bc)
}

type BlockChain struct {
	memPool []Transaction
	chain   []*Block
}

func NewBlockChain() *BlockChain {
	bc := new(BlockChain)
	//CreateBlock appends to the chain
	sum := sha3.Sum256([]byte("Genesis"))
	bc.CreateBlock(0, sum)
	return bc
}

func (bc *BlockChain) CreateBlock(nonce int, previousHash [32]byte) *Block {
	b := &Block{nonce: nonce, prevHash: previousHash}
	bc.chain = append(bc.chain, b)
	b.timestamp = time.Now().UnixNano()

	return b
}

//Formatting methods

func (bc *BlockChain) String() string {
	var sb strings.Builder
	div := strings.Repeat("=", 25)
	for i, block := range bc.chain {
		sb.WriteString(fmt.Sprintf("%s Block:%d %s", div, i, div))
		sb.WriteString(fmt.Sprintf(block.String()))
	}
	return sb.String()
}

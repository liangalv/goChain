package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
	"time"
)

func init() {
	log.SetPrefix("goChain Block: ")
}

func main() {
	//Set up the server
	lis, err := net.Listen("tcp", "9000")
	if err != nil {
		log.Fatalf("Failed to open port 9000: %v", err)
	}
	grpcServer := grpc.NewServer()
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to server grpcServer over port 9000: %v", err)
	}
	//check the network for any blockchain that being broadcasted, if so sync node's embedded db
	//else read from db and spin up bc state
	// bc := NewBlockChain()

}

type BlockChain struct {
	accounts []*Account
	memPool  []*Transaction
	chain    []*Block
}

// Todo: we are only taking a slice of the underlying blockchain, this genesis block should only be called once
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

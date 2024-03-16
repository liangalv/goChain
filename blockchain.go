package main

import (
	"fmt"
	"github.com/liangalv/goChain/core"
	"github.com/liangalv/goChain/core/services"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
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
	//Instantiate services
	//TODO: separate networking concerns to network.go
	ts := services.TransactionService{}

	grpcServer := grpc.NewServer()
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve grpcServer over port 9000: %v", err)
	}
	//check the network for any blockchain that being broadcasted, if so sync node's embedded db
	//else read from db and spin up bc state
	// bc := NewBlockChain()

}

type BlockChain struct {
	accounts []*core.Account
	memPool  []*core.Transaction
	chain    []*core.Block
}

// Todo: we are only taking a slice of the underlying blockchain, this genesis block should only be called once
func NewBlockChain() *BlockChain {
	bc := new(BlockChain)
	//Generates an Empty Block as the prevHash of the Genesis Block
	bc.CreateBlock((&core.Block{}).Hash(), make([]*core.Transaction, 0, 1000))
	return bc
}

func (bc *BlockChain) CreateBlock(prevHash [32]byte, trans []*core.Transaction) *core.Block {
	b := core.NewBlock(prevHash, trans)
	bc.chain = append(bc.chain, b)
	return b
}

// Helper Methods
func (bc *BlockChain) LastBlock() *core.Block {
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

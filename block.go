package main

import (
	"fmt"
	"strings"
	"time"
)

type Block struct {
	nonce        int
	prevHash     [32]byte
	timestamp    int64
	transactions []string
}

func (b Block) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\nNonce: %d\n", b.nonce))
	sb.WriteString(fmt.Sprintf("PrevHash: %s\n", b.prevHash))
	sb.WriteString(fmt.Sprintf("Timestamp: %d\n", b.timestamp))
	sb.WriteString("Transactions:\n")

	for i, transaction := range b.transactions {
		sb.WriteString(fmt.Sprintf("\t%d: %s\n", i+1, transaction))
	}
	sb.WriteString(fmt.Sprintln(""))

	return sb.String()
}

func NewBlock(nonce int, prevHash [32]byte) *Block {
	return &Block{
		timestamp: time.Now().UnixNano(),
		nonce:     nonce,
		prevHash:  prevHash,
	}
}

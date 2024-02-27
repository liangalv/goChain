package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/sha3"
)

type Block struct {
	nonce        int
	prevHash     [32]byte
	timestamp    int64
	transactions []*Transaction
	gasLimit     uint32
}

func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha3.Sum256(m)
}

func NewBlock(nonce int, prevHash [32]byte) *Block {
	return &Block{
		timestamp: time.Now().UnixNano(),
		nonce:     nonce,
		prevHash:  prevHash,
	}
}

// TODO
func (b *Block) BuildBlock(mp *MemPool) {
}

// Stringify: Stringer interface implementation
func (b Block) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("\nNonce: %d\n", b.nonce))
	sb.WriteString(fmt.Sprintf("PrevHash: %s\n", hex.EncodeToString(b.prevHash[:])))
	sb.WriteString(fmt.Sprintf("Timestamp: %d\n", b.timestamp))
	sb.WriteString("Transactions:\n")

	for i, transaction := range b.transactions {
		sb.WriteString(fmt.Sprintf("\t%d: \n %s\n", i+1, transaction))
	}
	sb.WriteString(fmt.Sprintln(""))

	return sb.String()
}

// JSON Encoding: this implements the json.Marshaller interface
// MarshalJSON needs to overriden as json.Marshal does not encode private fields
func (b *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Nonce        int            `json:"nonce"`
		PrevHash     [32]byte       `json:"previous_hash"`
		Timestamp    int64          `json:"timestamp"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    b.timestamp,
		Nonce:        b.nonce,
		PrevHash:     b.prevHash,
		Transactions: b.transactions,
	})
}

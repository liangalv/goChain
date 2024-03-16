package core

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strings"
	"time"
)

const (
	//TODO: make the gaslimit dynamic based on load
	GASLIMIT = 30 * 1000000 //30 million
)

type Block struct {
	//Header
	ID           [32]byte
	timestamp    int64
	parentHash   [32]byte
	trieRootHash [32]byte
	gasLimit     uint32
	//Body
	transactions []*Transaction
}

// TODO: write proto file, change marshalling type for actual encoding
func (b *Block) Hash() [32]byte {
	m, _ := json.Marshal(b)
	return sha3.Sum256(m)
}

// TODO: we need to provide a trieRootHash
// TODO: we need to hash the headers of the block
func NewBlock(parentHash [32]byte, trans []*Transaction) *Block {
	return &Block{
		timestamp:    time.Now().UnixNano(),
		parentHash:   parentHash,
		transactions: trans,
		gasLimit:     GASLIMIT,
	}
}

// Stringify: Stringer interface implementation
func (b Block) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("ID: %v\n", b.ID))
	sb.WriteString(fmt.Sprintf("ParentHash: %s\n", hex.EncodeToString(b.parentHash[:])))
	sb.WriteString(fmt.Sprintf("Timestamp: %d\n", b.timestamp))
	sb.WriteString(fmt.Sprintf("Root Hash: %v\n", b.trieRootHash))
	sb.WriteString(fmt.Sprintf("Gas Limit: %d\n", b.gasLimit))
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
		ID           [32]byte       `json:"ID"`
		ParentHash   [32]byte       `json:"parent_hash"`
		TrieRootHash [32]byte       `json:"trie_root_hash"`
		Timestamp    int64          `json:"timestamp"`
		Gaslimit     uint32         `json:"gas_limit"`
		Transactions []*Transaction `json:"transactions"`
	}{

		ID:           b.ID,
		ParentHash:   b.parentHash,
		TrieRootHash: b.trieRootHash,
		Timestamp:    b.timestamp,
		Gaslimit:     b.gasLimit,
		Transactions: b.transactions,
	})
}

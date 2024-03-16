package core

import (
	"fmt"
	"github.com/liangalv/goChain/core/types"
	"golang.org/x/crypto/sha3"
	"google.golang.org/protobuf/proto"
	"strings"
	"time"
)

type Transaction struct {
	//112 Bytes per transaction
	timestamp       int64               //24 bytes
	ID              [32]byte            //32
	senderAddress   [AddressLength]byte //20
	receiverAddress [AddressLength]byte
	//this index is only specific to the application layer
	index int    //8
	value uint32 //4
	gas   uint32 //4
}

func NewTransaction(i int, v, g uint32, s, r [AddressLength]byte) *Transaction {
	//TODO: Remember in the consensus engine you need to ensure that you're rejecting transaction with duplicates
	trans := &Transaction{
		//apparently we need something called consensus based timestamp generation
		timestamp:       time.Now().Unix(),
		senderAddress:   s,
		receiverAddress: r,
		value:           v,
		gas:             g,
	}
	//TODO:error handling
	trans.ID = trans.hashTransaction()
	return trans
}

// Helper methods
func (t *Transaction) hashTransaction() [32]byte {
	data, _ := proto.Marshal(t.convertToTransactionPbMsg())
	return sha3.Sum256(data)
}

// Convert to protoreflect.Protomessage type for pb marshalling
func (t *Transaction) convertToTransactionPbMsg() *types.TransactionMsg {
	return &types.TransactionMsg{
		Timestamp:       t.timestamp,
		SenderAddress:   t.senderAddress[:],
		ReceiverAddress: t.receiverAddress[:],
		Value:           t.value,
		Gas:             t.gas,
	}
}

// Stringer Interface override
func (t *Transaction) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d coins: \n", t.value))
	return sb.String()
}

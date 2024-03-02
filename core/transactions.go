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
	index           int    //8
	value           uint32 //4
	gas             uint32
}

func NewTransaction(i int, v, g uint32, s, r [AddressLength]byte) *Transaction {
	//TODO: Remember in the consensus engine you need to ensure that you're rejecting transaction with duplicates
	trans := &Transaction{
		timestamp:       time.Now().Unix(),
		senderAddress:   s,
		receiverAddress: r,
		value:           v,
		gas:             g,
	}
	//TODO:error handling
	data, _ := proto.Marshal(trans.ConvertToTransactionPbMsg())
	trans.ID = sha3.Sum256(data)
	return trans
}

// Helper method to convert to protoreflect.Protomessage type for pb marshalling
func (t *Transaction) ConvertToTransactionPbMsg() *types.TransactionMsg {
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

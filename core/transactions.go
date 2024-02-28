package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"strings"
	"time"
)

type Transaction struct {
	timestamp       time.Time
	ID              [32]byte
	senderAddress   [AddressLength]byte
	receiverAddress [AddressLength]byte
	index           int
	value           uint32
	gas             uint32
}

func NewTransaction(i int, v, g uint32, s, r [AddressLength]byte) *Transaction {
	//TODO: Remember in the consensus engine you need to ensure that you're rejecting transaction with duplicates
	time := time.Now()
	trans := &Transaction{
		timestamp:       time,
		senderAddress:   s,
		receiverAddress: r,
		value:           v,
		gas:             g,
		index:           i,
	}
	//TODO: implement protobuf Marshalling, json marshalling is too much overhead
	data, _ := json.Marshal(*trans)
	trans.ID = sha3.Sum256(data)
	return trans
}

// Stringer Interface override
func (t *Transaction) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d coins: \n", t.value))
	return sb.String()
}

// We don't need this marshalling function anymore we're going to Marshal to Protobuf
func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress   [AddressLength]byte `json:"sender_address"`
		ReceiverAddress [AddressLength]byte `json:"receiver_address"`
		Value           uint32              `json:"value"`
		Gas             uint32              `json:"gas"`
	}{
		SenderAddress:   t.senderAddress,
		ReceiverAddress: t.receiverAddress,
		Value:           t.value,
		Gas:             t.gas,
	})

}

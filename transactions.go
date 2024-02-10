package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	//AddressLength is the expected length of the address
	addressLength = 20
)

type Transaction struct {
	senderAddress   [addressLength]byte
	receiverAddress [addressLength]byte
	value           uint8
}

func NewTransaction(sender [addressLength]byte, receiver [addressLength]byte, amount uint8) *Transaction {
	return &Transaction{
		senderAddress:   sender,
		receiverAddress: receiver,
		value:           amount,
	}
}

// Stringer Interface override
func (t *Transaction) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d coins: \n", t.value))
	sb.WriteString(fmt.Sprintf("%x -> %x", t.senderAddress, t.receiverAddress))
	return sb.String()
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderAddress   [addressLength]byte `json:"sender_address"`
		ReceiverAddress [addressLength]byte `json:"receiver_address"`
		Value           uint8               `json:"value"`
	}{
		SenderAddress:   t.senderAddress,
		ReceiverAddress: t.receiverAddress,
		Value:           t.value,
	})

}

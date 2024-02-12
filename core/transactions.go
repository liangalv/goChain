package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Transaction struct {
	ID              string
	senderAddress   [AddressLength]byte
	receiverAddress [AddressLength]byte
	value           uint8
	gas             uint8
	index           int
}

func NewTransaction(sender [AddressLength]byte, receiver [AddressLength]byte, amount uint8) *Transaction {
	//TODO: generate ID pseudorandomly and prevent collisions when hashing
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
		SenderAddress   [AddressLength]byte `json:"sender_address"`
		ReceiverAddress [AddressLength]byte `json:"receiver_address"`
		Value           uint8               `json:"value"`
		Gas             uint8               `json: gas`
	}{
		SenderAddress:   t.senderAddress,
		ReceiverAddress: t.receiverAddress,
		Value:           t.value,
		Gas:             t.gas,
	})

}

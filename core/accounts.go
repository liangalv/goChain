package main

import (
	"fmt"

	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/sha3"
)

const (
	//AddressLength is the expected length of the address
	AddressLength = 20
	bitSize       = 256
)

type Account struct {
	address [AddressLength]byte
	pubKey  *bip32.Key
	privKey *bip32.Key
}

func NewAccount(passPhrase string) *Account {
	privKey, pubKey, mnemonic := generateNewPrivPubKeyPair(passPhrase)
	address := deriveAddressFromPubKey(pubKey)
	fmt.Println(mnemonic)
	return &Account{
		address: address,
		privKey: privKey,
		pubKey:  pubKey,
	}
}

// Helper
func generateNewPrivPubKeyPair(passPhrase string) (priv *bip32.Key, pub *bip32.Key, mne string) {
	//Generate a new mneomic with 256 bitsize
	entropy, _ := bip39.NewEntropy(bitSize)
	mnemonic, _ := bip39.NewMnemonic(entropy)

	//Generate a Bip32 HD wallet with mnenomic and passPhrase
	seed := bip39.NewSeed(mnemonic, passPhrase)
	priv, _ = bip32.NewMasterKey(seed)
	pub = priv.PublicKey()

	return priv, pub, mnemonic
}

func deriveAddressFromPubKey(pubKey *bip32.Key) [20]byte {
	hash := sha3.Sum256(pubKey.Key)
	var address [20]byte
	copy(address[:], hash[12:])
	return address
}

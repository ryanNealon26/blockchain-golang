package main

import (
	"math/rand"
)

type Wallet struct {
	User      string
	PublicKey string
	History   []Transaction
}

type Transaction struct {
	PayerKey           string
	PayeeKey           string
	TransactionMessage string
}

// creates public key for sending and recieving messages via the blockchain
func setPublicKey(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

// initialize wallets prior to recieving message
func initWallet(PublicKey string) []Transaction {
	var newWallet []Transaction
	newWallet = append(newWallet, Transaction{PublicKey, "", "No Messages Recieved"})
	return newWallet
}

func errorWallet(id string) []Transaction {
	var newWallet []Transaction
	errorMessage := "No wallet with Public Key " + id + " was found."
	newWallet = append(newWallet, Transaction{"", "", errorMessage})
	return newWallet
}

func saveToWallet(walletList []*Wallet, publicKey string, transaction Transaction) {
	for i := 0; i < len(walletList); i++ {
		if walletList[i].PublicKey == publicKey {
			if len(walletList[i].History) == 1 && walletList[i].History[0].TransactionMessage == "No Messages Recieved" {
				walletList[i].History[0] = transaction
			} else {
				walletList[i].History = append(walletList[i].History, transaction)
			}
		}
	}
}
func searchWallet(walletList []*Wallet, publicKey string) *Wallet {
	for i := 0; i < len(walletList); i++ {
		if walletList[i].PublicKey == publicKey {
			return walletList[i]
		}
	}
	err := Wallet{"Error Wallet", "", errorWallet(publicKey)}
	return &err
}

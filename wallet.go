package main

import (
	"math/rand"
)

type Wallet struct {
	PublicKey string
	History   []Transaction
}

type Transaction struct {
	PayerKey           string
	Payee              string
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

func sendTransaction(senderWallet Wallet, recieverWallet *Wallet, message string) {
	messageRec := Transaction{senderWallet.PublicKey, recieverWallet.PublicKey, message}
	recieverWallet.History = append(recieverWallet.History, messageRec)
}

func saveToWallet(walletList []*Wallet, publicKey string, transaction Transaction) {
	for i := 0; i < len(walletList); i++ {
		if walletList[i].PublicKey == publicKey {
			walletList[i].History = append(walletList[i].History, transaction)
		}
	}
}
func searchWallet(walletList []*Wallet, publicKey string) *Wallet {
	var wallet *Wallet
	for i := 0; i < len(walletList); i++ {
		if walletList[i].PublicKey == publicKey {
			wallet = walletList[i]
		}
	}
	return wallet
}

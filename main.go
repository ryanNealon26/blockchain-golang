package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	hash      []byte
	data      string
	prevHash  []byte
	timeStamp string
}

type BlockChain struct {
	chain []Block
}

func createHash(data string) []byte {
	h := sha256.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	fmt.Printf("String Before Hash: %s\n", data)
	fmt.Printf("Hashed String: %x\n", bs)
	return bs
}
func addBlock(chain *BlockChain, data string) {
	timeStamp := time.Now().Format(time.RFC850)
	block := Block{createHash(data), data, getPrevHash(*chain), timeStamp}
	chain.chain = append(chain.chain, block)
}
func getPrevHash(chain BlockChain) []byte {
	size := len(chain.chain)
	return chain.chain[size-1].hash
}
func initChain() BlockChain {
	var chain BlockChain
	data := "Genesis Block"
	timeStamp := time.Now().Format(time.RFC850)
	genesis := Block{createHash(data), data, []byte{}, timeStamp}
	chain.chain = append(chain.chain, genesis)
	return chain

}
func main() {
	chain := initChain()
	addBlock(&chain, "First Block After Genesis")
	addBlock(&chain, "Second Block After Genesis")
	for i := 0; i < len(chain.chain); i++ {
		fmt.Println("Block #", i+1)
		fmt.Println("	Block Data: ", chain.chain[i].data)
		fmt.Printf("	Previous Hash: %x\n", chain.chain[i].prevHash)
		fmt.Printf("	Current Hash: %x\n", chain.chain[i].hash)
		fmt.Println("	Time Stamp: ", chain.chain[i].timeStamp)
	}
}

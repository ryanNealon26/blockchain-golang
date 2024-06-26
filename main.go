package main

import (
	"crypto/sha256"
	"fmt"
	"time"
)

type Block struct {
	hash      string
	data      string
	prevHash  string
	timeStamp string
}

type BlockChain struct {
	chain []Block
}

func createHash(data string) string {
	h := sha256.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	fmt.Printf("String Before Hash: %s\n", data)
	hexString := fmt.Sprintf("%x", bs)

	fmt.Printf("Hashed String: %s\n", hexString)

	return hexString
}
func addBlock(chain *BlockChain, data string) {
	timeStamp := time.Now().Format(time.RFC850)
	block := Block{createHash(data), data, getPrevHash(*chain), timeStamp}
	chain.chain = append(chain.chain, block)
}
func getPrevHash(chain BlockChain) string {
	size := len(chain.chain)
	return chain.chain[size-1].hash
}
func initChain() BlockChain {
	var chain BlockChain
	data := "Genesis Block"
	timeStamp := time.Now().Format(time.RFC850)
	genesis := Block{createHash(data), data, "", timeStamp}
	chain.chain = append(chain.chain, genesis)
	return chain

}
func isChainValid(chain BlockChain) bool {
	currentBlock := chain.chain
	for i := 1; i < len(currentBlock); i++ {
		next := currentBlock[i]
		if next.prevHash != currentBlock[i-1].hash {
			return false
		}

	}
	return true
}
func main() {
	chain := initChain()
	addBlock(&chain, "First Block After Genesis")
	addBlock(&chain, "Second Block After Genesis")
	if isChainValid(chain) {
		fmt.Println("Chain is valid")
		for i := 0; i < len(chain.chain); i++ {
			fmt.Println("Block #", i+1)
			fmt.Println("	Block Data: ", chain.chain[i].data)
			fmt.Println("	Previous Hash: ", chain.chain[i].prevHash)
			fmt.Println("	Current Hash: ", chain.chain[i].hash)
			fmt.Println("	Time Stamp: ", chain.chain[i].timeStamp)
		}
	} else {
		fmt.Println("Execution halted, invalid block in chain")
	}
}

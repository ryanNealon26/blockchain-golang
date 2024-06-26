package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Block struct {
	Hash      string
	Data      string
	PrevHash  string
	TimeStamp string
}

type BlockChain struct {
	Chain []Block
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
	chain.Chain = append(chain.Chain, block)
}
func getPrevHash(chain BlockChain) string {
	size := len(chain.Chain)
	return chain.Chain[size-1].Hash
}
func initChain() BlockChain {
	var chain BlockChain
	data := "Genesis Block"
	timeStamp := time.Now().Format(time.RFC850)
	genesis := Block{createHash(data), data, "", timeStamp}
	chain.Chain = append(chain.Chain, genesis)
	return chain

}
func isChainValid(chain BlockChain) bool {
	currentBlock := chain.Chain
	for i := 1; i < len(currentBlock); i++ {
		next := currentBlock[i]
		if next.PrevHash != currentBlock[i-1].Hash {
			return false
		}

	}
	return true
}

var chain = initChain()

func getBlockChain(c *gin.Context) {
	if isChainValid(chain) {
		c.IndentedJSON(http.StatusOK, chain.Chain)
	} else {
		c.IndentedJSON(http.StatusOK, "Error Found in blockchain")
	}
}
func postBlock(c *gin.Context) {
	var data struct {
		Data string `json:"data"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	addBlock(&chain, data.Data)
	c.JSON(http.StatusOK, gin.H{"message": "Block added successfully", "data": data.Data})

}
func main() {
	router := gin.Default()
	router.GET("/blockchain", getBlockChain)
	router.POST("/blockchain", postBlock)
	addBlock(&chain, "First Block After Genesis")
	addBlock(&chain, "Second Block After Genesis")
	if isChainValid(chain) {
		fmt.Println("Chain is valid")
		for i := 0; i < len(chain.Chain); i++ {
			fmt.Println("Block #", i+1)
			fmt.Println("	Block Data: ", chain.Chain[i].Data)
			fmt.Println("	Previous Hash: ", chain.Chain[i].PrevHash)
			fmt.Println("	Current Hash: ", chain.Chain[i].Hash)
			fmt.Println("	Time Stamp: ", chain.Chain[i].TimeStamp)
		}
	} else {
		fmt.Println("Execution halted, invalid block in chain")
	}
	router.Run("localhost:8080")
}

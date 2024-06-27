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
	previousHash := getPrevHash(*chain)
	hashString := data + "%" + previousHash + "%" + timeStamp
	block := Block{createHash(hashString), data, previousHash, timeStamp}
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
	hashString := data + "%" + "%" + timeStamp
	genesis := Block{createHash(hashString), data, "", timeStamp}
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
		fmt.Println("Chain is valid")
		for i := 0; i < len(chain.Chain); i++ {
			fmt.Println("Block #", i+1)
			fmt.Println("	Block Data: ", chain.Chain[i].Data)
			fmt.Println("	Previous Hash: ", chain.Chain[i].PrevHash)
			fmt.Println("	Current Hash: ", chain.Chain[i].Hash)
			fmt.Println("	Time Stamp: ", chain.Chain[i].TimeStamp)
		}
		c.IndentedJSON(http.StatusOK, chain.Chain)
	} else {
		fmt.Println("Execution halted, invalid block in chain")
		c.IndentedJSON(http.StatusOK, "Error Found in blockchain")
	}
}

var userKey1 = setPublicKey(25)
var userNum1 = Wallet{userKey1, initWallet(userKey1)}
var userKey2 = setPublicKey(25)
var userNum2 = Wallet{userKey2, initWallet(userKey2)}

// temporary array for storing public keys, will replace with db.
var keyArray = []*Wallet{&userNum1, &userNum2}

func getWallet(c *gin.Context) {
	publicId := c.Param("id")
	wallet := searchWallet(keyArray, publicId)
	c.IndentedJSON(http.StatusOK, wallet.History)
}

func postBlock(c *gin.Context) {
	var data struct {
		Data     string `json:"data"`
		PayerKey string `json:"payerKey"`
		PayeeKey string `json:"payeeKey"`
	}
	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	messageRec := Transaction{data.PayerKey, data.PayeeKey, data.Data}
	saveToWallet(keyArray, data.PayeeKey, messageRec)

	addBlock(&chain, data.Data)
	c.JSON(http.StatusOK, gin.H{"message": "Block added successfully", "data": data.Data})
}
func main() {
	fmt.Println(userKey1)
	fmt.Println(userKey2)
	router := gin.Default()
	router.GET("/blockchain", getBlockChain)
	router.GET("/wallet/:id", getWallet)
	router.POST("/blockchain", postBlock)
	addBlock(&chain, "First Block After Genesis")
	addBlock(&chain, "Second Block After Genesis")
	router.Run("localhost:8080")
}

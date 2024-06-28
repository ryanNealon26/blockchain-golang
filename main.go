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
	Nonce     int
}

type BlockChain struct {
	Chain []Block
}

func createHash(data string, nonce int) string {
	data = data + string(nonce)
	h := sha256.New()
	h.Write([]byte(data))
	bs := h.Sum(nil)
	hexString := fmt.Sprintf("%x", bs)
	return hexString
}
func addBlock(chain *BlockChain, data string) {
	timeStamp := time.Now().Format(time.RFC850)
	previousHash := getPrevHash(*chain)
	hashString := data + "%" + previousHash + "%" + timeStamp + "%"
	block := Block{createHash(hashString, 0), data, previousHash, timeStamp, 0}
	proofOfWork(&block)
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
	genesis := Block{createHash(hashString, 0), data, "", timeStamp, 0}
	chain.Chain = append(chain.Chain, genesis)
	return chain

}
func proofOfWork(block *Block) {
	difficulty := 4
	target := "0000"
	for block.Hash[:difficulty] != target {
		block.Nonce += 1
		hashString := block.Data + "%" + block.PrevHash + "%" + block.TimeStamp + "%"
		block.Hash = createHash(hashString, block.Nonce)
	}
	fmt.Println("Blocked Mined with Nonce: " + string(block.Nonce))
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

// temporary array for storing public keys, will replace with db.
var keyArray []*Wallet

func getWallet(c *gin.Context) {
	publicId := c.Param("id")
	wallet := searchWallet(keyArray, publicId)
	c.IndentedJSON(http.StatusOK, wallet.History)
}

//function for posting user data, initializes crypto wallets

func postUserData(c *gin.Context) {
	var user struct {
		Username string `json:"username"`
	}
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key := setPublicKey(25)
	userData := Wallet{user.Username, key, initWallet(key)}
	keyArray = append(keyArray, &userData)
	fmt.Println("Public Key: " + userData.PublicKey)
	c.JSON(http.StatusOK, gin.H{"message": "UserAdded to blockchain", "data": user.Username})
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
	router := gin.Default()
	router.GET("/blockchain", getBlockChain)
	router.GET("/wallet/:id", getWallet)
	router.POST("/blockchain", postBlock)
	router.POST("/wallet", postUserData)
	addBlock(&chain, "First Block After Genesis")
	addBlock(&chain, "Second Block After Genesis")
	router.Run("localhost:8080")
}

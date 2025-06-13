package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Index     int
	Timestamp string
	Data      string
	PrevHash  string
	Hash      string
	Nonce     int
}

type Blockchain struct {
	Blocks     []*Block
	Difficulty int
}

func CalculateHash(index int, timestamp string, data string, prevhash string, nonce int) string {
	record := strconv.Itoa(index) + timestamp + data + prevhash + strconv.Itoa(nonce)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (b *Block) MineBlock(difficulty int) {
	target := strings.Repeat("0", difficulty)
	for {
		b.Nonce++
		b.Hash = CalculateHash(b.Index, b.Timestamp, b.Data, b.PrevHash, b.Nonce)
		if strings.HasPrefix(b.Hash, target) {
			break
		}
	}
}

func NewBlock(index int, data string, prevHash string, difficulty int) *Block {
	timestamp := time.Now().String()
	block := &Block{
		Index:     index,
		Timestamp: timestamp,
		Data:      data,
		PrevHash:  prevHash,
		Nonce:     0,
	}
	block.MineBlock(difficulty)
	return block
}

func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(prevBlock.Index+1, data, prevBlock.Hash, bc.Difficulty)
	bc.Blocks = append(bc.Blocks, newBlock)
}

func GenesisBlock(difficulty int) *Block {
	timestamp := time.Now().String()
	block := &Block{
		Index:     0,
		Timestamp: timestamp,
		Data:      "Genesis Block",
		PrevHash:  "",
		Nonce:     0,
	}
	block.MineBlock(difficulty)
	return block
}

func NewBlockchain(difficulty int) *Blockchain {
	return &Blockchain{
		Blocks:     []*Block{GenesisBlock(difficulty)},
		Difficulty: difficulty,
	}
}

func (bc *Blockchain) Validate() bool {
	for i := 1; i < len(bc.Blocks); i++ {
		currentBlock := bc.Blocks[i]
		prevBlock := bc.Blocks[i-1]

		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}

		if currentBlock.Hash != CalculateHash(
			currentBlock.Index,
			currentBlock.Timestamp,
			currentBlock.Data,
			currentBlock.PrevHash,
			currentBlock.Nonce,
		) {
			return false
		}
		target := strings.Repeat("0", bc.Difficulty)
		if !strings.HasPrefix(currentBlock.Hash, target) {
			return false
		}
	}

	return true
}

func main() {
	//difficulty = the number of zeros at the beginning
	difficulty := 4
	bc := NewBlockchain(difficulty)

	//Testing------
	numBlocks := 15
	fmt.Printf("Mining %d blocks with difficulty %d...\n", numBlocks, difficulty)

	for i := 0; i < numBlocks; i++ {
		val := strconv.Itoa(i)
		bc.AddBlock(val)
		fmt.Printf("Mined block %d\n", i+1)
	}

	start := len(bc.Blocks) - 5
	if start < 0 {
		start = 0
	}
	//Testing------

	for _, block := range bc.Blocks {
		fmt.Printf("Index: %d\n", block.Index)
		fmt.Printf("Timestamp: %s\n", block.Timestamp)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("PrevHash: %s\n", block.PrevHash)
		fmt.Printf("Hash: %s\n", block.Hash)
		fmt.Printf("Nonce: %d\n\n", block.Nonce)

		if !strings.HasPrefix(block.Hash, strings.Repeat("0", difficulty)) {
			fmt.Printf("WARNING: Hash does not meet difficulty requirement!\n")
		}
	}

	fmt.Println("Blockchain valid:", bc.Validate())
}

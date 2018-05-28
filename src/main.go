package main

import "fmt"
import "strconv"

import . "lib/blockchain"
import . "lib/block"

func main() {
	bc := NewBlockchain()

	bc.AddBlock("Sample Txn 1")
	bc.AddBlock("Sample Txn 2")
	bc.AddBlock("Sample Txn 3")
	bc.AddBlock("Sample Txn 4")
	bc.AddBlock("Sample Txn 5")

	for _, block := range bc.Blocks {
		fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		// validate the nonce
		pow := NewProofOfWork(block)
		fmt.Printf("Valid: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}

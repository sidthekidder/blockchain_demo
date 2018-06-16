package main

import "fmt"
import "strconv"
import "flag"
import "os"

import . "lib/blockchain"
import . "lib/block"

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	addBlockData := addBlockCmd.String("data", "", "Block Data")

	switch os.Args[1] {
		case "addblock":
			err := addBlockCmd.Parse(os.Args[2:])
			fmt.Println(err)
		case "printchain":
			err := printChainCmd.Parse(os.Args[2:])
			fmt.Println(err)
		default:
			cli.printUsage()
			os.Exit(1)
	}
	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}

func (cli *CLI) addBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Println("Success!")
}

func (cli *CLI) printChain() {
	bci := cli.bc.Iterator()

	for {
		block := bci.Next()

		fmt.Println("Previous block hash: %x\n", block.PrevBlockHash)
		fmt.Println("Data: %s\n", block.Data)
		fmt.Println("Hash: %s\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func main() {
	bc := NewBlockchain()
	defer bc.Db.Close()

	cli := CLI{bc}
	cli.Run()

	// bc.AddBlock("Sample Txn 1")
	// bc.AddBlock("Sample Txn 2")
	// bc.AddBlock("Sample Txn 3")
	// bc.AddBlock("Sample Txn 4")
	// bc.AddBlock("Sample Txn 5")

	// for _, block := range bc.Blocks {
	// 	fmt.Printf("Prev hash: %x\n", block.PrevBlockHash)
	// 	fmt.Printf("Data: %s\n", block.Data)
	// 	fmt.Printf("Hash: %x\n", block.Hash)

	// 	// validate the nonce
	// 	pow := NewProofOfWork(block)
	// 	fmt.Printf("Valid: %s\n", strconv.FormatBool(pow.Validate()))
	// 	fmt.Println()
	// }
}

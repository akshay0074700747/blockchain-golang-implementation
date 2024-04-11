package main

import (
	"os"

	"github.com/akshay0074700747/blockchain-GO/chain"
	"github.com/akshay0074700747/blockchain-GO/cli"
)

func main() {
	// chain := chain.InitBlockChain()

	// chain.AddBlock("First Block after Genesis")
	// chain.AddBlock("Second Block after Genesis")
	// chain.AddBlock("Third Block after Genesis")

	// for _, blockk := range chain.Blocks {

	// 	fmt.Printf("Previous Hash: %x\n", blockk.PrevHash)
	// 	fmt.Printf("Data in Block: %s\n", blockk.Data)
	// 	fmt.Printf("Hash: %x\n", blockk.Hash)

	// 	pow := block.NewProofofWork(blockk)
	// 	fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	// 	fmt.Println()

	// }
	defer os.Exit(0)
	chain := chain.InitBlockChain()
	defer chain.DB.Close()

	cli := cli.NewCli(chain)
	cli.Run()
}

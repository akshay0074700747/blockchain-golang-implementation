package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/akshay0074700747/blockchain-GO/block"
	"github.com/akshay0074700747/blockchain-GO/chain"
)

type FuncSig func()

type CommandLine struct {
	Blockchain *chain.BlockChain
	CallMap    map[string]FuncSig
}

func NewCli(chain *chain.BlockChain) *CommandLine {
	return &CommandLine{
		Blockchain: chain,
		CallMap:    make(map[string]FuncSig),
	}
}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println(" add -block BLOCK_DATA - add a block to the chain")
	fmt.Println(" print - Prints the blocks in the chain")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) AddBlock(data string) {
	cli.Blockchain.AddBlock(data)
	fmt.Println("Added Block!")
}

func (cli *CommandLine) PrintChain() {
	iter := cli.Blockchain.Iterator()

	for {
		blockk := iter.Next()

		fmt.Printf("Prev. hash: %x\n", blockk.PrevHash)
		fmt.Printf("Data: %s\n", blockk.Data)
		fmt.Printf("Hash: %x\n", blockk.Hash)
		pow := block.NewProofofWork(blockk)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(blockk.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	addBlockData := addBlockCmd.String("block", "", "Block data")

	cli.CallMap["add"] = func() {
		err := addBlockCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			runtime.Goexit()
		}
		if *addBlockData == "" {
			addBlockCmd.Usage()
			runtime.Goexit()
		}
		cli.AddBlock(*addBlockData)
	}
	cli.CallMap["print"] = func() {
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			runtime.Goexit()
		}
		cli.PrintChain()
	}

	if call := cli.CallMap[os.Args[1]]; call == nil {
		cli.PrintUsage()
		runtime.Goexit()
	} else {
		call()
	}

	// switch os.Args[1] {
	// case "add":
	// 	err := addBlockCmd.Parse(os.Args[2:])
	// 	blockchain.Handle(err)

	// case "print":
	// 	err := printChainCmd.Parse(os.Args[2:])
	// 	blockchain.Handle(err)

	// default:
	// 	cli.printUsage()
	// 	runtime.Goexit()
	// }

	// if addBlockCmd.Parsed() {
	// 	if *addBlockData == "" {
	// 		addBlockCmd.Usage()
	// 		runtime.Goexit()
	// 	}
	// 	cli.addBlock(*addBlockData)
	// }

	// if printChainCmd.Parsed() {
	// 	cli.printChain()
	// }
}

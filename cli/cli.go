package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"

	"github.com/akshay0074700747/blockchain-GO/block"
	"github.com/akshay0074700747/blockchain-GO/chain"
	"github.com/akshay0074700747/blockchain-GO/transactions"
)

type FuncSig func()

type CommandLine struct {
	CallMap map[string]FuncSig
}

func NewCli() *CommandLine {
	return &CommandLine{
		CallMap: make(map[string]FuncSig),
	}
}

func (cli *CommandLine) PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println(" getbalance -address ADDRESS - get the balance for an address")
	fmt.Println(" createblockchain -address ADDRESS creates a blockchain and sends genesis reward to address")
	fmt.Println(" printchain - Prints the blocks in the chain")
	fmt.Println(" send -from FROM -to TO -amount AMOUNT - Send amount of coins")
}

func (cli *CommandLine) ValidateArgs() {
	if len(os.Args) < 2 {
		cli.PrintUsage()
		runtime.Goexit()
	}
}

func (cli *CommandLine) printChain() {
	chain := chain.ContinueBlockChain("")
	defer chain.DB.Close()
	iter := chain.Iterator()

	for {
		blockk := iter.Next()

		fmt.Printf("Prev. hash: %x\n", blockk.PrevHash)
		fmt.Printf("Hash: %x\n", blockk.Hash)
		pow := block.NewProofofWork(blockk)
		fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(blockk.PrevHash) == 0 {
			break
		}
	}
}

func (cli *CommandLine) createBlockChain(address string) {
	chain := chain.InitBlockChain(address)
	chain.DB.Close()
	fmt.Println("Finished!")
}

func (cli *CommandLine) getBalance(address string) {
	chain := chain.ContinueBlockChain(address)
	defer chain.DB.Close()

	balance := 0
	UTXOs := chain.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s: %d\n", address, balance)
}

func (cli *CommandLine) send(from, to string, amount int) {
	chainn := chain.ContinueBlockChain(from)
	defer chainn.DB.Close()

	tx := chain.NewTransaction(from, to, amount, chainn)
	chainn.AddBlock([]*transactions.Transaction{tx})
	fmt.Println("Success!")
}

func (cli *CommandLine) Run() {
	cli.ValidateArgs()

	// addBlockCmd := flag.NewFlagSet("add", flag.ExitOnError)
	// printChainCmd := flag.NewFlagSet("print", flag.ExitOnError)
	// addBlockData := addBlockCmd.String("block", "", "Block data")
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	getBalanceAddress := getBalanceCmd.String("address", "", "The address to get balance for")
	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send genesis block reward to")
	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	cli.CallMap["getbalance"] = func() {
		err := getBalanceCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			runtime.Goexit()
		}
		if *getBalanceAddress == "" {
			getBalanceCmd.Usage()
			runtime.Goexit()
		}
		cli.getBalance(*getBalanceAddress)
	}
	cli.CallMap["printchain"] = func() {
		err := printChainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			runtime.Goexit()
		}
		cli.printChain()
	}
	cli.CallMap["createblockchain"] = func() {
		err := createBlockchainCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			runtime.Goexit()
		}
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			runtime.Goexit()
		}
		cli.createBlockChain(*createBlockchainAddress)
	}
	cli.CallMap["send"] = func() {
		err := sendCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Println(err)
			runtime.Goexit()
		}
		if *sendFrom == "" {
			sendCmd.Usage()
			runtime.Goexit()
		}
		if *sendTo == "" {
			sendCmd.Usage()
			runtime.Goexit()
		}
		if *sendAmount == 0 {
			sendCmd.Usage()
			runtime.Goexit()
		}
		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

	if call := cli.CallMap[os.Args[1]]; call == nil {
		cli.PrintUsage()
		runtime.Goexit()
	} else {
		call()
	}
}

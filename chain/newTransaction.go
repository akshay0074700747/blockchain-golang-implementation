package chain

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/akshay0074700747/blockchain-GO/transactions"
)

func NewTransaction(from, to string, amount int, chain *BlockChain) *transactions.Transaction {
	var inputs []transactions.TxInput
	var outputs []transactions.TxOutput

	acc, validOutputs := chain.FindSpendableOutputs(from, amount)

	if acc < amount {
		fmt.Println("the account balance is less than the amount you want to send")
		os.Exit(1)
	}

	for txid, outs := range validOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			fmt.Println(err)
		}

		for _, out := range outs {
			inputs = append(inputs, transactions.TxInput{
				ID:  txID,
				Out: out,
				Sig: from,
			})
		}
	}

	outputs = append(outputs, transactions.TxOutput{
		Value:  amount,
		PubKey: to,
	})

	if acc > amount {
		outputs = append(outputs, transactions.TxOutput{
			Value:  acc - amount,
			PubKey: from,
		})
	}

	tx := transactions.Transaction{
		ID: nil,
		Inputs: inputs,
		Outputs: outputs,
	}
	tx.SetID()

	return &tx
}

package transactions

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"os"
)

type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

type TxOutput struct {
	Value  int
	PubKey string
}

type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

func CoinBaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Coins to %s", to)
	}
	tx := &Transaction{
		ID: nil,
		Inputs: []TxInput{
			{ID: []byte{},
				Out: -1,
				Sig: data},
		},
		Outputs: []TxOutput{
			{Value: 100,
				PubKey: to},
		},
	}
	if err := tx.SetID(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return tx
}

func (tx *Transaction) SetID() error {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	var hash [32]byte

	if err := encoder.Encode(tx); err != nil {
		return err
	}

	hash = sha256.Sum256(buff.Bytes())
	tx.ID = hash[:]
	return nil
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Sig == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PubKey == data
}

//--------------------------------------<> NewTransaction is moved to the chain package due to the cyclic dependency issues

// func NewTransaction()  {
	
// }
//-------------------------------------------------------------------------------------------------------------------------


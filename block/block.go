package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"

	"github.com/akshay0074700747/blockchain-GO/transactions"
)

type Block struct {
	Data     []*transactions.Transaction
	PrevHash []byte
	Hash     []byte
	Nonce    int
}

// func (b *Block) DeriveHash() {
// 	summedHash := bytes.Join([][]byte{b.Data,b.PrevHash},[]byte{})
// 	hash := sha256.Sum256(summedHash)
// 	b.Hash = hash[:]
// }

func NewBlock(data []*transactions.Transaction, prevHash []byte) *Block {
	block := new(Block)
	block.Data = data
	block.PrevHash = prevHash
	pow := NewProofofWork(block)
	nonce, hash := pow.Compute()
	block.Nonce, block.Hash = nonce, hash[:]
	return block
}

func Genesis(coinbase *transactions.Transaction) *Block {
	return NewBlock([]*transactions.Transaction{coinbase}, []byte{})
}

func (b *Block) Serialize() ([]byte, error) {

	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	if err := encoder.Encode(b); err != nil {
		return nil, err
	}
	return buff.Bytes(), nil

}

func (b *Block) HashTransactions() []byte {
	var res [32]byte
	var ress [][]byte

	for _, v := range b.Data {
		ress = append(ress, v.ID)
	}

	res = sha256.Sum256(bytes.Join(ress, []byte{}))
	return res[:]
}

func Deserialize(data []byte) (*Block, error) {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(data))
	if err := decoder.Decode(&block); err != nil {
		return nil, err
	}
	return &block, nil
}

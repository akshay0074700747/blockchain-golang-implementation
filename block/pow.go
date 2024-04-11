package block

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"os"

	"github.com/akshay0074700747/blockchain-GO/util"
)

const (
	Difficulty = 18
)

type ProofofWork struct {
	Block  *Block
	Target *big.Int
}

func NewProofofWork(Block *Block) *ProofofWork {
	pow := new(ProofofWork)
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Difficulty))
	pow.Target = target
	pow.Block = Block
	return pow
}

func (pow *ProofofWork) InitData(nonce int) []byte {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			os.Exit(1)
		}
	}()
	nonceBytes, err := util.ToHex(int64(nonce))
	if err != nil {
		panic(err)
	}
	difficultyBytes, err := util.ToHex(Difficulty)
	if err != nil {
		panic(err)
	}
	data := bytes.Join([][]byte{pow.Block.HashTransactions()[:], pow.Block.PrevHash, nonceBytes, difficultyBytes}, []byte{})
	return data
}

func (pow *ProofofWork) Compute() (int, []byte) {
	var intHash big.Int
	var data []byte
	var hash [32]byte
	nonce := 0
	for nonce < math.MaxInt64 {
		data = pow.InitData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		intHash.SetBytes(hash[:])

		if intHash.Cmp(pow.Target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofofWork) Validate() bool {
	var intHash big.Int
	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])
	return intHash.Cmp(pow.Target) == -1
}

package chain

import (
	"encoding/hex"
	"fmt"
	"os"
	"runtime"

	"github.com/akshay0074700747/blockchain-GO/block"
	"github.com/akshay0074700747/blockchain-GO/transactions"
	"github.com/dgraph-io/badger"
)

const (
	DBpath      = "/tmp/blocks"
	DBFile      = "/tmp/blocks/MANIFEST"
	GenesisData = "First Transaction from Genesis"
)

type BlockChain struct {
	// Blocks []*block.Block
	LastHash []byte
	DB       *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func BadgerExists() bool {
	if _, err := os.Stat(DBFile); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		os.Exit(1)
	}
	return true
}

func NewBlockChain() *BlockChain {
	return new(BlockChain)
}

func (chain *BlockChain) AddBlock(data []*transactions.Transaction) {
	// lastBlock := chain.Blocks[len(chain.Blocks)-1]
	// chain.Blocks = append(chain.Blocks, block.NewBlock([]byte(data), lastBlock.Hash))
	if chain.LastHash == nil || len(chain.LastHash) == 0 {
		chain.DB.View(func(txn *badger.Txn) error {
			item, err := txn.Get([]byte("lh"))
			if err != nil {
				fmt.Println("error at getting the lh...", err)
				return err
			}
			chain.LastHash, err = item.ValueCopy(nil)
			if err != nil {
				fmt.Println("error at getting value from item...", err)
				return err
			}
			return nil
		})
	}

	blockk := block.NewBlock(data, chain.LastHash)
	serialized, err := blockk.Serialize()
	if err != nil {
		fmt.Println("error at serializing...", err)
		os.Exit(1)
	}
	chain.DB.Update(func(txn *badger.Txn) error {
		if err = txn.Set(blockk.Hash, serialized); err != nil {
			fmt.Println("error happened at setting the block...", err)
			return err
		}
		if err = txn.Set([]byte("lh"), blockk.Hash); err != nil {
			fmt.Println("error happened at setting the lh...", err)
			return err
		}
		return nil
	})

	chain.LastHash = blockk.Hash
}

func InitBlockChain(address string) *BlockChain {
	// return &BlockChain{Blocks: []*block.Block{block.Genesis()}}

	if BadgerExists() {
		fmt.Println("BlockChain Already Exists...")
		runtime.Goexit()
	}

	var lastHash []byte
	db, err := OpenBadger()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = db.Update(func(txn *badger.Txn) error {
		// if item, err := txn.Get([]byte("lh")); err != nil {
		// 	if err == badger.ErrKeyNotFound {
		// 		fmt.Println("No Existing BlockChain found...")
		// 		genesis := block.Genesis()
		// 		serialized, err := genesis.Serialize()
		// 		if err != nil {
		// 			fmt.Println("error upon serializing the genesis block...", err.Error())
		// 			return err
		// 		}
		// 		fmt.Println("Creating a BlockChain in DB...")
		// 		if err = txn.Set(genesis.Hash, serialized); err != nil {
		// 			fmt.Println("error upon setting the transaction...", err.Error())
		// 			return err
		// 		}
		// 		if err = txn.Set([]byte("lh"), genesis.Hash); err != nil {
		// 			fmt.Println("error upon setting the lh..", err.Error())
		// 			return err
		// 		}
		// 		lastHash = genesis.Hash
		// 	} else {
		// 		return err
		// 	}
		// } else {
		// 	lastHash, err = item.ValueCopy(nil)
		// 	if err != nil {
		// 		fmt.Println("error at getting value from item...", err)
		// 		return err
		// 	}
		// }

		cbtx := transactions.CoinBaseTx(address, GenesisData)
		genesis := block.Genesis(cbtx)
		fmt.Println("Genesis created")
		serialized, err := genesis.Serialize()
		if err != nil {
			fmt.Println("error at serializing genesis...", err)
			return err
		}
		if err = txn.Set(genesis.Hash, serialized); err != nil {
			fmt.Println("error at setting hash genesis...", err)
			return err
		}
		if err = txn.Set([]byte("lh"), genesis.Hash); err != nil {
			fmt.Println("error at setting lh...", err)
			return err
		}

		lastHash = genesis.Hash

		return nil

	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &BlockChain{
		DB:       db,
		LastHash: lastHash,
	}
}

func ContinueBlockChain(address string) *BlockChain {
	if !BadgerExists() {
		fmt.Println("No existing blockchain found, create one!")
		runtime.Goexit()
	}

	var lastHash []byte

	db, err := OpenBadger()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		if err != nil {
			fmt.Println("error at getting lh...", err)
			return err
		}
		lastHash, err = item.ValueCopy(nil)
		if err != nil {
			fmt.Println("error at getting hash from item...", err)
			return err
		}

		return nil
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &BlockChain{
		DB:       db,
		LastHash: lastHash,
	}
}

func OpenBadger() (*badger.DB, error) {
	opts := badger.DefaultOptions(DBpath)
	return badger.Open(opts)
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.DB}

	return iter
}

func (iter *BlockChainIterator) Next() *block.Block {
	var blockk *block.Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		if err != nil {
			fmt.Println(err)
		}
		encodedBlock, err := item.ValueCopy(nil)
		if err != nil {
			fmt.Println(err)
		}
		blockk, err = block.Deserialize(encodedBlock)
		if err != nil {
			fmt.Println(err)
		}
		return err
	})
	if err != nil {
		os.Exit(1)
	}

	iter.CurrentHash = blockk.PrevHash

	return blockk
}

func (chain *BlockChain) FindUnspentTransactions(address string) []transactions.Transaction {
	var unspendTxs []transactions.Transaction
	var spendTxs map[string][]int

	iter := chain.Iterator()
	for {
		blockk := iter.Next()
		for _, transaction := range blockk.Data {
			txID := hex.EncodeToString(transaction.ID)
			Outputs:
			for outID, out := range transaction.Outputs {
				if spendTxs[txID] != nil {
					for _, spentOut := range spendTxs[txID] {
						if spentOut == outID {
							continue Outputs
						}
					}
				}
				if out.CanBeUnlocked(address) {
					unspendTxs = append(unspendTxs, *transaction)
				}
			}
			if !transaction.IsCoinbase() {
				for _, in := range transaction.Inputs {
					if in.CanUnlock(address) {
						inTxID := hex.EncodeToString(in.ID)
						spendTxs[inTxID] = append(spendTxs[inTxID], in.Out)
					}
				}
			}
		}
		if len(blockk.PrevHash) == 0 {
			break
		}
	}
	return unspendTxs
}

func (chain *BlockChain) FindUTXO(address string) []transactions.TxOutput {
	var UTXOs []transactions.TxOutput
	unspentTransactions := chain.FindUnspentTransactions(address)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Outputs {
			if out.CanBeUnlocked(address) {
				UTXOs = append(UTXOs, out)
			}
		}
	}
	return UTXOs
}

func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOuts := make(map[string][]int)
	unspentTxs := chain.FindUnspentTransactions(address)
	accumulated := 0

Work:
	for _, tx := range unspentTxs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Outputs {
			if out.CanBeUnlocked(address) && accumulated < amount {
				accumulated += out.Value
				unspentOuts[txID] = append(unspentOuts[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOuts
}
package blockchain

import . "lib/block"
import "github.com/boltdb/bolt"
import "log"

type Blockchain struct {
	tip []byte
	Db *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	Db 			*bolt.DB
}

const dbFile = "blockchain.db"
const blocksBucket = "blocks"

func (bc *Blockchain) AddBlock(data string) {
	var prevHash []byte

	bc.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		prevHash = b.Get([]byte("l"))

		return nil
	})

	newBlock := NewBlock(data, prevHash)

	bc.Db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		b.Put(newBlock.Hash, newBlock.Serialize())
		b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genesis := NewGenesisBlock()
			b, _ := tx.CreateBucket([]byte(blocksBucket))

			b.Put(genesis.Hash, genesis.Serialize())
			b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.tip, bc.Db}
	return bci
}

func (bci *BlockchainIterator) Next() *Block {
	var block *Block
	bci.Db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		encodedBlock := b.Get(bci.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	bci.currentHash = block.PrevBlockHash

	return block
}

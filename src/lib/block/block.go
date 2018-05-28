package block

import (
	"time"
	"bytes"
	"encoding/binary"
	"log"
	"crypto/sha256"
	"math/big"
	"math"
	"fmt"
)

const targetBits = 20

type Block struct {
	Timestamp			int64
	Data 				[]byte
	PrevBlockHash		[]byte
	Hash 				[]byte
	Nonce 				int
}

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), []byte(prevBlockHash), []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing %s \n", pow.block.Data)
	count := 0
	for nonce < math.MaxInt64 {
		count += 1
		if count == 200 {
			fmt.Printf("Tried %d hashes", count)
			fmt.Printf("TRYING HASH!! %x\n", hashInt)
		}
		data := bytes.Join(
				[][]byte{
					pow.block.PrevBlockHash,
					pow.block.Data,
					IntToHex(pow.block.Timestamp),
					IntToHex(int64(targetBits)),
					IntToHex(int64(nonce)),
				},
				[]byte{},
		)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}
	fmt.Printf("Mining Done.\n")

	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := bytes.Join(
			[][]byte{
				pow.block.PrevBlockHash,
				pow.block.Data,
				IntToHex(pow.block.Timestamp),
				IntToHex(int64(targetBits)),
				IntToHex(int64(pow.block.Nonce)),
			},
			[]byte{},
	)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}


// IntToHex converts an int64 to a byte array
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

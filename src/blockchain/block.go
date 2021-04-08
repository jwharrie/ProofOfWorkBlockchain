package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
)

type Block struct {
	PrevHash   []byte
	Generation uint64
	Difficulty uint8
	Data       string
	Proof      uint64
	Hash       []byte
}

// Create new initial (generation 0) block.
func Initial(difficulty uint8) Block {
	var prevHash string
	for i := 0; i < 32; i++ {
		prevHash += "\x00"
	}
	return Block{[]byte(prevHash), uint64(0), difficulty, "", uint64(0), make([]byte, 32)}
}

// Create new block to follow this block, with provided data.
func (prev_block Block) Next(data string) Block {
	return Block{prev_block.Hash, prev_block.Generation + uint64(1), prev_block.Difficulty, data, uint64(0), make([]byte, 32)}
}

// Calculate the block's hash.
func (blk Block) CalcHash() []byte {
	hash := hex.EncodeToString(blk.PrevHash)
	hash += ":" + strconv.FormatUint(blk.Generation, 10)
	hash += ":" + strconv.FormatUint(uint64(blk.Difficulty), 10)
	hash += ":" + blk.Data
	hash += ":" + strconv.FormatUint(blk.Proof, 10)
	sum := sha256.Sum256([]byte(hash))
	return sum[:]
}

// Is this block's hash valid?
func (blk Block) ValidHash() bool {
	nBytes := blk.Difficulty / 8
	nBits := blk.Difficulty % 8
	var currByte uint8 = 31
	for currByte >= (32 - nBytes) {
		if (blk.Hash[currByte]) != byte(0) {
			return false
		}
		currByte--
	}
	if blk.Hash[currByte]%(1<<nBits) != 0 || blk.Hash[currByte-1] == byte(0) {
		return false
	}
	return true
}

// Set the proof-of-work and calculate the block's "true" hash.
func (blk *Block) SetProof(proof uint64) {
	blk.Proof = proof
	blk.Hash = blk.CalcHash()
}

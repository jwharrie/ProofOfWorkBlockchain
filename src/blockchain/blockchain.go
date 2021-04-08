package blockchain

import(
	"bytes"
)

type Blockchain struct {
	Chain []Block
}

func (chain *Blockchain) Add(blk Block) {
	chain.Chain = append(chain.Chain, blk)
}

func (chain Blockchain) IsValid() bool {
	var chainDifficulty uint8
	var prevBlkHash []byte
	for i, blk := range chain.Chain {
		if i == 0 {
			var nullByte string
			for i := 0; i < 32; i++ {
				nullByte += "\x00"
			}
			if blk.Generation != 0 || bytes.Compare(blk.PrevHash, []byte(nullByte)) != 0 {
				return false
			}
			chainDifficulty = blk.Difficulty
			prevBlkHash = blk.Hash
			continue
		}
		if blk.Generation != uint64(i) {
			return false
		}
		if blk.Difficulty != chainDifficulty {
			return false
		}
		if bytes.Compare(blk.PrevHash, prevBlkHash) != 0 {
			return false
		}
		if bytes.Compare(blk.Hash, blk.CalcHash()) != 0 {
			return false
		}
		if !blk.ValidHash() {
			return false
		}
		prevBlkHash = blk.Hash
	}
	return true
}

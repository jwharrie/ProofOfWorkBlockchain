package blockchain

import (
	"work_queue"
)

type miningWorker struct {
	blk Block
	startValue uint64
	endValue uint64
}

type MiningResult struct {
	Proof uint64 // proof-of-work value, if found.
	Found bool   // true if valid proof-of-work was found.
}

func (mw miningWorker) Run() interface{} {
	for i := mw.startValue; i <= mw.endValue; i++ {
		mw.blk.SetProof(i)
		if mw.blk.ValidHash() {
			return MiningResult{i, true}
		}
	}
	return MiningResult{Found: false}
}

// Mine the range of proof values, by breaking up into chunks and checking
// "workers" chunks concurrently in a work queue. Should return shortly after a result
// is found.
func (blk Block) MineRange(start uint64, end uint64, workers uint64, chunks uint64) MiningResult {
	q := work_queue.Create(uint(workers), uint(chunks))
	chunkSize := (end - start) / chunks
	if chunkSize == 0 {
		chunkSize = 1
	}
	prevVal := start
	currVal := prevVal + chunkSize
	chunksSent := 0
	for currVal < end {
		q.Enqueue(miningWorker{blk, prevVal, currVal})
		prevVal = currVal + 1
		currVal = prevVal + chunkSize
		chunksSent++
	}
	q.Enqueue(miningWorker{blk, prevVal, end})
	chunksSent++
	for i := 0; i < chunksSent; i++ {
		r := <- q.Results
		mr := r.(MiningResult)
		if mr.Found {
			q.Shutdown()
			return mr
		}
	}
	q.Shutdown()
	return MiningResult{Found: false}
}

// Call .MineRange with some reasonable values that will probably find a result.
// Good enough for testing at least. Updates the block's .Proof and .Hash if successful.
func (blk *Block) Mine(workers uint64) bool {
	reasonableRangeEnd := uint64(4 * 1 << blk.Difficulty) // 4 * 2^(bits that must be zero)
	mr := blk.MineRange(0, reasonableRangeEnd, workers, 4321)
	if mr.Found {
		blk.SetProof(mr.Proof)
	}
	return mr.Found
}


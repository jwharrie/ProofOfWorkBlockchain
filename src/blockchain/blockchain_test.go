package blockchain

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"encoding/hex"
)

func TestBlockInitial(t *testing.T) {
	var prevHash string
	for i := 0; i < 32; i++ {
		prevHash += "\x00"
	}
	b := Block{[]byte(prevHash), uint64(0), uint8(16), "", uint64(0), make([]byte, 32)}
	b0 := Initial(16)
	if !reflect.DeepEqual(b, b0) {
		t.Error("Block.Initial failed.")
	}
}

func TestBlockNext(t *testing.T) {
	b0 := Initial(16)
	b := Block{b0.Hash, uint64(1), uint8(16), "message", uint64(0), make([]byte, 32)}
	b1 := b0.Next("message")
	if !reflect.DeepEqual(b, b1) {
		t.Error("Block.Next failed.")
	}
}


func TestBlockCalcHash(t *testing.T) {
	b0 := Initial(16)
	b0.SetProof(56231)
	b0Expected := "6c71ff02a08a22309b7dbbcee45d291d4ce955caa32031c50d941e3e9dbd0000"
	b0Actual := hex.EncodeToString(b0.CalcHash())
	assert.Equal(t, b0Expected, b0Actual)

	b1 := b0.Next("message")
	b1.SetProof(2159)
	b1Expected := "9b4417b36afa6d31c728eed7abc14dd84468fdb055d8f3cbe308b0179df40000"
	b1Actual := hex.EncodeToString(b1.CalcHash())
	assert.Equal(t, b1Expected, b1Actual)
}

func TestValidHash(t *testing.T) {
	b0 := Initial(16)
	assert.False(t, b0.ValidHash())
	b0.SetProof(56231)
	assert.True(t, b0.ValidHash())

	b0 = Initial(19)
	b0.SetProof(87745)
	b1 := b0.Next("hash example 1234")
	b1.SetProof(1407891)
	assert.True(t, b1.ValidHash())
	b1.SetProof(346082)
	assert.False(t, b1.ValidHash())
}

func TestMiningDifficulty7(t *testing.T) {
	b0 := Initial(7)
	b0.Mine(1)
	assert.Equal(t, uint64(385), b0.Proof)
	assert.Equal(t, "379bf2fb1a558872f09442a45e300e72f00f03f2c6f4dd29971f67ea4f3d5300", hex.EncodeToString(b0.Hash))

	b1 := b0.Next("this is an interesting message")
	b1.Mine(1)
	assert.Equal(t, uint64(20), b1.Proof)
	assert.Equal(t, "4a1c722d8021346fa2f440d7f0bbaa585e632f68fd20fed812fc944613b92500", hex.EncodeToString(b1.Hash))

	b2 := b1.Next("this is not interesting")
	b2.Mine(1)
	assert.Equal(t, uint64(40), b2.Proof)
	assert.Equal(t, "ba2f9bf0f9ec629db726f1a5fe7312eb76270459e3f5bfdc4e213df9e47cd380", hex.EncodeToString(b2.Hash))
}

func TestMiningDifficulty20(t *testing.T) {
	b0 := Initial(20)
	b0.Mine(1)
	assert.Equal(t, uint64(1209938), b0.Proof)
	assert.Equal(t, "19e2d3b3f0e2ebda3891979d76f957a5d51e1ba0b43f4296d8fb37c470600000", hex.EncodeToString(b0.Hash))

	b1 := b0.Next("this is an interesting message")
	b1.Mine(1)
	assert.Equal(t, uint64(989099), b1.Proof)
	assert.Equal(t, "a42b7e319ee2dee845f1eb842c31dac60a94c04432319638ec1b9f989d000000", hex.EncodeToString(b1.Hash))

	b2 := b1.Next("this is not interesting")
	b2.Mine(1)
	assert.Equal(t, uint64(1017262), b2.Proof)
	assert.Equal(t, "6c589f7a3d2df217fdb39cd969006bc8651a0a3251ffb50470cbc9a0e4d00000", hex.EncodeToString(b2.Hash))
}

func TestMiningMultipleWorkers(t *testing.T) {
	b0 := Initial(20)
	b0.Mine(2)
	assert.Equal(t, uint64(1209938), b0.Proof)
	assert.Equal(t, "19e2d3b3f0e2ebda3891979d76f957a5d51e1ba0b43f4296d8fb37c470600000", hex.EncodeToString(b0.Hash))

	b1 := Initial(20)
	b1.Mine(4)
	assert.Equal(t, uint64(1209938), b1.Proof)
	assert.Equal(t, "19e2d3b3f0e2ebda3891979d76f957a5d51e1ba0b43f4296d8fb37c470600000", hex.EncodeToString(b1.Hash))

	b2 := Initial(20)
	b2.Mine(10)
	assert.Equal(t, uint64(1209938), b2.Proof)
	assert.Equal(t, "19e2d3b3f0e2ebda3891979d76f957a5d51e1ba0b43f4296d8fb37c470600000", hex.EncodeToString(b2.Hash))

	b3 := Initial(20)
	b3.Mine(100)
	assert.Equal(t, uint64(1209938), b3.Proof)
	assert.Equal(t, "19e2d3b3f0e2ebda3891979d76f957a5d51e1ba0b43f4296d8fb37c470600000", hex.EncodeToString(b3.Hash))
}

func TestBlockChainNormal(t *testing.T) {
	var testChain Blockchain
	b0 := Initial(20)
	b0.Mine(10)
	testChain.Add(b0)
	assert.True(t, testChain.IsValid())
	
	b1 := b0.Next("this is an interesting message")
	b1.Mine(10)
	testChain.Add(b1)
	assert.True(t, testChain.IsValid())

	b2 := b1.Next("this is not interesting")
	b2.Mine(10)
	testChain.Add(b2)
	assert.True(t, testChain.IsValid())
}

func TestBlockChainFailure(t *testing.T) {
	var testChain Blockchain
	b0 := Initial(20)
	b0.Mine(10)
	testChain.Add(b0)

	b01 := Initial(7)
	b01.Mine(10)
	b1 := b01.Next("this is an interesting message")
	b1.Mine(10)
	b2 := b1.Next("this is not interesting")
	b2.Mine(10)

	testChain.Add(b2)
	assert.False(t, testChain.IsValid())
}


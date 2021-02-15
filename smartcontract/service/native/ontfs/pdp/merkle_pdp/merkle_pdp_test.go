
package merkle_pdp

import (
	"crypto/rand"
	"testing"

	" github.com/Daironode/aingle/smartcontract/service/native/ontfs/pdp/types"
)

func TestMerkleProof(t *testing.T) {
	var blocks []types.Block
	data := make([]byte, 256*1024)
	count := 128
	for i := 0; i < count; i++ {
		rand.Read(data)
		blocks = append(blocks, data)
	}

	rootHash, err := CalcRootHash(blocks)
	if err != nil {
		t.Fatal(err.Error())
	}

	t.Logf("rootHash: %v", rootHash)

	prf, err := MerkleProof(blocks, 1)
	if err != nil {
		t.Fatal(err.Error())
	}
	var prfLen int
	for _, v := range prf {
		prfLen += len(v)
	}

	t.Logf("prfLen: %v", len(prf))
	t.Logf("prfTotalLength: %v", prfLen)
	t.Logf("prf: %v", prf)

	if err := VerifyMerkleProof(prf, rootHash, 1); err != nil {
		t.Fatal(err.Error())
	}
}

func BenchmarkHash(b *testing.B) {
	var blocks []types.Block
	data := make([]byte, 256*1024)
	count := 1024
	for i := 0; i < count; i++ {
		rand.Read(data)
		blocks = append(blocks, data)
	}

	for i := 0; i < b.N; i++ {
		MerkleProof(blocks, 10)
	}
}

func BenchmarkVerify(b *testing.B) {
	var blocks []types.Block
	data := make([]byte, 256*1024)
	count := 1024
	for i := 0; i < count; i++ {
		rand.Read(data)
		blocks = append(blocks, data)
	}
	rootHash, err := CalcRootHash(blocks)
	if err != nil {
		b.Fatal(err.Error())
	}

	proof, err := MerkleProof(blocks, 10)
	if err != nil {
		b.Fatal(err.Error())
	}
	for i := 0; i < b.N; i++ {
		if err := VerifyMerkleProof(proof, rootHash, 10); err != nil {
			b.Fatal(err.Error())
		}
	}
}

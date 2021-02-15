 
package auth

import (
	"bytes"
	"testing"

	" github.com/Daironode/aingle/common"
)

func TestSerRoleFuncs(t *testing.T) {
	param := &roleFuncs{
		[]string{"foo1", "foo2"},
		//[]string{},
	}
	bf := common.NewZeroCopySink(nil)
	param.Serialization(bf)
	rd := common.NewZeroCopySource(bf.Bytes())
	param2 := new(roleFuncs)
	if err := param2.Deserialization(rd); err != nil {
		t.Fatal(err)
	}

	if len(param.funcNames) != len(param2.funcNames) {
		t.Fatalf("does not match")
	}
	for i := 0; i < len(param.funcNames); i++ {
		if param.funcNames[i] != param2.funcNames[i] {
			t.Fatalf("%s \t %s does not match", param.funcNames[i], param2.funcNames[i])
		}
	}
}

func TestSerAuthToken(t *testing.T) {
	param := &AuthToken{
		role:       []byte("role"),
		expireTime: 1000000,
		level:      2,
	}

	bf := common.NewZeroCopySink(nil)
	param.Serialization(bf)
	rd := common.NewZeroCopySource(bf.Bytes())
	param2 := new(AuthToken)
	if err := param2.Deserialization(rd); err != nil {
		t.Fatal(err)
	}

	if param.expireTime != param2.expireTime ||
		param.level != param2.level ||
		bytes.Compare(param.role, param2.role) != 0 {
		t.Fatalf("failed")
	}
}

func TestSerDelegateStatus(t *testing.T) {
	token := &AuthToken{
		role:       []byte("role"),
		expireTime: 1000000,
		level:      2,
	}
	s1 := &DelegateStatus{
		root:      []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		AuthToken: *token,
	}
	bf := common.NewZeroCopySink(nil)
	s1.Serialization(bf)
	rd := common.NewZeroCopySource(bf.Bytes())
	s2 := new(DelegateStatus)
	if err := s2.Deserialization(rd); err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(s1.root, s2.root) != 0 ||
		bytes.Compare(s1.role, s2.role) != 0 ||
		s1.expireTime != s2.expireTime || s1.level != s2.level {
		t.Fatalf("failed")
	}
}

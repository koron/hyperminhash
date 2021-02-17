package hyperminhash

import (
	"strconv"
	"testing"
)

func TestMarshal(t *testing.T) {
	sk := New()
	const N = 1000000
	for i := 0; i < N; i++ {
		sk.Add([]byte(strconv.Itoa(i)))
	}
	b, err := sk.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal failed: %s", err)
	}
	sk2 := New()
	if err := sk2.UnmarshalBinary(b); err != nil {
		t.Fatalf("unmarshal failed: %s", err)
	}
	exp := sk.Cardinality()
	act := sk2.Cardinality()
	if act != exp {
		t.Errorf("cardinalities unmatch: want=%d got=%d", exp, act)
	}
}

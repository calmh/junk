package prototest

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestMarshal(t *testing.T) {
	msg := Test{
		Indexes: []int32{1, 2, 3},
	}
	bs, err := msg.Marshal()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(hex.Dump(bs))

	var dst Test
	err = dst.Unmarshal(bs)
	if err != nil {
		t.Fatal(err)
	}
}

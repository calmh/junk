package jsonstream

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestStream(t *testing.T) {
	data := []byte(`[{"v":0,"foo":"bar,bac"},{"v":1},{"v":2},{"v":3},{"v":4}]`)
	str := New(bytes.NewReader(data))
	dec := json.NewDecoder(str)
	var res struct{ V int }
	for i := 0; i < 5; i++ {
		err := dec.Decode(&res)
		if err != nil {
			t.Fatal(err)
		}
		if res.V != i {
			t.Errorf("%d != %d", res.V, i)
		}
	}
}

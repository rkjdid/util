package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

type testFloatStruct struct {
	Float Float
}

func TestFloat_MarshalJSON(t *testing.T) {
	f := Float(1.0)
	expect := fmt.Sprintf("\"%v\"", 1.0)

	b, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}

	res := string(b)
	if expect != res {
		t.Errorf("json encode expected: %s, got: %s", expect, res)
	}
}

func TestFloat_UnmashalJSON(t *testing.T) {
	expect := Float(1.0)
	b := new(bytes.Buffer)
	b.WriteString("1.0")
	enc := json.NewDecoder(b)
	var f Float
	err := enc.Decode(&f)
	if err != nil {
		t.Fatal(err)
	}

	if expect != f {
		t.Fatalf("json decode expected: %s, got: %s", expect, f)
	}

	expect = Float(1)
	b = new(bytes.Buffer)
	b.WriteString("1")
	enc = json.NewDecoder(b)
	err = enc.Decode(&f)
	if err != nil {
		t.Fatal(err)
	}

	if expect != f {
		t.Fatalf("json decode expected: %s, got: %s", expect, f)
	}
}

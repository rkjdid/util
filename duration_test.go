package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type testDurationStruct struct {
	Duration Duration
}

func TestDuration_MarshalJSON(t *testing.T) {
	d := Duration(time.Second)
	expect := fmt.Sprintf("\"%s\"", time.Duration(d))

	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}

	res := string(b)
	if expect != res {
		t.Errorf("json encode expected: %s, got: %s", expect, res)
	}
}

func TestDuration_UnmashalJSON(t *testing.T) {
	expect := Duration(time.Second)
	b := new(bytes.Buffer)
	b.WriteString("\"1s\"")
	enc := json.NewDecoder(b)
	var d Duration
	err := enc.Decode(&d)
	if err != nil {
		t.Fatal(err)
	}

	if expect != d {
		t.Fatalf("json decode expected: %s, got: %s", expect, d)
	}

	expect = Duration(1)
	b = new(bytes.Buffer)
	b.WriteString("\"1\"")
	enc = json.NewDecoder(b)
	err = enc.Decode(&d)
	if err != nil {
		t.Fatal(err)
	}

	if expect != d {
		t.Fatalf("json decode expected: %s, got: %s", expect, d)
	}
}

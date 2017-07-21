package util

import (
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

type vtest struct {
	A string
	B string `toml:"-" json:"-"`
}

func TestReadWriteJson(t *testing.T) {
	v0 := vtest{A: "hello", B: "world"}

	buf := new(bytes.Buffer)
	err := WriteJson(v0, buf)
	if err != nil {
		t.Fatal(err)
	}

	var v1 vtest
	err = ReadJson(&v1, buf)
	if err != nil {
		t.Fatal(err)
	}

	if v1.A != v0.A {
		t.Error("got", v1.A, "expected", v0.A)
	}
	if v1.B == v0.B {
		t.Error("unexported fields shouldn't be encoded")
	}
}

func TestReadWriteToml(t *testing.T) {
	v0 := vtest{A: "hello", B: "world"}

	buf := new(bytes.Buffer)
	err := WriteToml(v0, buf)
	if err != nil {
		t.Fatal(err)
	}

	var v1 vtest
	err = ReadToml(&v1, buf)
	if err != nil {
		t.Fatal(err)
	}

	if v1.A != v0.A {
		t.Error("got", v1.A, "expected", v0.A)
	}
	if v1.B == v0.B {
		t.Error("unexported fields shouldn't be encoded", v1.B, v0.B)
	}
}

func TestReadWriteFile(t *testing.T) {
	tmp, err := ioutil.TempDir(os.TempDir(), "TestUtil")
	if err != nil {
		t.Fatal("ioutil.TempDir", err)
	}
	defer os.RemoveAll(tmp)

	v0 := vtest{A: "hello", B: "world"}
	tmpJson := path.Join(tmp, "v.json")
	tmpToml := path.Join(tmp, "v.toml")

	err = WriteJsonFile(v0, tmpJson)
	if err != nil {
		t.Fatal(err)
	}
	err = WriteTomlFile(v0, tmpToml)
	if err != nil {
		t.Fatal(err)
	}

	var v1 vtest
	err = ReadJsonFile(&v1, tmpJson)
	if err != nil {
		t.Fatal(err)
	}
	if v1.A != v0.A {
		t.Error("got", v1.A, "expected", v0.A)
	}
	if v1.B == v0.B {
		t.Error("unexported fields shouldn't be encoded")
	}

	v1 = vtest{}
	err = ReadTomlFile(&v1, tmpToml)
	if err != nil {
		t.Fatal(err)
	}
	if v1.A != v0.A {
		t.Error("got", v1.A, "expected", v0.A)
	}
	if v1.B == v0.B {
		t.Error("unexported fields shouldn't be encoded")
	}

	v1 = vtest{}
	err = ReadGenericFile(&v1, tmpJson)
	if err != nil {
		t.Fatal(err)
	}
	if v1.A != v0.A {
		t.Error("got", v1.A, "expected", v0.A)
	}
	if v1.B == v0.B {
		t.Error("unexported fields shouldn't be encoded")
	}

	v1 = vtest{}
	err = ReadGenericFile(&v1, tmpToml)
	if err != nil {
		t.Fatal(err)
	}
	if v1.A != v0.A {
		t.Error("got", v1.A, "expected", v0.A)
	}
	if v1.B == v0.B {
		t.Error("unexported fields shouldn't be encoded")
	}
}

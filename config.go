package util

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/rkjdid/errors"
	"io"
	"log"
	"os"
	"time"
)

// 2006 Jan _2 15:04:05
var yearStamp = "2/1/2006 15:04:05"

func TimeStampParse(val string) time.Time {
	t, err := time.Parse(yearStamp, val)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// "Generic" Read/Writers

// Read tries in turn to ReadJson, then ReadToml,
// if both fail it returns a detailed error
func ReadGeneric(v interface{}, rd io.Reader) (errs error) {
	err := ReadJson(v, rd)
	if err == nil {
		return nil
	}

	errs = errors.New(err)
	rdSeeker, ok := rd.(io.ReadSeeker)
	if !ok {
		return errors.Add(errs,
			fmt.Errorf("stopped after trying json, provided Reader isn't a ReadSeeker"))
	}

	// reset ReadSeeker to beginning
	_, err = rdSeeker.Seek(0, io.SeekStart)
	if err != nil {
		return errors.Add(errs, err)
	}

	err = ReadToml(v, rd)
	if err == nil {
		return nil
	}

	return errors.Add(errs, err)
}

func ReadJson(v interface{}, rd io.Reader) error {
	dec := json.NewDecoder(rd)
	err := dec.Decode(v)
	if err == io.EOF {
		err = nil
	}
	return err
}

func ReadToml(v interface{}, rd io.Reader) error {
	_, err := toml.DecodeReader(rd, v)
	if err == io.EOF {
		err = nil
	}
	return err
}

func ReadGenericFile(v interface{}, path string) error {
	return readFile(v, path, ReadGeneric)
}

func ReadJsonFile(v interface{}, path string) error {
	return readFile(v, path, ReadJson)
}

func ReadTomlFile(v interface{}, path string) error {
	return readFile(v, path, ReadToml)
}

func WriteToml(v interface{}, wr io.Writer) error {
	enc := toml.NewEncoder(wr)
	return enc.Encode(v)
}

func WriteJson(v interface{}, wr io.Writer) error {
	prettyJson, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return err
	}
	_, err = wr.Write(append(prettyJson, byte('\n')))
	return err
}

func WriteTomlFile(v interface{}, path string) error {
	return writeFile(v, path, WriteToml)
}

func WriteJsonFile(v interface{}, path string) error {
	return writeFile(v, path, WriteJson)
}

// writeFile opens file to path for writing and calls write(v) on fd
func writeFile(v interface{}, path string, write func(interface{}, io.Writer) error) error {
	fd, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	err = write(v, fd)
	if err != nil {
		return fmt.Errorf("error writing \"%s\": %s", path, err)
	}
	return nil
}

// readFile opens file to path for reading and calls read(v) on fd
func readFile(v interface{}, path string, read func(interface{}, io.Reader) error) error {
	fd, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer fd.Close()
	err = read(v, fd)
	if err != nil {
		return fmt.Errorf("error reading \"%s\": %s", path, err)
	}
	return nil
}

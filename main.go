package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
)

var input, output, delimiter string

func init() {
	flag.StringVar(&input, "i", "", "input json file path")
	flag.StringVar(&output, "o", "", "output csv file path")
	flag.StringVar(&delimiter, "delimiter", ",", "csv delimiter")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	i, err := os.Open(input)
	if err != nil {
		return err
	}
	o, err := os.OpenFile(output+"_tmp", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	d := json.NewDecoder(i)
	if err := expectToken(d, json.Delim('[')); err != nil {
		return err
	}
	columns := []string{}
	for d.More() {
		el := map[string]json.RawMessage{}
		if err := d.Decode(&el); err != nil {
			return err
		}
		for k := range el {
			if !slices.Contains(columns, k) {
				columns = append(columns, k)
			}
		}
		for i, c := range columns {
			if i == len(columns)-1 {
				fmt.Fprintf(o, "%s\n", el[c])
				break
			}
			fmt.Fprintf(o, "%s%s", el[c], delimiter)
		}
	}
	if err := expectToken(d, json.Delim(']')); err != nil {
		return err
	}
	// Prepend columns.
	if _, err := o.Seek(0, io.SeekStart); err != nil {
		return err
	}
	ofinal, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fmt.Fprintf(ofinal, "%s\n", strings.Join(columns, delimiter))
	if _, err = io.Copy(ofinal, o); err != nil {
		return err
	}
	if err := o.Close(); err != nil {
		return err
	}
	return os.Remove(output + "_tmp")
}

func expectToken(d *json.Decoder, expected json.Token) error {
	t, err := d.Token()
	if err != nil {
		return err
	}
	if t != expected {
		return fmt.Errorf("expected %v but got %v", expected, t)
	}
	return nil
}

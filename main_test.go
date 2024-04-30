package main

import (
	"os"
	"testing"
)

func TestObj(t *testing.T) {
	f, err := os.CreateTemp("", "json2csv*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(`[{
		"string": "value",
		"number": 123445,
		"array": [
			"a",
			"b"
		],
		"obj": {"key": "value"}
	}]`); err != nil {
		t.Fatal(err)
	}
	input = f.Name()
	output = f.Name() + ".csv"
	delimiter = ";"
	if err := run(); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	expected := `array;number;obj;string
["a","b"];123445;{"key": "value"};"value"
`
	if string(actual) != expected {
		t.Errorf("expected %q, actual %q", expected, actual)
	}
}

func TestColumns(t *testing.T) {
	f, err := os.CreateTemp("", "json2csv*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	if _, err := f.WriteString(`[{
		"col2": 1
	},{
		"col3": 2,
		"col1": 3
	},{
		"col3": 4
	}]`); err != nil {
		t.Fatal(err)
	}
	input = f.Name()
	output = f.Name() + ".csv"
	delimiter = ";"
	if err := run(); err != nil {
		t.Fatal(err)
	}
	actual, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	expected := `col2;col1;col3
1
;3;2
;;4
`
	if string(actual) != expected {
		t.Errorf("expected %q, actual %q", expected, actual)
	}
}

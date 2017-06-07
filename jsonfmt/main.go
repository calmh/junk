package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	close := flag.Bool("close", false, "Attempt to auto close truncated JSON")
	flag.Parse()

	bs, _ := ioutil.ReadAll(os.Stdin)
	if *close {
		bs = append(bs, closeJSON(bs)...)
	}

	buf := new(bytes.Buffer)
	if err := json.Indent(buf, bs, "", "  "); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	buf.WriteTo(os.Stdout)
}

func closeJSON(bs []byte) []byte {
	var stack []byte
	var prevC byte
	var inString bool
	for _, c := range bs {
		switch {
		case !inString && c == '"':
			inString = true
			stack = append([]byte{'"'}, stack...)
		case inString && c == '"' && prevC != '\\':
			inString = false
			stack = stack[1:]
		case !inString && c == '{':
			stack = append([]byte("}"), stack...)
		case !inString && c == '[':
			stack = append([]byte("]"), stack...)
		case !inString && c == '}' || c == ']':
			stack = stack[1:]
		}
		prevC = c
	}
	return stack
}

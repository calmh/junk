package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/hcl/hcl/printer"
	"github.com/hashicorp/hcl/json/parser"
)

func main() {
	bs, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ast, err := parser.Parse(bs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = printer.Fprint(os.Stdout, ast)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

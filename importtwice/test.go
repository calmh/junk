package main

import "fmt"
import "github.com/calmh/junk/importtwice/a"
import foo "github.com/calmh/junk/importtwice/a"

func main() {
	a.Var = 5
	a.Print()
	foo.Print()
	fmt.Println("done")
}

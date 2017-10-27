package a

import "fmt"

func init() {
	fmt.Println("a init()")
}

var Var int

func Print() {
	fmt.Println("Var =", Var)
}

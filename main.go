package main

import (
	"fmt"

	. "github.com/costowell/parsing-fun/common"
	"github.com/costowell/parsing-fun/llrecurse"
)

func test(parser Parser, input string) {
	fmt.Printf("Evaluating string '%s'... ", input)
	if parser.Parse(input) {
		fmt.Println("Success!")
	} else {
		fmt.Println("Failed")
	}
}

func main() {
	rules := []Rule{
		NewRule("a", Expr{Ref("a"), "a"}),
		NewRule("a", Expr{""}),
	}
	g, err := NewGrammar(rules)
	if err != nil {
		fmt.Printf("error creating grammar: %s", err.Error())
		return
	}
	fmt.Println(g)

	parser := llrecurse.New(g)
	test(parser, "a")
}

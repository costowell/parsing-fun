package main

import (
	"fmt"

	. "github.com/costowell/parsing-fun/common"
	"github.com/costowell/parsing-fun/earley"
)

func main() {
	rules := []Rule{
		NewRule("S", Expr{Ref("S"), "+", Ref("M")}),
		NewRule("S", Expr{Ref("M")}),
		NewRule("M", Expr{Ref("M"), "*", Ref("T")}),
		NewRule("M", Expr{Ref("T")}),
		NewRule("T", Expr{"1"}),
		NewRule("T", Expr{"2"}),
		NewRule("T", Expr{"3"}),
		NewRule("T", Expr{"4"}),
	}
	g, err := NewGrammar(rules)
	if err != nil {
		fmt.Printf("error creating grammar: %s", err.Error())
		return
	}
	fmt.Println(g)

	parser := earley.New(g)

	parser.Parse("2")
}

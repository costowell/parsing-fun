package main

import (
	"fmt"

	. "github.com/costowell/parsing-fun/common"
	"github.com/costowell/parsing-fun/earley"
)

func main() {
	_ = []Rule{
		NewRule("S", Expr{"a", Ref("S"), "a"}),
		NewRule("S", Expr{"b", Ref("S"), "b"}),
		NewRule("S", Expr{}),
		NewRule("S", Expr{"a"}),
		NewRule("S", Expr{"b"}),
	}
	rulesB := []Rule{
		NewRule("S0", Expr{Ref("Na"), Ref("A02")}),
		NewRule("S0", Expr{Ref("Nb"), Ref("A12")}),
		NewRule("S0", Expr{"a"}),
		NewRule("S0", Expr{"b"}),
		NewRule("A12", Expr{Ref("_S"), Ref("Nb")}),
		NewRule("A12", Expr{"b"}),
		NewRule("Nb", Expr{"b"}),
		NewRule("Na", Expr{"a"}),
		NewRule("_S", Expr{Ref("Na"), Ref("A02")}),
		NewRule("_S", Expr{Ref("Nb"), Ref("A12")}),
		NewRule("_S", Expr{"a"}),
		NewRule("_S", Expr{"b"}),
		NewRule("A02", Expr{Ref("_S"), Ref("Na")}),
		NewRule("A02", Expr{"a"}),
	}
	g, err := NewGrammar(rulesB)
	if err != nil {
		fmt.Println("error creating grammar: ", err.Error())
		return
	}
	fmt.Println(g)
	cnf, err := g.ToCNF()
	if err != nil {
		fmt.Println("error creating CNF grammar: ", err.Error())
		return
	}
	fmt.Println("-------------------------")
	fmt.Println(cnf.String())

	parser := earley.New(g)
	leftParse, err := parser.Parse("ababa")
	if err != nil {
		fmt.Println("error parsing CNF grammar: ", err.Error())
		return
	}
	fmt.Println("Left Parse: ", leftParse)
}

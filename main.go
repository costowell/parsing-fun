package main

import (
	"fmt"

	. "github.com/costowell/parsing-fun/common"
	"github.com/costowell/parsing-fun/earley"
)

func test(gram *Grammar, parser Parser, input string) {
	fmt.Printf("Evaluating string '%s'...\n", input)
	leftParse, err := parser.Parse(input)
	if err != nil {
		fmt.Printf("Failed to parse string: %s\n", err.Error())
		return
	}
	fmt.Println(leftParse)
	leftParseStr, err := gram.EvalLeftParse(leftParse)
	if err != nil {
		fmt.Printf("Verification failed: %s\n", err.Error())
		return
	}
	if input != leftParseStr {
		fmt.Println("Left parse does not generate input:")
		fmt.Printf("\tLeft Parse: %+v\n", leftParse)
		fmt.Printf("\t%s != %s\n", input, leftParseStr)
		return
	}
	fmt.Printf("Success! %s == %s\n", input, leftParseStr)
}

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
	test(g, parser, "2+3*4")
	test(g, parser, "1+2*3+4*1+2+3*4+1+2*3+4+1*2+3+4*2+1+3*4+2+1*2+3+4")
}

package common

import (
	"fmt"
)

// Variable represents any non-terminal in a grammar
type Variable string

// Symbol represents a RuleRef or string
type Symbol any

// Expr represents an array of Symbols
type Expr []Symbol

// Rule some sequence of Symbols
type Rule struct {
	Variable Variable
	Expr     Expr
}

// RuleRef is a placeholder for referencing a rule within an expression
type RuleRef struct {
	Variable Variable
}

type Grammar struct {
	rules map[Variable][]*Expr
	first Variable
}

func (g *Grammar) Rules() map[Variable][]*Expr {
	return g.rules
}

func (g *Grammar) GetRule(v Variable) []*Expr {
	return g.rules[v]
}

func (g *Grammar) FirstRule() Variable {
	return g.first
}

func (g *Grammar) String() string {
	var s string
	for v, exprs := range g.Rules() {
		s += fmt.Sprintf("%s -> %s", v, exprs[0])
		for _, expr := range exprs[1:] {
			s += fmt.Sprintf(" | %s", expr)
		}
		s += "\n"
	}
	return s
}

func (e *Expr) String() string {
	var s string
	for _, sym := range *e {
		ref, ok := sym.(RuleRef)
		if ok {
			s += fmt.Sprintf("%s ", ref.Variable)
			continue
		}
		str, ok := sym.(string)
		if ok {
			s += fmt.Sprintf("'%s' ", str)
			continue
		}
		return "FAILED_PARSE "
	}
	return s
}

// Ref returns a RuleRef for a given Variable
func Ref(v Variable) RuleRef {
	return RuleRef{
		Variable: v,
	}
}

func NewGrammar(rules []Rule) (*Grammar, error) {
	ruleMap := make(map[Variable][]*Expr, len(rules))

	// Init map
	for _, rule := range rules {
		if _, ok := ruleMap[rule.Variable]; !ok {
			ruleMap[rule.Variable] = make([]*Expr, 0)
		}
		ruleMap[rule.Variable] = append(ruleMap[rule.Variable], &rule.Expr)
	}

	g := &Grammar{
		rules: ruleMap,
		first: rules[0].Variable,
	}
	return g, nil
}

func NewRule(v Variable, expr Expr) Rule {
	return Rule{
		Variable: v,
		Expr:     expr,
	}
}

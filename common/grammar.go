package common

import (
	"fmt"
)

// Variable represents any non-terminal in a grammar
type Variable string

// Terminal represents any terminal in a grammar
type Terminal string

// Symbol represents a RuleRef or string
type Symbol any

// Expr represents an array of Symbols
type Expr []Symbol

// Rule some sequence of Symbols
type Rule struct {
	Variable Variable
	Expr     Expr
}

func (r *Rule) String() string {
	return fmt.Sprintf("%s -> %s", r.Variable, r.Expr.String())
}

// RuleRef is a placeholder for referencing a rule within an expression
type RuleRef struct {
	Variable Variable
}

type Grammar struct {
	Rules     []*Rule
	RulesMap  map[Variable][]*Expr
	Terminals []Terminal
	Variables []Variable
}

func (e *Expr) ApplyRule(rule *Rule) error {
	for i, sym := range *e {
		ref, ok := sym.(RuleRef)
		if !ok {
			continue
		}
		if rule.Variable != ref.Variable {
			return fmt.Errorf("Expected \"%s\", got \"%s\"", rule.Variable, ref.Variable)
		}

		newExpr := make(Expr, i)
		copy(newExpr, (*e)[:i])
		newExpr = append(newExpr, rule.Expr...)
		newExpr = append(newExpr, (*e)[i+1:]...)
		*e = newExpr
		return nil
	}
	return fmt.Errorf("Failed to find '%s'", rule.Variable)
}

func (g *Grammar) EvalLeftParse(leftParse []int) (string, error) {
	expr := Expr{Ref(g.StartVariable())}
	for _, ruleNum := range leftParse {
		rule := g.Rules[ruleNum]
		fmt.Print(expr, "->")
		if err := expr.ApplyRule(rule); err != nil {
			return "", fmt.Errorf("Failed to apply rule \"%s\" to \"%s\": %s", rule.String(), expr.String(), err.Error())
		}
		fmt.Printf("%s [%s]\n", expr, rule.String())
	}

	var str string
	for _, sym := range expr {
		switch v := sym.(type) {
		case string:
			str += v
		default:
			return "", fmt.Errorf("Incomplete left parse '%s'", expr.String())
		}
	}
	return str, nil
}

func (g *Grammar) StartVariable() Variable {
	return g.Rules[0].Variable
}

func (g *Grammar) String() string {
	var s string
	s += fmt.Sprintf("Terminals: %v\n", g.Terminals)
	s += fmt.Sprintf("Variables: %v\n", g.Variables)
	for v, exprs := range g.RulesMap {
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
	var variables []Variable
	var terminals []Terminal
	var rulesP []*Rule

	// Init map
	for _, rule := range rules {
		if _, ok := ruleMap[rule.Variable]; !ok {
			ruleMap[rule.Variable] = make([]*Expr, 0)
		}
		ruleMap[rule.Variable] = append(ruleMap[rule.Variable], &rule.Expr)
		rulesP = append(rulesP, &rule)

		variables = append(variables, rule.Variable)
		for _, sym := range rule.Expr {
			if term, ok := sym.(string); ok {
				terminals = append(terminals, Terminal(term))
			}
		}
	}

	g := &Grammar{
		RulesMap:  ruleMap,
		Rules:     rulesP,
		Variables: variables,
		Terminals: terminals,
	}
	return g, nil
}

func NewRule(v Variable, expr Expr) Rule {
	return Rule{
		Variable: v,
		Expr:     expr,
	}
}

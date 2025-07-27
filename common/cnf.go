package common

import (
	"fmt"
)

func _variableRemovalPermutations(expr Expr, variable Variable, offset int) []Expr {
	var exprs []Expr
	for i, sym := range expr[offset:] {
		if v, ok := sym.(RuleRef); ok {
			if v.Variable == variable {
				index := offset + i
				newExpr := make(Expr, len(expr)-1)
				copy(newExpr, expr[:index])
				copy(newExpr[index:], expr[index+1:])
				exprs = append(exprs, newExpr)
				exprs = append(exprs, _variableRemovalPermutations(newExpr, variable, index)...)
				exprs = append(exprs, _variableRemovalPermutations(expr, variable, index+1)...)
				break
			}
		}
	}
	return exprs
}

func VariableRemovalPermutations(expr Expr, variable Variable) []Expr {
	return _variableRemovalPermutations(expr, variable, 0)
}

func removeEpsilonProductions(rules *[]Rule) {
	ruleMap := make(map[Variable][]Rule)
	for _, rule := range *rules {
		if _, ok := ruleMap[rule.Variable]; !ok {
			ruleMap[rule.Variable] = make([]Rule, 0)
		}
		ruleMap[rule.Variable] = append(ruleMap[rule.Variable], rule)
	}
	for i, rule := range *rules {
		if len(rule.Expr) == 0 {
			*rules = append((*rules)[:i], (*rules)[i+1:]...)
			targetVar := rule.Variable
			if len(ruleMap[targetVar]) > 1 {
				for _, rule := range *rules {
					for _, expr := range VariableRemovalPermutations(rule.Expr, targetVar) {
						*rules = append(*rules, Rule{Variable: rule.Variable, Expr: expr})
					}
				}
			} else {
				for j, rule := range *rules {
					for i, sym := range rule.Expr {
						if v, ok := sym.(RuleRef); ok {
							if v.Variable == targetVar {
								newExpr := make(Expr, len(rule.Expr)-1)
								copy(newExpr, rule.Expr[:i])
								copy(newExpr[i:], rule.Expr[i+1:])
								(*rules)[j].Expr = newExpr
							}
						}
					}
				}
			}
		}
	}
}

func cnfTransformSymbol(sym Symbol) Variable {
	switch v := sym.(type) {
	case string:
		return Variable("N" + v)
	case Terminal:
		return Variable("N" + v.String())
	case RuleRef:
		return Variable("_" + v.Variable)
	case Variable:
		return Variable("_" + v.String())
	}

	return "INVALID SYMBOL"
}

// ToCNF converts a Grammar to an equivalent Grammar in Chomsky Normal Form
// Based off of the algorithm described here: https://en.wikipedia.org/wiki/Chomsky_normal_form#Converting_a_grammar_to_Chomsky_normal_form
func (g *Grammar) ToCNF() (*Grammar, error) {
	var rules []Rule

	// START: Eliminate the start symbol from right-hand sides
	// S0 -> _S
	rules = append(rules, NewRule("S0", Expr{Ref(cnfTransformSymbol(g.StartVariable()))}))

	// TERM: Eliminate rules with nonsolitary terminals
	for _, term := range g.Terminals.Data {
		rules = append(rules, NewRule(cnfTransformSymbol(term), Expr{term}))
	}

	// BIN: Eliminate right-hand sides with more than 2 nonterminals
	for i, rule := range g.Rules {
		// A -> X1 X2 X3... becomes _A -> _X1 _X2 _X3...
		r := rule.Copy()
		r.Variable = cnfTransformSymbol(rule.Variable)

		// I love Go type system :)
		expr := make([]RuleRef, len(r.Expr))
		for i, sym := range r.Expr {
			expr[i] = Ref(cnfTransformSymbol(sym))
			r.Expr[i] = expr[i]
		}

		if len(r.Expr) > 2 {
			for j := 1; j < len(expr); j++ {
				var ruleVar Variable
				var endSymbol Variable
				if j == len(r.Expr)-1 {
					endSymbol = expr[j].Variable
				} else {
					endSymbol = Variable(fmt.Sprintf("A%d%d", i, j+1))
				}

				if j == 1 {
					ruleVar = r.Variable
				} else {
					ruleVar = Variable(fmt.Sprintf("A%d%d", i, j))
				}

				rule := Rule{
					Variable: ruleVar,
					Expr: Expr{
						expr[j-1],
						Ref(endSymbol),
					},
				}

				rules = append(rules, rule)
			}
		} else {
			rules = append(rules, r)
		}
	}

	// DEL: Eliminate ε-rules
	removeEpsilonProductions(&rules)

	for _, rule := range rules {
		fmt.Println(rule.String())
	}

	// UNIT: Eliminate unit rules
	var removeIndices []int
	for i := 0; i < len(rules); i++ {
		rule := rules[i]
		if len(rule.Expr) != 1 {
			continue
		}
		lhsVar := rule.Variable
		if v, ok := rule.Expr[0].(RuleRef); ok {
			removeIndices = append(removeIndices, i)
			for _, rule := range rules {
				if rule.Variable == v.Variable {
					newExpr := make(Expr, len(rule.Expr))
					copy(newExpr, rule.Expr)
					rules = append(rules, Rule{Variable: lhsVar, Expr: newExpr})
				}
			}
		}
	}
	fmt.Println(removeIndices)
	for i, index := range removeIndices {
		rules = append(rules[:index-i], rules[index-i+1:]...)
	}

	return NewGrammar(rules)
}

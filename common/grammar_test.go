package common

import (
	"testing"
)

func TestNewGrammar(t *testing.T) {
	tests := []struct {
		name        string
		rules       []Rule
		expectError bool
	}{
		{
			name: "simple grammar",
			rules: []Rule{
				NewRule("A", Expr{"a"}),
			},
			expectError: false,
		},
		{
			name: "invalid ref",
			rules: []Rule{
				NewRule("A", Expr{Ref("B")}),
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g, err := NewGrammar(tt.rules)
			if err == nil && tt.expectError {
				t.Errorf("NewGrammar() expected error but got none")
			}
			if err != nil && !tt.expectError {
				t.Errorf("NewGrammar() unexpected error: %v", err)
			}
			if g == nil {
				return
			}
			if len(g.Rules) != len(tt.rules) || len(g.RulesMap) != len(tt.rules) {
				t.Errorf("NewGrammar() unexpected number of rules %v, wanted %v", len(g.Rules), len(tt.rules))
			}
		})
	}
}

func TestGrammarEvalLeftParse(t *testing.T) {
	rules := []Rule{
		NewRule("S", Expr{"a", Ref("B"), Ref("C")}),
		NewRule("B", Expr{"b"}),
		NewRule("C", Expr{"c"}),
	}

	gram, err := NewGrammar(rules)
	if err != nil {
		t.Fatalf("NewGrammar() error = %v", err)
	}

	tests := []struct {
		name        string
		leftParse   []int
		expected    string
		expectError bool
	}{
		{
			name:        "valid left parse",
			leftParse:   []int{0, 1, 2},
			expected:    "abc",
			expectError: false,
		},
		{
			name:        "unknown rule",
			leftParse:   []int{3},
			expected:    "",
			expectError: true,
		},
		{
			name:        "no leftmost variable matching rule",
			leftParse:   []int{0, 2, 1},
			expected:    "",
			expectError: true,
		},
		{
			name:        "incomplete parse",
			leftParse:   []int{0},
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty parse",
			leftParse:   []int{},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := gram.EvalLeftParse(tt.leftParse)

			if tt.expectError && err == nil {
				t.Errorf("EvalLeftParse() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("EvalLeftParse() unexpected error: %v", err)
			}
			if !tt.expectError && result != tt.expected {
				t.Errorf("EvalLeftParse() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestApplyRuleLeft(t *testing.T) {
	tests := []struct {
		name        string
		expr        Expr
		rule        *Rule
		expected    Expr
		expectError bool
	}{
		{
			name: "successful rule application",
			expr: Expr{"a", Ref("B"), "c"},
			rule: &Rule{
				Variable: "B",
				Expr:     Expr{"b1", "b2"},
			},
			expected:    Expr{"a", "b1", "b2", "c"},
			expectError: false,
		},
		{
			name: "rule not found",
			expr: Expr{"a", Ref("C"), "c"},
			rule: &Rule{
				Variable: "B",
				Expr:     Expr{"b1", "b2"},
			},
			expected:    Expr{},
			expectError: true,
		},
		{
			name: "leftmost rule not match",
			expr: Expr{"a", Ref("C"), Ref("B"), "c"},
			rule: &Rule{
				Variable: "B",
				Expr:     Expr{"b1", "b2"},
			},
			expected:    Expr{},
			expectError: true,
		},
		{
			name: "no rule refs in expr",
			expr: Expr{"a", "b", "c"},
			rule: &Rule{
				Variable: "B",
				Expr:     Expr{"b1", "b2"},
			},
			expected:    Expr{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalExpr := make(Expr, len(tt.expr))
			copy(originalExpr, tt.expr)

			err := tt.expr.ApplyRuleLeft(tt.rule)

			if tt.expectError && err == nil {
				t.Errorf("ApplyRule() expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("ApplyRule() unexpected error: %v", err)
			}

			if !tt.expectError {
				if len(tt.expr) != len(tt.expected) {
					t.Errorf("ApplyRule() result length = %v, want %v", len(tt.expr), len(tt.expected))
				}
				for i, sym := range tt.expr {
					if i < len(tt.expected) {
						if sym != tt.expected[i] {
							t.Errorf("ApplyRule() result[%d] = %v, want %v", i, sym, tt.expected[i])
						}
					}
				}
			}
		})
	}
}

package earley

import (
	"testing"

	. "github.com/costowell/parsing-fun/common"
)

func TestOperatorPrecedence(t *testing.T) {
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
		t.Errorf("NewGrammar() unexpected error: %v", err)
		return
	}
	parser := New(g)
	tests := []struct {
		name        string
		input       string
		expectError bool
	}{
		{
			name:        "single",
			input:       "1",
			expectError: false,
		},
		{
			name:        "binary",
			input:       "1*2",
			expectError: false,
		},
		{
			name:        "mixed precedence",
			input:       "2+3*4",
			expectError: false,
		},
		{
			name:        "looong",
			input:       "2*3+4*2*3+4+2*3+4*2*3+4+2*3+4+2*3+4*2*3+4+2*3+4*2*3+4+2*3+4",
			expectError: false,
		},
		{
			name:        "empty",
			input:       "",
			expectError: true,
		},
		{
			name:        "incomplete",
			input:       "2*3+",
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leftParse, err := parser.Parse(tt.input)
			if !tt.expectError && err != nil {
				t.Errorf("Parse() unexpected error: %v", err)
			}
			if tt.expectError && err == nil {
				t.Errorf("Parse() expected error, got none")
			}
			if err != nil {
				return
			}

			leftParseStr, err := g.EvalLeftParse(leftParse)
			if !tt.expectError && err != nil {
				t.Errorf("EvalLeftParse() unexpected error %v", err)
			}
			if tt.expectError && err == nil {
				t.Error("EvalLeftParse() expected error, got none")
			}
			if err != nil {
				return
			}

			if tt.input != leftParseStr {
				t.Errorf("EvalLeftParse() invalid left parse, expected %v, got %v", tt.input, leftParseStr)
				return
			}
		})
	}
}

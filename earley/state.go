package earley

import (
	"fmt"
	. "github.com/costowell/parsing-fun/common"
)

type State struct {
	variable       Variable
	rule           *Expr
	position       int
	originPosition int
}

const positionMarker = "•"

func (s *State) String() string {
	var ruleString string
	for i, sym := range *s.rule {
		if i == s.position {
			ruleString += positionMarker
		}
		switch s := sym.(type) {
		case string:
			ruleString += "'" + s + "'" + " "
		case RuleRef:
			ruleString += string(s.Variable) + " "
		}
	}
	if len(*s.rule) == s.position {
		ruleString += positionMarker
	}
	return fmt.Sprintf("(%s -> %s, %d)", s.variable, ruleString, s.originPosition)
}

func (s *State) IncrementPosition() State {
	return State{
		variable:       s.variable,
		rule:           s.rule,
		position:       s.position + 1,
		originPosition: s.originPosition,
	}
}

func (s *State) NextSym() Symbol {
	if s.IsComplete() {
		return nil
	}
	return (*s.rule)[s.position]
}

func (s *State) IsComplete() bool {
	return s.position >= len(*s.rule)
}

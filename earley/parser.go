package earley

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/costowell/parsing-fun/common"
)

const startSymbol = "_P"

type State struct {
	variable       Variable
	rule           *Expr
	position       int
	originPosition int
}

func (s *State) String() string {
	var ruleString string
	for i, sym := range *s.rule {
		if i == s.position {
			ruleString += "•"
		}
		switch s := sym.(type) {
		case string:
			ruleString += "'" + s + "'" + " "
		case RuleRef:
			ruleString += string(s.Variable) + " "
		}
	}
	if len(*s.rule) == s.position {
		ruleString += "•"
	}
	return fmt.Sprintf("(%s -> %s, %d)", s.variable, ruleString, s.originPosition)
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

type realParser struct {
	gram *Grammar
	S    []OrderedSet[State]
}

func (p *realParser) InsertState(k int, rule *Expr, startSymbol Variable, position int, originPosition int) {
	for i := len(p.S) - 1; i < k; i++ {
		p.S = append(p.S, NewOrderedSet[State]())
	}
	p.S[k].Insert(State{startSymbol, rule, position, originPosition})
}

func (p *realParser) Predict(k int, state State) {
	nextSym := state.NextSym()
	if nextSym == nil {
		return
	}
	ref, ok := nextSym.(RuleRef)
	if !ok {
		return
	}

	for _, production := range p.gram.GetRule(ref.Variable) {
		p.InsertState(k, production, ref.Variable, 0, k)
	}
}

func (p *realParser) Scan(k int, input string, state State) bool {
	nextSym := state.NextSym()
	if nextSym == nil {
		return false
	}
	ref, ok := nextSym.(string)
	if !ok {
		return false
	}
	_, found := strings.CutPrefix(input, ref)
	if found {
		p.InsertState(k+1, state.rule, state.variable, state.position+1, state.originPosition)
		return true
	}
	return false
}

func (p *realParser) Complete(k int, state State) {
	if !state.IsComplete() {
		return
	}
	for _, kState := range p.S[state.originPosition].Data {
		if ref, ok := kState.NextSym().(RuleRef); ok && ref.Variable == state.variable {
			p.InsertState(k, kState.rule, kState.variable, kState.position+1, kState.originPosition)
		}
	}
}

func (p *realParser) Parse(input string) error {
	p.S = make([]OrderedSet[State], 0)

	finalState := State{
		variable:       startSymbol,
		rule:           &Expr{Ref(p.gram.FirstRule())},
		position:       1,
		originPosition: 0,
	}

	// Add the first state: P -> *S
	p.InsertState(0, finalState.rule, finalState.variable, 0, 0)

	// maxK := len(p.input)
	for k := 0; k >= 0; k++ {
		if k >= len(p.S) {
			break
		}
		for i := 0; i < len(p.S[k].Data); i++ {
			state := p.S[k].Data[i]
			if state.IsComplete() {
				p.Complete(k, state)
			} else {
				switch state.NextSym().(type) {
				case string:
					if k >= len(input) {
						p.Scan(k, "", state)
					} else {
						p.Scan(k, input[k:], state)
					}
				case RuleRef:
					p.Predict(k, state)
				}
			}
		}
	}
	p.PrintState()
	if !p.S[len(p.S)-1].Contains(finalState) || len(p.S) != len(input)+1 {
		p.PrintState()
		return errors.New("State did not end with completion")
	}

	return nil
}

func (p *realParser) PrintState() {
	for k, set := range p.S {
		fmt.Printf("S(%d):\n", k)
		for _, state := range set.Data {
			fmt.Println(state.String())
		}
	}
}

func New(gram *Grammar) Parser {
	return &realParser{
		gram: gram,
	}
}

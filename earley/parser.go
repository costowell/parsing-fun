package earley

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/costowell/parsing-fun/common"
)

type realParser struct {
	gram *Grammar
	S    []OrderedSet[State]
}

func (p *realParser) InsertState(k int, state State) {
	for i := len(p.S) - 1; i < k; i++ {
		p.S = append(p.S, NewOrderedSet[State]())
	}
	p.S[k].Insert(state)
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
		p.InsertState(k, State{
			variable:       ref.Variable,
			rule:           production,
			position:       0,
			originPosition: k,
		})
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
		p.InsertState(k+1, state.IncrementPosition())
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
			p.InsertState(k, kState.IncrementPosition())
		}
	}
}

func (p *realParser) Parse(input string) error {
	p.S = make([]OrderedSet[State], 0)

	// _P -> •S
	startState := State{
		variable:       "_P",
		rule:           &Expr{Ref(p.gram.FirstRule())},
		position:       0,
		originPosition: 0,
	}
	// _P -> S•
	finalState := startState.IncrementPosition()

	// Add the first state
	p.InsertState(0, startState)

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

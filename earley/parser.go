package earley

import (
	"errors"
	"fmt"
	"strings"

	. "github.com/costowell/parsing-fun/common"
)

type realParser struct {
	gram      *Grammar
	S         []OrderedSet[State]
	BT        map[State][]*State
	ruleOrder map[*Expr]*int
}

func (p *realParser) InsertBT(state State, completedBy *State) {
	if _, ok := p.BT[state]; !ok {
		p.BT[state] = make([]*State, 0)
	}
	p.BT[state] = append(p.BT[state], completedBy)
}

func (p *realParser) InsertState(state State) {
	for i := len(p.S) - 1; i < state.k; i++ {
		p.S = append(p.S, NewOrderedSet[State]())
	}
	p.S[state.k].Insert(state)
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

	for _, production := range p.gram.RulesMap[ref.Variable] {
		p.InsertState(State{
			k:              k,
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
		s := state.IncrementPosition().IncrementK()
		for _, btState := range p.BT[state] {
			p.InsertBT(s, btState)
		}
		p.InsertState(s)
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
			newKState := kState.IncrementPosition()
			newKState.k = k
			for _, btState := range p.BT[kState] {
				p.InsertBT(newKState, btState)
			}
			p.InsertBT(newKState, &state)
			p.InsertState(newKState)
		}
	}
}

func (p *realParser) Parse(input string) ([]int, error) {
	p.S = make([]OrderedSet[State], 0)
	p.BT = make(map[State][]*State)

	// _P -> •S
	startState := State{
		k:              0,
		variable:       "_P",
		rule:           &Expr{Ref(p.gram.StartVariable())},
		position:       0,
		originPosition: 0,
	}
	// _P -> S•

	// Add the first state
	p.InsertState(startState)

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

	finalState := startState.IncrementPosition()
	finalState.k = len(p.S) - 1

	if !p.S[len(p.S)-1].Contains(finalState) || len(p.S) != len(input)+1 {
		return nil, errors.New("State did not end with completion")
	}

	return p.GenerateLeftParse(finalState), nil
}

func (p *realParser) GenerateLeftParse(state State) []int {
	var leftParse []int
	var i int

	if ruleNum, ok := p.ruleOrder[state.rule]; ok {
		leftParse = append(leftParse, *ruleNum)
	}

	states := p.BT[state]
	for _, sym := range *state.rule {
		switch v := sym.(type) {
		case string:
			fmt.Print(v)
		case RuleRef:
			leftParse = append(leftParse, p.GenerateLeftParse(*states[i])...)
			i++
		}
	}
	return leftParse
}

func (p *realParser) PrintState() {
	for k, set := range p.S {
		fmt.Printf("S(%d):\n", k)
		for _, state := range set.Data {
			fmt.Println(state.String())
		}
	}

	fmt.Println("Completed:")
	for key, value := range p.BT {
		fmt.Printf("%+v %+v\n", key.String(), value)
	}
}

func New(gram *Grammar) Parser {
	ruleOrder := make(map[*Expr]*int)
	for k, rule := range gram.Rules {
		ruleOrder[&rule.Expr] = &k
	}
	return &realParser{
		gram:      gram,
		ruleOrder: ruleOrder,
	}
}

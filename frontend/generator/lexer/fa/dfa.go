package fa

import (
	"github.com/uans3k/pl/infra"
)

type DFAEdge struct {
	Char      Char
	FromState int
	ToState   int
}

type DFA struct {
	nfa                    *NFA
	enableState            int
	StartState             int
	AcceptState2TokenTypes map[int]*TokenType
	State2Edges            [][]*DFAEdge
}

func NFA2DFA(nfa *NFA) *DFA {
	dfa := &DFA{
		nfa:                    nfa,
		AcceptState2TokenTypes: map[int]*TokenType{},
	}
	dfa.transform()
	dfa.nfa = nil
	return dfa
}

type char2StateSet struct {
	Key2StateSet map[string]infra.Set[int]
	Key2Char     map[string]Char
}

func (d *DFA) transform() {
	var (
		state2NFAStateSet      []infra.Set[int]
		nextState              int
		nextNFAStateSet        infra.Set[int]
		curState               int
		curNFAStateSet         infra.Set[int]
		openList               []int
		curStateDelta          *char2StateSet
		stateExist             bool
		acceptState2TokenTypes = map[int][]*TokenType{}
	)
	nextState = d.nextEnableState()
	nextNFAStateSet = d.eliminateClosure(infra.NewSet(d.nfa.StartState))
	state2NFAStateSet = append(state2NFAStateSet, nextNFAStateSet)
	openList = []int{nextState}
	for len(openList) > 0 {
		curState = openList[0]
		openList = openList[1:]
		curNFAStateSet = state2NFAStateSet[curState]
		curStateDelta = d.delta(curNFAStateSet)
		for key, stateSet := range curStateDelta.Key2StateSet {
			nextNFAStateSet = d.eliminateClosure(stateSet)
			if nextState, stateExist = d.exist(state2NFAStateSet, nextNFAStateSet); !stateExist {
				nextState = d.nextEnableState()
				if tokenTypes, accept := d.accept(nextNFAStateSet); accept {
					acceptState2TokenTypes[nextState] = tokenTypes
				}
				state2NFAStateSet = append(state2NFAStateSet, nextNFAStateSet)
				openList = append(openList, nextState)
			}
			d.State2Edges[curState] = append(d.State2Edges[curState], &DFAEdge{
				Char:      curStateDelta.Key2Char[key],
				FromState: curState,
				ToState:   nextState,
			})
		}
	}
	for acceptState, tokenTypes := range acceptState2TokenTypes {
		infra.SliceSort(tokenTypes, func(left, right *TokenType) bool {
			return left.Order < right.Order
		})
		d.AcceptState2TokenTypes[acceptState] = tokenTypes[0]
	}
}

func (d *DFA) accept(nfaStateSet infra.Set[int]) (tokenTypes []*TokenType, accept bool) {
	for nfaState := range nfaStateSet {
		if tokenType, ok := d.nfa.AcceptState2TokenType[nfaState]; ok {
			accept = ok
			tokenTypes = append(tokenTypes, tokenType)
		}
	}
	return
}

func (d *DFA) exist(state2NFAStateSet []infra.Set[int], nfaStateSet infra.Set[int]) (state int, exist bool) {
	var curNFAStateSet infra.Set[int]
	for state, curNFAStateSet = range state2NFAStateSet {
		if curNFAStateSet.Equal(nfaStateSet) {
			exist = true
			return
		}
	}
	exist = false
	return
}

func (d *DFA) nextEnableState() int {
	nextEnableState := d.enableState
	d.State2Edges = append(d.State2Edges, nil)
	d.enableState++
	return nextEnableState
}

func (d *DFA) eliminateClosure(stateSet infra.Set[int]) (closeSet infra.Set[int]) {
	openList := stateSet.Members()
	closeSet = infra.NewSet[int]()
	for len(openList) > 0 {
		state := openList[0]
		openList = openList[1:]
		closeSet.Add(state)
		edges := d.nfa.State2Edges[state]
		for _, edge := range edges {
			if edge.Char.Equal(CharEliminate) {
				exist := closeSet.AddIfNotExist(edge.ToState)
				if !exist {
					openList = append(openList, edge.ToState)
				}
			}
		}
	}
	return
}

func (d *DFA) delta(stateSet infra.Set[int]) (stateDelta *char2StateSet) {
	stateDelta = &char2StateSet{
		Key2StateSet: map[string]infra.Set[int]{},
		Key2Char:     map[string]Char{},
	}
	for state := range stateSet {
		edges := d.nfa.State2Edges[state]
		for _, edge := range edges {
			switch v := edge.Char.(type) {
			case *CharSingle, *CharRange:
				set, ok := stateDelta.Key2StateSet[v.Key()]
				if !ok {
					set = infra.NewSet[int]()
					stateDelta.Key2StateSet[v.Key()] = set
					stateDelta.Key2Char[v.Key()] = v
				}
				set.Add(edge.ToState)
			}
		}
	}
	return
}

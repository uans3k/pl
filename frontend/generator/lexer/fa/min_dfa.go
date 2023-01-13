package fa

import (
	"github.com/uans3k/pl/infra"
)

type MinDFAEdge struct {
	Chars   []Char
	ToState int
}

type SplitStates struct {
	State2DFAStates [][]int
	DFAState2State  map[int]int
}

func (s *SplitStates) Append(t *SplitStates) {
	shift := len(s.State2DFAStates)
	for _, dfaStates := range t.State2DFAStates {
		s.State2DFAStates = append(s.State2DFAStates, dfaStates)
		for _, dfaState := range dfaStates {
			s.DFAState2State[dfaState] = t.DFAState2State[dfaState] + shift
		}
	}
}

type MinDFA struct {
	dfa                   *DFA
	StartState            int
	TokenTypes            []*TokenType
	AcceptState2TokenType map[int]*TokenType
	State2Edges           [][]*MinDFAEdge
}

func DFA2MinDFA(dfa *DFA) *MinDFA {
	minDFA := &MinDFA{
		dfa:                   dfa,
		AcceptState2TokenType: map[int]*TokenType{},
	}
	minDFA.transform()
	minDFA.dfa = nil
	return minDFA
}

func (m *MinDFA) transform() {
	splitStates := m.hopcroft()
	acceptState2TokenTypes := map[int][]*TokenType{}
	m.State2Edges = make([][]*MinDFAEdge, len(splitStates.State2DFAStates))
	for state, dfaStates := range splitStates.State2DFAStates {
		for _, dfaState := range dfaStates {
			if dfaTokenType, ok := m.dfa.AcceptState2TokenTypes[dfaState]; ok {
				var tokenTypes []*TokenType
				if tokenTypes, ok = acceptState2TokenTypes[state]; !ok {
					tokenTypes = []*TokenType{dfaTokenType}
				} else {
					tokenTypes = append(tokenTypes, dfaTokenType)
				}
				acceptState2TokenTypes[state] = tokenTypes
			}
			stateEdges := m.State2Edges[state]
			dfaEdges := m.dfa.State2Edges[dfaState]
			for _, dfaEdge := range dfaEdges {
				stateEdges = m.appendStateEdges(splitStates, stateEdges, dfaEdge)
			}
			m.State2Edges[state] = stateEdges
		}
	}
	name2TokenType := map[string]*TokenType{}
	for acceptState, tokenTypes := range acceptState2TokenTypes {
		infra.SliceSort(tokenTypes, func(left, right *TokenType) bool {
			return left.Order < right.Order
		})
		tokenType := tokenTypes[0]
		m.AcceptState2TokenType[acceptState] = tokenType
		name2TokenType[tokenType.Name] = tokenType
	}
	m.TokenTypes = infra.MapSortV(name2TokenType, func(left, right *TokenType) bool {
		return left.Order < right.Order
	})
}

func (m *MinDFA) existChar(chars []Char, t Char) bool {
	for _, char := range chars {
		if char.Equal(t) {
			return true
		}
	}
	return false
}

func (m *MinDFA) appendStateEdges(splitStates *SplitStates, stateEdges []*MinDFAEdge, dfaEdge *DFAEdge) (appendedStateEdges []*MinDFAEdge) {
	toState := splitStates.DFAState2State[dfaEdge.ToState]
	for _, stateEdge := range stateEdges {
		if stateEdge.ToState == toState {
			if !m.existChar(stateEdge.Chars, dfaEdge.Char) {
				stateEdge.Chars = append(stateEdge.Chars, dfaEdge.Char)
			}
			appendedStateEdges = stateEdges
			return
		}
	}
	appendedStateEdges = append(stateEdges, &MinDFAEdge{
		Chars:   []Char{dfaEdge.Char},
		ToState: toState,
	})
	return
}

func (m *MinDFA) hopcroft() *SplitStates {
	var (
		nextSplitStates *SplitStates
		splittingStates *SplitStates
		split           bool
	)
	nextSplitStates, split = m.initSplitStates()
	for split {
		split = false
		splittingStates = nextSplitStates
		nextSplitStates = &SplitStates{
			DFAState2State: map[int]int{},
		}
		for splittingStateIndex := range splittingStates.State2DFAStates {
			curSplitStates, curSplit := m.split(splittingStates, splittingStateIndex)
			nextSplitStates.Append(curSplitStates)
			split = split || curSplit
		}
	}
	return splittingStates
}

func (m *MinDFA) initSplitStates() (splitStates *SplitStates, split bool) {
	tokenType2State := map[string]int{}
	splitStates = &SplitStates{
		State2DFAStates: make([][]int, 1),
		DFAState2State:  map[int]int{},
	}
	for dfaState, _ := range m.dfa.State2Edges {
		if tokenType, ok := m.dfa.AcceptState2TokenTypes[dfaState]; !ok {
			splitStates.State2DFAStates[0] = append(splitStates.State2DFAStates[0], dfaState)
			splitStates.DFAState2State[dfaState] = 0
		} else {
			split = true
			if state, ok := tokenType2State[tokenType.Name]; ok {
				splitStates.State2DFAStates[state] = append(splitStates.State2DFAStates[state], dfaState)
				splitStates.DFAState2State[dfaState] = state
			} else {
				splitStates.State2DFAStates = append(splitStates.State2DFAStates, []int{dfaState})
				state = len(splitStates.State2DFAStates) - 1
				splitStates.DFAState2State[dfaState] = state
				tokenType2State[tokenType.Name] = state
			}
		}
	}
	return
}

func (m *MinDFA) split(splittingStates *SplitStates, splittingStateIndex int) (splitStates *SplitStates, split bool) {
	var (
		closeSet  = map[string]Char{}
		splitChar Char
	)
	for _, dfaState := range splittingStates.State2DFAStates[splittingStateIndex] {
		for _, edge := range m.dfa.State2Edges[dfaState] {
			splitChar = edge.Char
			key := splitChar.Key()
			if _, ok := closeSet[key]; !ok {
				closeSet[key] = splitChar
				splitStates, split = m.splitByChar(splittingStates, splittingStateIndex, splitChar)
				if split {
					return
				}
			}
		}
	}
	splitStates = m.noSplit(splittingStates, splittingStateIndex)
	return
}

func (m *MinDFA) noSplit(splittingStates *SplitStates, splittingStatesIndex int) (nextSplitStates *SplitStates) {
	nextSplitStates = &SplitStates{
		DFAState2State: map[int]int{},
	}
	splittingDFAStates := splittingStates.State2DFAStates[splittingStatesIndex]
	nextSplitStates.State2DFAStates = append(nextSplitStates.State2DFAStates, splittingDFAStates)
	for _, dfaState := range splittingDFAStates {
		nextSplitStates.DFAState2State[dfaState] = 0
	}
	return
}

func (m *MinDFA) splitByChar(splittingStates *SplitStates, splittingStatesIndex int, splitChar Char) (nextSplitStates *SplitStates, split bool) {
	splittingStateIndex2NextSplitStatesIndex := map[int]int{}
	noCharStateIndex := -1
	nextSplitStates = &SplitStates{
		DFAState2State: map[int]int{},
	}

	for _, dfaState := range splittingStates.State2DFAStates[splittingStatesIndex] {
		splittingStateIndex := noCharStateIndex
		for _, edge := range m.dfa.State2Edges[dfaState] {
			if edge.Char.Equal(splitChar) {
				splittingStateIndex = splittingStates.DFAState2State[edge.ToState]
				m.appendNextSplitStates(nextSplitStates, dfaState, splittingStateIndex2NextSplitStatesIndex, splittingStateIndex)
				break
			}
		}
		if splittingStateIndex == noCharStateIndex {
			m.appendNextSplitStates(nextSplitStates, dfaState, splittingStateIndex2NextSplitStatesIndex, splittingStateIndex)
		}
	}
	split = len(nextSplitStates.State2DFAStates) > 1
	return
}

func (m *MinDFA) appendNextSplitStates(nextSplitStates *SplitStates, dfaState int, splittingStateIndex2NextSplitStateIndex map[int]int, splittingStateIndex int) {
	if nextSplitStateIndex, ok := splittingStateIndex2NextSplitStateIndex[splittingStateIndex]; ok {
		nextSplitStates.State2DFAStates[nextSplitStateIndex] = append(nextSplitStates.State2DFAStates[nextSplitStateIndex], dfaState)
		nextSplitStates.DFAState2State[dfaState] = nextSplitStateIndex
	} else {
		nextSplitStates.State2DFAStates = append(nextSplitStates.State2DFAStates, []int{dfaState})
		nextSplitStateIndex = len(nextSplitStates.State2DFAStates) - 1
		nextSplitStates.DFAState2State[dfaState] = nextSplitStateIndex
		splittingStateIndex2NextSplitStateIndex[splittingStateIndex] = nextSplitStateIndex
	}
}

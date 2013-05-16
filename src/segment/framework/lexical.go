package framework

import (
	"segment/dict"
	"sort"
)

const (
	Continue = 0
	Quit     = 1
	ElseQuit = 2
	End      = 3
)

// Lexical Function
const (
	None             = 0
	OutputIdentifier = 1
	DoSpace          = 2
	OutputSpace      = 3
	OutputNumeric    = 4
	OutputChinese    = 5
	Other            = 255
)

type State struct {
	NoFunction      bool
	Id              int
	Func            int
	IsQuitState     bool
	NextStateIdDict map[rune]int
	NextStateIds    []int
	ElseStateId     int
}

func NewState(id int, isQuit bool, function int, nextStateIdDict map[rune]int) (s *State) {
	s = &State{Id: id, IsQuitState: isQuit, Func: function, NextStateIdDict: nextStateIdDict}
	s.NoFunction = (s.Func == None)
	return
}

func NewStateNoDict(id int, isQuit bool, function int) (s *State) {
	return NewState(id, isQuit, function, nil)
}

func NewStateIdQuit(id int, isQuit bool) (s *State) {
	return NewState(id, isQuit, None, nil)
}

func NewStateNoFunc(id int, isQuit bool, nextStateIdDict map[rune]int) (s *State) {
	return NewState(id, isQuit, None, nextStateIdDict)
}

func NewStateId(id int) (s *State) {
	return NewState(id, false, None, nil)
}

func NewStateIdDict(id int, nextStateIdDict map[rune]int) (s *State) {
	return NewState(id, false, None, nextStateIdDict)
}

func (s *State) AddNextState(action rune, nextstate int) {
	if s.NextStateIdDict != nil {
		s.NextStateIdDict[action] = nextstate
	} else {
		if s.NextStateIds == nil {
			s.NextStateIds = make([]int, (int(action) + 1))
			for i := 0; i < len(s.NextStateIds); i++ {
				s.NextStateIds[i] = -1
			}
		} else {
			if len(s.NextStateIds) < int(action)+1 {
				old := s.NextStateIds
				s.NextStateIds = make([]int, int(action)+1)
				copy(s.NextStateIds, old)
				for i := len(old); i < len(s.NextStateIds); i++ {
					s.NextStateIds[i] = -1
				}
			}
		}
		s.NextStateIds[int(action)] = nextstate
	}
}

func (s *State) AddNextStateFromTo(beginAction rune, endAction rune, nextstate int) {
	for action := endAction; action >= beginAction; action-- {
		s.AddNextState(action, nextstate)
	}
}

type RuneSlice []rune

func (s RuneSlice) Len() int           { return len(s) }
func (s RuneSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s RuneSlice) Less(i, j int) bool { return s[i] < s[j] }

func (s *State) AddNextStateArr(actions []rune, nextstate int) {
	sort.Sort(RuneSlice(actions))
	for i := len(actions) - 1; i >= 0; i-- {
		s.AddNextState(actions[i], nextstate)
	}
}

func (s *State) AddElseState(nextstate int) {
	s.ElseStateId = nextstate
}

func (s *State) NextState(action rune) (nextstate int, isElseAction bool) {
	isElseAction = false
	nextstate = -1

	if action < 0 {
		isElseAction = true
		return s.ElseStateId, isElseAction
	}

	if s.NextStateIdDict != nil {
		if nextstate, ok := s.NextStateIdDict[action]; ok {
			if nextstate < 0 {
				isElseAction = true
				return s.ElseStateId, isElseAction
			} else {
				return nextstate, isElseAction
			}
		} else {
			isElseAction = true
			return s.ElseStateId, isElseAction
		}
	} else {
		if s.NextStateIds == nil {
			isElseAction = true
			return s.ElseStateId, isElseAction
		}

		if int(action) >= len(s.NextStateIds) {
			isElseAction = true
			return s.ElseStateId, isElseAction
		} else {
			nextstate = s.NextStateIds[action]
			if nextstate < 0 {
				isElseAction = true
				return s.ElseStateId, isElseAction
			} else {
				return nextstate, isElseAction
			}
		}
	}
	return
}

func (s *State) DoThings(action rune, dfa *Lexical) {
	switch s.Func {
	case OutputIdentifier:
		dfa.OutputToken = dict.NewWordInfoDefault()
		s.getTextElse(dfa)
		dfa.OutputToken.WordType = dict.TEnglish
	case OutputSpace:
		dfa.OutputToken = dict.NewWordInfoDefault()
		s.getTextElse(dfa)
		dfa.OutputToken.WordType = dict.TSpace
	case OutputNumeric:
		dfa.OutputToken = dict.NewWordInfoDefault()
		s.getTextElse(dfa)
		dfa.OutputToken.WordType = dict.TNumeric
	case OutputChinese:
		dfa.OutputToken = dict.NewWordInfoDefault()
		s.getTextElse(dfa)
		dfa.OutputToken.WordType = dict.TSimplifiedChinese
	case Other:
		dfa.OutputToken = dict.NewWordInfoDefault()
		s.getText(dfa)
		dfa.OutputToken.WordType = dict.TSymbol
	}
}

func (s *State) getTextElse(dfa *Lexical) {
	endIndex := dfa.CurrentToken
	dfa.OutputToken.Position = dfa.beginIndex
	dfa.OutputToken.Word = string(dfa.inputText[dfa.beginIndex:endIndex])
	dfa.beginIndex = endIndex
}

func (s *State) getText(dfa *Lexical) {
	endIndex := dfa.CurrentToken
	if endIndex == len(dfa.inputText) {
		dfa.OutputToken.Position = dfa.beginIndex
		dfa.OutputToken.Word = string(dfa.inputText[dfa.beginIndex:endIndex])
	} else {
		dfa.OutputToken.Position = dfa.beginIndex
		dfa.OutputToken.Word = string(dfa.inputText[dfa.beginIndex:(endIndex + 1)])
	}
	dfa.beginIndex = endIndex + 1
}

// static 
var states = []*State{}
var EofAction rune = 0
var s0 = addState(NewStateId(0))                        // Start state
var sother = addState(NewStateNoDict(255, true, Other)) // Start state

func init() {
	initDFAStates()
}

func initDFAStates() {
	initIdentifierStates()
	initSpaceStates()
	initNumericStates()
	initChineseStates()
	initOtherStates()
}

func initIdentifierStates() {
	s1 := addState(NewStateId(1))                             // Identifier begin state
	s2 := addState(NewStateNoDict(2, true, OutputIdentifier)) // Identifier quit state

	// s0 [_a-zA-Z] s1
	s0.AddNextState('_', s1.Id)
	s0.AddNextStateFromTo('a', 'z', s1.Id)
	s0.AddNextStateFromTo('A', 'Z', s1.Id)
	s0.AddNextStateFromTo('ａ', 'ｚ', s1.Id)
	s0.AddNextStateFromTo('Ａ', 'Ｚ', s1.Id)

	// s1 [_a-zA-Z0-9] s1
	s1.AddNextState('_', s1.Id)
	s1.AddNextStateFromTo('a', 'z', s1.Id)
	s1.AddNextStateFromTo('A', 'Z', s1.Id)
	s1.AddNextStateFromTo('0', '9', s1.Id)
	s1.AddNextStateFromTo('ａ', 'ｚ', s1.Id)
	s1.AddNextStateFromTo('Ａ', 'Ｚ', s1.Id)
	s1.AddNextStateFromTo('０', '９', s1.Id)

	// s1 ^[_z-zA-Z0-9] s2
	s1.AddElseState(s2.Id)
}

func initSpaceStates() {
	s3 := addState(NewStateIdQuit(3, false))             // Space begin state
	s4 := addState(NewStateNoDict(4, true, OutputSpace)) // Space quit state

	// s0 [ \t\r\n] s3
	s0.AddNextStateArr([]rune{' ', '\t', '\r', '\n'}, s3.Id)

	// s3 [ \t\r\n] s3
	s3.AddNextStateArr([]rune{' ', '\t', '\r', '\n'}, s3.Id)

	// s3 ^[ \t\r\n] s4
	s3.AddElseState(s4.Id)
}

func initNumericStates() {
	s5 := addState(NewStateIdQuit(5, false))               // Numberic begin state
	s6 := addState(NewStateIdQuit(6, false))               // Number dot state
	s7 := addState(NewStateNoDict(7, true, OutputNumeric)) // Number quit state

	// s0 [0-9] s5
	s0.AddNextStateFromTo('0', '9', s5.Id)
	s0.AddNextStateFromTo('０', '９', s5.Id)

	// s5 [0-9] s5
	s5.AddNextStateFromTo('0', '9', s5.Id)
	s5.AddNextStateFromTo('０', '９', s5.Id)

	// s5 [\.] s6
	s5.AddNextState('.', s6.Id)

	// s5 else s7 (integer)
	s5.AddElseState(s7.Id)

	// s6 [0-9] s6
	s6.AddNextStateFromTo('0', '9', s6.Id)
	s6.AddNextStateFromTo('０', '９', s6.Id)

	// s6 else s7 (float)
	s6.AddElseState(s7.Id)
}

func initChineseStates() {
	s8 := addState(NewStateIdQuit(8, false))               // Chinese begin state
	s9 := addState(NewStateNoDict(9, true, OutputChinese)) // Chinese quit state

	// s0 [4e00-9fa5] s5
	s0.AddNextStateFromTo('\u4e00', '\u9fa5', s8.Id)

	s8.AddNextStateFromTo('\u4e00', '\u9fa5', s8.Id)

	s8.AddElseState(s9.Id)
}

func initOtherStates() {
	s0.AddElseState(sother.Id)
}

func addState(state *State) *State {
	if state.Id >= len(states) {
		newLength := 0
		if state.Id < 2*len(states) {
			newLength = 2 * len(states)
		} else {
			newLength = state.Id + 1
		}
		oldStates := states
		states = make([]*State, newLength)
		copy(states, oldStates)
	}

	states[state.Id] = state
	return state
}

type Lexical struct {
	OldState     int
	CurrentState int
	QuitManuelly bool
	CurrentToken int

	beginIndex  int
	inputText   []rune
	OutputToken *dict.WordInfo
}

func NewLexical(runes []rune) *Lexical {
	return &Lexical{inputText: runes, OutputToken: nil}
}

func (l *Lexical) Input(action rune, token int) int {
	if len(states) == 0 {
		return Continue
	}

	l.CurrentToken = token

	if l.CurrentState == 0 && action == EofAction {
		return End
	}

	isElseAction := false
	l.OldState = l.CurrentState
	l.CurrentState, isElseAction = states[l.CurrentState].NextState(action)

	if !states[l.CurrentState].NoFunction {
		states[l.CurrentState].DoThings(action, l)
	}

	if !states[l.CurrentState].IsQuitState && !l.QuitManuelly {
		return Continue
	}

	l.QuitManuelly = false
	l.OldState = l.CurrentState
	l.CurrentState = 0
	if isElseAction {
		// Else action
		return ElseQuit
	}

	return Quit
}

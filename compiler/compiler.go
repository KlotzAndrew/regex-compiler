package compiler

func Re2post(re string) []rune {
	buffer := []rune{}

	ncat := 0
	nalt := 0

	type opstruct struct{ ncat int }

	opstack := []opstruct{}

	for _, ch := range re {
		switch ch {
		case '(':
			if ncat > 1 {
				ncat--
				buffer = append(buffer, '.')
			}

			opstack = append(opstack, opstruct{ncat: ncat})
			ncat = 0
		case ')':
			for ncat--; ncat > 0; ncat-- {
				buffer = append(buffer, '.')
			}

			op := opstack[len(opstack)-1]
			opstack = opstack[:len(opstack)-1]
			ncat = op.ncat

			ncat++
		case '|':
			for ncat--; ncat > 0; ncat-- {
				buffer = append(buffer, '.')
			}
			nalt++
		case '*', '+', '?':
			buffer = append(buffer, ch)
		default:
			if ncat > 1 {
				ncat--
				buffer = append(buffer, '.')
			}
			buffer = append(buffer, ch)
			ncat++
		}
	}

	for ncat--; ncat > 0; ncat-- {
		buffer = append(buffer, '.')
	}

	for ; nalt > 0; nalt-- {
		buffer = append(buffer, '|')
	}

	return buffer
}

type State struct {
	c         int
	s         rune
	out, out1 *State
	lastlist  int
}

type Frag struct {
	start *State
	out   []**State
}

type FragmentStack struct {
	frags []Frag
}

func NewFragmentStack() *FragmentStack { return &FragmentStack{} }
func (s *FragmentStack) Push(f Frag)   { s.frags = append(s.frags, f) }
func (s *FragmentStack) Pop() Frag {
	fs, f := s.frags[:len(s.frags)-1], s.frags[len(s.frags)-1]
	s.frags = fs
	return f
}

var Match = 1
var Split = 2
var matchstate = State{c: Match}

func patch(out []**State, s *State) {
	for _, p := range out {
		*p = s
	}
}

func Post2nfa(pf []rune) *State {
	stack := NewFragmentStack()
	j := len(pf)

	for i := 0; i < j; i++ {
		switch pf[i] {
		default:
			s := State{s: pf[i], out: nil, out1: nil}
			stack.Push(Frag{&s, []**State{&s.out}})
		case '.':
			fragTwo := stack.Pop()
			fragOne := stack.Pop()

			patch(fragOne.out, fragTwo.start)
			stack.Push(Frag{fragOne.start, fragTwo.out})
		case '*':
			frag := stack.Pop()
			s := State{c: Split, out: frag.start}
			patch(frag.out, &s)

			stack.Push(Frag{&s, []**State{&s.out1}})
		case '|':
			fragTwo := stack.Pop()
			fragOne := stack.Pop()

			s := State{c: Split, out: fragOne.start, out1: fragTwo.start}
			stack.Push(Frag{&s, append(fragOne.out, fragTwo.out...)})
		}
	}

	e := stack.Pop()

	patch(e.out, &matchstate)
	return e.start
}

func addState(list []*State, s *State, lid *int) []*State {
	s.lastlist = *lid
	if s.c == Split {
		list = addState(list, s.out, lid)
		list = addState(list, s.out1, lid)
	}
	list = append(list, s)
	return list
}

func step(clist []*State, ch rune, nlist []*State, lid *int) []*State {
	nlid := *lid + 1
	lid = &nlid

	nlist = nlist[:0]
	for _, s := range clist {
		if s.s == ch {
			nlist = addState(nlist, s.out, lid)
		}
	}
	return nlist
}

func isMatch(l []*State) bool {
	for _, s := range l {
		if s == &matchstate {
			return true
		}
	}
	return false
}

func MatchRe(start *State, s string) bool {
	var clist, nlist []*State

	lid := 0
	clist = addState(clist, start, &lid)

	for _, ch := range s {
		nlist = step(clist, ch, nlist, &lid)
		clist, nlist = nlist, clist
	}

	return isMatch(clist)
}

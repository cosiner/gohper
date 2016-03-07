package state

type State struct {
	state int

	cbs []func(now, old int)
}

func (s *State) Change(state int) int {
	orig := state
	s.state = state

	for _, cb := range s.cbs {
		cb(state, orig)
	}
	return orig
}

func (s *State) OnChange(fn func(int, int), replace bool) {
	if replace {
		s.cbs = s.cbs[:0]
	}
	s.cbs = append(s.cbs, fn)
}

func (s *State) Curr() int {
	return s.state
}

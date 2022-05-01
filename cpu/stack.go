package cpu

type stack struct {
	n int
	s []uint16
}

func (s *stack) push(pc uint16) {
	if s.n >= 16 {
		panic("stack overflow")
	}
	s.s = append(s.s, pc)
	s.n = len(s.s)
}

func (s *stack) pop() uint16 {
	item := s.s[s.n-1]
	s.s = s.s[:s.n-1]
	s.n = len(s.s)
	return item
}

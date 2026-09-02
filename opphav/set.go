package opphav

const maxSize = 256

type Set[M any] struct {
	domain  string
	mint    func(name string, index uint8) M
	members []M
	table   map[string]M
}

func NewSignSet[M Sign](domain string, wrap func(SignEntry) M) *Set[M] {
	return newSet(
		domain,
		func(name string, index uint8) M {
			return wrap(SignEntry{entry{name: name, index: index}})
		},
	)
}

func NewCategorySet[M Category](domain string, wrap func(CategoryEntry) M) *Set[M] {
	return newSet(
		domain,
		func(name string, index uint8) M {
			return wrap(CategoryEntry{entry{name: name, index: index}})
		},
	)
}

func NewSpectrumSet[M Spectrum](domain string, wrap func(SpectrumEntry) M) *Set[M] {
	return newSet(
		domain,
		func(name string, index uint8) M {
			return wrap(SpectrumEntry{entry{name: name, index: index}})
		},
	)
}

func newSet[M any](domain string, mint func(string, uint8) M) *Set[M] {
	if domain == "" {
		panic("opphav: a set needs a domain")
	}

	return &Set[M]{
		domain: domain,
		mint:   mint,
		table:  make(map[string]M),
	}
}

func (s *Set[M]) Define(name string) M {
	if name == "" {
		panic("opphav: member name must not be empty")
	}

	if len(s.members) == maxSize {
		panic("opphav: set is full")
	}

	qualified := s.domain + "." + name

	if _, taken := s.table[qualified]; taken {
		panic("opphav: duplicate member " + qualified)
	}

	member := s.mint(qualified, uint8(len(s.members)))

	s.members = append(s.members, member)
	s.table[qualified] = member

	return member
}

func (s *Set[M]) Lookup(qualified string) (M, bool) {
	member, ok := s.table[qualified]

	return member, ok
}

func (s *Set[M]) Members() []M {
	out := make([]M, len(s.members))
	copy(out, s.members)

	return out
}

func (s *Set[M]) At(index int) M {
	return s.members[index]
}

func (s *Set[M]) Size() int {
	return len(s.members)
}

func (s *Set[M]) Domain() string {
	return s.domain
}

package main

type Set struct {
	content map[string] bool
}

func NewSet() *Set {
	return &Set{
		content: make(map[string] bool),
	}
}

func (s *Set) Add(entry string)  {
	s.content[entry] = true
}

func (s *Set) Remove(entry string)  {
	delete(s.content, entry)
}

func (s *Set) Contains(entry string) bool {
	return s.content[entry]
}

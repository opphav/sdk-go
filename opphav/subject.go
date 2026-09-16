package opphav

import (
	"encoding/json"
	"strings"
)

type Subject struct {
	name     Name
	instance string
}

func NewSubject(name Name) Subject {
	if name.IsZero() {
		panic("opphav: subject name must not be zero")
	}

	return Subject{name: name}
}

func NewSubjectWithInstance(name Name, instance string) Subject {
	if name.IsZero() {
		panic("opphav: subject name must not be zero")
	}

	if instance == "" {
		panic("opphav: subject instance must not be empty")
	}

	return Subject{name: name, instance: instance}
}

func ParseSubject(s string) Subject {
	if s == "" {
		return Subject{}
	}

	name, instance, hasAt := strings.Cut(s, "@")
	if !hasAt || instance == "" {
		return NewSubject(NewName(name))
	}

	return NewSubjectWithInstance(NewName(name), instance)
}

func (s Subject) Name() Name {
	return s.name
}

func (s Subject) Instance() string {
	return s.instance
}

func (s Subject) IsZero() bool {
	return s.name.IsZero()
}

func (s Subject) String() string {
	if s.IsZero() {
		return ""
	}

	if s.instance == "" {
		return s.name.Path()
	}

	return s.name.Path() + "@" + s.instance
}

func (s Subject) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *Subject) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	*s = ParseSubject(str)
	return nil
}

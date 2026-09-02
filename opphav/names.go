package opphav

import (
	"encoding/json"
	"strings"
	"unique"
)

type node struct {
	part   string
	path   string
	parent unique.Handle[node]
}

type Name struct {
	handle unique.Handle[node]
}

func NewName(path string) Name {
	return Name{}.Extend(path)
}

func (n Name) IsZero() bool {
	return n.handle == (unique.Handle[node]{})
}

func (n Name) Extend(path string) Name {
	if path == "" {
		return n
	}

	current := n

	for {
		seg, rest, found := strings.Cut(path, ".")

		if seg == "" {
			panic("opphav: name segment must not be empty")
		}

		var full string
		if current.IsZero() {
			full = seg
		} else {
			full = current.handle.Value().path + "." + seg
		}

		current = Name{handle: unique.Make(node{seg, full, current.handle})}

		if !found {
			break
		}

		path = rest
	}

	return current
}

func (n Name) Path() string {
	if n.IsZero() {
		return ""
	}

	return n.handle.Value().path
}

func (n Name) Part() string {
	if n.IsZero() {
		return ""
	}

	return n.handle.Value().part
}

func (n Name) Parent() Name {
	if n.IsZero() {
		return Name{}
	}

	return Name{handle: n.handle.Value().parent}
}

func (n Name) String() string {
	return n.Path()
}

func (n Name) MarshalJSON() ([]byte, error) {
	return json.Marshal(n.Path())
}

func (n *Name) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*n = NewName(v)

	return nil
}

package field

import (
	"bytes"
	"fmt"
	"strconv"
)

// Path represents the path from some root to a particular field.
type Path struct {
	name   string
	index  string
	parent *Path
}

// NewPath creates a root Path object.
func NewPath(name string, moreNames ...string) *Path {
	r := &Path{name: name, parent: nil}
	for _, anotheName := range moreNames {
		r = &Path{name: anotheName, parent: r}
	}
	return r
}

// Root returns the root element of this Path.
func (p *Path) Root() *Path {
	for ; p.parent != nil; p = p.parent {
		//noop
	}
	return p
}

// Child creates a new Path that is a child of the method receiver.
func (p *Path) Child(name string, moreNames ...string) *Path {
	r := NewPath(name, moreNames...)
	r.Root().parent = p
	return r
}

// Index indicates that the previous Path is to be subscripted by an int.
// This sets the same underlying value as Key.
func (p *Path) Index(index int) *Path {
	return &Path{index: strconv.Itoa(index), parent: p}
}

// Key indicates that the previous Path is to be subscripted by a string.
// This sets the same underlying value as Index.
func (p *Path) Key(key string) *Path {
	return &Path{index: key, parent: p}
}

func (p *Path) String() string {
	if p == nil {
		return "<nil>"
	}
	elems := []*Path{}
	for ; p != nil; p = p.parent {
		elems = append(elems, p)
	}

	buf := bytes.NewBuffer(nil)
	for i := range elems {
		p := elems[len(elems)-1-i]
		if p.parent != nil && len(p.name) > 0 {
			buf.WriteString(".")
		}
		if len(p.name) > 0 {
			buf.WriteString(p.name)
		} else {
			fmt.Fprintf(buf, "[%s]", p.index)
		}
	}
	return buf.String()
}
